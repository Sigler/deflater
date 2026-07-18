<script lang="ts">
  import { SvelteSet } from "svelte/reactivity";
  import mascot from "./assets/mascot-512.png";
  import { api } from "./lib/api";
  import { computeChanges, initialSelection } from "./lib/changes";
  import { S } from "./lib/i18n";
  import AlertsBanner from "./lib/components/AlertsBanner.svelte";
  import ApplyBar from "./lib/components/ApplyBar.svelte";
  import CategoryNav from "./lib/components/CategoryNav.svelte";
  import CategorySection from "./lib/components/CategorySection.svelte";
  import Header from "./lib/components/Header.svelte";
  import MaintenanceCard from "./lib/components/MaintenanceCard.svelte";
  import Modal from "./lib/components/Modal.svelte";
  import ScanBar from "./lib/components/ScanBar.svelte";
  import type { FixResult, Report } from "./lib/types";

  let report = $state<Report | null>(null);
  let loadError = $state("");
  let selection = $state(new SvelteSet<string>());
  let applying = $state(false);
  let progressText = $state("");
  let showElevateModal = $state(false);
  let doneMessage = $state("");
  let doneWarn = $state(false);
  let failures = $state<FixResult[]>([]);
  let maintenancePendingElevation = $state(false);
  let watcherPendingElevation = $state(false);

  const changes = $derived(
    report ? computeChanges(report.fixes, selection) : { enable: [], disable: [] },
  );
  const changeCount = $derived(changes.enable.length + changes.disable.length);
  const pendingIds = $derived(new Set([...changes.enable, ...changes.disable]));
  const inPlaceCount = $derived(
    report?.fixes.filter((f) => f.status === "on" || f.status === "removed").length ?? 0,
  );

  // Keep the backend informed so closing the window can warn about
  // staged changes that were never applied.
  $effect(() => {
    if (report) void api.setDirty(applying ? 0 : changeCount);
  });

  // Section navigation: sidenav when wide, tabs when narrow, scrollspy
  // highlighting whichever section is under the sticky bars.
  let scroller = $state<HTMLDivElement | null>(null);
  let activeSection = $state("");

  const navItems = $derived(
    report
      ? [
          ...report.categories.map((c) => ({
            id: c,
            label: S.categories[c as keyof typeof S.categories]?.nav ?? c,
          })),
          { id: "maintenance", label: S.nav.maintenance },
        ]
      : [],
  );

  function updateActive() {
    if (!scroller || navItems.length === 0) return;
    const y = scroller.scrollTop + 150;
    let current = navItems[0].id;
    for (const item of navItems) {
      const el = document.getElementById(`sec-${item.id}`);
      if (el && el.offsetTop <= y) current = item.id;
    }
    // At the end of a real scroll the last section may never reach the
    // spy line; scrolled to the bottom means the last section is active.
    // Only when there is actually something to scroll.
    const scrollable = scroller.scrollHeight > scroller.clientHeight + 8;
    if (scrollable && scroller.scrollTop + scroller.clientHeight >= scroller.scrollHeight - 8) {
      current = navItems[navItems.length - 1].id;
    }
    activeSection = current;
  }

  function jump(id: string) {
    document.getElementById(`sec-${id}`)?.scrollIntoView({ behavior: "smooth", block: "start" });
  }

  $effect(() => {
    if (report) updateActive();
  });

  function byCategory(cat: string) {
    return report?.fixes.filter((f) => f.category === cat) ?? [];
  }

  async function load() {
    loadError = "";
    try {
      const r = await api.getReport();
      report = r;
      selection = new SvelteSet(initialSelection(r.fixes, r.managed));
      // If we were relaunched elevated to finish a staged apply, claim it
      // (consume-on-read, so it can never fire twice) and run it.
      const pending = await api.takePending();
      if (pending) {
        progressText = S.apply.resuming;
        await runApply(pending.enable, pending.disable);
      }
    } catch (e) {
      loadError = `${e}`;
    }
  }

  // Progress events carry a phase: "start" fires before a fix's work, so
  // the label names the fix actually being worked on (not the last one).
  api.onApplyProgress((raw) => {
    const res = raw as FixResult;
    if (res.phase !== "start") return;
    const title = S.fixes[res.id as keyof typeof S.fixes]?.title ?? res.id;
    progressText = `${S.apply.applying} ${title}`;
  });

  async function runApply(enable: string[], disable: string[]) {
    applying = true;
    doneMessage = "";
    doneWarn = false;
    failures = [];
    try {
      const outcome = await api.apply(enable, disable);
      if (outcome.needsElevation) {
        showElevateModal = true;
        return;
      }
      const failed = (outcome.results ?? []).filter((r) => !r.ok);
      failures = failed;
      if (outcome.saveWarning) {
        doneWarn = true;
        doneMessage = S.apply.saveWarning;
      } else {
        doneWarn = failed.length > 0;
        doneMessage = failed.length > 0 ? S.apply.doneSomeFailed(failed.length) : S.apply.doneBody;
      }
      const r = await api.getReport();
      report = r;
      // Reflect reality, but keep failed enables selected so the user can
      // retry them without hunting each one down again.
      const next = initialSelection(r.fixes, r.managed);
      for (const f of failed) if (enable.includes(f.id)) next.add(f.id);
      selection = new SvelteSet(next);
      maintenancePendingElevation = false;
      watcherPendingElevation = false;
    } catch (e) {
      doneWarn = true;
      doneMessage = `${S.apply.applyError} ${e}`;
    } finally {
      applying = false;
      progressText = "";
    }
  }

  function apply() {
    void runApply(changes.enable, changes.disable);
  }

  async function confirmElevate() {
    showElevateModal = false;
    try {
      // Snapshot the change set the modal described, so a background
      // toggle can't alter what gets elevated.
      await api.saveAndElevate([...elevateSnapshot.enable], [...elevateSnapshot.disable]);
    } catch {
      // UAC declined; nothing was changed and nothing stays queued.
    }
  }

  // The change set captured when the elevate modal opened.
  let elevateSnapshot = $state<{ enable: string[]; disable: string[] }>({ enable: [], disable: [] });
  $effect(() => {
    if (showElevateModal) elevateSnapshot = { enable: changes.enable, disable: changes.disable };
  });

  function reset() {
    if (report) selection = new SvelteSet(initialSelection(report.fixes, report.managed));
    doneMessage = "";
    failures = [];
  }

  function toggleFix(id: string, value: boolean) {
    if (value) selection.add(id);
    else selection.delete(id);
    doneMessage = "";
  }

  async function setMaintenance(on: boolean) {
    if (!report) return;
    report.maintenance = on;
    try {
      const res = await api.setMaintenance(on);
      maintenancePendingElevation = res.needsElevation;
    } catch (e) {
      report.maintenance = !on; // roll the toggle back to reality
      maintenancePendingElevation = false;
      doneWarn = true;
      doneMessage = `${e}`;
    }
  }

  async function setWatcher(on: boolean) {
    if (!report) return;
    report.watcher = on;
    try {
      const res = await api.setWatcher(on);
      watcherPendingElevation = res.needsElevation;
    } catch (e) {
      report.watcher = !on;
      watcherPendingElevation = false;
      doneWarn = true;
      doneMessage = `${e}`;
    }
  }

  async function removeAlertPackage(pkg: string) {
    await api.removePackage(pkg);
    if (report) report.alerts = report.alerts.filter((a) => a.package !== pkg);
  }

  async function dismissAlerts() {
    await api.dismissAlerts();
    if (report) report.alerts = [];
  }

  void load();
