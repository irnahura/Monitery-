import type { FormEvent } from "react";
import { Field, MessageLine } from "../components/common";
import type { AuthMode, Message } from "../types";

export function AuthScreen(props: {
  mode: AuthMode;
  form: { name: string; email: string; password: string; confirm: string };
  message: Message;
  setMode: (mode: AuthMode) => void;
  setForm: (form: { name: string; email: string; password: string; confirm: string }) => void;
  submit: (event: FormEvent) => void;
  startDemo: () => void;
}) {
  const isLogin = props.mode === "login";
  return (
    <div className="auth-shell">
      <section className="auth-panel">
        <div>
          <span className="brand-mark">Pulse</span>
          <span className="eyebrow">{isLogin ? "Welcome back" : "Create account"}</span>
          <h1>{isLogin ? "Sign in to your monitor workspace." : "Start a small, intentional monitor workspace."}</h1>
          <p className="muted-copy">{isLogin ? "Minimal auth chrome, direct validation, quiet visual weight." : "Only the fields required by the expansion brief."}</p>
        </div>
        <form className="form-stack" onSubmit={props.submit} noValidate>
          {!isLogin && <Field label="Full name"><input value={props.form.name} placeholder="Mayank Sharma" onChange={(event) => props.setForm({ ...props.form, name: event.target.value })} /></Field>}
          <Field label="Email"><input type="email" value={props.form.email} placeholder="you@company.com" onChange={(event) => props.setForm({ ...props.form, email: event.target.value })} /></Field>
          <Field label="Password"><input type="password" value={props.form.password} placeholder={isLogin ? "Enter password" : "At least 8 characters"} onChange={(event) => props.setForm({ ...props.form, password: event.target.value })} /></Field>
          {!isLogin && <Field label="Confirm password"><input type="password" value={props.form.confirm} placeholder="Repeat password" onChange={(event) => props.setForm({ ...props.form, confirm: event.target.value })} /></Field>}
          {isLogin && <label className="check-row"><input type="checkbox" /><span>Remember me</span></label>}
          <MessageLine message={props.message} />
          <button className="button button-primary full-width" type="submit">{isLogin ? "Login" : "Register"}</button>
          <button className="button button-secondary full-width" type="button" onClick={props.startDemo}>Preview demo</button>
        </form>
        <div className="auth-links">
          <button className="link-button" type="button" onClick={() => props.setMode(isLogin ? "register" : "login")}>{isLogin ? "Create account" : "Back to login"}</button>
          {isLogin && <span className="muted">Forgot password</span>}
        </div>
      </section>
    </div>
  );
}
