package main

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const createTableSQL = `
CREATE TABLE IF NOT EXISTS pm_activity_events (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  event_id CHAR(64) NOT NULL,
  user_wallet VARCHAR(64) NOT NULL,
  proxy_wallet VARCHAR(64) NOT NULL DEFAULT '',
  timestamp_ms BIGINT NOT NULL DEFAULT 0,
  event_time DATETIME(3) NULL,
  condition_id VARCHAR(128) NOT NULL DEFAULT '',
  activity_type VARCHAR(32) NOT NULL DEFAULT '',
  size DECIMAL(38,18) NULL,
  usdc_size DECIMAL(38,18) NULL,
  transaction_hash VARCHAR(128) NOT NULL DEFAULT '',
  price DECIMAL(38,18) NULL,
  asset VARCHAR(128) NOT NULL DEFAULT '',
  side VARCHAR(16) NOT NULL DEFAULT '',
  outcome_index INT NULL,
  title VARCHAR(1024) NOT NULL DEFAULT '',
  slug VARCHAR(255) NOT NULL DEFAULT '',
  event_slug VARCHAR(255) NOT NULL DEFAULT '',
  outcome VARCHAR(128) NOT NULL DEFAULT '',
  source_offset INT NOT NULL DEFAULT 0,
  source_index INT NOT NULL DEFAULT 0,
  pulled_at DATETIME(3) NOT NULL,
  raw_json JSON NOT NULL,
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (id),
  UNIQUE KEY uq_event_id (event_id),
  KEY idx_user_ts (user_wallet, timestamp_ms DESC),
  KEY idx_slug_ts (slug, timestamp_ms DESC),
  KEY idx_tx_hash (transaction_hash)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
`

type MySQLStore struct {
	db *sql.DB
}

type EventQuery struct {
	Limit        int
	Offset       int
	Slug         string
	ActivityType string
	Side         string
}

func NewMySQLStore(dsn string, maxIdle, maxOpen int, maxLife time.Duration) (*MySQLStore, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("open mysql failed: %w", err)
	}
	db.SetMaxIdleConns(maxIdle)
	db.SetMaxOpenConns(maxOpen)
	db.SetConnMaxLifetime(maxLife)
	if err = db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping mysql failed: %w", err)
	}
	return &MySQLStore{db: db}, nil
}

func (s *MySQLStore) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *MySQLStore) EnsureSchema(ctx context.Context) error {
	if s == nil || s.db == nil {
		return fmt.Errorf("mysql store is nil")
	}
	_, err := s.db.ExecContext(ctx, createTableSQL)
	if err != nil {
		return fmt.Errorf("create table failed: %w", err)
	}
	return nil
}

func (s *MySQLStore) InsertIgnoreBatch(ctx context.Context, events []ActivityEvent) (int, int, error) {
	if len(events) == 0 {
		return 0, 0, nil
	}
	const maxBatch = 200
	insertedTotal := 0
	for start := 0; start < len(events); start += maxBatch {
		end := start + maxBatch
		if end > len(events) {
			end = len(events)
		}
		inserted, err := s.insertBatchOnce(ctx, events[start:end])
		if err != nil {
			return insertedTotal, len(events) - insertedTotal, err
		}
		insertedTotal += inserted
	}
	return insertedTotal, len(events) - insertedTotal, nil
}

