// A tiny reactive toast queue. Components push toasts to confirm what
// just happened; transient ones fade on their own, sticky ones stay until
// dismissed (used for failures and "needs a restart" outcomes). Every
// apply spawns a fresh toast, so there's never the old-vs-new ambiguity
// of a single reused banner.

export type ToastKind = "success" | "warn" | "info";

export interface Toast {
  id: number;
  kind: ToastKind;
  message: string;
  // Optional secondary lines (e.g. per-fix failures, a restart hint).
  detail?: string[];
  // Sticky toasts never auto-dismiss; the user closes them.
  sticky?: boolean;
}

let seq = 0;

// Mutated in place (push/splice) so the exported reference stays stable
// across the module boundary, which keeps Svelte's reactivity happy.
export const toasts = $state<Toast[]>([]);

const TRANSIENT_MS = 5000;

export function pushToast(t: Omit<Toast, "id">): number {
  const id = ++seq;
  toasts.push({ id, ...t });
  if (!t.sticky) {
    setTimeout(() => dismissToast(id), TRANSIENT_MS);
  }
  return id;
}

export function dismissToast(id: number): void {
  const i = toasts.findIndex((t) => t.id === id);
  if (i >= 0) toasts.splice(i, 1);
}

// Clear everything, e.g. when the user resets or starts a fresh apply.
export function clearToasts(): void {
  toasts.splice(0, toasts.length);
}
