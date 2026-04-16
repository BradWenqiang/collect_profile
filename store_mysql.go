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
  market_tag VARCHAR(24) NOT NULL DEFAULT '',
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
  KEY idx_market_tag_ts (market_tag, timestamp_ms DESC),
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
	MarketTag    string
	ActivityType string
	Side         string
}

type SlugSummaryQuery struct {
	Limit   int
	Offset  int
	Keyword string
	Tag     string
}

type SlugSummary struct {
	MarketTag        string `json:"market_tag"`
	Slug             string `json:"slug"`
	EventCount       int64  `json:"event_count"`
	FirstTimestampMs int64  `json:"first_timestamp_ms"`
	LastTimestampMs  int64  `json:"last_timestamp_ms"`
	BuyCount         int64  `json:"buy_count"`
	SellCount        int64  `json:"sell_count"`
	UpCount          int64  `json:"up_count"`
	DownCount        int64  `json:"down_count"`
}

type StrategyGroupQuery struct {
	Symbol   string
	StartSec int64
	EndSec   int64
	Limit    int
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
	if err = s.ensureMarketTagSchema(ctx); err != nil {
		return err
	}
	return nil
}

func (s *MySQLStore) ensureMarketTagSchema(ctx context.Context) error {
	if _, err := s.db.ExecContext(ctx, "ALTER TABLE pm_activity_events ADD COLUMN market_tag VARCHAR(24) NOT NULL DEFAULT '' AFTER slug"); err != nil {
		if !strings.Contains(strings.ToLower(err.Error()), "duplicate column name") {
			return fmt.Errorf("ensure market_tag column failed: %w", err)
		}
	}
	if _, err := s.db.ExecContext(ctx, "ALTER TABLE pm_activity_events ADD INDEX idx_market_tag_ts (market_tag, timestamp_ms DESC)"); err != nil {
		msg := strings.ToLower(err.Error())
		if !strings.Contains(msg, "duplicate key name") && !strings.Contains(msg, "already exists") {
			return fmt.Errorf("ensure market_tag index failed: %w", err)
		}
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
	valueSQL := "(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	sqlBuilder := strings.Builder{}
	sqlBuilder.WriteString("INSERT IGNORE INTO pm_activity_events (")
	sqlBuilder.WriteString("event_id,user_wallet,proxy_wallet,timestamp_ms,event_time,condition_id,activity_type,size,usdc_size,transaction_hash,price,asset,side,outcome_index,title,slug,market_tag,event_slug,outcome,source_offset,source_index,pulled_at,raw_json,created_at")
	sqlBuilder.WriteString(") VALUES ")
	args := make([]interface{}, 0, len(events)*24)
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
			normalizeMarketTag(e.MarketTag),
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

func (s *MySQLStore) CountEvents(ctx context.Context, q EventQuery) (int64, error) {
	sqlText, args := buildEventQueryBase("SELECT COUNT(1)", q)
	var total int64
	err := s.db.QueryRowContext(ctx, sqlText, args...).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("count events failed: %w", err)
	}
	return total, nil
}

func (s *MySQLStore) QueryEvents(ctx context.Context, q EventQuery) ([]ActivityEvent, error) {
	if q.Limit <= 0 {
		q.Limit = 100
	}
	if q.Offset < 0 {
		q.Offset = 0
	}

	sqlText, args := buildEventQueryBase("SELECT event_id,user_wallet,proxy_wallet,timestamp_ms,event_time,condition_id,activity_type,size,usdc_size,transaction_hash,price,asset,side,outcome_index,title,slug,market_tag,event_slug,outcome,raw_json,source_offset,source_index,pulled_at", q)
	sqlText += " ORDER BY timestamp_ms DESC, id DESC LIMIT ? OFFSET ?"
	args = append(args, q.Limit, q.Offset)

	rows, err := s.db.QueryContext(ctx, sqlText, args...)
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
			&e.MarketTag,
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

func (s *MySQLStore) QuerySlugSummaries(ctx context.Context, q SlugSummaryQuery) ([]SlugSummary, int64, error) {
	if q.Limit <= 0 {
		q.Limit = 12
	}
	if q.Offset < 0 {
		q.Offset = 0
	}

	whereSQL, args := buildSlugWhere(q.Keyword, q.Tag)
	countSQL := "SELECT COUNT(1) FROM (SELECT market_tag, slug FROM pm_activity_events" + whereSQL + " GROUP BY market_tag, slug) t"
	var total int64
	if err := s.db.QueryRowContext(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count slugs failed: %w", err)
	}
	if total == 0 {
		return []SlugSummary{}, 0, nil
	}

	listSQL := "SELECT market_tag, slug, " +
		"COUNT(1) AS event_count, " +
		"COALESCE(MIN(timestamp_ms), 0) AS first_ts, " +
		"COALESCE(MAX(timestamp_ms), 0) AS last_ts, " +
		"SUM(CASE WHEN side = 'buy' THEN 1 ELSE 0 END) AS buy_count, " +
		"SUM(CASE WHEN side = 'sell' THEN 1 ELSE 0 END) AS sell_count, " +
		"SUM(CASE WHEN (LOWER(outcome) LIKE '%up%' OR LOWER(outcome) = 'yes' OR LOWER(outcome) LIKE '%long%' OR LOWER(outcome) LIKE '%higher%') THEN 1 ELSE 0 END) AS up_count, " +
		"SUM(CASE WHEN (LOWER(outcome) LIKE '%down%' OR LOWER(outcome) = 'no' OR LOWER(outcome) LIKE '%short%' OR LOWER(outcome) LIKE '%lower%') THEN 1 ELSE 0 END) AS down_count " +
		"FROM pm_activity_events" + whereSQL + " GROUP BY market_tag, slug ORDER BY CASE market_tag WHEN 'btc' THEN 1 WHEN 'eth' THEN 2 WHEN 'sol' THEN 3 ELSE 9 END, last_ts DESC LIMIT ? OFFSET ?"

	listArgs := make([]interface{}, 0, len(args)+2)
	listArgs = append(listArgs, args...)
	listArgs = append(listArgs, q.Limit, q.Offset)

	rows, err := s.db.QueryContext(ctx, listSQL, listArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("query slug summaries failed: %w", err)
	}
	defer func() { _ = rows.Close() }()

	items := make([]SlugSummary, 0, q.Limit)
	for rows.Next() {
		var row SlugSummary
		if err = rows.Scan(
			&row.MarketTag,
			&row.Slug,
			&row.EventCount,
			&row.FirstTimestampMs,
			&row.LastTimestampMs,
			&row.BuyCount,
			&row.SellCount,
			&row.UpCount,
			&row.DownCount,
		); err != nil {
			return nil, 0, fmt.Errorf("scan slug summary failed: %w", err)
		}
		items = append(items, row)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate slug summaries failed: %w", err)
	}
	return items, total, nil
}

func (s *MySQLStore) QueryEventsByStrategyGroup(ctx context.Context, q StrategyGroupQuery) ([]ActivityEvent, error) {
	q.Symbol = strings.ToLower(strings.TrimSpace(q.Symbol))
	if q.Symbol == "" {
		return nil, fmt.Errorf("empty symbol")
	}
	if q.EndSec <= q.StartSec {
		return nil, fmt.Errorf("invalid window: start=%d end=%d", q.StartSec, q.EndSec)
	}
	if q.Limit <= 0 {
		q.Limit = 4000
	}
	if q.Limit > 10000 {
		q.Limit = 10000
	}

	querySQL := "SELECT event_id,user_wallet,proxy_wallet,timestamp_ms,event_time,condition_id,activity_type,size,usdc_size,transaction_hash,price,asset,side,outcome_index,title,slug,market_tag,event_slug,outcome,raw_json,source_offset,source_index,pulled_at " +
		"FROM pm_activity_events WHERE slug <> '' " +
		"AND market_tag = ? " +
		"AND (slug LIKE '%-5m-%' OR slug LIKE '%-15m-%') " +
		"AND CAST(SUBSTRING_INDEX(slug, '-', -1) AS UNSIGNED) > ? " +
		"AND CAST(SUBSTRING_INDEX(slug, '-', -1) AS UNSIGNED) <= ? " +
		"ORDER BY timestamp_ms DESC, id DESC LIMIT ?"

	rows, err := s.db.QueryContext(ctx, querySQL, normalizeMarketTag(q.Symbol), q.StartSec, q.EndSec, q.Limit)
	if err != nil {
		return nil, fmt.Errorf("query strategy group events failed: %w", err)
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
			&e.MarketTag,
			&e.EventSlug,
			&e.Outcome,
			&e.RawJSON,
			&e.SourceOffset,
			&e.SourceIndex,
			&e.PulledAt,
		); err != nil {
			return nil, fmt.Errorf("scan strategy group event failed: %w", err)
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
		return nil, fmt.Errorf("iterate strategy group events failed: %w", err)
	}
	return out, nil
}

func buildEventQueryBase(selectClause string, q EventQuery) (string, []interface{}) {
	builder := strings.Builder{}
	builder.WriteString(selectClause)
	builder.WriteString(" FROM pm_activity_events WHERE 1=1")

	args := make([]interface{}, 0, 4)
	if q.Slug != "" {
		builder.WriteString(" AND slug = ?")
		args = append(args, strings.TrimSpace(q.Slug))
	}
	if q.MarketTag != "" {
		builder.WriteString(" AND market_tag = ?")
		args = append(args, normalizeMarketTag(q.MarketTag))
	}
	if q.ActivityType != "" {
		builder.WriteString(" AND activity_type = ?")
		args = append(args, strings.ToUpper(strings.TrimSpace(q.ActivityType)))
	}
	if q.Side != "" {
		builder.WriteString(" AND side = ?")
		args = append(args, strings.ToLower(strings.TrimSpace(q.Side)))
	}
	return builder.String(), args
}

func buildSlugWhere(keyword, tag string) (string, []interface{}) {
	builder := strings.Builder{}
	builder.WriteString(" WHERE slug <> ''")
	args := make([]interface{}, 0, 2)

	tag = normalizeMarketTag(tag)
	if tag != "" && tag != "all" {
		builder.WriteString(" AND market_tag = ?")
		args = append(args, tag)
	}

	keyword = strings.TrimSpace(keyword)
	if keyword != "" {
		like := "%" + keyword + "%"
		builder.WriteString(" AND (slug LIKE ? OR title LIKE ?)")
		args = append(args, like, like)
	}
	return builder.String(), args
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

func normalizeMarketTag(tag string) string {
	tag = strings.ToLower(strings.TrimSpace(tag))
	switch tag {
	case "btc", "bitcoin":
		return "btc"
	case "eth", "ethereum":
		return "eth"
	case "sol", "solana":
		return "sol"
	case "all":
		return "all"
	case "":
		return ""
	default:
		return "other"
	}
}
