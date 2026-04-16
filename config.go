package main

import (
	"fmt"
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

func loadConfig() (*Config, error) {
	cfg := &Config{
		ListenAddr:        ":18202",
		DataAPIBaseURL:    "https://data-api.polymarket.com",
		TargetUser:        "0x89b5cdaaa4866c1e738406712012a630b4078beb",
		PageLimit:         1000,
		FastInterval:      30 * time.Second,
		BackfillInterval:  300 * time.Second,
		RequestGap:        800 * time.Millisecond,
		HTTPTimeout:       12 * time.Second,
		FetchRetry:        3,
		StartupBackfill:   true,
		BackfillOffsets:   []int{0, 1000, 2000, 3000},
		MySQLDSN:          "root:root@tcp(127.0.0.1:3306)/pm?charset=utf8mb4&parseTime=true&loc=Local",
		MySQLMaxIdle:      4,
		MySQLMaxOpen:      16,
		MySQLMaxLifeMin:   30 * time.Minute,
		DefaultQueryLimit: 100,
		MaxQueryLimit:     500,
	}

	if cfg.TargetUser == "" {
		return nil, fmt.Errorf("hardcoded target user is empty")
	}
	if cfg.MySQLDSN == "" {
		return nil, fmt.Errorf("hardcoded mysql dsn is empty")
	}
	return cfg, nil
}
