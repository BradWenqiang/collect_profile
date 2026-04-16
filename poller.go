package main

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"
	"time"
)

type syncJob struct {
	Mode   string
	Reason string
}

type PollerStats struct {
	StartedAt          time.Time `json:"started_at"`
	LastJobMode        string    `json:"last_job_mode"`
	LastJobReason      string    `json:"last_job_reason"`
	LastJobStartedAt   time.Time `json:"last_job_started_at"`
	LastJobFinishedAt  time.Time `json:"last_job_finished_at"`
	LastJobDurationMs  int64     `json:"last_job_duration_ms"`
	LastSuccessAt      time.Time `json:"last_success_at"`
	LastErrorAt        time.Time `json:"last_error_at"`
	LastError          string    `json:"last_error"`
	LastFetchedRows    int       `json:"last_fetched_rows"`
	LastInsertedRows   int       `json:"last_inserted_rows"`
	LastDuplicateRows  int       `json:"last_duplicate_rows"`
	LastNewestTS       int64     `json:"last_newest_ts"`
	LastOldestTS       int64     `json:"last_oldest_ts"`
	QueueDropped       int64     `json:"queue_dropped"`
	QueueLen           int       `json:"queue_len"`
	Running            bool      `json:"running"`
	TotalFetchedRows   int64     `json:"total_fetched_rows"`
	TotalInsertedRows  int64     `json:"total_inserted_rows"`
	TotalDuplicateRows int64     `json:"total_duplicate_rows"`
}

type Poller struct {
	cfg    *Config
	client *ActivityClient
	store  *MySQLStore

	jobCh chan syncJob

	mu    sync.Mutex
	stats PollerStats
}

func NewPoller(cfg *Config, client *ActivityClient, store *MySQLStore) *Poller {
	return &Poller{
		cfg:    cfg,
		client: client,
		store:  store,
		jobCh:  make(chan syncJob, 8),
		stats: PollerStats{
			StartedAt: time.Now().UTC(),
		},
	}
}

func (p *Poller) Run(ctx context.Context) {
	if p == nil {
		return
	}
	go p.worker(ctx)
	if p.cfg.StartupBackfill {
		p.Enqueue(syncJob{Mode: "backfill", Reason: "startup"})
	}
	p.Enqueue(syncJob{Mode: "fast", Reason: "startup"})

	fastTicker := time.NewTicker(p.cfg.FastInterval)
	backfillTicker := time.NewTicker(p.cfg.BackfillInterval)
	defer fastTicker.Stop()
	defer backfillTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-fastTicker.C:
			p.Enqueue(syncJob{Mode: "fast", Reason: "ticker"})
		case <-backfillTicker.C:
			p.Enqueue(syncJob{Mode: "backfill", Reason: "ticker"})
		}
	}
}

func (p *Poller) Enqueue(job syncJob) bool {
	if p == nil {
		return false
	}
	select {
	case p.jobCh <- job:
		return true
	default:
		p.mu.Lock()
		p.stats.QueueDropped++
		p.mu.Unlock()
		return false
	}
}

func (p *Poller) Snapshot() PollerStats {
	p.mu.Lock()
	defer p.mu.Unlock()
	s := p.stats
	s.QueueLen = len(p.jobCh)
	return s
}

func (p *Poller) worker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case job := <-p.jobCh:
			p.runJob(ctx, job)
		}
	}
}

