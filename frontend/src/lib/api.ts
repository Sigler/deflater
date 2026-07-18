// Thin typed wrapper around the generated Wails bindings, so components
// import one module and tests can mock it.

import {
  Apply,
  DismissAlerts,
  GetReport,
  OpenLogFolder,
  RemovePackage,
  SaveAndElevate,
  SetDirty,
  SetMaintenance,
  SetWatcher,
} from "../../wailsjs/go/main/App";
import { BrowserOpenURL, EventsOn } from "../../wailsjs/runtime/runtime";
import type { ApplyOutcome, Report } from "./types";

export const api = {
  getReport: () => GetReport() as Promise<Report>,
  apply: (enable: string[], disable: string[]) => Apply(enable, disable) as Promise<ApplyOutcome>,
  saveAndElevate: (enable: string[], disable: string[]) => SaveAndElevate(enable, disable),
  setMaintenance: (on: boolean) => SetMaintenance(on) as Promise<boolean>,
  setWatcher: (on: boolean) => SetWatcher(on) as Promise<boolean>,
  setDirty: (n: number) => SetDirty(n),
  dismissAlerts: () => DismissAlerts(),
  removePackage: (name: string) => RemovePackage(name),
  openLogFolder: () => OpenLogFolder(),
  onApplyProgress: (cb: (result: unknown) => void) => EventsOn("apply:progress", cb),
  openStoreSearch: (query: string) =>
    BrowserOpenURL(`ms-windows-store://search/?query=${encodeURIComponent(query)}`),
};