func (s *MySQLStore) insertBatchOnce(ctx context.Context, events []ActivityEvent) (int, error) {
	if len(events) == 0 {
		return 0, nil
	}
	valueSQL := "(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	sqlBuilder := strings.Builder{}
	sqlBuilder.WriteString("INSERT IGNORE INTO pm_activity_events (")
	sqlBuilder.WriteString("event_id,user_wallet,proxy_wallet,timestamp_ms,event_time,condition_id,activity_type,size,usdc_size,transaction_hash,price,asset,side,outcome_index,title,slug,event_slug,outcome,source_offset,source_index,pulled_at,raw_json,created_at")
	sqlBuilder.WriteString(") VALUES ")
	args := make([]interface{}, 0, len(events)*23)
	now := time.Now().UTC()
	for i := range events {
		if i > 0 {
			sqlBuilder.WriteByte(',')
		}
		sqlBuilder.WriteString(valueSQL)
		e := events[i]
		args = append(args,
			e.EventID,
			e.UserWallet,
			e.ProxyWallet,
			e.TimestampMs,
			nullableTime(e.EventTime),
			e.ConditionID,
			e.ActivityType,
			nullableDecimal(e.Size),
			nullableDecimal(e.USDCSize),
			e.TransactionHash,
			nullableDecimal(e.Price),
			e.Asset,
			e.Side,
			nullableInt(e.OutcomeIndex),
			e.Title,
			e.Slug,
			e.EventSlug,
			e.Outcome,
			e.SourceOffset,
			e.SourceIndex,
			e.PulledAt.UTC(),
			e.RawJSON,
			now,
		)
	}
	res, err := s.db.ExecContext(ctx, sqlBuilder.String(), args...)
	if err != nil {
		return 0, fmt.Errorf("insert batch failed: %w", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("read rows affected failed: %w", err)
	}
	return int(affected), nil
}

func (s *MySQLStore) Count(ctx context.Context) (int64, error) {
	var total int64
	err := s.db.QueryRowContext(ctx, "SELECT COUNT(1) FROM pm_activity_events").Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("count failed: %w", err)
	}
	return total, nil
}

func (s *MySQLStore) QueryEvents(ctx context.Context, q EventQuery) ([]ActivityEvent, error) {
	base := "SELECT event_id,user_wallet,proxy_wallet,timestamp_ms,event_time,condition_id,activity_type,size,usdc_size,transaction_hash,price,asset,side,outcome_index,title,slug,event_slug,outcome,raw_json,source_offset,source_index,pulled_at FROM pm_activity_events WHERE 1=1"
	args := make([]interface{}, 0, 8)
	if q.Slug != "" {
		base += " AND slug = ?"
		args = append(args, q.Slug)
	}
	if q.ActivityType != "" {
		base += " AND activity_type = ?"
		args = append(args, strings.ToUpper(strings.TrimSpace(q.ActivityType)))
	}
	if q.Side != "" {
		base += " AND side = ?"
		args = append(args, strings.ToLower(strings.TrimSpace(q.Side)))
	}
	base += " ORDER BY timestamp_ms DESC, id DESC LIMIT ? OFFSET ?"
	args = append(args, q.Limit, q.Offset)

	rows, err := s.db.QueryContext(ctx, base, args...)
	if err != nil {
		return nil, fmt.Errorf("query events failed: %w", err)
	}
	defer func() { _ = rows.Close() }()

	out := make([]ActivityEvent, 0, q.Limit)
	for rows.Next() {
		var e ActivityEvent
		var eventTime sql.NullTime
		var size sql.NullString
		var usdcSize sql.NullString
		var price sql.NullString
		var outcomeIndex sql.NullInt64
		if err = rows.Scan(
			&e.EventID,
			&e.UserWallet,
			&e.ProxyWallet,
			&e.TimestampMs,
			&eventTime,
			&e.ConditionID,
			&e.ActivityType,
			&size,
			&usdcSize,
			&e.TransactionHash,
			&price,
			&e.Asset,
			&e.Side,
			&outcomeIndex,
			&e.Title,
			&e.Slug,
			&e.EventSlug,
			&e.Outcome,
			&e.RawJSON,
			&e.SourceOffset,
			&e.SourceIndex,
			&e.PulledAt,
		); err != nil {
			return nil, fmt.Errorf("scan event failed: %w", err)
		}
		if eventTime.Valid {
			e.EventTime = eventTime.Time.UTC()
		}
		if size.Valid {
			e.Size = size.String
		}
		if usdcSize.Valid {
			e.USDCSize = usdcSize.String
		}
		if price.Valid {
			e.Price = price.String
		}
		if outcomeIndex.Valid {
			tmp := int(outcomeIndex.Int64)
			e.OutcomeIndex = &tmp
		}
		out = append(out, e)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rows failed: %w", err)
	}
	return out, nil
}

func nullableDecimal(raw string) interface{} {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	return raw
}

func nullableTime(t time.Time) interface{} {
	if t.IsZero() {
		return nil
	}
	return t.UTC()
}

func nullableInt(v *int) interface{} {
	if v == nil {
		return nil
	}
	return *v
}
