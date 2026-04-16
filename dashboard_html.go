package main

const dashboardHTML = `<!doctype html>
<html lang="zh-CN">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>PM Activity Console</title>
  <style>
    @import url("https://fonts.googleapis.com/css2?family=Space+Grotesk:wght@400;500;700&family=IBM+Plex+Mono:wght@400;500&display=swap");
    :root {
      --bg: #f5f7ff;
      --panel: #ffffff;
      --ink: #182039;
      --muted: #5f6782;
      --line: #d9def0;
      --ok: #0a8f5d;
      --warn: #d17d00;
      --bad: #bb2a31;
      --accent: #2d5bff;
      --accent-soft: #dfe8ff;
      --shadow: 0 10px 30px rgba(30, 47, 110, 0.12);
    }
    * { box-sizing: border-box; }
    body {
      margin: 0;
      font-family: "Space Grotesk", "Avenir Next", "Segoe UI", sans-serif;
      color: var(--ink);
      background:
        radial-gradient(1200px 550px at 80% -10%, #d9e4ff 0%, transparent 65%),
        radial-gradient(1100px 500px at -10% 5%, #ffe5cc 0%, transparent 60%),
        var(--bg);
      min-height: 100vh;
      line-height: 1.45;
    }
    .page {
      width: min(1200px, calc(100vw - 32px));
      margin: 22px auto 34px;
    }
    .head {
      display: flex;
      justify-content: space-between;
      align-items: center;
      gap: 12px;
      flex-wrap: wrap;
      margin-bottom: 14px;
    }
    .title {
      margin: 0;
      font-size: clamp(24px, 4vw, 32px);
      font-weight: 700;
      letter-spacing: -0.03em;
    }
    .sub {
      margin: 2px 0 0;
      color: var(--muted);
      font-size: 14px;
    }
    .badge {
      border: 1px solid var(--line);
      border-radius: 999px;
      padding: 7px 12px;
      background: #fff;
      font-size: 13px;
      color: var(--muted);
      font-family: "IBM Plex Mono", Menlo, monospace;
    }
    .grid {
      display: grid;
      grid-template-columns: repeat(4, minmax(0, 1fr));
      gap: 10px;
      margin-bottom: 12px;
    }
    .card, .panel {
      background: var(--panel);
      border: 1px solid var(--line);
      border-radius: 14px;
      box-shadow: var(--shadow);
    }
    .card { padding: 12px 14px; }
    .card .k { font-size: 12px; color: var(--muted); }
    .card .v {
      margin-top: 4px;
      font-family: "IBM Plex Mono", Menlo, monospace;
      font-size: 18px;
      font-weight: 500;
      word-break: break-word;
    }
    .panel { padding: 14px; margin-top: 10px; }
    .panel h2 {
      margin: 0 0 10px;
      font-size: 17px;
      letter-spacing: -0.01em;
    }
    .row {
      display: grid;
      grid-template-columns: repeat(6, minmax(0, 1fr));
      gap: 10px;
    }
    label {
      display: block;
      font-size: 12px;
      color: var(--muted);
      margin-bottom: 4px;
    }
    input, select, button {
      border-radius: 10px;
      border: 1px solid var(--line);
      background: #fff;
      color: var(--ink);
      font: inherit;
      padding: 9px 10px;
    }
    input, select { width: 100%; }
    button {
      cursor: pointer;
      transition: transform .08s ease, background .15s ease;
      font-weight: 600;
    }
    button:active { transform: translateY(1px); }
    .btn-primary {
      background: var(--accent);
      color: #fff;
      border-color: transparent;
    }
    .btn-soft {
      background: var(--accent-soft);
      color: #213ea8;
      border-color: transparent;
    }
    .btn-warning {
      background: #fff1dd;
      color: #7d4a00;
      border-color: transparent;
    }
    .ops {
      display: flex;
      gap: 8px;
      flex-wrap: wrap;
      margin-top: 10px;
      align-items: center;
    }
    .hint {
      font-size: 12px;
      color: var(--muted);
      margin-left: auto;
    }
    .mono {
      font-family: "IBM Plex Mono", Menlo, monospace;
      font-size: 12px;
      color: var(--muted);
      margin-top: 8px;
      word-break: break-all;
    }
    .notice {
      margin-top: 8px;
      padding: 8px 10px;
      border-radius: 10px;
      font-size: 13px;
      display: none;
    }
    .notice.show { display: block; }
    .notice.ok { background: #e5f9ee; color: var(--ok); }
    .notice.err { background: #ffe7e8; color: var(--bad); }
    .table-wrap {
      overflow: auto;
      border: 1px solid var(--line);
      border-radius: 12px;
      margin-top: 10px;
      background: #fff;
    }
    table {
      width: 100%;
      border-collapse: collapse;
      min-width: 980px;
    }
    th, td {
      border-bottom: 1px solid #eef1fb;
      padding: 8px 10px;
      text-align: left;
      font-size: 13px;
      vertical-align: top;
    }
    th {
      position: sticky;
      top: 0;
      z-index: 1;
      background: #f8faff;
      font-size: 12px;
      color: #4b567a;
      text-transform: uppercase;
      letter-spacing: 0.02em;
    }
    td code {
      font-family: "IBM Plex Mono", Menlo, monospace;
      font-size: 12px;
      color: #33468f;
    }
    .status-grid {
      display: grid;
      grid-template-columns: repeat(3, minmax(0, 1fr));
      gap: 8px;
    }
    .status-item {
      padding: 8px;
      background: #f7f9ff;
      border-radius: 10px;
      border: 1px solid #ebeffc;
    }
    .status-item .s-k { font-size: 12px; color: var(--muted); }
    .status-item .s-v {
      margin-top: 2px;
      font-family: "IBM Plex Mono", Menlo, monospace;
      font-size: 13px;
      word-break: break-word;
    }
    @media (max-width: 1024px) {
      .grid { grid-template-columns: repeat(2, minmax(0, 1fr)); }
      .row { grid-template-columns: repeat(2, minmax(0, 1fr)); }
      .status-grid { grid-template-columns: 1fr; }
      .hint { width: 100%; margin-left: 0; }
    }
    @media (max-width: 640px) {
      .page { width: calc(100vw - 20px); margin: 12px auto 20px; }
      .grid { grid-template-columns: 1fr; }
      .row { grid-template-columns: 1fr; }
      .panel { padding: 10px; }
    }
  </style>
</head>
<body>
  <main class="page">
    <header class="head">
      <div>
        <h1 class="title">PM Activity Console</h1>
        <p class="sub">采集状态看板 + 手动同步 + 事件查询</p>
      </div>
      <div class="badge" id="badge-health">service: checking</div>
    </header>

    <section class="grid">
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
      <h2>服务状态</h2>
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
        <button class="btn-primary" id="btn-refresh-status">刷新状态</button>
        <button class="btn-soft" id="btn-sync-fast">触发 Fast</button>
        <button class="btn-warning" id="btn-sync-backfill">触发 Backfill</button>
        <label style="display:flex;align-items:center;gap:6px;margin:0 0 0 6px;color:#445078;font-size:13px;">
          <input id="auto-refresh" type="checkbox" checked style="width:auto;padding:0;margin:0;">
          自动刷新(10s)
        </label>
        <span class="hint" id="status-config">listen=- / wallet=- / limit=-</span>
      </div>
      <div class="notice" id="notice"></div>
    </section>

    <section class="panel">
      <h2>事件查询</h2>
      <div class="row">
        <div>
          <label for="q-limit">limit</label>
          <input id="q-limit" type="number" min="1" max="500" value="100" />
        </div>
        <div>
          <label for="q-offset">offset</label>
          <input id="q-offset" type="number" min="0" value="0" />
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
          <input id="q-side" type="text" placeholder="buy / sell ..." />
        </div>
        <div style="display:flex;align-items:flex-end;">
          <button class="btn-primary" id="btn-query" style="width:100%;">查询事件</button>
        </div>
      </div>
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
            <tr><td colspan="10" style="text-align:center;color:#667;">加载中...</td></tr>
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

      function escapeHTML(v) {
        return String(v == null ? "" : v)
          .replace(/&/g, "&amp;")
          .replace(/</g, "&lt;")
          .replace(/>/g, "&gt;")
          .replace(/"/g, "&quot;")
          .replace(/'/g, "&#39;");
      }

      function showNotice(msg, isError) {
        notice.textContent = msg;
        notice.className = "notice show " + (isError ? "err" : "ok");
        setTimeout(function () {
          notice.className = "notice";
        }, 3000);
      }

      function fmtTime(raw) {
        if (!raw) return "-";
        var d = new Date(raw);
        if (isNaN(d.getTime())) return "-";
        return d.toLocaleString();
      }

      function fmtMS(ms) {
        if (!ms || ms <= 0) return "-";
        var d = new Date(ms);
        if (isNaN(d.getTime())) return "-";
        return d.toLocaleString();
      }

      function shortTx(tx) {
        if (!tx) return "-";
        var s = String(tx);
        if (s.length <= 14) return s;
        return s.slice(0, 8) + "..." + s.slice(-6);
      }

      async function request(path, options) {
        var res = await fetch(path, options || {});
        var body = await res.json();
        if (!res.ok || body.code !== 20000) {
          throw new Error((body && body.message) || ("HTTP " + res.status));
        }
        return body.data;
      }

      function readEventQuery() {
        var p = new URLSearchParams();
        var limit = document.getElementById("q-limit").value.trim();
        var offset = document.getElementById("q-offset").value.trim();
        var slug = document.getElementById("q-slug").value.trim();
        var type = document.getElementById("q-type").value.trim();
        var side = document.getElementById("q-side").value.trim();
        if (limit) p.set("limit", limit);
        if (offset) p.set("offset", offset);
        if (slug) p.set("slug", slug);
        if (type) p.set("type", type);
        if (side) p.set("side", side);
        return p.toString();
      }

      async function loadStatus() {
        try {
          var data = await request("/api/v1/status");
          var poller = data.poller || {};
          var cfg = data.config || {};
          badgeHealth.textContent = "service: online";
          badgeHealth.style.color = "#0a8f5d";

          document.getElementById("v-stored").textContent = String(data.stored_events || 0);
          document.getElementById("v-last").textContent = String(poller.last_fetched_rows || 0) + " / " + String(poller.last_inserted_rows || 0);
          document.getElementById("v-total").textContent = String(poller.total_fetched_rows || 0) + " / " + String(poller.total_inserted_rows || 0) + " / " + String(poller.total_duplicate_rows || 0);
          document.getElementById("v-queue").textContent = String(poller.queue_len || 0) + " / " + (poller.running ? "running" : "idle");

          document.getElementById("s-success").textContent = fmtTime(poller.last_success_at);
          document.getElementById("s-error").textContent = poller.last_error ? fmtTime(poller.last_error_at) + " | " + poller.last_error : "-";
          document.getElementById("s-ts").textContent = fmtMS(poller.last_newest_ts) + " / " + fmtMS(poller.last_oldest_ts);

          document.getElementById("status-config").textContent =
            "listen=" + (cfg.listen_addr || "-") +
            " / wallet=" + (cfg.user_wallet || "-") +
            " / limit=" + (cfg.page_limit || "-");
        } catch (err) {
          badgeHealth.textContent = "service: offline";
          badgeHealth.style.color = "#bb2a31";
          showNotice("状态刷新失败: " + err.message, true);
        }
      }

      function renderRows(items) {
        var tbody = document.getElementById("rows");
        if (!Array.isArray(items) || items.length === 0) {
          tbody.innerHTML = "<tr><td colspan=\"10\" style=\"text-align:center;color:#667;\">暂无数据</td></tr>";
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
            "<td>" + escapeHTML(e.price || "-") + "</td>" +
            "<td>" + escapeHTML(e.usdc_size || "-") + "</td>" +
            "<td>" + escapeHTML(e.size || "-") + "</td>" +
            "<td><code>" + escapeHTML(shortTx(e.transaction_hash)) + "</code></td>" +
            "<td>" + escapeHTML(String(e.source_offset || 0) + "/" + String(e.source_index || 0)) + "</td>" +
            "</tr>";
        }).join("");
        tbody.innerHTML = html;
      }

      async function loadEvents() {
        var query = readEventQuery();
        var endpoint = "/api/v1/events" + (query ? ("?" + query) : "");
        document.getElementById("query-meta").textContent = "query: " + endpoint;
        try {
          var data = await request(endpoint);
          renderRows(data.items || []);
        } catch (err) {
          showNotice("事件查询失败: " + err.message, true);
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

      function resetAutoRefresh() {
        if (timer) {
          clearInterval(timer);
          timer = null;
        }
        if (autoRefreshBox.checked) {
          timer = setInterval(function () {
            loadStatus();
            loadEvents();
          }, 10000);
        }
      }

      document.getElementById("btn-refresh-status").addEventListener("click", function () {
        loadStatus();
      });
      document.getElementById("btn-query").addEventListener("click", function () {
        loadEvents();
      });
      document.getElementById("btn-sync-fast").addEventListener("click", function () {
        triggerSync("fast");
      });
      document.getElementById("btn-sync-backfill").addEventListener("click", function () {
        triggerSync("backfill");
      });
      autoRefreshBox.addEventListener("change", resetAutoRefresh);

      loadStatus();
      loadEvents();
      resetAutoRefresh();
    })();
  </script>
</body>
</html>
`
