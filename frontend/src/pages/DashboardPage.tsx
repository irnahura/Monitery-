import { InfoRow, StatCard, StatusPill } from "../components/common";
import type { HealthCheck, Monitor, Route, Summary } from "../types";
import { formatDate, relativeTime } from "../utils/format";

export function DashboardPage(props: {
  monitors: Monitor[];
  history: HealthCheck[];
  summary: Summary;
  selected: Monitor | null;
  openDetail: (monitor: Monitor) => void;
  setRoute: (route: Route) => void;
}) {
  const activeCount = props.monitors.filter((monitor) => monitor.latest_is_up).length;
  const responses = props.monitors.map((monitor) => monitor.latest_response_time_ms).filter((value): value is number => typeof value === "number" && value > 0).sort((a, b) => a - b);
  const median = responses.length ? responses[Math.floor(responses.length / 2)] : 0;
  return (
    <>
      <header className="page-header">
        <div>
          <span className="eyebrow">Monitor fleet</span>
          <h1>Dashboard</h1>
          <p className="muted-copy">A table-first view built for a few dozen URLs, not hundreds.</p>
        </div>
        <div className="header-actions">
          <button className="button button-secondary" onClick={() => props.setRoute("settings")}>Settings</button>
          <button className="button button-primary" onClick={() => props.setRoute("add")}>Add monitor</button>
        </div>
      </header>
      <section className="stats-grid">
        <StatCard label="Active monitors" value={String(props.monitors.length)} copy={`${activeCount} responding, ${props.monitors.length - activeCount} waiting or down.`} />
        <StatCard label="Median response" value={median ? `${median}ms` : "-"} copy="Across the latest successful checks." />
        <StatCard label="Availability" value={`${props.summary.availability_percent.toFixed(2)}%`} copy="Current selected monitor history." />
      </section>
      <section className="card">
        <div className="card-header">
          <div>
            <h2>Monitors</h2>
            <p className="muted-copy">Rows link into the retained detail view.</p>
          </div>
          <span className="pill pill-neutral">Last refresh {relativeTime(new Date())}</span>
        </div>
        <div className="monitor-table" role="table" aria-label="Monitors">
          <div className="table-head" role="row"><span>Name</span><span>Status</span><span>Response</span><span>Last checked</span></div>
          {props.monitors.map((monitor) => <div className="table-row" key={monitor.id} role="row" onClick={() => props.openDetail(monitor)}><div><strong>{monitor.name}</strong><span>{monitor.url}</span></div><StatusPill monitor={monitor} /><span className="mono">{monitor.latest_is_up && monitor.latest_response_time_ms ? `${monitor.latest_response_time_ms}ms` : "-"}</span><span className="mono muted">{monitor.last_checked_at ? relativeTime(new Date(monitor.last_checked_at)) : "waiting"}</span></div>)}
        </div>
      </section>
      <section className="split-grid">
        <article className="card">
          <div className="card-header">
            <div>
              <h2>Recent incidents</h2>
              <p className="muted-copy">Keep the product honest with simple status history.</p>
            </div>
          </div>
          <div className="stack-list">
            {props.history.filter((check) => !check.is_up).slice(0, 2).map((check) => <div className="list-row" key={check.id}><div><strong>{props.selected?.name || "Monitor"} was down</strong><span>Status {check.status_code || "unreachable"} on scheduler pass.</span></div><span className="mono muted">{formatDate(check.checked_at)}</span></div>)}
            {!props.history.some((check) => !check.is_up) && <InfoRow title="No incidents in recent history" copy="Latest retained checks are inside the healthy range." />}
          </div>
        </article>
        <article className="card">
          <div className="card-header">
            <div>
              <h2>Quick actions</h2>
              <p className="muted-copy">Minimal operator shortcuts, no extra tooling panel.</p>
            </div>
          </div>
          <div className="action-group">
            <button className="button button-secondary" onClick={() => props.setRoute("add")}>Register another URL</button>
            <button className="button button-secondary" onClick={() => props.setRoute("empty")}>View first-run empty state</button>
            <button className="button button-secondary" onClick={() => props.setRoute("settings")}>Review API keys</button>
          </div>
        </article>
      </section>
    </>
  );
}
