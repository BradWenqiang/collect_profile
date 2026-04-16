package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	ListenAddr       string
	DataAPIBaseURL   string
	TargetUser       string
	PageLimit        int
	FastInterval     time.Duration
	BackfillInterval time.Duration
	RequestGap       time.Duration
	HTTPTimeout      time.Duration
	FetchRetry       int
	StartupBackfill  bool
	BackfillOffsets  []int

	MySQLDSN        string
	MySQLMaxIdle    int
	MySQLMaxOpen    int
	MySQLMaxLifeMin time.Duration

	DefaultQueryLimit int
	MaxQueryLimit     int
}

func loadConfigFromEnv() (*Config, error) {
	cfg := &Config{
		ListenAddr:        envString("PM_ACTIVITY_LISTEN_ADDR", ":18202"),
		DataAPIBaseURL:    envString("PM_ACTIVITY_API_BASE", "https://data-api.polymarket.com"),
		TargetUser:        strings.TrimSpace(os.Getenv("PM_ACTIVITY_USER")),
		PageLimit:         envInt("PM_ACTIVITY_PAGE_LIMIT", 1000),
		FastInterval:      time.Duration(envInt("PM_ACTIVITY_FAST_INTERVAL_SEC", 30)) * time.Second,
		BackfillInterval:  time.Duration(envInt("PM_ACTIVITY_BACKFILL_INTERVAL_SEC", 300)) * time.Second,
		RequestGap:        time.Duration(envInt("PM_ACTIVITY_REQUEST_GAP_MS", 800)) * time.Millisecond,
		HTTPTimeout:       time.Duration(envInt("PM_ACTIVITY_HTTP_TIMEOUT_SEC", 12)) * time.Second,
		FetchRetry:        envInt("PM_ACTIVITY_FETCH_RETRY", 3),
		StartupBackfill:   envBool("PM_ACTIVITY_STARTUP_BACKFILL", true),
		BackfillOffsets:   parseOffsets(envString("PM_ACTIVITY_BACKFILL_OFFSETS", "0,1000,2000,3000")),
		MySQLDSN:          strings.TrimSpace(os.Getenv("PM_ACTIVITY_MYSQL_DSN")),
		MySQLMaxIdle:      envInt("PM_ACTIVITY_MYSQL_MAX_IDLE", 4),
		MySQLMaxOpen:      envInt("PM_ACTIVITY_MYSQL_MAX_OPEN", 16),
		MySQLMaxLifeMin:   time.Duration(envInt("PM_ACTIVITY_MYSQL_MAX_LIFE_MIN", 30)) * time.Minute,
		DefaultQueryLimit: envInt("PM_ACTIVITY_QUERY_DEFAULT_LIMIT", 100),
		MaxQueryLimit:     envInt("PM_ACTIVITY_QUERY_MAX_LIMIT", 500),
	}

	if cfg.TargetUser == "" {
		return nil, fmt.Errorf("PM_ACTIVITY_USER is required")
	}
	if cfg.MySQLDSN == "" {
		return nil, fmt.Errorf("PM_ACTIVITY_MYSQL_DSN is required")
	}
	if cfg.PageLimit <= 0 {
		cfg.PageLimit = 1000
	}
	if cfg.PageLimit > 1000 {
		cfg.PageLimit = 1000
	}
	if cfg.FetchRetry <= 0 {
		cfg.FetchRetry = 3
	}
	if len(cfg.BackfillOffsets) == 0 {
		cfg.BackfillOffsets = []int{0, 1000, 2000, 3000}
	}
	if cfg.FastInterval < 10*time.Second {
		cfg.FastInterval = 10 * time.Second
	}
	if cfg.BackfillInterval < 60*time.Second {
		cfg.BackfillInterval = 60 * time.Second
	}
	if cfg.RequestGap < 200*time.Millisecond {
		cfg.RequestGap = 200 * time.Millisecond
	}
	if cfg.HTTPTimeout < 3*time.Second {
		cfg.HTTPTimeout = 3 * time.Second
	}
	if cfg.DefaultQueryLimit <= 0 {
		cfg.DefaultQueryLimit = 100
	}
	if cfg.MaxQueryLimit < cfg.DefaultQueryLimit {
		cfg.MaxQueryLimit = cfg.DefaultQueryLimit
	}

	return cfg, nil
}

func parseOffsets(raw string) []int {
	parts := strings.Split(raw, ",")
	seen := make(map[int]struct{}, len(parts))
	out := make([]int, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		v, err := strconv.Atoi(p)
		if err != nil || v < 0 {
			continue
		}
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	sort.Ints(out)
	return out
}

func envString(key, fallback string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback
	}
	return v
}

func envInt(key string, fallback int) int {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return v
}

func envBool(key string, fallback bool) bool {
	raw := strings.ToLower(strings.TrimSpace(os.Getenv(key)))
	if raw == "" {
		return fallback
	}
	switch raw {
	case "1", "true", "yes", "y", "on":
		return true
	case "0", "false", "no", "n", "off":
		return false
	default:
		return fallback
	}
}
