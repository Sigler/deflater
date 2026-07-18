// Shapes shared across the UI. These mirror the Go structs in app.go;
// the wailsjs bindings carry the actual data across the bridge.

export type FixKind = "switch" | "app-junk" | "app-might" | "onedrive";

export type FixStatus = "on" | "off" | "partial" | "removed" | "installed" | "unknown";

export type ProfileId = "light-touch" | "clean-sweep" | "full-deflate";

export interface RegOp {
  hive: string;
  path: string;
  name: string;
  value: number;
  revert: string;
}

export type Refresh = "none" | "explorer" | "signout" | "reboot";

export interface FixState {
  id: string;
  category: string;
  kind: FixKind;
  caution: boolean;
  profiles: string[];
  reg?: RegOp[];
  appx?: string[];
  status: FixStatus;
  refresh: Refresh;
}

export interface Alert {
  package: string;
  seen: string;
}

export interface Pending {
  enable: string[];
  disable: string[];
  removeTasks?: string[];
  token: string;
  created: string;
}

export interface ForeignTask {
  name: string;
  tool: string;
  note: string;
}

export interface UpdateInfo {
  available: boolean;
  current: string;
  latest: string;
  url: string;
}

export interface Report {
  version: string;
  elevated: boolean;
  categories: string[];
  fixes: FixState[];
  managed: string[];
  maintenance: boolean;
  watcher: boolean;
  alerts: Alert[];
  taskMismatch: boolean;
  conflictingTasks: ForeignTask[];
  pending: Pending | null;
}

export interface FixResult {
  id: string;
  ok: boolean;
  error?: string;
  status: FixStatus;
  phase: "start" | "done";
}

export interface ApplyOutcome {
  needsElevation: boolean;
  results: FixResult[] | null;
  saveWarning?: string;
  refresh?: Refresh;
}

export interface ToggleResult {
  saved: boolean;
  needsElevation: boolean;
}
