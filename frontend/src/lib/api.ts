// Thin typed wrapper around the generated Wails bindings, so components
// import one module and tests can mock it.

import {
  Apply,
  CheckUpdate,
  DismissAlerts,
  GetReport,
  OpenLogFolder,
  OpenRecallFolder,
  RecallSnapshots,
  RemoveConflictingTasks,
  RemovePackage,
  RestartExplorer,
  SaveAndElevate,
  SetDirty,
  SetMaintenance,
  SetWatcher,
  StageTaskRemovalAndElevate,
  TakePending,
} from "../../wailsjs/go/main/App";
import { BrowserOpenURL, EventsOn } from "../../wailsjs/runtime/runtime";
import type {
  ApplyOutcome,
  Pending,
  RecallInfo,
  Report,
  ToggleResult,
  UpdateInfo,
} from "./types";

export const api = {
  getReport: () => GetReport() as Promise<Report>,
  apply: (enable: string[], disable: string[]) => Apply(enable, disable) as Promise<ApplyOutcome>,
  saveAndElevate: (enable: string[], disable: string[]) => SaveAndElevate(enable, disable),
  takePending: () => TakePending() as Promise<Pending | null>,
  setMaintenance: (on: boolean) => SetMaintenance(on) as Promise<ToggleResult>,
  setWatcher: (on: boolean) => SetWatcher(on) as Promise<ToggleResult>,
  setDirty: (n: number) => SetDirty(n),
  dismissAlerts: () => DismissAlerts(),
  removePackage: (name: string) => RemovePackage(name),
  removeConflictingTasks: (names: string[]) => RemoveConflictingTasks(names),
  stageTaskRemovalAndElevate: (name: string) => StageTaskRemovalAndElevate(name),
  restartExplorer: () => RestartExplorer(),
  checkUpdate: () => CheckUpdate() as Promise<UpdateInfo>,
  openUrl: (url: string) => BrowserOpenURL(url),
  recallSnapshots: () => RecallSnapshots() as Promise<RecallInfo>,
  openRecallFolder: () => OpenRecallFolder(),
  openLogFolder: () => OpenLogFolder(),
  onApplyProgress: (cb: (result: unknown) => void) => EventsOn("apply:progress", cb),
  openStorePage: (productId: string) =>
    BrowserOpenURL(`ms-windows-store://pdp/?ProductId=${productId}`),
  openStoreSearch: (query: string) =>
    BrowserOpenURL(`ms-windows-store://search/?query=${encodeURIComponent(query)}`),
};
