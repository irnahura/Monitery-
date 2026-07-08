import type { APIKey, HealthCheck, Monitor, User } from "../types";

export const demoUser: User = { id: 1, name: "Mayank Sharma", email: "mayank@pulse.dev" };
export const demoMonitors: Monitor[] = [
  { id: 1, name: "Example", url: "https://example.com", latest_status_code: 200, latest_response_time_ms: 84, latest_is_up: true, last_checked_at: new Date().toISOString() },
  { id: 2, name: "Docs API", url: "https://api.pulse.local/health", latest_status_code: 204, latest_response_time_ms: 61, latest_is_up: true, last_checked_at: new Date(Date.now() - 24000).toISOString() },
  { id: 3, name: "Billing webhook", url: "https://billing.pulse.local/webhook", latest_status_code: 0, latest_response_time_ms: null, latest_is_up: false, last_checked_at: new Date(Date.now() - 51000).toISOString() }
];
export const demoHistory: HealthCheck[] = [
  { id: 1, status_code: 200, response_time_ms: 84, is_up: true, checked_at: new Date().toISOString() },
  { id: 2, status_code: 204, response_time_ms: 61, is_up: true, checked_at: new Date(Date.now() - 60000).toISOString() },
  { id: 3, status_code: 0, response_time_ms: 0, is_up: false, checked_at: new Date(Date.now() - 120000).toISOString() },
  { id: 4, status_code: 200, response_time_ms: 78, is_up: true, checked_at: new Date(Date.now() - 180000).toISOString() }
];
export const demoAPIKeys: APIKey[] = [
  { id: 1, name: "CLI integration", created_at: "2026-07-01T00:00:00Z", expires_at: "2026-10-01T00:00:00Z" }
];
