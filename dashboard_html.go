package main

const dashboardHTML = `<!doctype html>
<html lang="zh-CN">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>PM Match Analyzer</title>
  <style>
    @import url("https://fonts.googleapis.com/css2?family=Plus+Jakarta+Sans:wght@400;500;600;700;800&family=JetBrains+Mono:wght@400;500&display=swap");
    :root {
      --bg: #eef4ea;
      --paper: #ffffff;
      --paper-soft: #f8fbf3;
      --ink: #132232;
      --muted: #546071;
      --line: #d4ddcc;
      --line-strong: #b7c5aa;
      --accent: #0c7a53;
      --accent-2: #1963d3;
      --warn: #b65d1d;
      --bad: #b73a39;
      --good: #1e8a59;
      --chip: #e6f7ef;
      --shadow: 0 10px 28px rgba(15, 38, 28, 0.12);
      --radius: 14px;
    }

    * { box-sizing: border-box; }

    body {
      margin: 0;
      color: var(--ink);
      font-family: "Plus Jakarta Sans", "PingFang SC", "Hiragino Sans GB", "Microsoft YaHei", sans-serif;
      background:
        radial-gradient(900px 380px at 90% -10%, #dff0ff 0%, transparent 65%),
        radial-gradient(1000px 420px at -10% 0%, #ffe9cf 0%, transparent 62%),
        radial-gradient(1100px 500px at 50% 120%, #d8f0df 0%, transparent 70%),
        var(--bg);
      min-height: 100vh;
      line-height: 1.45;
    }

    .page {
      width: min(1400px, calc(100vw - 28px));
      margin: 14px auto 24px;
      animation: rise .35s ease;
    }

    .head {
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 12px;
      margin-bottom: 10px;
      flex-wrap: wrap;
    }

    .title {
      margin: 0;
      font-size: clamp(22px, 3.4vw, 34px);
      letter-spacing: -0.03em;
      font-weight: 800;
    }

    .sub {
      margin: 4px 0 0;
      color: var(--muted);
      font-size: 13px;
    }

    .badge {
      border: 1px solid var(--line);
      border-radius: 999px;
      padding: 8px 12px;
      background: #fff;
      color: #2d3f59;
      font-family: "JetBrains Mono", Menlo, monospace;
      font-size: 12px;
    }

    .cards {
      display: grid;
      grid-template-columns: repeat(4, minmax(0, 1fr));
      gap: 10px;
      margin-bottom: 10px;
    }

    .card,
    .panel {
      border: 1px solid var(--line);
      background: rgba(255, 255, 255, 0.9);
      border-radius: var(--radius);
      box-shadow: var(--shadow);
      backdrop-filter: blur(5px);
    }

    .card { padding: 12px; }

    .card .k {
      font-size: 12px;
      color: var(--muted);
    }

    .card .v {
      margin-top: 6px;
      font-family: "JetBrains Mono", Menlo, monospace;
      font-size: 18px;
      font-weight: 500;
      word-break: break-word;
    }

    .panel {
      padding: 12px;
      margin-bottom: 10px;
    }

    .panel h2 {
      margin: 0;
      font-size: 16px;
      letter-spacing: -0.01em;
    }

    .muted {
      color: var(--muted);
      font-size: 12px;
    }

    .row {
      display: grid;
      grid-template-columns: repeat(6, minmax(0, 1fr));
      gap: 9px;
      margin-top: 10px;
    }

    label {
      display: block;
      font-size: 12px;
      color: var(--muted);
      margin-bottom: 4px;
    }

    input,
    button,
    select {
      width: 100%;
      border: 1px solid var(--line);
      border-radius: 10px;
      padding: 8px 10px;
      background: #fff;
      color: var(--ink);
      font: inherit;
      min-height: 36px;
    }

    button {
      cursor: pointer;
      font-weight: 600;
      transition: transform .08s ease, filter .12s ease;
    }

    button:hover { filter: brightness(0.97); }
    button:active { transform: translateY(1px); }

    .btn-primary {
      border-color: transparent;
      background: linear-gradient(135deg, #10875b, #0d6f8a);
      color: #fff;
    }

    .btn-soft {
      border-color: transparent;
      background: #e6f0ff;
      color: #184596;
    }

    .btn-warn {
      border-color: transparent;
      background: #ffe9d8;
      color: #7a4200;
    }

    .btn-light {
      background: #f6f8f2;
      color: #2d4358;
    }

    .ops {
      margin-top: 9px;
      display: flex;
      align-items: center;
      gap: 8px;
      flex-wrap: wrap;
    }

    .ops .inline-check {
      display: inline-flex;
      align-items: center;
      gap: 6px;
      margin: 0;
      color: #385168;
      font-size: 13px;
      width: auto;
    }

    .ops .inline-check input {
      width: auto;
      min-height: 0;
      padding: 0;
      margin: 0;
    }

    .status-grid {
      display: grid;
      grid-template-columns: repeat(3, minmax(0, 1fr));
      gap: 8px;
      margin-top: 10px;
    }

    .status-item {
      border: 1px solid #e7ede0;
      border-radius: 10px;
      background: #f7faf3;
      padding: 9px;
    }

    .status-item .s-k {
      font-size: 11px;
      color: var(--muted);
      text-transform: uppercase;
      letter-spacing: .03em;
    }

    .status-item .s-v {
      margin-top: 4px;
      font-family: "JetBrains Mono", Menlo, monospace;
      font-size: 12px;
      color: #25364a;
      word-break: break-word;
    }

    .notice {
      margin-top: 8px;
      border-radius: 10px;
      padding: 8px 10px;
      font-size: 13px;
      display: none;
    }

    .notice.show { display: block; }
    .notice.ok {
      background: #e4f8ee;
      color: #1b7f53;
    }

    .notice.err {
      background: #ffe7e7;
      color: #a53838;
    }

    .split {
      display: grid;
      grid-template-columns: minmax(320px, 450px) minmax(0, 1fr);
      gap: 10px;
      align-items: start;
    }

    .toolbar {
      display: grid;
      grid-template-columns: 1.2fr .8fr .8fr auto;
      gap: 8px;
      margin-top: 8px;
      align-items: end;
    }

    .toolbar-2 {
      margin-top: 8px;
      display: grid;
      grid-template-columns: repeat(3, minmax(0, 1fr));
      gap: 8px;
    }

    .mono {
      margin-top: 7px;
      font-family: "JetBrains Mono", Menlo, monospace;
      font-size: 12px;
      color: #4f6077;
      word-break: break-all;
    }

    .table-wrap {
      margin-top: 8px;
      border: 1px solid var(--line);
      border-radius: 12px;
      overflow: auto;
      background: #fff;
    }

    table {
      width: 100%;
      border-collapse: collapse;
      min-width: 760px;
    }

    th,
    td {
      border-bottom: 1px solid #edf1e8;
      padding: 8px 10px;
      text-align: left;
      font-size: 12px;
      vertical-align: top;
    }

    th {
      position: sticky;
      top: 0;
      background: #f2f7ee;
      z-index: 1;
      color: #485b73;
      text-transform: uppercase;
      letter-spacing: .04em;
      font-size: 11px;
    }

    td code {
      font-family: "JetBrains Mono", Menlo, monospace;
      color: #2c4f88;
      font-size: 11px;
    }

    .btn-link {
      border: 1px solid #bdd4bf;
      border-radius: 8px;
      background: #f2faef;
      color: #1f6d45;
      padding: 5px 9px;
      width: auto;
      min-height: 0;
      font-size: 12px;
      white-space: nowrap;
    }

    .row-selected {
      background: #f3faf5;
    }

    .tag-group-row td {
      background: #edf5ea;
      color: #2d5a42;
      font-weight: 700;
      letter-spacing: .01em;
      text-transform: uppercase;
      font-size: 11px;
      border-top: 1px solid #d7e5cf;
    }

    .kpi-grid {
      margin-top: 9px;
      display: grid;
      grid-template-columns: repeat(4, minmax(0, 1fr));
      gap: 8px;
    }

    .kpi {
      border: 1px solid #dce6d0;
      border-radius: 10px;
      background: #f7fbf0;
      padding: 8px;
    }

    .kpi .k {
      color: var(--muted);
      font-size: 11px;
    }

    .kpi .v {
      margin-top: 4px;
      font-family: "JetBrains Mono", Menlo, monospace;
      font-size: 14px;
      font-weight: 500;
      word-break: break-word;
    }

    .tone-up { color: var(--good); }
    .tone-down { color: var(--bad); }
    .tone-warn { color: var(--warn); }

    .analysis-box {
      margin-top: 8px;
      border: 1px solid #e5ecdb;
      border-radius: 10px;
      background: #f9fcf5;
      padding: 10px;
      font-size: 12px;
      color: #2c3f53;
    }

    .analysis-box .line {
      margin-top: 4px;
    }

    .combo-grid {
      margin-top: 9px;
      display: grid;
      grid-template-columns: repeat(2, minmax(0, 1fr));
      gap: 8px;
    }

    .combo-card {
      border: 1px solid #dfe8d5;
      border-radius: 10px;
      background: #fff;
      overflow: hidden;
    }

    .combo-head {
      padding: 7px 9px;
      border-bottom: 1px solid #e9f0e2;
      background: #f7fbf2;
      display: flex;
      align-items: baseline;
      justify-content: space-between;
      gap: 6px;
    }

    .combo-head .name {
      font-size: 12px;
      font-weight: 700;
      letter-spacing: .01em;
    }

    .combo-head .meta {
      font-size: 11px;
      color: var(--muted);
      font-family: "JetBrains Mono", Menlo, monospace;
    }

    .combo-table table {
      min-width: 600px;
    }

    .empty-cell {
      text-align: center;
      color: #607082;
      padding: 10px;
    }

    .wave-wrap {
      margin-top: 9px;
      border: 1px solid #dae5cf;
      border-radius: 10px;
      background: #fff;
      overflow: auto;
    }

    .wave-wrap table {
      min-width: 920px;
    }

    .pager {
      margin-top: 8px;
      display: flex;
      align-items: center;
      gap: 8px;
      flex-wrap: wrap;
    }

    .pager .meta {
      font-size: 12px;
      color: #41576f;
      font-family: "JetBrains Mono", Menlo, monospace;
    }

    .pager-meta-line {
      margin-top: 4px;
    }

    .slug-table {
      min-width: 0;
    }

    .slug-table th:first-child,
    .slug-table td:first-child {
      width: 260px;
      max-width: 260px;
    }

    .slug-table td:first-child code {
      display: inline-block;
      max-width: 248px;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
      vertical-align: top;
    }

    @keyframes rise {
      from {
        opacity: .4;
        transform: translateY(8px);
      }
      to {
        opacity: 1;
        transform: translateY(0);
      }
    }

    @media (max-width: 1200px) {
      .cards { grid-template-columns: repeat(2, minmax(0, 1fr)); }
      .kpi-grid { grid-template-columns: repeat(2, minmax(0, 1fr)); }
      .split { grid-template-columns: 1fr; }
    }

    @media (max-width: 880px) {
      .row,
      .toolbar,
      .toolbar-2,
      .status-grid,
      .combo-grid { grid-template-columns: 1fr; }

      .page {
        width: calc(100vw - 16px);
        margin-top: 8px;
      }
    }
  </style>
</head>
<body>
  <main class="page">
    <header class="head">
      <div>
        <h1 class="title">Polymarket 比赛分析台</h1>
        <p class="sub">分页查询 + 按 slug 比赛视角展开 + buy/sell/up/down 四组合分析 + 波段时间判断</p>
      </div>
      <div class="badge" id="badge-health">service: checking</div>
    </header>

    <section class="cards">
      <article class="card">
        <div class="k">Stored Events</div>
        <div class="v" id="v-stored">-</div>
      </article>
      <article class="card">
        <div class="k">Last Fetched / Inserted</div>
        <div class="v" id="v-last">-</div>
      </article>
      <article class="card">
        <div class="k">Total Fetched / Inserted / Dup</div>
        <div class="v" id="v-total">-</div>
      </article>
      <article class="card">
        <div class="k">Queue / Running</div>
        <div class="v" id="v-queue">-</div>
      </article>
    </section>

    <section class="panel">
      <h2>服务状态与同步</h2>
      <div class="status-grid">
        <div class="status-item">
          <div class="s-k">Last Success</div>
          <div class="s-v" id="s-success">-</div>
        </div>
        <div class="status-item">
          <div class="s-k">Last Error</div>
          <div class="s-v" id="s-error">-</div>
        </div>
        <div class="status-item">
          <div class="s-k">Newest / Oldest TS</div>
          <div class="s-v" id="s-ts">-</div>
        </div>
      </div>
      <div class="ops">
        <button class="btn-primary" id="btn-refresh-status" style="width:auto;">刷新状态</button>
        <button class="btn-soft" id="btn-sync-fast" style="width:auto;">触发 Fast</button>
        <button class="btn-warn" id="btn-sync-backfill" style="width:auto;">触发 Backfill</button>
        <label class="inline-check">
          <input id="auto-refresh" type="checkbox" checked>
          自动刷新(10s)
        </label>
        <span class="muted" id="status-config">listen=- / wallet=- / limit=-</span>
      </div>
      <div class="notice" id="notice"></div>
    </section>

    <section class="panel">
        <h2>比赛列表（按 slug 分页）</h2>
        <div class="toolbar">
          <div>
            <label for="slug-keyword">搜索 slug / title</label>
            <input id="slug-keyword" type="text" placeholder="输入关键字" />
          </div>
          <div>
            <label for="slug-page-size">每页</label>
            <input id="slug-page-size" type="number" min="1" max="100" value="12" />
          </div>
          <div>
            <label for="slug-tag">Tag</label>
            <select id="slug-tag">
              <option value="all">全部</option>
              <option value="btc">BTC</option>
              <option value="eth">ETH</option>
              <option value="sol">SOL</option>
              <option value="other">OTHER</option>
            </select>
          </div>
          <div style="display:flex;align-items:flex-end;">
            <button class="btn-primary" id="btn-slug-search">查询比赛</button>
          </div>
        </div>
        <div class="pager">
          <button class="btn-light" id="btn-slug-prev" style="width:auto;">上一页</button>
          <button class="btn-light" id="btn-slug-next" style="width:auto;">下一页</button>
        </div>
        <div class="pager-meta-line"><span class="meta" id="slug-page-meta">-</span></div>
        <div class="table-wrap">
          <table class="slug-table">
            <thead>
              <tr>
                <th>slug</th>
                <th>tag</th>
                <th>events</th>
                <th>buy/sell</th>
                <th>up/down</th>
                <th>latest</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody id="slug-rows">
              <tr><td colspan="7" class="empty-cell">加载中...</td></tr>
            </tbody>
          </table>
        </div>
    </section>

    <section class="panel">
        <h2>比赛展开分析</h2>
        <div class="toolbar-2">
          <div>
            <label for="detail-slug">当前 slug</label>
            <input id="detail-slug" type="text" placeholder="请点击左侧比赛" readonly />
          </div>
          <div>
            <label for="detail-wave-gap">波段间隔(秒)</label>
            <input id="detail-wave-gap" type="number" min="10" max="900" value="120" />
          </div>
          <div style="display:flex;align-items:flex-end;">
            <button class="btn-soft" id="detail-refresh">刷新当前比赛</button>
          </div>
        </div>

        <div class="mono" id="detail-note">请选择左侧 slug，系统会自动按比赛聚合并做分析。</div>
        <div class="mono" id="detail-group-meta">分组: -</div>
        <div class="kpi-grid" id="detail-kpis"></div>

        <div class="analysis-box" id="detail-summary">
          还没有比赛数据。
        </div>

        <div class="combo-grid">
          <section class="combo-card">
            <div class="combo-head">
              <span class="name">Buy + Up</span>
              <span class="meta" id="meta-buy-up">-</span>
            </div>
            <div class="combo-table table-wrap" style="margin-top:0;border:none;border-radius:0;">
              <table>
                <thead>
                  <tr>
                    <th>time</th>
                    <th>price</th>
                    <th>qty</th>
                    <th>usdc</th>
                    <th>outcome</th>
                    <th>tx</th>
                  </tr>
                </thead>
                <tbody id="combo-buy-up"><tr><td colspan="6" class="empty-cell">暂无数据</td></tr></tbody>
              </table>
            </div>
          </section>

          <section class="combo-card">
            <div class="combo-head">
              <span class="name">Buy + Down</span>
              <span class="meta" id="meta-buy-down">-</span>
            </div>
            <div class="combo-table table-wrap" style="margin-top:0;border:none;border-radius:0;">
              <table>
                <thead>
                  <tr>
                    <th>time</th>
                    <th>price</th>
                    <th>qty</th>
                    <th>usdc</th>
                    <th>outcome</th>
                    <th>tx</th>
                  </tr>
                </thead>
                <tbody id="combo-buy-down"><tr><td colspan="6" class="empty-cell">暂无数据</td></tr></tbody>
              </table>
            </div>
          </section>

          <section class="combo-card">
            <div class="combo-head">
              <span class="name">Sell + Up</span>
              <span class="meta" id="meta-sell-up">-</span>
            </div>
            <div class="combo-table table-wrap" style="margin-top:0;border:none;border-radius:0;">
              <table>
                <thead>
                  <tr>
                    <th>time</th>
                    <th>price</th>
                    <th>qty</th>
                    <th>usdc</th>
                    <th>outcome</th>
                    <th>tx</th>
                  </tr>
                </thead>
                <tbody id="combo-sell-up"><tr><td colspan="6" class="empty-cell">暂无数据</td></tr></tbody>
              </table>
            </div>
          </section>

          <section class="combo-card">
            <div class="combo-head">
              <span class="name">Sell + Down</span>
              <span class="meta" id="meta-sell-down">-</span>
            </div>
            <div class="combo-table table-wrap" style="margin-top:0;border:none;border-radius:0;">
              <table>
                <thead>
                  <tr>
                    <th>time</th>
                    <th>price</th>
                    <th>qty</th>
                    <th>usdc</th>
                    <th>outcome</th>
                    <th>tx</th>
                  </tr>
                </thead>
                <tbody id="combo-sell-down"><tr><td colspan="6" class="empty-cell">暂无数据</td></tr></tbody>
              </table>
            </div>
          </section>
        </div>

        <h2 style="margin-top:12px;">波段时间分析（先后顺序）</h2>
        <div class="wave-wrap">
          <table>
            <thead>
              <tr>
                <th>#</th>
                <th>start</th>
                <th>end</th>
                <th>span(s)</th>
                <th>records</th>
                <th>先后关系</th>
                <th>首尾组合</th>
                <th>Buy配对盈亏</th>
                <th>Sell配对盈亏</th>
                <th>总估算</th>
              </tr>
            </thead>
            <tbody id="wave-rows">
              <tr><td colspan="10" class="empty-cell">暂无波段数据</td></tr>
            </tbody>
          </table>
        </div>
    </section>

    <section class="panel">
      <h2>事件明细查询（支持分页）</h2>
      <div class="row">
        <div>
          <label for="q-limit">page_size</label>
          <input id="q-limit" type="number" min="1" max="500" value="100" />
        </div>
        <div>
          <label for="q-page">page</label>
          <input id="q-page" type="number" min="1" value="1" />
        </div>
        <div>
          <label for="q-slug">slug</label>
          <input id="q-slug" type="text" placeholder="可留空" />
        </div>
        <div>
          <label for="q-type">type</label>
          <input id="q-type" type="text" placeholder="BUY / SELL / TRADE..." />
        </div>
        <div>
          <label for="q-side">side</label>
          <input id="q-side" type="text" placeholder="buy / sell" />
        </div>
        <div style="display:flex;align-items:flex-end;">
          <button class="btn-primary" id="btn-query">查询</button>
        </div>
      </div>

      <div class="pager">
        <button class="btn-light" id="btn-events-prev" style="width:auto;">上一页</button>
        <button class="btn-light" id="btn-events-next" style="width:auto;">下一页</button>
        <button class="btn-light" id="btn-query-reset" style="width:auto;">重置过滤</button>
      </div>
      <div class="pager-meta-line"><span class="meta" id="events-page-meta">-</span></div>
      <div class="mono" id="query-meta">query: /api/v1/events</div>

      <div class="table-wrap">
        <table>
          <thead>
            <tr>
              <th>event_time</th>
              <th>type</th>
              <th>side</th>
              <th>slug</th>
              <th>outcome</th>
              <th>price</th>
              <th>usdc_size</th>
              <th>size</th>
              <th>tx</th>
              <th>source</th>
            </tr>
          </thead>
          <tbody id="rows">
            <tr><td colspan="10" class="empty-cell">加载中...</td></tr>
          </tbody>
        </table>
      </div>
    </section>
  </main>

  <script>
    (function () {
      var notice = document.getElementById("notice");
      var badgeHealth = document.getElementById("badge-health");
      var autoRefreshBox = document.getElementById("auto-refresh");
      var timer = null;

      var state = {
        slugPage: 1,
        slugTotalPages: 0,
        slugPageSize: 12,
        slugKeyword: "",
        slugTag: "all",
        slugMode: "cursor",
        slugCursor: "",
        slugNextCursor: "",
        slugHasMore: false,
        slugCursorStack: [""],
        slugFilterKey: "",
        selectedSlug: "",
        selectedStrategy: null,
        eventPage: 1,
        eventTotalPages: 0,
        eventPageSize: 100
      };

      function escapeHTML(v) {
        return String(v == null ? "" : v)
          .replace(/&/g, "&amp;")
          .replace(/</g, "&lt;")
          .replace(/>/g, "&gt;")
          .replace(/\"/g, "&quot;")
          .replace(/'/g, "&#39;");
      }

      function parseNumber(v) {
        var n = Number(v);
        if (!Number.isFinite(n)) return 0;
        return n;
      }

      function toInt(v, fallback) {
        var n = parseInt(v, 10);
        if (!Number.isFinite(n) || n <= 0) return fallback;
        return n;
      }

      function clamp(n, min, max) {
        if (n < min) return min;
        if (n > max) return max;
        return n;
      }

      function fmtTime(raw) {
        if (!raw) return "-";
        var d = new Date(raw);
        if (isNaN(d.getTime())) return "-";
        return d.toLocaleString();
      }

      function normalizeEpochMs(v) {
        var ts = Number(v || 0);
        if (!Number.isFinite(ts) || ts <= 0) return 0;
        if (ts < 1000000000000) return ts * 1000;
        return ts;
      }

      function fmtMS(ms) {
        var v = normalizeEpochMs(ms);
        if (!v || v <= 0) return "-";
        var d = new Date(v);
        if (isNaN(d.getTime())) return "-";
        return d.toLocaleString();
      }

      function fmtNum(v, digits) {
        if (!Number.isFinite(v)) return "-";
        var d = Number.isFinite(digits) ? digits : 2;
        return Number(v).toFixed(d);
      }

      function fmtNumFromRaw(raw) {
        var s = String(raw == null ? "" : raw).trim();
        if (!s) return "-";
        var n = Number(s);
        if (!Number.isFinite(n)) return s;
        return fmtNum(n, 2);
      }

      function shortTx(tx) {
        if (!tx) return "-";
        var s = String(tx);
        if (s.length <= 14) return s;
        return s.slice(0, 8) + "..." + s.slice(-6);
      }

      function showNotice(msg, isError) {
        notice.textContent = msg;
        notice.className = "notice show " + (isError ? "err" : "ok");
        setTimeout(function () {
          notice.className = "notice";
        }, 3200);
      }

      async function request(path, options) {
        var res = await fetch(path, options || {});
        var body = await res.json();
        if (!res.ok || !body || body.code !== 20000) {
          throw new Error((body && body.message) || ("HTTP " + res.status));
        }
        return body.data;
      }

      function parseStrategyGroupBySlug(slug, fallbackTag) {
        var text = String(slug || "").trim().toLowerCase();
        if (!text) return null;
        var parts = text.split("-");
        if (parts.length < 2) return null;

        var symbol = parts[0].replace(/[^a-z0-9]/g, "");
        if (fallbackTag && fallbackTag !== "all") {
          symbol = normalizeTag(fallbackTag);
        }
        if (!symbol) return null;

        var closeSec = parseInt(parts[parts.length - 1], 10);
        if (!Number.isFinite(closeSec) || closeSec <= 0) return null;

        var closeMs = normalizeEpochMs(closeSec);
        var closeSecNorm = Math.floor(closeMs / 1000);
        var startSec = Math.floor((closeSecNorm - 1) / 900) * 900;
        var endSec = startSec + 900;

        return {
          symbol: symbol,
          startSec: startSec,
          endSec: endSec
        };
      }

      function fmtTimeRangeSec(startSec, endSec) {
        if (!startSec || !endSec || endSec <= startSec) return "-";
        var start = new Date(startSec * 1000);
        var end = new Date(endSec * 1000);
        return start.toLocaleString() + " ~ " + end.toLocaleTimeString();
      }

      async function fetchEventsByStrategyGroup(strategy) {
        var params = new URLSearchParams();
        params.set("symbol", strategy.symbol);
        params.set("start_sec", String(strategy.startSec));
        params.set("end_sec", String(strategy.endSec));
        params.set("limit", "5000");

        var data = await request("/api/v1/events/strategy-group?" + params.toString());
        var items = Array.isArray(data.items) ? data.items : [];
        var slugMap = {};
        var i;
        for (i = 0; i < items.length; i++) {
          var s = String(items[i].slug || "").trim();
          if (s) slugMap[s] = true;
        }
        return {
          items: items,
          total: Number(data.count || items.length),
          truncated: false,
          strategy: strategy,
          slugs: Object.keys(slugMap).sort()
        };
      }

      function getEventTimestampMs(e) {
        var ts = normalizeEpochMs(e && e.timestamp_ms);
        if (Number.isFinite(ts) && ts > 0) return ts;
        if (e && e.event_time) {
          var d = new Date(e.event_time);
          if (!isNaN(d.getTime())) return d.getTime();
        }
        return 0;
      }

      function normalizeSide(raw) {
        var s = String(raw || "").trim().toLowerCase();
        if (!s) return "";
        if (s.indexOf("buy") >= 0 || s.indexOf("买") >= 0) return "buy";
        if (s.indexOf("sell") >= 0 || s.indexOf("卖") >= 0) return "sell";
        return "";
      }

      function normalizeDirection(e) {
        var raw = String((e && e.outcome) || "").trim().toLowerCase();
        if (!raw && e && e.outcome_index != null) {
          var idxOnly = Number(e.outcome_index);
          if (idxOnly === 0) return "up";
          if (idxOnly === 1) return "down";
        }

        var upWords = ["up", "yes", "long", "higher", "bull", "上涨", "看涨", "涨"];
        var downWords = ["down", "no", "short", "lower", "bear", "下跌", "看跌", "跌"];

        var i;
        for (i = 0; i < upWords.length; i++) {
          if (raw.indexOf(upWords[i]) >= 0) return "up";
        }
        for (i = 0; i < downWords.length; i++) {
          if (raw.indexOf(downWords[i]) >= 0) return "down";
        }

        var idx = Number(e && e.outcome_index);
        if (idx === 0) return "up";
        if (idx === 1) return "down";
        return "";
      }

      function deriveQty(e, price) {
        var q = parseNumber(e && e.size);
        if (q > 0) return q;
        var usdc = parseNumber(e && e.usdc_size);
        if (usdc > 0 && price > 0) {
          var converted = usdc / price;
          if (Number.isFinite(converted) && converted > 0) return converted;
        }
        return 1;
      }

      function normalizeRecord(e) {
        var side = normalizeSide(e && e.side);
        var direction = normalizeDirection(e);
        var combo = (side && direction) ? (side + "_" + direction) : "";
        var price = parseNumber(e && e.price);
        var qty = deriveQty(e, price);
        return {
          raw: e,
          ts: getEventTimestampMs(e),
          side: side,
          direction: direction,
          combo: combo,
          price: price,
          qty: qty,
          usdc: parseNumber(e && e.usdc_size)
        };
      }

      function comboLabel(combo) {
        if (combo === "buy_up") return "buy-up";
        if (combo === "buy_down") return "buy-down";
        if (combo === "sell_up") return "sell-up";
        if (combo === "sell_down") return "sell-down";
        return "unknown";
      }

      function pairLegs(upLegs, downLegs, side) {
        var eps = 1e-9;
        var ups = upLegs.slice().sort(function (a, b) { return a.ts - b.ts; }).map(function (x) {
          return { item: x, remain: x.qty > 0 ? x.qty : 0 };
        });
        var downs = downLegs.slice().sort(function (a, b) { return a.ts - b.ts; }).map(function (x) {
          return { item: x, remain: x.qty > 0 ? x.qty : 0 };
        });

        var totalUpQty = 0;
        var totalDownQty = 0;
        var k;
        for (k = 0; k < upLegs.length; k++) totalUpQty += upLegs[k].qty;
        for (k = 0; k < downLegs.length; k++) totalDownQty += downLegs[k].qty;

        var i = 0;
        var j = 0;
        var pairedQty = 0;
        var pairedPnl = 0;
        var samples = [];

        while (i < ups.length && j < downs.length) {
          if (ups[i].remain <= eps) { i += 1; continue; }
          if (downs[j].remain <= eps) { j += 1; continue; }

          var matched = Math.min(ups[i].remain, downs[j].remain);
          if (matched <= eps) break;

          var priceSum = ups[i].item.price + downs[j].item.price;
          var edgePerUnit = side === "buy" ? (1 - priceSum) : (priceSum - 1);
          var pnl = edgePerUnit * matched;

          pairedQty += matched;
          pairedPnl += pnl;

          if (samples.length < 12) {
            samples.push({
              matchedQty: matched,
              edgePerUnit: edgePerUnit,
              pnl: pnl,
              order: ups[i].item.ts <= downs[j].item.ts ? "up->down" : "down->up",
              upTs: ups[i].item.ts,
              downTs: downs[j].item.ts,
              upPrice: ups[i].item.price,
              downPrice: downs[j].item.price
            });
          }

          ups[i].remain -= matched;
          downs[j].remain -= matched;
        }

        var base = Math.max(totalUpQty, totalDownQty);
        var symmetryRate = base > eps ? (pairedQty / base) : 0;

        return {
          totalUpQty: totalUpQty,
          totalDownQty: totalDownQty,
          pairedQty: pairedQty,
          pairedPnl: pairedPnl,
          avgEdgePerUnit: pairedQty > eps ? (pairedPnl / pairedQty) : 0,
          symmetryRate: symmetryRate,
          unpairedUpQty: Math.max(0, totalUpQty - pairedQty),
          unpairedDownQty: Math.max(0, totalDownQty - pairedQty),
          samples: samples
        };
      }

      function buildWaves(records, gapSec) {
        var list = records.slice().filter(function (r) { return r.ts > 0 && r.combo; });
        if (!list.length) return [];
        list.sort(function (a, b) { return a.ts - b.ts; });

        var gapMs = Math.max(10, gapSec || 120) * 1000;
        var waves = [];
        var current = null;
        var idx;

        for (idx = 0; idx < list.length; idx++) {
          var rec = list[idx];
          if (!current || rec.ts - current.endTs > gapMs) {
            if (current) waves.push(current);
            current = {
              startTs: rec.ts,
              endTs: rec.ts,
              events: [rec]
            };
          } else {
            current.events.push(rec);
            current.endTs = rec.ts;
          }
        }
        if (current) waves.push(current);

        return waves.map(function (w, i) {
          var counts = {
            buy_up: 0,
            buy_down: 0,
            sell_up: 0,
            sell_down: 0
          };

          var x;
          for (x = 0; x < w.events.length; x++) {
            if (counts[w.events[x].combo] != null) counts[w.events[x].combo] += 1;
          }

          var first = w.events[0];
          var last = w.events[w.events.length - 1];

          var firstUpTs = 0;
          var firstDownTs = 0;
          for (x = 0; x < w.events.length; x++) {
            if (!firstUpTs && w.events[x].direction === "up") firstUpTs = w.events[x].ts;
            if (!firstDownTs && w.events[x].direction === "down") firstDownTs = w.events[x].ts;
          }

          var orderHint = "方向不明";
          if (firstUpTs && firstDownTs) {
            orderHint = firstUpTs <= firstDownTs ? "先上后下" : "先下后上";
          } else if (firstUpTs) {
            orderHint = "只有上腿";
          } else if (firstDownTs) {
            orderHint = "只有下腿";
          }

          var buyP = pairLegs(
            w.events.filter(function (r) { return r.combo === "buy_up"; }),
            w.events.filter(function (r) { return r.combo === "buy_down"; }),
            "buy"
          );
          var sellP = pairLegs(
            w.events.filter(function (r) { return r.combo === "sell_up"; }),
            w.events.filter(function (r) { return r.combo === "sell_down"; }),
            "sell"
          );

          return {
            index: i + 1,
            startTs: w.startTs,
            endTs: w.endTs,
            spanSec: Math.max(0, Math.round((w.endTs - w.startTs) / 1000)),
            records: w.events.length,
            orderHint: orderHint,
            firstCombo: comboLabel(first.combo),
            lastCombo: comboLabel(last.combo),
            buyPnl: buyP.pairedPnl,
            sellPnl: sellP.pairedPnl,
            totalPnl: buyP.pairedPnl + sellP.pairedPnl,
            comboCounts: counts
          };
        });
      }

      function analyzeSlugEvents(events, waveGapSec) {
        var normalized = events.map(normalizeRecord);
        var buckets = {
          buy_up: [],
          buy_down: [],
          sell_up: [],
          sell_down: []
        };
        var unknown = [];

        var i;
        for (i = 0; i < normalized.length; i++) {
          var r = normalized[i];
          if (buckets[r.combo]) {
            buckets[r.combo].push(r);
          } else {
            unknown.push(r);
          }
        }

        Object.keys(buckets).forEach(function (k) {
          buckets[k].sort(function (a, b) { return a.ts - b.ts; });
        });

        var buyPair = pairLegs(buckets.buy_up, buckets.buy_down, "buy");
        var sellPair = pairLegs(buckets.sell_up, buckets.sell_down, "sell");

        var totalPairedQty = buyPair.pairedQty + sellPair.pairedQty;
        var totalPnl = buyPair.pairedPnl + sellPair.pairedPnl;

        var verdict = "样本不足（没有可配对上下腿）";
        if (totalPairedQty > 0) {
          if (totalPnl > 0.000001) verdict = "偏赚钱";
          else if (totalPnl < -0.000001) verdict = "偏赔钱";
          else verdict = "接近平衡";
        }

        var riskParts = [];
        if (buyPair.symmetryRate < 0.9 && (buyPair.totalUpQty > 0 || buyPair.totalDownQty > 0)) {
          riskParts.push("Buy 组合单边风险较高");
        }
        if (sellPair.symmetryRate < 0.9 && (sellPair.totalUpQty > 0 || sellPair.totalDownQty > 0)) {
          riskParts.push("Sell 组合单边风险较高");
        }
        if (!riskParts.length) {
          riskParts.push("对称度较好，接近可配对状态");
        }

        var knownRecords = [];
        Object.keys(buckets).forEach(function (k) {
          knownRecords = knownRecords.concat(buckets[k]);
        });
        var waves = buildWaves(knownRecords, waveGapSec);

        return {
          buckets: buckets,
          unknown: unknown,
          buyPair: buyPair,
          sellPair: sellPair,
          totalPnl: totalPnl,
          totalPairedQty: totalPairedQty,
          verdict: verdict,
          riskText: riskParts.join("；"),
          waves: waves
        };
      }

      function renderStatusSummary(data) {
        var poller = data.poller || {};
        var cfg = data.config || {};

        badgeHealth.textContent = "service: online";
        badgeHealth.style.color = "#1e8a59";

        document.getElementById("v-stored").textContent = String(data.stored_events || 0);
        document.getElementById("v-last").textContent = String(poller.last_fetched_rows || 0) + " / " + String(poller.last_inserted_rows || 0);
        document.getElementById("v-total").textContent = String(poller.total_fetched_rows || 0) + " / " + String(poller.total_inserted_rows || 0) + " / " + String(poller.total_duplicate_rows || 0);
        document.getElementById("v-queue").textContent = String(poller.queue_len || 0) + " / " + (poller.running ? "running" : "idle");

        document.getElementById("s-success").textContent = fmtTime(poller.last_success_at);
        document.getElementById("s-error").textContent = poller.last_error ? (fmtTime(poller.last_error_at) + " | " + poller.last_error) : "-";
        document.getElementById("s-ts").textContent = fmtMS(poller.last_newest_ts) + " / " + fmtMS(poller.last_oldest_ts);

        document.getElementById("status-config").textContent =
          "listen=" + (cfg.listen_addr || "-") +
          " / wallet=" + (cfg.user_wallet || "-") +
          " / limit=" + (cfg.page_limit || "-");
      }

      async function loadStatus() {
        try {
          var data = await request("/api/v1/status");
          renderStatusSummary(data);
        } catch (err) {
          badgeHealth.textContent = "service: offline";
          badgeHealth.style.color = "#b73a39";
          showNotice("状态刷新失败: " + err.message, true);
        }
      }

      async function triggerSync(mode) {
        try {
          var data = await request("/api/v1/sync/once", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ mode: mode })
          });
          showNotice("已提交 " + mode + " 同步，queued=" + String(!!data.queued), false);
          await loadStatus();
        } catch (err) {
          showNotice("触发同步失败: " + err.message, true);
        }
      }

      function normalizeTag(raw) {
        var t = String(raw || "").trim().toLowerCase();
        if (!t) return "other";
        if (t === "btc" || t === "bitcoin") return "btc";
        if (t === "eth" || t === "ethereum") return "eth";
        if (t === "sol" || t === "solana") return "sol";
        return "other";
      }

      function prettyTag(tag) {
        var t = normalizeTag(tag);
        if (t === "btc") return "BTC";
        if (t === "eth") return "ETH";
        if (t === "sol") return "SOL";
        return "OTHER";
      }

      function readSlugFilters() {
        var pageSize = clamp(toInt(document.getElementById("slug-page-size").value, state.slugPageSize || 12), 1, 100);
        state.slugPageSize = pageSize;
        state.slugKeyword = document.getElementById("slug-keyword").value.trim();
        state.slugTag = String(document.getElementById("slug-tag").value || "all").trim().toLowerCase();
      }

      function buildSlugFilterKey() {
        return [String(state.slugPageSize || 12), state.slugKeyword || "", state.slugTag || "all"].join("|");
      }

      function resetSlugCursorPaging() {
        state.slugPage = 1;
        state.slugCursor = "";
        state.slugNextCursor = "";
        state.slugHasMore = false;
        state.slugCursorStack = [""];
      }

      function renderSlugRows(items) {
        var tbody = document.getElementById("slug-rows");
        if (!items || !items.length) {
          tbody.innerHTML = "<tr><td colspan=\"7\" class=\"empty-cell\">暂无比赛</td></tr>";
          return;
        }
        var order = ["btc", "eth", "sol", "other"];
        var grouped = {};
        var i;
        for (i = 0; i < items.length; i++) {
          var row = items[i];
          var tag = normalizeTag(row.market_tag);
          if (!grouped[tag]) grouped[tag] = [];
          grouped[tag].push(row);
        }

        var html = "";
        order.forEach(function (tag) {
          var rows = grouped[tag] || [];
          if (!rows.length) return;
          html += "<tr class=\"tag-group-row\"><td colspan=\"7\">" + escapeHTML(prettyTag(tag)) + " · " + rows.length + " 场</td></tr>";
          rows.forEach(function (s) {
            var slug = String(s.slug || "");
            var selected = state.selectedSlug && state.selectedSlug === slug ? "row-selected" : "";
            var actionSlug = encodeURIComponent(slug);
            var actionTag = encodeURIComponent(normalizeTag(s.market_tag));
            html += "<tr class=\"" + selected + "\">" +
              "<td><code>" + escapeHTML(slug || "-") + "</code></td>" +
              "<td>" + escapeHTML(prettyTag(s.market_tag)) + "</td>" +
              "<td>" + escapeHTML(String(s.event_count || 0)) + "</td>" +
              "<td>" + escapeHTML(String(s.buy_count || 0) + " / " + String(s.sell_count || 0)) + "</td>" +
              "<td>" + escapeHTML(String(s.up_count || 0) + " / " + String(s.down_count || 0)) + "</td>" +
              "<td>" + escapeHTML(fmtMS(s.last_timestamp_ms || 0)) + "</td>" +
              "<td><button class=\"btn-link btn-open-slug\" data-slug=\"" + actionSlug + "\" data-tag=\"" + actionTag + "\">展开</button></td>" +
              "</tr>";
          });
        });
        tbody.innerHTML = html || "<tr><td colspan=\"7\" class=\"empty-cell\">暂无比赛</td></tr>";
      }

      function renderSlugPager(data) {
        var mode = String((data && data.mode) || "").trim().toLowerCase();
        if (mode === "cursor") {
          var cPage = Number(state.slugPage || 1);
          var cCount = Number(data && data.count);
          if (!Number.isFinite(cCount) || cCount < 0) {
            cCount = Array.isArray(data && data.items) ? data.items.length : 0;
          }
          var cHasMore = !!(data && data.has_more);
          if (!Number.isFinite(cPage) || cPage <= 0) cPage = 1;
          document.getElementById("slug-page-meta").textContent = "第 " + cPage + " 页（游标分页），本页 " + cCount + " 场" + (cHasMore ? "，可继续下一页" : "，已到末页");
          document.getElementById("btn-slug-prev").disabled = cPage <= 1;
          document.getElementById("btn-slug-next").disabled = !cHasMore;
          return;
        }

        var page = Number(data.page || 1);
        var total = Number(data.total || 0);
        var pages = Number(data.total_pages || 0);
        if (!Number.isFinite(page) || page <= 0) page = 1;
        if (!Number.isFinite(pages) || pages < 0) pages = 0;

        state.slugPage = page;
        state.slugTotalPages = pages;

        document.getElementById("slug-page-meta").textContent = "第 " + page + " / " + Math.max(1, pages) + " 页，共 " + total + " 场";
        document.getElementById("btn-slug-prev").disabled = page <= 1;
        document.getElementById("btn-slug-next").disabled = pages <= 0 || page >= pages;
      }

      async function loadSlugs(action) {
        readSlugFilters();
        var nextAction = String(action || "refresh").trim().toLowerCase();
        var filterKey = buildSlugFilterKey();

        if (!state.slugFilterKey || state.slugFilterKey !== filterKey || nextAction === "reset") {
          state.slugFilterKey = filterKey;
          resetSlugCursorPaging();
          nextAction = "refresh";
        }

        var targetPage = state.slugPage || 1;
        var cursor = state.slugCursor || "";
        if (nextAction === "next") {
          if (!state.slugHasMore || !state.slugNextCursor) {
            renderSlugPager({ mode: "cursor", count: 0, has_more: false, items: [] });
            return;
          }
          targetPage += 1;
          cursor = state.slugNextCursor;
        } else if (nextAction === "prev") {
          if (targetPage <= 1) {
            targetPage = 1;
            cursor = "";
          } else {
            targetPage -= 1;
            cursor = state.slugCursorStack[targetPage - 1] || "";
          }
        } else {
          cursor = state.slugCursorStack[targetPage - 1] || state.slugCursor || "";
        }

        var params = new URLSearchParams();
        params.set("mode", "cursor");
        params.set("page_size", String(state.slugPageSize));
        if (cursor) params.set("cursor", cursor);
        if (state.slugKeyword) params.set("keyword", state.slugKeyword);
        if (state.slugTag && state.slugTag !== "all") params.set("tag", state.slugTag);

        try {
          var data = await request("/api/v1/slugs?" + params.toString());
          var items = Array.isArray(data.items) ? data.items : [];
          state.slugMode = "cursor";
          state.slugPage = targetPage;
          state.slugCursor = cursor;
          state.slugNextCursor = String(data.next_cursor || "").trim();
          state.slugHasMore = !!data.has_more;
          if (!Array.isArray(state.slugCursorStack) || !state.slugCursorStack.length) {
            state.slugCursorStack = [""];
          }
          state.slugCursorStack[targetPage - 1] = cursor;
          state.slugCursorStack = state.slugCursorStack.slice(0, targetPage);
          if (state.slugHasMore && state.slugNextCursor) {
            state.slugCursorStack[targetPage] = state.slugNextCursor;
          }
          renderSlugRows(items);
          renderSlugPager({
            mode: "cursor",
            count: Number(data.count || items.length),
            has_more: state.slugHasMore,
            items: items
          });
        } catch (err) {
          showNotice("slug 查询失败: " + err.message, true);
        }
      }

      async function fetchAllEventsBySlug(slug) {
        var limit = 200;
        var maxPages = 30;
        var all = [];
        var total = 0;
        var page = 1;

        while (page <= maxPages) {
          var params = new URLSearchParams();
          params.set("slug", slug);
          params.set("limit", String(limit));
          params.set("page", String(page));

          var data = await request("/api/v1/events?" + params.toString());
          var items = Array.isArray(data.items) ? data.items : [];

          if (page === 1) {
            total = Number(data.total || 0);
          }

          all = all.concat(items);

          if (!items.length || all.length >= total) {
            break;
          }
          page += 1;
        }

        return {
          items: all,
          total: total,
          truncated: all.length < total,
          slugs: slug ? [slug] : []
        };
      }

      function renderComboBody(tbodyID, list) {
        var tbody = document.getElementById(tbodyID);
        if (!tbody) return;
        if (!list || !list.length) {
          tbody.innerHTML = "<tr><td colspan=\"6\" class=\"empty-cell\">暂无</td></tr>";
          return;
        }

        var sorted = list.slice().sort(function (a, b) { return b.ts - a.ts; });
        var showLimit = 80;
        var show = sorted.slice(0, showLimit);

        var html = show.map(function (r) {
          return "<tr>" +
            "<td>" + escapeHTML(fmtMS(r.ts)) + "</td>" +
            "<td>" + escapeHTML(fmtNum(r.price, 2)) + "</td>" +
            "<td>" + escapeHTML(fmtNum(r.qty, 2)) + "</td>" +
            "<td>" + escapeHTML(fmtNum(r.usdc, 2)) + "</td>" +
            "<td>" + escapeHTML(String((r.raw && r.raw.outcome) || "-")) + "</td>" +
            "<td><code>" + escapeHTML(shortTx(r.raw && r.raw.transaction_hash)) + "</code></td>" +
            "</tr>";
        }).join("");

        if (sorted.length > show.length) {
          html += "<tr><td colspan=\"6\" class=\"empty-cell\">仅展示最新 " + show.length + " 条，共 " + sorted.length + " 条</td></tr>";
        }

        tbody.innerHTML = html;
      }

      function renderWaveRows(waves) {
        var tbody = document.getElementById("wave-rows");
        if (!waves || !waves.length) {
          tbody.innerHTML = "<tr><td colspan=\"10\" class=\"empty-cell\">暂无波段数据</td></tr>";
          return;
        }

        var html = waves.map(function (w) {
          return "<tr>" +
            "<td>" + escapeHTML(String(w.index)) + "</td>" +
            "<td>" + escapeHTML(fmtMS(w.startTs)) + "</td>" +
            "<td>" + escapeHTML(fmtMS(w.endTs)) + "</td>" +
            "<td>" + escapeHTML(String(w.spanSec)) + "</td>" +
            "<td>" + escapeHTML(String(w.records)) + "</td>" +
            "<td>" + escapeHTML(w.orderHint) + "</td>" +
            "<td>" + escapeHTML(w.firstCombo + " -> " + w.lastCombo) + "</td>" +
            "<td>" + escapeHTML(fmtNum(w.buyPnl, 2)) + "</td>" +
            "<td>" + escapeHTML(fmtNum(w.sellPnl, 2)) + "</td>" +
            "<td>" + escapeHTML(fmtNum(w.totalPnl, 2)) + "</td>" +
            "</tr>";
        }).join("");

        tbody.innerHTML = html;
      }

      function renderDetail(slug, fetched, analyzed, strategyInfo) {
        document.getElementById("detail-slug").value = slug || "";

        var totalRecords = fetched.items.length;
        var buy = analyzed.buyPair;
        var sell = analyzed.sellPair;

        var combinedSymmetry = 0;
        if (buy.symmetryRate > 0 || sell.symmetryRate > 0) {
          combinedSymmetry = (buy.symmetryRate + sell.symmetryRate) / ((buy.symmetryRate > 0 && sell.symmetryRate > 0) ? 2 : 1);
          if (!Number.isFinite(combinedSymmetry)) combinedSymmetry = 0;
        }

        var pnlClass = analyzed.totalPnl >= 0 ? "tone-up" : "tone-down";
        var symmetryClass = combinedSymmetry >= 0.9 ? "tone-up" : "tone-warn";

        var kpis = [
          { k: "总记录", v: String(totalRecords), c: "" },
          { k: "可配对数量", v: fmtNum(analyzed.totalPairedQty, 2), c: "" },
          { k: "估算总盈亏", v: fmtNum(analyzed.totalPnl, 2), c: pnlClass },
          { k: "平均对称度", v: fmtNum(combinedSymmetry * 100, 2) + "%", c: symmetryClass }
        ];

        document.getElementById("detail-kpis").innerHTML = kpis.map(function (x) {
          return "<div class=\"kpi\">" +
            "<div class=\"k\">" + escapeHTML(x.k) + "</div>" +
            "<div class=\"v " + escapeHTML(x.c || "") + "\">" + escapeHTML(x.v) + "</div>" +
            "</div>";
        }).join("");

        var lines = [];
        lines.push("结论: " + analyzed.verdict + "；" + analyzed.riskText);
        lines.push("Buy 配对公式: (1 - (up_price + down_price)) * 配对数量。当前配对数量=" + fmtNum(buy.pairedQty, 2) + "，估算=" + fmtNum(buy.pairedPnl, 2) + "，对称度=" + fmtNum(buy.symmetryRate * 100, 2) + "%");
        lines.push("Sell 配对公式: ((up_price + down_price) - 1) * 配对数量。当前配对数量=" + fmtNum(sell.pairedQty, 2) + "，估算=" + fmtNum(sell.pairedPnl, 2) + "，对称度=" + fmtNum(sell.symmetryRate * 100, 2) + "%");
        lines.push("示例: 如果 buy-up=0.05 且 buy-down=0.92，则 1-(0.05+0.92)=0.03，理论上这对是赚钱的（前提是数量可配对）。");
        lines.push("未知方向记录: " + analyzed.unknown.length + " 条（outcome 无法识别为 up/down 时计入）。");

        document.getElementById("detail-summary").innerHTML = lines.map(function (x) {
          return "<div class=\"line\">" + escapeHTML(x) + "</div>";
        }).join("");

        var truncatedText = fetched.truncated ? "注意: 当前只分析了前 " + fetched.items.length + " 条，库内总量约 " + fetched.total + " 条（已触发保护上限）。" : "分析覆盖记录: " + fetched.items.length + " / " + fetched.total;
        document.getElementById("detail-note").textContent = truncatedText;

        var groupText = "分组: 单场 slug";
        if (strategyInfo) {
          groupText = "分组: " + strategyInfo.symbol.toUpperCase() + " | " + fmtTimeRangeSec(strategyInfo.startSec, strategyInfo.endSec);
        }
        if (Array.isArray(fetched.slugs) && fetched.slugs.length) {
          groupText += " | 本组盘口: " + fetched.slugs.join(", ");
        }
        document.getElementById("detail-group-meta").textContent = groupText;

        renderComboBody("combo-buy-up", analyzed.buckets.buy_up);
        renderComboBody("combo-buy-down", analyzed.buckets.buy_down);
        renderComboBody("combo-sell-up", analyzed.buckets.sell_up);
        renderComboBody("combo-sell-down", analyzed.buckets.sell_down);

        document.getElementById("meta-buy-up").textContent = "count=" + analyzed.buckets.buy_up.length;
        document.getElementById("meta-buy-down").textContent = "count=" + analyzed.buckets.buy_down.length;
        document.getElementById("meta-sell-up").textContent = "count=" + analyzed.buckets.sell_up.length;
        document.getElementById("meta-sell-down").textContent = "count=" + analyzed.buckets.sell_down.length;

        renderWaveRows(analyzed.waves);
      }

      async function loadSlugDetail(slug, marketTag) {
        if (!slug) return;

        var strategyInfo = parseStrategyGroupBySlug(slug, marketTag);
        state.selectedStrategy = strategyInfo;

        document.getElementById("detail-slug").value = slug;
        document.getElementById("detail-note").textContent = "正在加载比赛数据并分析...";
        document.getElementById("detail-summary").innerHTML = "<div class=\"line\">加载中...</div>";
        document.getElementById("detail-group-meta").textContent = "分组: 计算中...";

        try {
          var fetched = null;
          if (strategyInfo) {
            fetched = await fetchEventsByStrategyGroup(strategyInfo);
          } else {
            fetched = await fetchAllEventsBySlug(slug);
          }
          var waveGapSec = clamp(toInt(document.getElementById("detail-wave-gap").value, 120), 10, 900);
          var analyzed = analyzeSlugEvents(fetched.items, waveGapSec);
          renderDetail(slug, fetched, analyzed, strategyInfo);
        } catch (err) {
          document.getElementById("detail-note").textContent = "分析失败";
          document.getElementById("detail-group-meta").textContent = "分组: -";
          showNotice("加载比赛详情失败: " + err.message, true);
        }
      }

      async function selectSlug(slug, marketTag) {
        if (!slug) return;
        state.selectedSlug = slug;
        document.getElementById("q-slug").value = slug;
        await loadSlugDetail(slug, marketTag);
        await loadEvents(1);
        await loadSlugs("refresh");
      }

      function readEventFilters() {
        var limit = clamp(toInt(document.getElementById("q-limit").value, state.eventPageSize || 100), 1, 500);
        var page = clamp(toInt(document.getElementById("q-page").value, state.eventPage || 1), 1, 999999);

        return {
          limit: limit,
          page: page,
          slug: document.getElementById("q-slug").value.trim(),
          type: document.getElementById("q-type").value.trim(),
          side: document.getElementById("q-side").value.trim()
        };
      }

      function renderEventRows(items) {
        var tbody = document.getElementById("rows");
        if (!Array.isArray(items) || !items.length) {
          tbody.innerHTML = "<tr><td colspan=\"10\" class=\"empty-cell\">暂无数据</td></tr>";
          return;
        }

        var html = items.map(function (e) {
          var eventTime = e.timestamp_ms ? fmtMS(e.timestamp_ms) : fmtTime(e.event_time);
          return "<tr>" +
            "<td>" + escapeHTML(eventTime) + "</td>" +
            "<td>" + escapeHTML(e.activity_type || "-") + "</td>" +
            "<td>" + escapeHTML(e.side || "-") + "</td>" +
            "<td>" + escapeHTML(e.slug || "-") + "</td>" +
            "<td>" + escapeHTML(e.outcome || "-") + "</td>" +
            "<td>" + escapeHTML(fmtNumFromRaw(e.price)) + "</td>" +
            "<td>" + escapeHTML(fmtNumFromRaw(e.usdc_size)) + "</td>" +
            "<td>" + escapeHTML(fmtNumFromRaw(e.size)) + "</td>" +
            "<td><code>" + escapeHTML(shortTx(e.transaction_hash)) + "</code></td>" +
            "<td>" + escapeHTML(String(e.source_offset || 0) + "/" + String(e.source_index || 0)) + "</td>" +
            "</tr>";
        }).join("");

        tbody.innerHTML = html;
      }

      function renderEventPager(data) {
        var page = Number(data.page || 1);
        var pages = Number(data.total_pages || 0);
        var total = Number(data.total || 0);

        if (!Number.isFinite(page) || page <= 0) page = 1;
        if (!Number.isFinite(pages) || pages < 0) pages = 0;

        state.eventPage = page;
        state.eventTotalPages = pages;

        document.getElementById("q-page").value = String(page);
        document.getElementById("events-page-meta").textContent = "第 " + page + " / " + Math.max(1, pages) + " 页，共 " + total + " 条";
        document.getElementById("btn-events-prev").disabled = page <= 1;
        document.getElementById("btn-events-next").disabled = pages <= 0 || page >= pages;
      }

      async function loadEvents(pageOverride) {
        var f = readEventFilters();
        if (pageOverride != null) {
          f.page = clamp(toInt(pageOverride, 1), 1, 999999);
        }

        state.eventPageSize = f.limit;

        var params = new URLSearchParams();
        params.set("limit", String(f.limit));
        params.set("page", String(f.page));
        if (f.slug) params.set("slug", f.slug);
        if (f.type) params.set("type", f.type);
        if (f.side) params.set("side", f.side);

        var endpoint = "/api/v1/events?" + params.toString();
        document.getElementById("query-meta").textContent = "query: " + endpoint;

        try {
          var data = await request(endpoint);
          renderEventRows(data.items || []);
          renderEventPager(data);
        } catch (err) {
          showNotice("事件查询失败: " + err.message, true);
        }
      }

      function resetAutoRefresh() {
        if (timer) {
          clearInterval(timer);
          timer = null;
        }
        if (autoRefreshBox.checked) {
          timer = setInterval(function () {
            loadStatus();
            loadSlugs("refresh");
          }, 10000);
        }
      }

      document.getElementById("btn-refresh-status").addEventListener("click", function () {
        loadStatus();
      });
      document.getElementById("btn-sync-fast").addEventListener("click", function () {
        triggerSync("fast");
      });
      document.getElementById("btn-sync-backfill").addEventListener("click", function () {
        triggerSync("backfill");
      });
      autoRefreshBox.addEventListener("change", resetAutoRefresh);

      document.getElementById("btn-slug-search").addEventListener("click", function () {
        loadSlugs("reset");
      });
      document.getElementById("btn-slug-prev").addEventListener("click", function () {
        loadSlugs("prev");
      });
      document.getElementById("btn-slug-next").addEventListener("click", function () {
        loadSlugs("next");
      });
      document.getElementById("slug-keyword").addEventListener("keydown", function (ev) {
        if (ev.key === "Enter") {
          ev.preventDefault();
          loadSlugs("reset");
        }
      });
      document.getElementById("slug-tag").addEventListener("change", function () {
        loadSlugs("reset");
      });

      document.getElementById("slug-rows").addEventListener("click", function (ev) {
        var btn = ev.target.closest(".btn-open-slug");
        if (!btn) return;
        var slug = decodeURIComponent(btn.getAttribute("data-slug") || "");
        var marketTag = decodeURIComponent(btn.getAttribute("data-tag") || "");
        if (!slug) return;
        selectSlug(slug, marketTag);
      });

      document.getElementById("detail-refresh").addEventListener("click", function () {
        if (!state.selectedSlug) {
          showNotice("请先在左侧选择一个 slug", true);
          return;
        }
        loadSlugDetail(state.selectedSlug, state.selectedStrategy && state.selectedStrategy.symbol);
      });

      document.getElementById("btn-query").addEventListener("click", function () {
        loadEvents(readEventFilters().page);
      });
      document.getElementById("btn-events-prev").addEventListener("click", function () {
        loadEvents(Math.max(1, (state.eventPage || 1) - 1));
      });
      document.getElementById("btn-events-next").addEventListener("click", function () {
        loadEvents((state.eventPage || 1) + 1);
      });
      document.getElementById("btn-query-reset").addEventListener("click", function () {
        document.getElementById("q-type").value = "";
        document.getElementById("q-side").value = "";
        if (!state.selectedSlug) {
          document.getElementById("q-slug").value = "";
        }
        document.getElementById("q-page").value = "1";
        loadEvents(1);
      });

      loadStatus();
      loadSlugs("reset");
      loadEvents(1);
      resetAutoRefresh();
    })();
  </script>
</body>
</html>
`