</script>

{#if report === null}
  <div class="loading">
    <img class="loading-mascot" src={mascot} alt="" draggable="false" />
    {#if loadError}
      <p>{S.app.loadFailed}</p>
      <p class="hint">{loadError}</p>
      <button type="button" class="primary" onclick={() => load()}>{S.app.retry}</button>
    {:else}
      <div class="spinner" aria-hidden="true"></div>
      <p>{S.app.loading}</p>
      <p class="hint">{S.app.loadingHint}</p>
    {/if}
  </div>
{:else}
  <div class="page" bind:this={scroller} onscroll={updateActive}>
    <div class="shell">
      <aside class="sidenav">
        <CategoryNav items={navItems} active={activeSection} variant="side" onjump={jump} />
      </aside>
      <main>
        <Header />
        <AlertsBanner
          alerts={report.alerts}
          onremove={removeAlertPackage}
          ondismiss={dismissAlerts}
        />

        {#if doneMessage}
          <div class="done" class:warn={doneWarn}>
            <p>{doneMessage}</p>
            {#each failures as f (f.id)}
              <p class="fail">{S.fixes[f.id as keyof typeof S.fixes]?.title ?? f.id}: {f.error}</p>
            {/each}
            {#if !doneWarn && failures.length === 0 && !report.maintenance}
              <p class="tip">{S.apply.doneMaintenanceTip}</p>
            {/if}
          </div>
        {/if}

        <div class="stickytop">
          <ScanBar inPlace={inPlaceCount} total={report.fixes.length} />
          <div class="tabsrow">
            <CategoryNav items={navItems} active={activeSection} variant="tabs" onjump={jump} />
          </div>
        </div>

        {#each report.categories as cat (cat)}
          <CategorySection
            id={cat}
            fixes={byCategory(cat)}
            {selection}
            pending={pendingIds}
            {applying}
            ontoggle={toggleFix}
          />
        {/each}

        <div class="anchor" id="sec-maintenance">
          <div class="secheader">
            <h2>{S.nav.maintenance}</h2>
            <p>{S.maintenance.sectionBlurb}</p>
          </div>
          {#if report.taskMismatch}
            <div class="done warn"><p>{S.maintenance.mismatch}</p></div>
          {/if}
          <MaintenanceCard
            maintenance={report.maintenance}
            watcher={report.watcher}
            maintenancePending={maintenancePendingElevation}
            watcherPending={watcherPendingElevation}
            onmaintenance={setMaintenance}
            onwatcher={setWatcher}
          />
        </div>

        <footer>
          <button type="button" class="link" onclick={() => api.openLogFolder()}>
            {S.footer.logs}
          </button>
          <span>{S.footer.assurance}</span>
          <span class="stamp">{S.footer.version(report.version)}</span>
        </footer>
      </main>
    </div>

    <ApplyBar {changeCount} {applying} {progressText} onapply={apply} onreset={reset} />
  </div>
{/if}

{#if showElevateModal}
  <Modal title={S.apply.elevateTitle} oncancel={() => (showElevateModal = false)}>
    {#snippet children()}
      <p>{S.apply.elevateBody}</p>
    {/snippet}
    {#snippet actions()}
      <button type="button" class="ghost" onclick={() => (showElevateModal = false)}>
        {S.apply.elevateCancel}
      </button>
      <button type="button" class="primary" onclick={confirmElevate}>
        {S.apply.elevateConfirm}
      </button>
    {/snippet}
  </Modal>
{/if}

<style>
  .page {
    height: 100vh;
    display: flex;
    flex-direction: column;
    overflow-y: auto;
  }
  .shell {
    flex: 1;
    display: grid;
    grid-template-columns: minmax(0, 1fr);
    justify-content: center;
  }
  .sidenav {
    display: none;
  }
  main {
    width: min(860px, 100%);
    margin: 0 auto;
    padding: 24px 24px 40px;
    display: grid;
    gap: 22px;
    align-content: start;
  }
  @media (min-width: 1220px) {
    .shell {
      grid-template-columns: 190px minmax(0, 860px);
      gap: 28px;
    }
    .sidenav {
      display: block;
      padding-top: 128px;
    }
    .sidenav :global(nav) {
      position: sticky;
      top: 24px;
    }
    main {
      width: 100%;
      margin: 0;
    }
    .tabsrow {
      display: none;
    }
  }
  .stickytop {
    position: sticky;
    top: 0;
    z-index: 10;
    background: color-mix(in srgb, var(--bg-window) 90%, transparent);
    backdrop-filter: blur(12px);
    border-bottom: 1px solid var(--stroke);
  }
  .anchor {
    scroll-margin-top: 96px;
    display: grid;
    gap: 10px;
    /* Same crate separation as the category sections. */
    margin-top: 18px;
  }
  /* Matches the category section headers in CategorySection.svelte. */
  .secheader {
    display: flex;
    align-items: baseline;
    gap: 12px;
    padding: 0 2px;
  }
  .secheader h2 {
    font-size: 13px;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.09em;
    color: var(--text);
  }
  .secheader p {
    color: var(--text-faint);
    font-size: 12.5px;
  }
  .loading {
    height: 100vh;
    display: grid;
    place-content: center;
    justify-items: center;
    gap: 12px;
    color: var(--text-dim);
  }
  .loading-mascot {
    width: 320px;
    height: 320px;
    margin-bottom: 4px;
    user-select: none;
  }
  .loading .hint {
    font-size: 12px;
    color: var(--text-faint);
  }
  .spinner {
    width: 28px;
    height: 28px;
    border-radius: 50%;
    border: 3px solid var(--stroke-strong);
    border-top-color: var(--coral);
    animation: spin 0.9s linear infinite;
  }
  @keyframes spin {
    to {
      transform: rotate(360deg);
    }
  }
  .done {
    background: var(--sage-soft);
    border: 1px solid rgba(132, 191, 164, 0.35);
    border-radius: var(--r-card);
    padding: 12px 16px;
    font-size: 13px;
    color: var(--text);
    display: grid;
    gap: 6px;
  }
  .done.warn {
    background: var(--gold-soft);
    border-color: rgba(214, 164, 76, 0.35);
  }
  .fail {
    font-size: 12px;
    color: var(--text-dim);
  }
  .tip {
    font-size: 12px;
    color: var(--text-dim);
  }
  footer {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 12px;
    padding-top: 6px;
    font-size: 12px;
    color: var(--text-faint);
  }
  .link {
    color: var(--text-dim);
    text-decoration: underline;
    text-underline-offset: 3px;
  }
  .link:hover {
    color: var(--text);
  }
  .stamp {
    text-transform: uppercase;
    letter-spacing: 0.07em;
    font-size: 10.5px;
  }
  .primary {
    padding: 8px 18px;
    border-radius: var(--r-control);
    background: linear-gradient(180deg, var(--coral-bright), var(--coral));
    color: #241511;
    font-weight: 600;
    box-shadow:
      inset 0 1px 0 rgba(255, 255, 255, 0.25),
      0 1px 3px rgba(0, 0, 0, 0.35);
  }
  .primary:hover {
    background: linear-gradient(180deg, #ff8f70, var(--coral-bright));
  }
  .primary:active {
    transform: translateY(1px);
  }
  .ghost {
    padding: 8px 14px;
    border-radius: var(--r-control);
    border: 1px solid var(--stroke-strong);
    color: var(--text-dim);
  }
  .ghost:hover {
    color: var(--text);
  }
</style>
