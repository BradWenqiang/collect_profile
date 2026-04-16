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
		if limit > cfg.MaxQueryLimit {
			limit = cfg.MaxQueryLimit
		}
		offset := parsePositiveInt(c.Query("offset"), 0)

		rows, err := store.QueryEvents(ctx, EventQuery{
			Limit:        limit,
			Offset:       offset,
			Slug:         strings.TrimSpace(c.Query("slug")),
			ActivityType: strings.TrimSpace(c.Query("type")),
			Side:         strings.TrimSpace(c.Query("side")),
		})
		if err != nil {
			writeErr(c, http.StatusInternalServerError, err.Error())
			return
		}
		writeOK(c, map[string]interface{}{
			"limit":  limit,
			"offset": offset,
			"items":  rows,
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