func (p *Poller) runJob(ctx context.Context, job syncJob) {
	startedAt := time.Now().UTC()
	p.mu.Lock()
	p.stats.Running = true
	p.stats.LastJobMode = job.Mode
	p.stats.LastJobReason = job.Reason
	p.stats.LastJobStartedAt = startedAt
	p.mu.Unlock()

	fetched, inserted, duplicated, newestTS, oldestTS, err := p.syncOnce(ctx, job)
	finishedAt := time.Now().UTC()

	p.mu.Lock()
	defer p.mu.Unlock()
	p.stats.Running = false
	p.stats.LastJobFinishedAt = finishedAt
	p.stats.LastJobDurationMs = finishedAt.Sub(startedAt).Milliseconds()
	p.stats.LastFetchedRows = fetched
	p.stats.LastInsertedRows = inserted
	p.stats.LastDuplicateRows = duplicated
	p.stats.LastNewestTS = newestTS
	p.stats.LastOldestTS = oldestTS
	p.stats.TotalFetchedRows += int64(fetched)
	p.stats.TotalInsertedRows += int64(inserted)
	p.stats.TotalDuplicateRows += int64(duplicated)
	if err != nil {
		p.stats.LastErrorAt = finishedAt
		p.stats.LastError = err.Error()
		log.Printf("[pm_activity] sync failed mode=%s reason=%s fetched=%d inserted=%d dup=%d err=%v", job.Mode, job.Reason, fetched, inserted, duplicated, err)
		return
	}
	p.stats.LastSuccessAt = finishedAt
	p.stats.LastError = ""
	log.Printf("[pm_activity] sync ok mode=%s reason=%s fetched=%d inserted=%d dup=%d newest=%d oldest=%d", job.Mode, job.Reason, fetched, inserted, duplicated, newestTS, oldestTS)
}

func (p *Poller) syncOnce(ctx context.Context, job syncJob) (int, int, int, int64, int64, error) {
	offsets := []int{0}
	if strings.EqualFold(job.Mode, "backfill") {
		offsets = append(offsets[:0], p.cfg.BackfillOffsets...)
		sort.Ints(offsets)
	}
	allFetched := 0
	allInserted := 0
	allDup := 0
	var newestTS int64
	var oldestTS int64

	for i, offset := range offsets {
		rows, err := p.fetchWithRetry(ctx, p.cfg.PageLimit, offset)
		if err != nil {
			return allFetched, allInserted, allDup, newestTS, oldestTS, fmt.Errorf("fetch offset=%d failed: %w", offset, err)
		}
		allFetched += len(rows)
		pulledAt := time.Now().UTC()
		events := make([]ActivityEvent, 0, len(rows))
		for idx := range rows {
			event, err := parseActivityEvent(rows[idx], p.cfg.TargetUser, offset, idx, pulledAt)
			if err != nil {
				continue
			}
			events = append(events, event)
			if event.TimestampMs > 0 {
				if newestTS == 0 || event.TimestampMs > newestTS {
					newestTS = event.TimestampMs
				}
				if oldestTS == 0 || event.TimestampMs < oldestTS {
					oldestTS = event.TimestampMs
				}
			}
		}
		inserted, dup, err := p.store.InsertIgnoreBatch(ctx, events)
		if err != nil {
			return allFetched, allInserted, allDup, newestTS, oldestTS, fmt.Errorf("insert offset=%d failed: %w", offset, err)
		}
		allInserted += inserted
		allDup += dup

		// 数据尾页提前退出，避免无意义请求。
		if strings.EqualFold(job.Mode, "backfill") && len(rows) < p.cfg.PageLimit {
			break
		}
		if i < len(offsets)-1 {
			select {
			case <-ctx.Done():
				return allFetched, allInserted, allDup, newestTS, oldestTS, ctx.Err()
			case <-time.After(p.cfg.RequestGap):
			}
		}
	}

	return allFetched, allInserted, allDup, newestTS, oldestTS, nil
}

func (p *Poller) fetchWithRetry(ctx context.Context, limit, offset int) ([]RawActivity, error) {
	var lastErr error
	for i := 0; i < p.cfg.FetchRetry; i++ {
		rows, err := p.client.FetchPage(ctx, limit, offset)
		if err == nil {
			return rows, nil
		}
		lastErr = err
		backoff := time.Duration(i+1) * 700 * time.Millisecond
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(backoff):
		}
	}
	return nil, lastErr
}
