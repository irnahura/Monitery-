import React from "react";
import { NavButton } from "./common";
import type { Route } from "../types";

export function Sidebar({ route, setRoute, logout }: { route: Route; setRoute: (route: Route) => void; logout: () => void }) {
  return <aside className="sidebar"><button className="brand-mark link-button" type="button" onClick={() => setRoute("dashboard")}>Pulse</button><nav className="nav-list" aria-label="Primary"><NavButton active={route === "dashboard" || route === "empty"} onClick={() => setRoute("dashboard")}>Dashboard</NavButton><NavButton active={route === "add"} onClick={() => setRoute("add")}>Add Monitor</NavButton><NavButton active={route === "settings"} onClick={() => setRoute("settings")}>Settings</NavButton><NavButton active={false} onClick={logout}>Logout</NavButton></nav><div className="sidebar-note"><span className="eyebrow">Cadence</span><p>Frontend polls every 5s. Backend checks run roughly every 60s.</p></div></aside>;
}
