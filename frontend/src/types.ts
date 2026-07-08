export type Route = "dashboard" | "add" | "detail" | "settings" | "empty";
export type AuthMode = "login" | "register";
export type Message = { text: string; type?: "error" | "success" };
export type User = { id: number; name: string; email: string };
export type Monitor = { id: number; name: string; url: string; latest_status_code?: number | null; latest_response_time_ms?: number | null; latest_is_up?: boolean | null; last_checked_at?: string | null };
export type HealthCheck = { id: number; status_code: number; response_time_ms: number; is_up: boolean; checked_at: string };
export type Summary = { availability_percent: number; sla_percent: number };
export type APIKey = { id: number; name: string; created_at: string; expires_at?: string | null };
