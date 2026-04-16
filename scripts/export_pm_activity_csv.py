#!/usr/bin/env python3
"""
Export MySQL table data to CSV in batches.

Default table is `pm_activity_events`.
"""

from __future__ import annotations

import argparse
import csv
import datetime as dt
import os
import sys
import time

try:
    import pymysql
except Exception:
    print("[error] missing dependency: pymysql")
    print("install (ubuntu): sudo apt-get install -y python3-pymysql")
    raise


def parse_args() -> argparse.Namespace:
    p = argparse.ArgumentParser(description="Export pm_activity_events to CSV")
    p.add_argument("--host", default="127.0.0.1")
    p.add_argument("--port", type=int, default=3306)
    p.add_argument("--user", default="root")
    p.add_argument("--password", default="root")
    p.add_argument("--database", default="pm")
    p.add_argument("--table", default="pm_activity_events")
    p.add_argument("--output", default="", help="CSV output path")
    p.add_argument("--batch-size", type=int, default=50000)
    p.add_argument(
        "--tag",
        default="",
        help="Optional market_tag filter: btc/eth/sol/other",
    )
    p.add_argument("--start-id", type=int, default=0, help="Optional id lower bound")
    p.add_argument("--end-id", type=int, default=0, help="Optional id upper bound")
    return p.parse_args()


def q_ident(name: str) -> str:
    return "`" + name.replace("`", "``") + "`"


def normalize_tag(raw: str) -> str:
    t = (raw or "").strip().lower()
    if t in ("btc", "bitcoin"):
        return "btc"
    if t in ("eth", "ethereum"):
        return "eth"
    if t in ("sol", "solana"):
        return "sol"
    if t in ("other",):
        return "other"
    return ""


def pick_output_path(args: argparse.Namespace) -> str:
    if args.output:
        return args.output
    ts = dt.datetime.now().strftime("%Y%m%d_%H%M%S")
    tag = normalize_tag(args.tag)
    suffix = f"_{tag}" if tag else ""
    return f"{args.table}{suffix}_{ts}.csv"


def get_columns(cur, table: str) -> list[str]:
    cur.execute(f"SHOW COLUMNS FROM {q_ident(table)}")
    rows = cur.fetchall()
    cols = [str(r[0]) for r in rows]
    if not cols:
        raise RuntimeError(f"table {table} has no columns")
    if "id" not in cols:
        raise RuntimeError(f"table {table} has no id column, cannot batch by id")
    return cols


def build_base_where(tag: str, start_id: int, end_id: int) -> tuple[str, list[object]]:
    where_parts: list[str] = ["1=1"]
    args: list[object] = []

    t = normalize_tag(tag)
    if t:
        where_parts.append("market_tag = %s")
        args.append(t)
    if start_id > 0:
        where_parts.append("id >= %s")
        args.append(start_id)
    if end_id > 0:
        where_parts.append("id <= %s")
        args.append(end_id)
    return " WHERE " + " AND ".join(where_parts), args


def get_id_span(cur, table: str, base_where: str, base_args: list[object]) -> tuple[int, int, int]:
    sql = (
        "SELECT COALESCE(MIN(id),0), COALESCE(MAX(id),0), COUNT(1) "
        f"FROM {q_ident(table)}" + base_where
    )
    cur.execute(sql, tuple(base_args))
    min_id, max_id, total = cur.fetchone()
    return int(min_id or 0), int(max_id or 0), int(total or 0)


def to_csv_value(v):
    if v is None:
        return ""
    if isinstance(v, (dt.datetime, dt.date, dt.time)):
        return v.isoformat(sep=" ")
    if isinstance(v, (bytes, bytearray, memoryview)):
        try:
            return bytes(v).decode("utf-8", errors="replace")
        except Exception:
            return repr(bytes(v))
    return v


def export_rows(conn, args: argparse.Namespace) -> int:
    table = args.table.strip()
    if not table:
        raise RuntimeError("empty table")
    if args.batch_size <= 0:
        raise RuntimeError("batch-size must be > 0")

    out_path = pick_output_path(args)
    os.makedirs(os.path.dirname(out_path) or ".", exist_ok=True)

    with conn.cursor() as cur:
        cols = get_columns(cur, table)
        base_where, base_args = build_base_where(args.tag, args.start_id, args.end_id)
        min_id, max_id, total = get_id_span(cur, table, base_where, base_args)

        if total <= 0:
            print("[ok] no rows matched, writing header only")
            with open(out_path, "w", newline="", encoding="utf-8-sig") as f:
                writer = csv.writer(f)
                writer.writerow(cols)
            print("[ok] csv:", os.path.abspath(out_path))
            return 0

        print(
            f"[info] rows={total}, id_range=[{min_id},{max_id}], batch_size={args.batch_size}, output={out_path}"
        )

        col_sql = ",".join(q_ident(c) for c in cols)
        written = 0
        t0 = time.time()

        with open(out_path, "w", newline="", encoding="utf-8-sig") as f:
            writer = csv.writer(f)
            writer.writerow(cols)

            for start in range(min_id, max_id + 1, args.batch_size):
                end = min(start + args.batch_size - 1, max_id)
                sql = (
                    f"SELECT {col_sql} FROM {q_ident(table)}"
                    + base_where
                    + " AND id BETWEEN %s AND %s ORDER BY id ASC"
                )
                qargs = list(base_args) + [start, end]
                cur.execute(sql, tuple(qargs))
                rows = cur.fetchall()
                if rows:
                    for row in rows:
                        writer.writerow([to_csv_value(v) for v in row])
                    written += len(rows)
                print(f"[batch] id=[{start},{end}] rows={len(rows)} written={written}/{total}")

        cost = time.time() - t0
        print(f"[done] wrote={written}, elapsed={cost:.2f}s")
        print("[ok] csv:", os.path.abspath(out_path))
        return written


def main() -> int:
    args = parse_args()
    conn = pymysql.connect(
        host=args.host,
        port=args.port,
        user=args.user,
        password=args.password,
        database=args.database,
        charset="utf8mb4",
        autocommit=True,
    )
    try:
        export_rows(conn, args)
        return 0
    except Exception as exc:
        print("[error]", exc)
        return 1
    finally:
        conn.close()


if __name__ == "__main__":
    raise SystemExit(main())
