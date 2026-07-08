import type { Route } from "../types";

export function EmptyStatePage({ setRoute }: { setRoute: (route: Route) => void }) {
  return (
    <>
      <header className="page-header">
        <div>
          <span className="eyebrow">First-run state</span>
          <h1>No monitors yet</h1>
          <p className="muted-copy">The message stays direct and product-facing.</p>
        </div>
      </header>
      <section className="card empty-card">
        <div className="empty-signal" aria-hidden="true" />
        <h2>Add one to start watching its pulse.</h2>
        <p>The first scheduler pass will append the initial check. After that, the row moves from waiting into either a healthy or down state without a refresh.</p>
        <div className="button-row center">
          <button className="button button-primary" onClick={() => setRoute("add")}>Add monitor</button>
          <button className="button button-secondary" onClick={() => setRoute("dashboard")}>Back to dashboard</button>
        </div>
      </section>
    </>
  );
}
