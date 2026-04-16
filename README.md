# collect_profile

面向 Polymarket `activity` 接口的独立采集服务，目标是稳定抓取最近窗口数据、去重落库，并提供后续页面可直接接入的 API。

## 已验证的接口行为

基于 `2026-04-16` 的实测：

1. `limit` 实际上限是 `1000`（传 `2000` 也只返回 `1000`）。
2. 分页支持 `offset`，但历史上限是 `3000`；`offset=4000` 会报错：
   - `max historical activity offset of 3000 exceeded`
3. 因此可获取的最近窗口约 `4000` 条（`offset=0/1000/2000/3000`）。
4. `transactionHash` 不是唯一键（同一 tx 可有多行）。

## 轮询策略（默认）

1. 快速轮询：每 `30s` 拉 `offset=0,limit=1000`
2. 回补轮询：每 `5min` 拉 `offset=0,1000,2000,3000`
3. 页间间隔：默认 `800ms`，防止太激进触发限流
4. 每次请求失败会自动重试（默认 `3` 次，指数回退）

这套策略来自你当前采样结果：

- `/tmp/pm_activity_probe_20260416_005450.csv`
  - 334 秒内 oldest timestamp 前移 708 秒（约 `127s/min`）
- `/tmp/pm_activity_churn_20260416_010638.csv`
  - 15~20 秒粒度下，平均新增约 `6` 条，峰值可到 `46` 条

## 去重方案

去重键不使用 `transactionHash`，改为：

- `event_id = sha256(user_wallet + "|" + canonical_raw_json)`

其中 `canonical_raw_json` 会对 key 排序并稳定序列化，确保同一事件重复拉取时哈希一致。

数据库层使用 `UNIQUE KEY (event_id)` + `INSERT IGNORE`，天然幂等。

## 数据库选型

建议：

1. 采集主链路：`MySQL`
   - 优点：强唯一键、写入幂等简单、状态接口查询方便
2. 分析看板：`ClickHouse`（后续可从 MySQL 同步）
   - 优点：宽表聚合、时间窗口分析、明细扫描更便宜

当前服务先落 MySQL（可靠采集），并已给出 ClickHouse 建表模板：

- `schema/mysql_pm_activity.sql`
- `schema/clickhouse_pm_activity.sql`

## 服务 API（Hertz）

1. `GET /`：内置 HTML 控制台（状态看板 + 手动同步 + 事件查询）
2. `GET /healthz`
3. `GET /api/v1/status`：轮询状态、累计抓取、累计去重、库内总行数
4. `POST /api/v1/sync/once`：手动触发（`mode=fast|backfill`）
5. `GET /api/v1/events`：查询最近明细（支持 `limit/page/slug/tag/type/side`）
6. `GET /api/v1/slugs`：按 `slug + market_tag` 聚合分页（支持 `page/page_size/keyword/tag`）
7. `GET /api/v1/events/strategy-group`：按策略分组窗口查询（`symbol/start_sec/end_sec`）

## Market Tag

采集入库时会自动打 `market_tag`：

- `btc` / `eth` / `sol` / `other`
- 表字段已带索引：`idx_market_tag_ts (market_tag, timestamp_ms DESC)`

## 存量数据回刷（Python）

脚本：`scripts/backfill_market_tag.py`

用途：

1. 自动补齐 `market_tag` 列和索引（老库兼容）
2. 按批次重刷存量数据 tag（BTC/ETH/SOL/OTHER）

示例：

```bash
python3 -m pip install pymysql
python3 scripts/backfill_market_tag.py --host 127.0.0.1 --port 3306 --user root --password root --database pm
```

## 运行配置（写死代码）

当前版本不读取环境变量，配置直接写死在 `config.go` 的 `loadConfig()` 里。

核心固定值：

1. `listen_addr`: `:18202`
2. `target_user`: `0x89b5cdaaa4866c1e738406712012a630b4078beb`
3. `mysql_dsn`: `root:root@tcp(127.0.0.1:3306)/pm?charset=utf8mb4&parseTime=true&loc=Local`
4. `page_limit`: `1000`
5. `fast_interval`: `30s`
6. `backfill_interval`: `300s`
7. `backfill_offsets`: `0,1000,2000,3000`

## 启动示例

```bash
cd collect_profile
go run .
```
