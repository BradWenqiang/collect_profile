#!/usr/bin/env python3
"""
Backfill market_tag for pm_activity_events.

What it does:
1) Ensure column `market_tag` exists.
2) Ensure index `idx_market_tag_ts (market_tag, timestamp_ms DESC)` exists.
3) Recompute market_tag in batches from slug/title/outcome for existing rows.

Tag rules:
- btc: contains btc / bitcoin
- eth: contains eth / ethereum
- sol: contains sol / solana
- other: fallback
"""

from __future__ import annotations

import argparse
import sys
import time


try:
    import pymysql
except Exception as exc:  # pragma: no cover
    print("[error] missing dependency: pymysql")
    print("install: python3 -m pip install pymysql")
    raise


def parse_args() -> argparse.Namespace:
    p = argparse.ArgumentParser(description="Backfill market_tag for pm_activity_events")
    p.add_argument("--host", default="127.0.0.1")
    p.add_argument("--port", type=int, default=3306)
    p.add_argument("--user", default="root")
    p.add_argument("--password", default="root")
    p.add_argument("--database", default="pm")
    p.add_argument("--table", default="pm_activity_events")
    p.add_argument("--batch-size", type=int, default=50000)
    p.add_argument("--dry-run", action="store_true", help="Do not commit changes")
    return p.parse_args()


def q_ident(name: str) -> str:
    return "`" + name.replace("`", "``") + "`"


def ensure_column(cur, db: str, table: str) -> None:
    cur.execute(
        """
        SELECT COUNT(1)
        FROM information_schema.columns
        WHERE table_schema=%s AND table_name=%s AND column_name='market_tag'
        """,
        (db, table),
    )
    exists = int(cur.fetchone()[0] or 0)
    if exists:
        print("[ok] column market_tag exists")
        return

    sql = f"ALTER TABLE {q_ident(table)} ADD COLUMN market_tag VARCHAR(24) NOT NULL DEFAULT '' AFTER slug"
    print("[exec]", sql)
    cur.execute(sql)
    print("[ok] column market_tag added")


def ensure_index(cur, db: str, table: str) -> None:
    cur.execute(
        """
        SELECT COUNT(1)
        FROM information_schema.statistics
        WHERE table_schema=%s AND table_name=%s AND index_name='idx_market_tag_ts'
        """,
        (db, table),
    )
    exists = int(cur.fetchone()[0] or 0)
    if exists:
        print("[ok] index idx_market_tag_ts exists")
        return

    sql = f"ALTER TABLE {q_ident(table)} ADD INDEX idx_market_tag_ts (market_tag, timestamp_ms DESC)"
    print("[exec]", sql)
    cur.execute(sql)
    print("[ok] index idx_market_tag_ts added")


def backfill(cur, table: str, batch_size: int, dry_run: bool) -> int:
    cur.execute(f"SELECT COALESCE(MIN(id),0), COALESCE(MAX(id),0), COUNT(1) FROM {q_ident(table)}")
    min_id, max_id, total = cur.fetchone()
    min_id = int(min_id or 0)
    max_id = int(max_id or 0)
    total = int(total or 0)

    if total == 0:
        print("[ok] no rows, nothing to backfill")
        return 0

    print(f"[info] rows={total}, id_range=[{min_id}, {max_id}], batch_size={batch_size}")

    update_sql = f"""
    UPDATE {q_ident(table)}
    SET market_tag = CASE
      WHEN (
        LOWER(SUBSTRING_INDEX(slug, '-', 1)) IN ('btc', 'bitcoin') OR
        LOWER(slug) LIKE '%btc%' OR LOWER(slug) LIKE '%bitcoin%' OR
        LOWER(title) LIKE '%btc%' OR LOWER(title) LIKE '%bitcoin%' OR
        LOWER(outcome) LIKE '%btc%' OR LOWER(outcome) LIKE '%bitcoin%'
      ) THEN 'btc'
      WHEN (
        LOWER(SUBSTRING_INDEX(slug, '-', 1)) IN ('eth', 'ethereum') OR
        LOWER(slug) LIKE '%eth%' OR LOWER(slug) LIKE '%ethereum%' OR
        LOWER(title) LIKE '%eth%' OR LOWER(title) LIKE '%ethereum%' OR
        LOWER(outcome) LIKE '%eth%' OR LOWER(outcome) LIKE '%ethereum%'
      ) THEN 'eth'
      WHEN (
        LOWER(SUBSTRING_INDEX(slug, '-', 1)) IN ('sol', 'solana') OR
        LOWER(slug) LIKE '%sol%' OR LOWER(slug) LIKE '%solana%' OR
        LOWER(title) LIKE '%sol%' OR LOWER(title) LIKE '%solana%' OR
        LOWER(outcome) LIKE '%sol%' OR LOWER(outcome) LIKE '%solana%'
      ) THEN 'sol'
      ELSE 'other'
    END
    WHERE id BETWEEN %s AND %s
    """

    touched = 0
    t0 = time.time()
    for start in range(min_id, max_id + 1, batch_size):
        end = min(start + batch_size - 1, max_id)
        cur.execute(update_sql, (start, end))
        changed = int(cur.rowcount or 0)
        touched += changed
        print(f"[batch] id=[{start},{end}] changed={changed}")
        if dry_run:
            pass
        else:
            cur.connection.commit()

    cost = time.time() - t0
    print(f"[done] changed={touched}, elapsed={cost:.2f}s, dry_run={dry_run}")
    return touched


def main() -> int:
    args = parse_args()

    conn = pymysql.connect(
        host=args.host,
        port=args.port,
        user=args.user,
        password=args.password,
        database=args.database,
        charset="utf8mb4",
        autocommit=False,
    )

    try:
        with conn.cursor() as cur:
            ensure_column(cur, args.database, args.table)
            ensure_index(cur, args.database, args.table)
            conn.commit()
            if args.dry_run:
                print("[info] --dry-run only affects data backfill. DDL (column/index) is already committed by MySQL.")

            changed = backfill(cur, args.table, args.batch_size, args.dry_run)
            if args.dry_run:
                conn.rollback()
                print("[info] dry-run rollback complete")
            else:
                conn.commit()

            print("[ok] finished, changed rows:", changed)
        return 0
    except Exception as exc:
        conn.rollback()
        print("[error]", exc)
        return 1
    finally:
        conn.close()


if __name__ == "__main__":
    raise SystemExit(main())
