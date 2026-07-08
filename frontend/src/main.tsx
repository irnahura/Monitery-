import { useEffect, useState } from "react";
import { createRoot } from "react-dom/client";
import "./styles.css";
import { demoAPIKeys, demoHistory, demoMonitors, demoUser } from "./data/demo";
import { Sidebar } from "./components/Sidebar";
import { AddMonitorPage } from "./pages/AddMonitorPage";
import { AuthScreen } from "./pages/AuthScreen";
import { DashboardPage } from "./pages/DashboardPage";
import { EmptyStatePage } from "./pages/EmptyStatePage";
import { ServiceDetailPage } from "./pages/ServiceDetailPage";
import { SettingsPage } from "./pages/SettingsPage";
import type { APIKey, AuthMode, HealthCheck, Message, Monitor, Route, Summary, User } from "./types";

const API_BASE = import.meta.env.VITE_API_BASE_URL || "";
const DEMO_TOKEN = "demo-preview-token";

function App() {
  const [token, setToken] = useState(() => localStorage.getItem("token") || "");
  const [route, setRoute] = useState<Route>("dashboard");
  const [authMode, setAuthMode] = useState<AuthMode>("login");
  const [authForm, setAuthForm] = useState({ name: "", email: "", password: "", confirm: "" });
  const [monitorForm, setMonitorForm] = useState({ name: "", url: "", interval_seconds: 60, timeout_seconds: 10 });
  const [settingsTab, setSettingsTab] = useState("profile");
  const [message, setMessage] = useState<Message>({ text: "" });
  const [apiMessage, setApiMessage] = useState<Message>({ text: "" });
  const [newKey, setNewKey] = useState("");
  const [user, setUser] = useState<User | null>(null);
  const [monitors, setMonitors] = useState<Monitor[]>([]);
  const [selected, setSelected] = useState<Monitor | null>(null);
  const [history, setHistory] = useState<HealthCheck[]>([]);
  const [summary, setSummary] = useState<Summary>({ availability_percent: 100, sla_percent: 100 });
  const [apiKeys, setApiKeys] = useState<APIKey[]>([]);

  const headers = { Authorization: `Bearer ${token}`, "Content-Type": "application/json" };
  const isDemo = token === DEMO_TOKEN;

  useEffect(() => {
    if (!token) return;
    void loadData();
    const timer = window.setInterval(() => void loadData(false), 5000);
    return () => window.clearInterval(timer);
  }, [token]);

  async function request(path: string, options: RequestInit = {}) {
    const response = await fetch(`${API_BASE}${path}`, options);
    if (!response.ok) {
      let error = "Request failed";
      try { error = (await response.json()).error || error; } catch { error = response.statusText || error; }
      throw new Error(error);
    }
    if (response.status === 204) return null;
    return response.json();
  }

  async function loadData(selectFirst = true) {
    if (isDemo) {
      setUser(demoUser);
      setMonitors((current) => current.length ? current : demoMonitors);
      setApiKeys((current) => current.length ? current : demoAPIKeys);
      const nextSelected = selected || demoMonitors[0];
      if (selectFirst && nextSelected) {
        setSelected(nextSelected);
        setHistory(demoHistory);
        setSummary({ availability_percent: 99.92, sla_percent: 99.9 });
      }
      return;
    }
    try {
      const [profile, monitorRows, keyRows] = await Promise.all([
        request("/auth/profile", { headers }),
        request("/monitors", { headers }),
        request("/apikeys", { headers })
      ]);
      setUser(profile);
      setMonitors(monitorRows);
      setApiKeys(keyRows);
      const nextSelected = selected ? monitorRows.find((monitor: Monitor) => monitor.id === selected.id) : monitorRows[0];
      if (nextSelected && selectFirst) {
        setSelected(nextSelected);
        await loadHistory(nextSelected);
      }
      if (!monitorRows.length) {
        setSelected(null);
        setHistory([]);
      }
    } catch (error) {
      if (error instanceof Error && /authentication|unauthorized/i.test(error.message)) logout();
    }
  }

  async function loadHistory(monitor: Monitor) {
    if (isDemo) {
      setSelected(monitor);
      setHistory(demoHistory);
      setSummary({ availability_percent: 99.92, sla_percent: 99.9 });
      return;
    }
    const result = await request(`/monitors/${monitor.id}/history?limit=20`, { headers });
    setHistory(result.history || []);
    setSummary(result.summary || { availability_percent: 100, sla_percent: 100 });
  }

  async function submitAuth(event: React.FormEvent) {
    event.preventDefault();
    setMessage({ text: "" });
    if (authMode === "register" && authForm.password !== authForm.confirm) {
      setMessage({ text: "Passwords do not match.", type: "error" });
      return;
    }
    try {
      const payload = authMode === "login" ? { email: authForm.email, password: authForm.password } : { name: authForm.name, email: authForm.email, password: authForm.password };
      const result = await request(`/auth/${authMode}`, { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify(payload) });
      localStorage.setItem("token", result.token);
      setToken(result.token);
      setRoute("dashboard");
    } catch (error) {
      setMessage({ text: error instanceof Error ? error.message : "Authentication failed.", type: "error" });
    }
  }

  async function createMonitor(event: React.FormEvent) {
    event.preventDefault();
    setMessage({ text: "" });
    if (isDemo) {
      const monitor: Monitor = {
        id: Date.now(),
        name: monitorForm.name,
        url: monitorForm.url,
        latest_status_code: null,
        latest_response_time_ms: null,
        latest_is_up: null,
        last_checked_at: null
      };
      setMonitors((current) => [monitor, ...current]);
      setSelected(monitor);
      setMonitorForm({ name: "", url: "", interval_seconds: 60, timeout_seconds: 10 });
      setRoute("dashboard");
      return;
    }
    try {
      await request("/monitors", { method: "POST", headers, body: JSON.stringify({ ...monitorForm, method: "GET", follow_redirects: true, validate_ssl: true, retry_count: 1 }) });
      setMonitorForm({ name: "", url: "", interval_seconds: 60, timeout_seconds: 10 });
      setMessage({ text: "Monitor created. Returning to dashboard...", type: "success" });
      await loadData();
      setRoute("dashboard");
    } catch (error) {
      setMessage({ text: error instanceof Error ? error.message : "Could not create monitor.", type: "error" });
    }
  }

  async function deleteMonitor(id: number) {
    if (isDemo) {
      setMonitors((current) => current.filter((monitor) => monitor.id !== id));
      setSelected(null);
      setHistory([]);
      return;
    }
    await request(`/monitors/${id}`, { method: "DELETE", headers });
    setSelected(null);
    setHistory([]);
    await loadData();
  }

  async function createAPIKey() {
    if (isDemo) {
      const key = `pk_live_demo_${Math.random().toString(16).slice(2, 7)}`;
      setNewKey(key);
      setApiKeys((current) => [{ id: Date.now(), name: "Dashboard key", created_at: new Date().toISOString(), expires_at: null }, ...current]);
      setApiMessage({ text: "New API key generated.", type: "success" });
      return;
    }
    const result = await request("/apikeys", { method: "POST", headers, body: JSON.stringify({ name: "Dashboard key" }) });
    setNewKey(result.key);
    setApiMessage({ text: "New API key generated.", type: "success" });
    await loadData(false);
  }

  async function deleteAPIKey(id: number) {
    if (isDemo) {
      setApiKeys((current) => current.filter((key) => key.id !== id));
      setApiMessage({ text: "Key revoked.", type: "success" });
      return;
    }
    await request(`/apikeys/${id}`, { method: "DELETE", headers });
    setApiMessage({ text: "Key revoked.", type: "success" });
    await loadData(false);
  }

  function logout() {
    localStorage.removeItem("token");
    setToken("");
    setUser(null);
    setRoute("dashboard");
  }

  function startDemo() {
    localStorage.setItem("token", DEMO_TOKEN);
    setToken(DEMO_TOKEN);
    setUser(demoUser);
    setMonitors(demoMonitors);
    setSelected(demoMonitors[0]);
    setHistory(demoHistory);
    setSummary({ availability_percent: 99.92, sla_percent: 99.9 });
    setApiKeys(demoAPIKeys);
    setRoute("dashboard");
  }

  function openDetail(monitor: Monitor) {
    setSelected(monitor);
    void loadHistory(monitor);
    setRoute("detail");
  }

  if (!token) return <AuthScreen mode={authMode} form={authForm} message={message} setMode={setAuthMode} setForm={setAuthForm} submit={submitAuth} startDemo={startDemo} />;

  const activeRoute = route === "dashboard" && monitors.length === 0 ? "empty" : route;
  return (
    <div className="app-shell">
      <Sidebar route={activeRoute} setRoute={setRoute} logout={logout} />
      <main className="page">
        {activeRoute === "dashboard" && <DashboardPage monitors={monitors} history={history} summary={summary} selected={selected} openDetail={openDetail} setRoute={setRoute} />}
        {activeRoute === "empty" && <EmptyStatePage setRoute={setRoute} />}
        {activeRoute === "add" && <AddMonitorPage form={monitorForm} setForm={setMonitorForm} submit={createMonitor} message={message} setRoute={setRoute} />}
        {activeRoute === "detail" && <ServiceDetailPage monitor={selected || monitors[0]} history={history} summary={summary} setRoute={setRoute} deleteMonitor={deleteMonitor} />}
        {activeRoute === "settings" && <SettingsPage user={user} tab={settingsTab} setTab={setSettingsTab} apiKeys={apiKeys} newKey={newKey} apiMessage={apiMessage} createAPIKey={createAPIKey} deleteAPIKey={deleteAPIKey} />}
      </main>
    </div>
  );
}

createRoot(document.getElementById("root")!).render(<App />);
