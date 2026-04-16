CREATE TABLE IF NOT EXISTS polymarket.pm_activity_events
(
  event_id FixedString(64),
  user_wallet String,
  proxy_wallet String,
  timestamp_ms Int64,
  event_time DateTime64(3, 'UTC'),
  condition_id String,
  activity_type LowCardinality(String),
  size Decimal(38, 18),
  usdc_size Decimal(38, 18),
  transaction_hash String,
  price Decimal(38, 18),
  asset String,
  side LowCardinality(String),
  outcome_index Nullable(Int32),
  title String,
  slug String,
  event_slug String,
  outcome String,
  source_offset Int32,
  source_index Int32,
  pulled_at DateTime64(3, 'UTC'),
  raw_json String,
  created_at DateTime64(3, 'UTC') DEFAULT now64(3)
)
ENGINE = ReplacingMergeTree(created_at)
PARTITION BY toDate(event_time)
ORDER BY (user_wallet, timestamp_ms, event_id)
SETTINGS index_granularity = 8192;

ALTER TABLE polymarket.pm_activity_events
ADD INDEX IF NOT EXISTS idx_slug slug TYPE bloom_filter GRANULARITY 8;

ALTER TABLE polymarket.pm_activity_events
ADD INDEX IF NOT EXISTS idx_tx transaction_hash TYPE bloom_filter GRANULARITY 8;
