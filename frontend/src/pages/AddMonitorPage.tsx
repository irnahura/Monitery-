import type { FormEvent } from "react";
import { Field, InfoRow, MessageLine } from "../components/common";
import type { Message, Route } from "../types";

export function AddMonitorPage(props: {
  form: { name: string; url: string; interval_seconds: number; timeout_seconds: number };
  setForm: (form: { name: string; url: string; interval_seconds: number; timeout_seconds: number }) => void;
  submit: (event: FormEvent) => void;
  message: Message;
  setRoute: (route: Route) => void;
}) {
  return (
    <>
      <header className="page-header">
        <div>
          <span className="eyebrow">Create monitor</span>
          <h1>Add a URL</h1>
          <p className="muted-copy">The flow stays short: name, URL, validate, return.</p>
        </div>
      </header>
      <div className="form-layout">
        <section className="card">
          <form className="form-stack" onSubmit={props.submit} noValidate>
            <Field label="Name"><input value={props.form.name} placeholder="Example (up)" onChange={(event) => props.setForm({ ...props.form, name: event.target.value })} /></Field>
            <Field label="URL"><input type="url" value={props.form.url} placeholder="https://example.com" onChange={(event) => props.setForm({ ...props.form, url: event.target.value })} /></Field>
            <Field label="Interval seconds"><input type="number" min="60" value={props.form.interval_seconds} onChange={(event) => props.setForm({ ...props.form, interval_seconds: Number(event.target.value) })} /></Field>
            <MessageLine message={props.message} />
            <div className="button-row">
              <button className="button button-secondary" type="button" onClick={() => props.setRoute("dashboard")}>Cancel</button>
              <button className="button button-primary" type="submit">Add monitor</button>
            </div>
          </form>
        </section>
        <aside className="card detail-side">
          <div className="card-header">
            <div>
              <h2>What happens next</h2>
              <p className="muted-copy">This mirrors the design doc and expansion plan.</p>
            </div>
          </div>
          <div className="stack-list">
            <InfoRow title="First check queues automatically" copy="The next scheduler pass inserts the first health row." />
            <InfoRow title="Failures still create history" copy="Down checks show - for response time instead of fake numbers." />
            <InfoRow title="Detail view stays simple" copy="Availability, SLA, and health history remain table-first." />
          </div>
        </aside>
      </div>
    </>
  );
}
