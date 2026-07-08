import { InfoRow, Meta, StatCard, StatusPill } from "../components/common";
import type { HealthCheck, Monitor, Route, Summary } from "../types";
import { formatDate, relativeTime } from "../utils/format";

export function ServiceDetailPage(props: {
  monitor?: Monitor;
  history: HealthCheck[];
  summary: Summary;
  setRoute: (route: Route) => void;
  deleteMonitor: (id: number) => Promise<void>;
}) {
  if (!props.monitor) return <div />;
  const monitor = props.monitor;
  return (
    <>
      <header className="page-header">
        <div>
          <span className="eyebrow">Service inspection</span>
          <h1>{monitor.name}</h1>
          <p className="muted-copy">{monitor.url}</p>
        </div>
        <div className="header-actions">
          <button className="button button-secondary" onClick={() => props.setRoute("dashboard")}>Back to dashboard</button>
          <button className="button button-tertiary" onClick={() => void props.deleteMonitor(monitor.id)}>Delete monitor</button>
        </div>
      </header>
      <section className="stats-grid">
        <StatCard label="Availability" value={`${props.summary.availability_percent.toFixed(2)}%`} copy="Trailing retained uptime." />
        <StatCard label="SLA target" value={`${props.summary.sla_percent.toFixed(2)}%`} copy="Current history clears the target." />
        <StatCard label="Latest response" value={monitor.latest_response_time_ms ? `${monitor.latest_response_time_ms}ms` : "-"} copy={`Status ${monitor.latest_status_code || "waiting"} - ${monitor.last_checked_at ? relativeTime(new Date(monitor.last_checked_at)) : "not checked yet"}.`} />
      </section>
      <section className="split-grid">
        <article className="card">
          <div className="card-header">
            <div>
              <h2>Latest observation</h2>
              <p className="muted-copy">Joined from the retained monitor response.</p>
            </div>
            <StatusPill monitor={monitor} />
          </div>
          <div className="meta-grid">
            <Meta label="Checked at" value={monitor.last_checked_at ? formatDate(monitor.last_checked_at) : "waiting"} />
            <Meta label="Response time" value={monitor.latest_response_time_ms ? `${monitor.latest_response_time_ms}ms` : "-"} />
            <Meta label="Failure rule" value="Timeout or non-2xx/3xx -> DOWN" />
          </div>
        </article>
        <article className="card">
          <div className="card-header">
            <div>
              <h2>Recent note</h2>
              <p className="muted-copy">The product keeps the explanation short and operational.</p>
            </div>
          </div>
          <p className="body-copy">{monitor.latest_is_up === false ? "This service is currently down. The next scheduler pass will record whether it recovered." : "This service is currently inside the healthy range. Recent checks remain available below."}</p>
        </article>
      </section>
      <section className="card">
        <div className="card-header">
          <div>
            <h2>Recent health check history</h2>
            <p className="muted-copy">Simple table, no heavy analytics shell.</p>
          </div>
        </div>
        <div className="history-table" role="table" aria-label="Recent health checks">
          <div className="table-head" role="row"><span>Checked at</span><span>Status</span><span>Response</span></div>
          {props.history.map((check) => <div className="table-row static-row" role="row" key={check.id}><div><strong>{formatDate(check.checked_at)}</strong><span>{check.is_up ? "Healthy response" : "Failed response"}</span></div><span className={`pill ${check.is_up ? "pill-up" : "pill-down"}`}>{check.status_code || "DOWN"}</span><span className="mono">{check.is_up ? `${check.response_time_ms}ms` : "-"}</span></div>)}
          {!props.history.length && <InfoRow title="No checks yet" copy="The scheduler will add the first health row on its next pass." />}
        </div>
      </section>
    </>
  );
}
