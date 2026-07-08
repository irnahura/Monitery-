import React from "react";
import type { Message, Monitor } from "../types";

export function Field({ label, children }: { label: string; children: React.ReactNode }) {
  return <label className="field"><span>{label}</span>{children}</label>;
}

export function MessageLine({ message }: { message: Message }) {
  return <div className={`message ${message.type || ""}`} aria-live="polite">{message.text}</div>;
}

export function NavButton({ active, onClick, children }: { active: boolean; onClick: () => void; children: React.ReactNode }) {
  return <button className={`nav-item ${active ? "active" : ""}`} type="button" onClick={onClick}>{children}</button>;
}

export function StatCard({ label, value, copy }: { label: string; value: string; copy: string }) {
  return <article className="card stat-card"><span className="card-label">{label}</span><strong>{value}</strong><p>{copy}</p></article>;
}

export function StatusPill({ monitor }: { monitor: Monitor }) {
  if (monitor.latest_is_up == null) return <span className="pill pill-neutral">WAITING</span>;
  return <span className={`pill ${monitor.latest_is_up ? "pill-up" : "pill-down"}`}>{monitor.latest_is_up ? `UP - ${monitor.latest_status_code || ""}` : "DOWN"}</span>;
}

export function InfoRow({ title, copy }: { title: string; copy: string }) {
  return <div className="list-row"><div><strong>{title}</strong><span>{copy}</span></div></div>;
}

export function Meta({ label, value }: { label: string; value: string }) {
  return <div><span className="card-label">{label}</span><strong className="mono-block">{value}</strong></div>;
}
