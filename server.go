package main

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func NewHTTPServer(cfg *Config, poller *Poller, store *MySQLStore) *server.Hertz {
	h := server.Default(server.WithHostPorts(cfg.ListenAddr))

	h.GET("/", func(ctx context.Context, c *app.RequestContext) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(dashboardHTML))
	})

	h.GET("/healthz", func(ctx context.Context, c *app.RequestContext) {
		writeOK(c, map[string]interface{}{
			"ok": true,
		})
	})

	h.GET("/api/v1/status", func(ctx context.Context, c *app.RequestContext) {
		total, err := store.Count(ctx)
		if err != nil {
			writeErr(c, http.StatusInternalServerError, err.Error())
			return
		}
		stats := poller.Snapshot()
		writeOK(c, map[string]interface{}{
			"config": map[string]interface{}{
				"listen_addr":       cfg.ListenAddr,
				"user_wallet":       cfg.TargetUser,
				"page_limit":        cfg.PageLimit,
				"fast_interval_sec": int(cfg.FastInterval.Seconds()),
				"backfill_interval": int(cfg.BackfillInterval.Seconds()),
				"request_gap_ms":    cfg.RequestGap.Milliseconds(),
				"backfill_offsets":  cfg.BackfillOffsets,
			},
			"poller":        stats,
			"stored_events": total,
		})
	})

	h.POST("/api/v1/sync/once", func(ctx context.Context, c *app.RequestContext) {
		var req struct {
			Mode string `json:"mode"`
		}
		if len(c.Request.Body()) > 0 {
			if err := c.BindJSON(&req); err != nil {
				writeErr(c, http.StatusBadRequest, err.Error())
				return
			}
		}
		mode := strings.ToLower(strings.TrimSpace(req.Mode))
		if mode == "" {
			mode = "backfill"
		}
		if mode != "fast" && mode != "backfill" {
			writeErr(c, http.StatusBadRequest, "mode must be fast/backfill")
			return
		}
		ok := poller.Enqueue(syncJob{
			Mode:   mode,
			Reason: "api",
		})
		writeOK(c, map[string]interface{}{
			"queued": ok,
			"mode":   mode,
		})
	})

	h.GET("/api/v1/events", func(ctx context.Context, c *app.RequestContext) {
		limit := parsePositiveInt(c.Query("limit"), cfg.DefaultQueryLimit)
		if limit <= 0 {
			limit = cfg.DefaultQueryLimit
		}
		if limit > cfg.MaxQueryLimit {
			limit = cfg.MaxQueryLimit
		}

		offset := parsePositiveInt(c.Query("offset"), -1)
		page := parsePositiveInt(c.Query("page"), 1)
		if page <= 0 {
			page = 1
		}
		if offset < 0 {
			offset = (page - 1) * limit
		}

		query := EventQuery{
			Limit:        limit,
			Offset:       offset,
			Slug:         strings.TrimSpace(c.Query("slug")),
			MarketTag:    strings.TrimSpace(c.Query("tag")),
			ActivityType: strings.TrimSpace(c.Query("type")),
			Side:         strings.TrimSpace(c.Query("side")),
		}

		total, err := store.CountEvents(ctx, query)
		if err != nil {
			writeErr(c, http.StatusInternalServerError, err.Error())
			return
		}
		rows, err := store.QueryEvents(ctx, query)
		if err != nil {
			writeErr(c, http.StatusInternalServerError, err.Error())
			return
		}

		writeOK(c, map[string]interface{}{
			"limit":       limit,
			"offset":      offset,
			"page":        page,
			"page_size":   limit,
			"total":       total,
			"total_pages": calcTotalPages(total, limit),
			"items":       rows,
		})
	})

	h.GET("/api/v1/events/strategy-group", func(ctx context.Context, c *app.RequestContext) {
		symbol := strings.TrimSpace(c.Query("symbol"))
		startSec := parsePositiveInt64(c.Query("start_sec"), 0)
		endSec := parsePositiveInt64(c.Query("end_sec"), 0)
		limit := parsePositiveInt(c.Query("limit"), 4000)
		if symbol == "" {
			writeErr(c, http.StatusBadRequest, "symbol is required")
			return
		}
		if startSec <= 0 || endSec <= startSec {
			writeErr(c, http.StatusBadRequest, "start_sec/end_sec is invalid")
			return
		}

		rows, err := store.QueryEventsByStrategyGroup(ctx, StrategyGroupQuery{
			Symbol:   symbol,
			StartSec: startSec,
			EndSec:   endSec,
			Limit:    limit,
		})
		if err != nil {
			writeErr(c, http.StatusInternalServerError, err.Error())
			return
		}

		writeOK(c, map[string]interface{}{
			"symbol":    strings.ToLower(symbol),
			"start_sec": startSec,
			"end_sec":   endSec,
			"items":     rows,
			"count":     len(rows),
		})
	})

	h.GET("/api/v1/slugs", func(ctx context.Context, c *app.RequestContext) {
		pageSize := parsePositiveInt(c.Query("page_size"), 12)
		if pageSize <= 0 {
			pageSize = 12
		}
		if pageSize > 100 {
			pageSize = 100
		}
		page := parsePositiveInt(c.Query("page"), 1)
		if page <= 0 {
			page = 1
		}
		offset := (page - 1) * pageSize

		keyword := strings.TrimSpace(c.Query("keyword"))
		tag := strings.TrimSpace(c.Query("tag"))
		items, total, err := store.QuerySlugSummaries(ctx, SlugSummaryQuery{
			Limit:   pageSize,
			Offset:  offset,
			Keyword: keyword,
			Tag:     tag,
		})
		if err != nil {
			writeErr(c, http.StatusInternalServerError, err.Error())
			return
		}
		writeOK(c, map[string]interface{}{
			"page":        page,
			"page_size":   pageSize,
			"offset":      offset,
			"total":       total,
			"total_pages": calcTotalPages(total, pageSize),
			"keyword":     keyword,
			"tag":         tag,
			"items":       items,
		})
	})

	return h
}

func parsePositiveInt(raw string, fallback int) int {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return fallback
	}
	v, err := strconv.Atoi(raw)
	if err != nil || v < 0 {
		return fallback
	}
	return v
}

func parsePositiveInt64(raw string, fallback int64) int64 {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return fallback
	}
	v, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || v < 0 {
		return fallback
	}
	return v
}

func calcTotalPages(total int64, pageSize int) int {
	if total <= 0 || pageSize <= 0 {
		return 0
	}
	return int((total + int64(pageSize) - 1) / int64(pageSize))
}

func writeOK(c *app.RequestContext, data interface{}) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    20000,
		"message": "success",
		"data":    data,
	})
}

func writeErr(c *app.RequestContext, status int, msg string) {
	c.JSON(status, map[string]interface{}{
		"code":    status,
		"message": msg,
		"data":    nil,
	})
}
