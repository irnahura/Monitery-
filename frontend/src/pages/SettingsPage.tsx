import { Field, MessageLine } from "../components/common";
import type { APIKey, Message, User } from "../types";
import { dateOnly } from "../utils/format";

export function SettingsPage(props: {
  user: User | null;
  tab: string;
  setTab: (tab: string) => void;
  apiKeys: APIKey[];
  newKey: string;
  apiMessage: Message;
  createAPIKey: () => Promise<void>;
  deleteAPIKey: (id: number) => Promise<void>;
}) {
  const tabs = [["profile", "Profile"], ["password", "Change Password"], ["api", "API Keys"], ["notifications", "Email Notifications"]];
  return (
    <>
      <header className="page-header">
        <div>
          <span className="eyebrow">Account settings</span>
          <h1>Settings</h1>
          <p className="muted-copy">Profile, password, API keys, and notifications in one place.</p>
        </div>
      </header>
      <div className="settings-layout">
        <nav className="settings-tabs" aria-label="Settings sections">{tabs.map(([id, label]) => <button key={id} className={`settings-tab ${props.tab === id ? "active" : ""}`} type="button" onClick={() => props.setTab(id)}>{label}</button>)}</nav>
        <div className="settings-panels">
          {props.tab === "profile" && <section className="card settings-panel active"><div className="card-header"><div><h2>Profile</h2><p className="muted-copy">Update core account details.</p></div></div><div className="form-stack inline-form"><Field label="Full name"><input readOnly value={props.user?.name || ""} /></Field><Field label="Email"><input readOnly value={props.user?.email || ""} /></Field><button className="button button-primary" type="button">Save changes</button></div></section>}
          {props.tab === "password" && <section className="card settings-panel active"><div className="card-header"><div><h2>Change password</h2><p className="muted-copy">Keep the action explicit and lightweight.</p></div></div><div className="form-stack inline-form"><Field label="Current password"><input type="password" /></Field><Field label="New password"><input type="password" /></Field><Field label="Confirm password"><input type="password" /></Field><button className="button button-primary" type="button">Update password</button></div></section>}
          {props.tab === "api" && <section className="card settings-panel active"><div className="card-header"><div><h2>API keys</h2><p className="muted-copy">Management stays inside settings, per the expansion plan.</p></div><button className="button button-secondary" type="button" onClick={() => void props.createAPIKey()}>Generate new key</button></div><MessageLine message={props.apiMessage} />{props.newKey && <input className="key-output mono" readOnly value={props.newKey} onFocus={(event) => event.currentTarget.select()} />}<div className="history-table api-table" role="table" aria-label="API keys"><div className="table-head" role="row"><span>Key name</span><span>Created</span><span>Expires</span><span>Actions</span></div>{props.apiKeys.map((key) => <div className="table-row static-row api-row" role="row" key={key.id}><div><strong>{key.name}</strong><span className="mono">Stored securely</span></div><span className="mono">{dateOnly(key.created_at)}</span><span className="mono">{key.expires_at ? dateOnly(key.expires_at) : "Never"}</span><div className="table-actions"><button className="button button-tertiary" type="button" onClick={() => void props.deleteAPIKey(key.id)}>Revoke key</button></div></div>)}</div></section>}
          {props.tab === "notifications" && <section className="card settings-panel active"><div className="card-header"><div><h2>Email notifications</h2><p className="muted-copy">Basic operational preferences only.</p></div></div><div className="form-stack inline-form"><label className="check-row"><input type="checkbox" defaultChecked /><span>Enable email notifications</span></label><label className="check-row"><input type="checkbox" defaultChecked /><span>Notify when service goes down</span></label><label className="check-row"><input type="checkbox" defaultChecked /><span>Notify when service recovers</span></label><div className="button-row"><button className="button button-secondary" type="button">Send test email</button></div></div></section>}
        </div>
      </div>
    </>
  );
}
