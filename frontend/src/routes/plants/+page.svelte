<script>
  import { onMount } from 'svelte';
  import {
    createPlant,
    updatePlant,
    deletePlant,
    listPlants,
    getCompanionData,
    listJournalEntries,
    scanSeedPacket
  } from '$lib/api';

  let query = '';
  let filterSeason = '';
  let filterStarred = false;
  let plants = [];
  let loading = true;
  let error = null;
  let formError = null;
  let isSaving = false;
  let editingPlantId = null;
  let confirmDeletePlantId = null;

  let name = '';
  let variety = '';
  let season = '';
  let spacing = '';
  let notes = '';
  let daysToMaturity = '';
  let sowDepth = '';
  let whenToSow = '';
  let daysToEmerge = '';
  let seedSpacing = '';
  let rowSpacing = '';
  let thinning = '';
  let showModal = false;

  let showScanModal = false;
  let scanFile = null;
  let scanLoading = false;
  let scanError = null;
  let scanResult = null;
  let scanCompany = 'botanical_interests';
  let scanPreviewUrl = '';
  let helperEnabled = false;
  let helperBox = null;
  let helperDragging = false;
  let helperStart = null;
  let helperCopied = false;
  let helperAttribute = 'name';
  let helperMappings = {};
  let helperMappingsCopied = false;
  let helperLockField = true;
  let helperContainer;
  const helperFieldLabels = {
    name: 'Name',
    variety: 'Variety',
    season: 'Season',
    spacing: 'Spacing',
    specs: 'Specs block',
    sow: 'When to sow block',
    emerge: 'Days to emerge',
    seed_spacing: 'Seed spacing',
    row_spacing: 'Row spacing',
    days_to_maturity: 'Days to maturity',
    depth: 'Sow depth',
    when_to_sow: 'When to sow',
    notes: 'Notes'
  };
  let selectedPlantId = '';
  let companionData = null;
  let companionLoading = false;
  let companionError = null;
  let journalByPlant = {};
  let journalLoading = false;

  const PAGE_SIZE = 25;
  let currentPage = 1;

  $: seasons = [...new Set(plants.map((p) => p.season).filter(Boolean))].sort();

  $: filteredPlants = plants.filter((p) => {
    const q = query.toLowerCase().trim();
    const matchesText =
      !q ||
      p.name.toLowerCase().includes(q) ||
      (p.variety || '').toLowerCase().includes(q);
    const matchesSeason = !filterSeason || p.season === filterSeason;
    const matchesStarred = !filterStarred || p.is_starred;
    return matchesText && matchesSeason && matchesStarred;
  });

  $: totalPages = Math.max(1, Math.ceil(filteredPlants.length / PAGE_SIZE));
  $: safePage = Math.min(Math.max(1, currentPage), totalPages);
  $: pagedPlants = filteredPlants.slice((safePage - 1) * PAGE_SIZE, safePage * PAGE_SIZE);

  async function loadPlants() {
    loading = true;
    error = null;
    try {
      plants = await listPlants({ limit: 500 });
    } catch (err) {
      error = err instanceof Error ? err.message : 'Unable to load plants.';
    } finally {
      loading = false;
    }
  }

  async function toggleStar(plant) {
    const newVal = !plant.is_starred;
    try {
      await updatePlant(plant.id, { is_starred: newVal });
      plants = plants.map((p) => (p.id === plant.id ? { ...p, is_starred: newVal } : p));
    } catch (err) {
      error = err instanceof Error ? err.message : 'Unable to update star.';
    }
  }

  async function submitPlant() {
    formError = null;
    if (!name.trim()) {
      formError = 'Plant name is required.';
      return;
    }
    if (daysToMaturity && Number.isNaN(Number(daysToMaturity))) {
      formError = 'Days to maturity must be a number.';
      return;
    }

    isSaving = true;
    try {
      const payload = {
        name: name.trim(),
        variety: variety.trim(),
        season: season.trim(),
        spacing: spacing.trim(),
        notes: notes.trim(),
        days_to_maturity: daysToMaturity ? Number(daysToMaturity) : null,
        sow_depth: sowDepth.trim(),
        when_to_sow: whenToSow.trim(),
        days_to_emerge: daysToEmerge.trim(),
        seed_spacing: seedSpacing.trim(),
        row_spacing: rowSpacing.trim(),
        thinning: thinning.trim()
      };
      if (editingPlantId) {
        const updated = await updatePlant(editingPlantId, payload);
        plants = plants.map((p) => (p.id === updated.id ? updated : p));
      } else {
        const created = await createPlant(payload);
        plants = [created, ...plants.filter((p) => p.id !== created.id)];
        await loadPlants();
      }
      resetPlantForm();
    } catch (err) {
      formError = err instanceof Error ? err.message : 'Unable to save plant.';
    } finally {
      isSaving = false;
    }
  }

  function resetPlantForm() {
    name = '';
    variety = '';
    season = '';
    spacing = '';
    notes = '';
    daysToMaturity = '';
    sowDepth = '';
    whenToSow = '';
    daysToEmerge = '';
    seedSpacing = '';
    rowSpacing = '';
    thinning = '';
    editingPlantId = null;
    formError = null;
  }

  function openEditModal(plant) {
    editingPlantId = plant.id;
    name = plant.name || '';
    variety = plant.variety || '';
    season = plant.season || '';
    spacing = plant.spacing || '';
    notes = plant.notes || '';
    daysToMaturity = plant.days_to_maturity != null ? String(plant.days_to_maturity) : '';
    sowDepth = plant.sow_depth || '';
    whenToSow = plant.when_to_sow || '';
    daysToEmerge = plant.days_to_emerge || '';
    seedSpacing = plant.seed_spacing || '';
    rowSpacing = plant.row_spacing || '';
    thinning = plant.thinning || '';
    formError = null;
    showModal = true;
  }

  async function handleDeletePlant(plantId) {
    try {
      await deletePlant(plantId);
      plants = plants.filter((p) => p.id !== plantId);
      if (selectedPlantId === plantId) selectedPlantId = '';
    } catch (err) {
      error = err instanceof Error ? err.message : 'Unable to delete plant.';
    }
  }

  async function handleScan() {
    scanError = null;
    if (!scanFile) {
      scanError = 'Select an image to scan.';
      return;
    }
    scanLoading = true;
    try {
      scanResult = await scanSeedPacket(scanFile, scanCompany);
      name = scanResult.name || '';
      variety = scanResult.variety || '';
      season = scanResult.season || '';
      spacing = scanResult.spacing || '';
      notes = scanResult.notes || '';
      daysToMaturity = scanResult.days_to_maturity ? String(scanResult.days_to_maturity) : '';
      sowDepth = scanResult.sow_depth || '';
      whenToSow = scanResult.when_to_sow || '';
      daysToEmerge = scanResult.days_to_emerge || '';
      seedSpacing = scanResult.seed_spacing || '';
      rowSpacing = scanResult.row_spacing || '';
      thinning = scanResult.thinning || '';
    } catch (err) {
      scanError = err instanceof Error ? err.message : 'Unable to scan seed packet.';
    } finally {
      scanLoading = false;
    }
  }

  function clamp01(value) {
    if (Number.isNaN(value)) return 0;
    return Math.max(0, Math.min(1, value));
  }

  function handleHelperPointerDown(event) {
    if (!helperEnabled || !helperContainer) return;
    event.preventDefault();
    const rect = helperContainer.getBoundingClientRect();
    const x = clamp01((event.clientX - rect.left) / rect.width);
    const y = clamp01((event.clientY - rect.top) / rect.height);
    helperStart = { x, y };
    helperBox = { x, y, w: 0, h: 0 };
    helperDragging = true;
  }

  function handleHelperPointerMove(event) {
    if (!helperEnabled || !helperDragging || !helperStart || !helperContainer) return;
    event.preventDefault();
    const rect = helperContainer.getBoundingClientRect();
    const x = clamp01((event.clientX - rect.left) / rect.width);
    const y = clamp01((event.clientY - rect.top) / rect.height);
    const left = Math.min(helperStart.x, x);
    const top = Math.min(helperStart.y, y);
    const right = Math.max(helperStart.x, x);
    const bottom = Math.max(helperStart.y, y);
    helperBox = { x: left, y: top, w: right - left, h: bottom - top };
  }

  function handleHelperPointerUp() {
    if (!helperEnabled) return;
    helperDragging = false;
    helperStart = null;
    if (helperBox && helperAttribute) {
      const finalized = { ...helperBox };
      helperMappings = { ...helperMappings, [helperAttribute]: finalized };
      if (!helperLockField) {
        helperAttribute = 'name';
      }
    }
    helperBox = null;
  }

  async function copyHelperBox() {
    if (!helperBox) return;
    const payload = {
      x: Number(helperBox.x.toFixed(3)),
      y: Number(helperBox.y.toFixed(3)),
      w: Number(helperBox.w.toFixed(3)),
      h: Number(helperBox.h.toFixed(3))
    };
    try {
      await navigator.clipboard.writeText(JSON.stringify(payload));
      helperCopied = true;
      setTimeout(() => {
        helperCopied = false;
      }, 1200);
    } catch (err) {
      helperCopied = false;
    }
  }

  async function copyHelperMappings() {
    if (!helperMappings || Object.keys(helperMappings).length === 0) return;
    try {
      await navigator.clipboard.writeText(
        JSON.stringify(
          Object.fromEntries(
            Object.entries(helperMappings).map(([key, value]) => [
              key,
              {
                x: Number(value.x.toFixed(3)),
                y: Number(value.y.toFixed(3)),
                w: Number(value.w.toFixed(3)),
                h: Number(value.h.toFixed(3))
              }
            ])
          )
        )
      );
      helperMappingsCopied = true;
      setTimeout(() => {
        helperMappingsCopied = false;
      }, 1200);
    } catch (err) {
      helperMappingsCopied = false;
    }
  }

  function clearHelperBox() {
    helperBox = null;
    helperDragging = false;
    helperStart = null;
  }

  function clearHelperMapping() {
    if (!helperAttribute) return;
    const next = { ...helperMappings };
    delete next[helperAttribute];
    helperMappings = next;
  }

  $: if (selectedPlantId) {
    const selected = plants.find((plant) => plant.id === selectedPlantId);
    if (selected) {
      loadCompanionData(selected.name);
    }
  }

  onMount(loadPlants);
  onMount(loadJournalEntries);

  async function loadCompanionData(name) {
    companionLoading = true;
    companionError = null;
    try {
      companionData = await getCompanionData(name);
    } catch (err) {
      companionData = null;
      companionError = err instanceof Error ? err.message : 'Unable to load companion data.';
    } finally {
      companionLoading = false;
    }
  }

  async function loadJournalEntries() {
    journalLoading = true;
    try {
      const entries = await listJournalEntries({ limit: 200 });
      const grouped = {};
      entries.forEach((entry) => {
        (entry.plants || []).forEach((plant) => {
          if (!grouped[plant.id]) {
            grouped[plant.id] = [];
          }
          grouped[plant.id].push(entry);
        });
      });
      journalByPlant = grouped;
    } catch (err) {
      journalByPlant = {};
    } finally {
      journalLoading = false;
    }
  }
</script>

<section class="mx-auto w-full max-w-6xl px-6 py-12">
  <div class="flex flex-wrap items-end justify-between gap-4">
    <div>
      <p class="text-xs uppercase tracking-[0.4em] text-earth-terracotta">Registry</p>
      <h1 class="mt-3 text-3xl font-semibold">Plant Registry</h1>
      <p class="mt-2 text-sm text-earth-text/70">
        Keep reference data for your core crops, spacing, and seasonal notes.
      </p>
    </div>
    <div class="flex flex-col gap-3">
      <div class="min-w-[280px]">
        <label class="text-xs uppercase tracking-[0.3em] text-earth-terracotta">Search</label>
        <input
          class="mt-2 w-full rounded-2xl border border-earth-border bg-white px-3 py-2 text-sm"
          placeholder="Filter by name or variety"
          bind:value={query}
          on:input={() => (currentPage = 1)}
        />
      </div>
      {#if seasons.length > 0}
        <div class="flex flex-wrap gap-2">
          <button
            class="rounded-full border px-3 py-1 text-xs font-semibold transition-colors {filterStarred ? 'border-amber-500 bg-amber-500 text-white' : 'border-earth-border text-earth-text/70'}"
            on:click={() => { filterStarred = !filterStarred; currentPage = 1; }}
          >
            ★ Starred
          </button>
          <span class="border-r border-earth-border"></span>
          <button
            class="rounded-full border px-3 py-1 text-xs font-semibold transition-colors {!filterSeason ? 'border-earth-forest bg-earth-forest text-earth-bg' : 'border-earth-border text-earth-text/70'}"
            on:click={() => { filterSeason = ''; currentPage = 1; }}
          >
            All
          </button>
          {#each seasons as s}
            <button
              class="rounded-full border px-3 py-1 text-xs font-semibold transition-colors {filterSeason === s ? 'border-earth-forest bg-earth-forest text-earth-bg' : 'border-earth-border text-earth-text/70'}"
              on:click={() => { filterSeason = filterSeason === s ? '' : s; currentPage = 1; }}
            >
              {s}
            </button>
          {/each}
        </div>
      {/if}
    </div>
  </div>

  <div class="mt-8 grid gap-6 lg:grid-cols-[1fr]">
    <div class="rounded-3xl border border-earth-border bg-white">
      <div class="border-b border-earth-border px-6 py-4">
        <div class="flex flex-wrap items-center justify-between gap-4">
          <div>
            <h2 class="text-lg font-semibold">Registry</h2>
            <p class="mt-1 text-sm text-earth-text/70">
              {#if loading}
                Loading plants…
              {:else if filteredPlants.length === plants.length}
                {plants.length} plants
              {:else}
                {filteredPlants.length} of {plants.length} plants
              {/if}
            </p>
          </div>
          <button
            class="rounded-full bg-earth-forest px-5 py-2.5 text-sm font-semibold text-earth-bg"
            on:click={() => {
              resetPlantForm();
              showModal = true;
            }}
          >
            Add plant
          </button>
          <button
            class="rounded-full border border-earth-border px-5 py-2.5 text-sm font-semibold"
            on:click={() => {
              scanError = null;
              scanFile = null;
              scanResult = null;
              scanPreviewUrl = '';
              showScanModal = true;
            }}
          >
            Scan seed packet
          </button>
        </div>
      </div>

      {#if error}
        <div class="mx-6 mt-4 rounded-2xl border border-rose-200 bg-rose-50 p-3 text-sm text-rose-700">
          {error}
        </div>
      {/if}

      <div class="hidden overflow-x-auto lg:block">
        <table class="w-full text-left text-sm">
          <thead class="bg-earth-bg/80 text-xs uppercase tracking-[0.2em] text-earth-terracotta">
            <tr>
              <th class="px-2 py-3 w-8"></th>
              <th class="px-4 py-3">Plant</th>
              <th class="px-4 py-3">Variety</th>
              <th class="px-4 py-3">Season</th>
              <th class="px-4 py-3">Spacing</th>
              <th class="px-4 py-3">Emerge</th>
              <th class="px-4 py-3">Seed spacing</th>
              <th class="px-4 py-3">Row spacing</th>
              <th class="px-4 py-3">Thinning</th>
              <th class="px-4 py-3">Days</th>
              <th class="px-4 py-3">Sow depth</th>
              <th class="px-4 py-3">When to sow</th>
              <th class="px-4 py-3">Journal</th>
              <th class="px-4 py-3"></th>
            </tr>
          </thead>
          <tbody>
            {#if loading}
              <tr>
                <td class="px-4 py-4 text-earth-text/60" colspan="14">Loading plants…</td>
              </tr>
            {:else if plants.length === 0}
              <tr>
                <td class="px-4 py-4 text-earth-text/60" colspan="14">No plants yet.</td>
              </tr>
            {:else if filteredPlants.length === 0}
              <tr>
                <td class="px-4 py-4 text-earth-text/60" colspan="14">No plants match your filters.</td>
              </tr>
            {:else}
            {#each pagedPlants as plant}
              <tr class="border-t border-earth-border">
                <td class="px-2 py-3 text-center">
                  <button
                    class="text-lg leading-none transition-colors {plant.is_starred ? 'text-amber-400' : 'text-earth-border hover:text-amber-300'}"
                    on:click={() => toggleStar(plant)}
                    title={plant.is_starred ? 'Unstar' : 'Star'}
                  >★</button>
                </td>
                <td class="px-4 py-3 font-semibold">{plant.name}</td>
                <td class="px-4 py-3">{plant.variety || '—'}</td>
                <td class="px-4 py-3">{plant.season || '—'}</td>
                <td class="px-4 py-3">{plant.spacing || '—'}</td>
                <td class="px-4 py-3">{plant.days_to_emerge || '—'}</td>
                <td class="px-4 py-3">{plant.seed_spacing || '—'}</td>
                <td class="px-4 py-3">{plant.row_spacing || '—'}</td>
                <td class="px-4 py-3">{plant.thinning || '—'}</td>
                <td class="px-4 py-3">{plant.days_to_maturity ?? '—'}</td>
                <td class="px-4 py-3">{plant.sow_depth || '—'}</td>
                <td class="px-4 py-3">{plant.when_to_sow || '—'}</td>
                <td class="px-4 py-3">
                  {#if journalLoading}
                    <span class="text-earth-text/60">Loading…</span>
                  {:else if (journalByPlant[plant.id] ?? []).length === 0}
                    <span class="text-earth-text/60">—</span>
                  {:else}
                    <div class="space-y-1">
                      <div class="text-xs text-earth-text/60">
                        {(journalByPlant[plant.id] ?? []).length} entries
                      </div>
                      <div class="text-sm text-earth-text/80">
                        {(journalByPlant[plant.id] ?? [])[0]?.text?.slice(0, 80)}
                        {(journalByPlant[plant.id] ?? [])[0]?.text?.length > 80 ? '…' : ''}
                      </div>
                    </div>
                  {/if}
                </td>
                <td class="px-4 py-3">
                  <div class="flex gap-2">
                    <button
                      class="rounded-full border border-earth-border px-3 py-1 text-xs"
                      on:click={() => openEditModal(plant)}
                    >
                      Edit
                    </button>
                    <button
                      class="rounded-full border border-earth-border px-3 py-1 text-xs text-rose-700"
                      on:click={() => (confirmDeletePlantId = plant.id)}
                    >
                      Delete
                    </button>
                  </div>
                </td>
              </tr>
            {/each}
          {/if}
          </tbody>
        </table>
      </div>
      <div class="block space-y-4 p-4 lg:hidden">
        {#if loading}
          <div class="rounded-2xl border border-earth-border bg-earth-bg/70 p-4 text-sm text-earth-text/60">
            Loading plants…
          </div>
        {:else if plants.length === 0}
          <div class="rounded-2xl border border-earth-border bg-earth-bg/70 p-4 text-sm text-earth-text/60">
            No plants yet.
          </div>
        {:else if filteredPlants.length === 0}
          <div class="rounded-2xl border border-earth-border bg-earth-bg/70 p-4 text-sm text-earth-text/60">
            No plants match your filters.
          </div>
        {:else}
          {#each pagedPlants as plant}
            <div class="rounded-2xl border border-earth-border bg-white p-4">
              <div class="flex items-start justify-between gap-4">
                <div class="flex items-center gap-2">
                  <button
                    class="text-lg leading-none transition-colors {plant.is_starred ? 'text-amber-400' : 'text-earth-border'}"
                    on:click={() => toggleStar(plant)}
                  >★</button>
                  <div>
                    <p class="text-sm font-semibold">{plant.name}</p>
                    <p class="text-xs text-earth-text/60">
                      {plant.variety || '—'} · {plant.season || '—'}
                    </p>
                  </div>
                </div>
                <span class="rounded-full border border-earth-border px-3 py-1 text-xs">
                  {plant.spacing || '—'}
                </span>
              </div>
                <div class="mt-3 grid gap-2 text-xs text-earth-text/70">
                  <div>Days to emerge: {plant.days_to_emerge || '—'}</div>
                  <div>Seed spacing: {plant.seed_spacing || '—'}</div>
                  <div>Row spacing: {plant.row_spacing || '—'}</div>
                  <div>Thinning: {plant.thinning || '—'}</div>
                  <div>Days to maturity: {plant.days_to_maturity ?? '—'}</div>
                  <div>Sow depth: {plant.sow_depth || '—'}</div>
                  <div>When to sow: {plant.when_to_sow || '—'}</div>
                </div>
              <div class="mt-3 text-xs text-earth-text/70">
                Journal: {(journalByPlant[plant.id] ?? []).length || '—'}
              </div>
              <div class="mt-3 flex gap-2">
                <button
                  class="rounded-full border border-earth-border px-3 py-1 text-xs"
                  on:click={() => openEditModal(plant)}
                >
                  Edit
                </button>
                <button
                  class="rounded-full border border-earth-border px-3 py-1 text-xs text-rose-700"
                  on:click={() => (confirmDeletePlantId = plant.id)}
                >
                  Delete
                </button>
              </div>
            </div>
          {/each}
        {/if}
      </div>

      {#if totalPages > 1}
        <div class="flex flex-wrap items-center justify-between gap-4 border-t border-earth-border px-6 py-4">
          <p class="text-xs text-earth-text/60">
            Page {safePage} of {totalPages} · {filteredPlants.length} plants
          </p>
          <div class="flex items-center gap-2">
            <button
              class="rounded-full border border-earth-border px-4 py-1.5 text-xs font-semibold disabled:opacity-40"
              disabled={safePage <= 1}
              on:click={() => (currentPage = safePage - 1)}
            >
              Prev
            </button>
            {#each Array.from({ length: totalPages }, (_, i) => i + 1) as page}
              {#if totalPages <= 7 || page === 1 || page === totalPages || Math.abs(page - safePage) <= 1}
                <button
                  class="rounded-full border px-3 py-1.5 text-xs font-semibold transition-colors {page === safePage ? 'border-earth-forest bg-earth-forest text-earth-bg' : 'border-earth-border text-earth-text/70'}"
                  on:click={() => (currentPage = page)}
                >
                  {page}
                </button>
              {:else if Math.abs(page - safePage) === 2}
                <span class="px-1 text-xs text-earth-text/40">…</span>
              {/if}
            {/each}
            <button
              class="rounded-full border border-earth-border px-4 py-1.5 text-xs font-semibold disabled:opacity-40"
              disabled={safePage >= totalPages}
              on:click={() => (currentPage = safePage + 1)}
            >
              Next
            </button>
          </div>
        </div>
      {/if}
    </div>
  </div>

  <div class="mt-8 rounded-3xl border border-earth-border bg-earth-surface p-6">
    <div class="flex flex-wrap items-center justify-between gap-4">
      <div>
        <h2 class="text-lg font-semibold">Companion guidance</h2>
        <p class="mt-1 text-sm text-earth-text/70">
          Select a plant to see companion planting notes from your dataset.
        </p>
      </div>
      <div class="min-w-[240px]">
        <label class="text-xs uppercase tracking-[0.3em] text-earth-terracotta">Plant</label>
        <select
          class="mt-2 w-full rounded-2xl border border-earth-border bg-white px-3 py-2 text-sm"
          bind:value={selectedPlantId}
        >
          <option value="">Select plant</option>
          {#each plants as plant}
            <option value={plant.id}>
              {plant.name}{plant.variety ? ` · ${plant.variety}` : ''}
            </option>
          {/each}
        </select>
      </div>
    </div>

    {#if companionLoading}
      <div class="mt-4 rounded-2xl border border-earth-border bg-white/70 p-4 text-sm text-earth-text/70">
        Loading companion data…
      </div>
    {:else if companionError}
      <div class="mt-4 rounded-2xl border border-rose-200 bg-rose-50 p-4 text-sm text-rose-700">
        {companionError}
      </div>
    {:else if companionData}
      <div class="mt-4 grid gap-4 md:grid-cols-2">
        <div class="rounded-2xl border border-earth-border bg-white p-4">
          <p class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">Helps</p>
          <p class="mt-2 text-sm text-earth-text/80">
            {companionData.helps || '—'}
          </p>
        </div>
        <div class="rounded-2xl border border-earth-border bg-white p-4">
          <p class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">Helped by</p>
          <p class="mt-2 text-sm text-earth-text/80">
            {companionData.helped_by || '—'}
          </p>
        </div>
        <div class="rounded-2xl border border-earth-border bg-white p-4">
          <p class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">Attracts</p>
          <p class="mt-2 text-sm text-earth-text/80">
            {companionData.attracts || '—'}
          </p>
        </div>
        <div class="rounded-2xl border border-earth-border bg-white p-4">
          <p class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">Repels</p>
          <p class="mt-2 text-sm text-earth-text/80">
            {companionData.repels || '—'}
          </p>
        </div>
        <div class="rounded-2xl border border-earth-border bg-white p-4">
          <p class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">Avoid</p>
          <p class="mt-2 text-sm text-earth-text/80">
            {companionData.avoid || '—'}
          </p>
        </div>
        <div class="rounded-2xl border border-earth-border bg-white p-4">
          <p class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">Notes</p>
          <p class="mt-2 text-sm text-earth-text/80 whitespace-pre-line">
            {companionData.comments || '—'}
          </p>
        </div>
      </div>
    {:else}
      <div class="mt-4 rounded-2xl border border-earth-border bg-white/70 p-4 text-sm text-earth-text/70">
        Select a plant to view companion guidance.
      </div>
    {/if}
  </div>

  {#if showModal}
    <div class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 px-4 py-6">
      <div class="w-full max-w-xl max-h-[90vh] overflow-y-auto rounded-3xl border border-earth-border bg-earth-surface p-6">
        <div class="flex items-center justify-between">
          <h2 class="text-lg font-semibold">{editingPlantId ? 'Edit plant' : 'Add plant'}</h2>
          <button
            class="rounded-full border border-earth-border px-3 py-1 text-xs"
            on:click={() => { resetPlantForm(); showModal = false; }}
          >
            Close
          </button>
        </div>
        <div class="mt-4 grid gap-4 md:grid-cols-2">
          <div>
            <label class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">Name</label>
            <input
              class="mt-2 w-full rounded-2xl border border-earth-border bg-white px-3 py-2 text-sm"
              placeholder="Tomato"
              bind:value={name}
            />
          </div>
          <div>
            <label class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">Variety</label>
            <input
              class="mt-2 w-full rounded-2xl border border-earth-border bg-white px-3 py-2 text-sm"
              placeholder="San Marzano"
              bind:value={variety}
            />
          </div>
          <div>
            <label class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">Season</label>
            <input
              class="mt-2 w-full rounded-2xl border border-earth-border bg-white px-3 py-2 text-sm"
              placeholder="Warm"
              bind:value={season}
            />
          </div>
          <div>
            <label class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">Spacing</label>
            <input
              class="mt-2 w-full rounded-2xl border border-earth-border bg-white px-3 py-2 text-sm"
              placeholder="18-24 in"
              bind:value={spacing}
            />
          </div>
          <div>
            <label class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">
              Days to maturity
            </label>
            <input
              class="mt-2 w-full rounded-2xl border border-earth-border bg-white px-3 py-2 text-sm"
              placeholder="75"
              bind:value={daysToMaturity}
            />
          </div>
          <div>
            <label class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">
              Sow depth
            </label>
            <input
              class="mt-2 w-full rounded-2xl border border-earth-border bg-white px-3 py-2 text-sm"
              placeholder="1/4 in"
              bind:value={sowDepth}
            />
          </div>
          <div>
            <label class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">
              Days to emerge
            </label>
            <input
              class="mt-2 w-full rounded-2xl border border-earth-border bg-white px-3 py-2 text-sm"
              placeholder="5-10 days"
              bind:value={daysToEmerge}
            />
          </div>
          <div>
            <label class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">
              Seed spacing
            </label>
            <input
              class="mt-2 w-full rounded-2xl border border-earth-border bg-white px-3 py-2 text-sm"
              placeholder='2 seeds every 8"'
              bind:value={seedSpacing}
            />
          </div>
          <div class="md:col-span-2">
            <label class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">
              Row spacing
            </label>
            <input
              class="mt-2 w-full rounded-2xl border border-earth-border bg-white px-3 py-2 text-sm"
              placeholder='18"'
              bind:value={rowSpacing}
            />
          </div>
          <div class="md:col-span-2">
            <label class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">
              Thinning
            </label>
            <input
              class="mt-2 w-full rounded-2xl border border-earth-border bg-white px-3 py-2 text-sm"
              placeholder='Thin to 1 every 8"'
              bind:value={thinning}
            />
          </div>
          <div class="md:col-span-2">
            <label class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">Notes</label>
            <textarea
              class="mt-2 h-24 w-full rounded-2xl border border-earth-border bg-white p-3 text-sm"
              placeholder="Succession sow every 2 weeks."
              bind:value={notes}
            ></textarea>
          </div>
        </div>

        {#if formError}
          <div class="mt-4 rounded-2xl border border-rose-200 bg-rose-50 p-3 text-sm text-rose-700">
            {formError}
          </div>
        {/if}

        <div class="mt-4 flex flex-wrap gap-3">
          <button
            class="rounded-full border border-earth-border px-5 py-2.5 text-sm font-semibold"
            on:click={() => { resetPlantForm(); showModal = false; }}
          >
            Cancel
          </button>
          <button
            class="rounded-full bg-earth-forest px-5 py-2.5 text-sm font-semibold text-earth-bg"
            on:click={async () => {
              await submitPlant();
              if (!formError) showModal = false;
            }}
            disabled={isSaving}
          >
            {isSaving ? 'Saving…' : editingPlantId ? 'Update plant' : 'Save plant'}
          </button>
        </div>
      </div>
    </div>
  {/if}

  {#if showScanModal}
    <div class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 px-4 py-6">
      <div class="w-full max-w-xl max-h-[90vh] overflow-y-auto rounded-3xl border border-earth-border bg-earth-surface p-6">
        <div class="flex items-center justify-between">
          <h2 class="text-lg font-semibold">Scan seed packet</h2>
          <button
            class="rounded-full border border-earth-border px-3 py-1 text-xs"
            on:click={() => (showScanModal = false)}
          >
            Close
          </button>
        </div>
        <div class="mt-4">
          <label class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">Image</label>
          <input
            class="mt-2 w-full rounded-2xl border border-earth-border bg-white px-3 py-2 text-sm"
            type="file"
            accept="image/*"
            on:change={(event) => {
              const file = event.currentTarget.files?.[0] ?? null;
              scanFile = file;
              scanPreviewUrl = file ? URL.createObjectURL(file) : '';
            }}
          />
        </div>
        <div class="mt-4">
          <label class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">
            Seed company
          </label>
          <select
            class="mt-2 w-full rounded-2xl border border-earth-border bg-white px-3 py-2 text-sm"
            bind:value={scanCompany}
          >
            <option value="botanical_interests">Botanical Interests</option>
            <option value="">Generic</option>
          </select>
        </div>
        <div class="mt-4 flex flex-wrap items-center gap-3 text-xs text-earth-text/70">
          <label class="flex items-center gap-2">
            <input
              type="checkbox"
              class="h-4 w-4 rounded border-earth-border text-earth-forest"
              bind:checked={helperEnabled}
              on:change={() => {
                if (!helperEnabled) {
                  helperBox = null;
                  helperDragging = false;
                  helperStart = null;
                  helperMappings = {};
                }
              }}
            />
            Box helper
          </label>
          <span>Drag on the image to measure a rectangle.</span>
        </div>
        {#if helperEnabled}
          <div class="mt-3 grid gap-3 text-xs text-earth-text/70 md:grid-cols-[220px_1fr]">
            <div>
              <label class="text-[10px] uppercase tracking-[0.3em] text-earth-terracotta">
                Map to field
              </label>
              <select
                class="mt-2 w-full rounded-2xl border border-earth-border bg-white px-3 py-2 text-xs"
                bind:value={helperAttribute}
              >
                <option value="name">Name</option>
                <option value="variety">Variety</option>
                <option value="season">Season</option>
                <option value="spacing">Spacing</option>
                <option value="specs">Specs block</option>
                <option value="sow">When to sow block</option>
                <option value="emerge">Days to emerge</option>
                <option value="seed_spacing">Seed spacing</option>
                <option value="row_spacing">Row spacing</option>
                <option value="thinning">Thinning</option>
                <option value="days_to_maturity">Days to maturity</option>
                <option value="depth">Sow depth</option>
                <option value="when_to_sow">When to sow</option>
                <option value="notes">Notes</option>
              </select>
            </div>
            <div class="flex flex-wrap items-end gap-2">
              <label class="flex items-center gap-2 text-[11px]">
                <input
                  type="checkbox"
                  class="h-4 w-4 rounded border-earth-border text-earth-forest"
                  bind:checked={helperLockField}
                />
                Lock field
              </label>
              <button
                class="rounded-full border border-earth-border px-3 py-1 text-xs font-semibold"
                type="button"
                on:click={clearHelperBox}
              >
                Clear box
              </button>
              <button
                class="rounded-full border border-earth-border px-3 py-1 text-xs font-semibold"
                type="button"
                on:click={clearHelperMapping}
              >
                Clear mapping
              </button>
              <button
                class="rounded-full border border-earth-border px-3 py-1 text-xs font-semibold"
                type="button"
                on:click={copyHelperMappings}
                disabled={Object.keys(helperMappings).length === 0}
              >
                {helperMappingsCopied ? 'Copied' : 'Copy mappings'}
              </button>
            </div>
          </div>
        {/if}
        {#if scanError}
          <div class="mt-4 rounded-2xl border border-rose-200 bg-rose-50 p-3 text-sm text-rose-700">
            {scanError}
          </div>
        {/if}
        <div class="mt-4 flex flex-wrap gap-3">
          <button
            class="rounded-full border border-earth-border px-5 py-2.5 text-sm font-semibold"
            on:click={() => (showScanModal = false)}
          >
            Cancel
          </button>
          <button
            class="rounded-full bg-earth-forest px-5 py-2.5 text-sm font-semibold text-earth-bg"
            on:click={handleScan}
            disabled={scanLoading}
          >
            {scanLoading ? 'Scanning…' : 'Scan packet'}
          </button>
        </div>

        {#if scanResult}
          <div class="mt-6 rounded-2xl border border-earth-border bg-white p-4">
            <p class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">
              Confirm details
            </p>
            <p class="mt-2 text-sm text-earth-text/70">
              Adjust any fields before saving to the registry.
            </p>
            <div class="mt-4 grid gap-4 md:grid-cols-2">
              <div class="md:col-span-2">
                <div
                  class="relative overflow-hidden rounded-2xl border border-earth-border bg-earth-bg/60"
                  class:touch-none={helperEnabled}
                  bind:this={helperContainer}
                  on:pointerdown={handleHelperPointerDown}
                  on:pointermove={handleHelperPointerMove}
                  on:pointerup={handleHelperPointerUp}
                  on:pointerleave={handleHelperPointerUp}
                >
                  {#if scanPreviewUrl}
                    <img src={scanPreviewUrl} alt="Seed packet preview" class="w-full object-contain" />
                    {#if helperEnabled}
                      {#each Object.entries(helperMappings) as [key, box]}
                        <div
                          class="absolute border-2 border-emerald-500/70"
                          style={`left:${box.x * 100}%;top:${box.y * 100}%;width:${box.w * 100}%;height:${box.h * 100}%;`}
                          title={key}
                        ></div>
                        <div
                          class="absolute rounded-full bg-emerald-600 px-2 py-0.5 text-[10px] font-semibold text-white shadow"
                          style={`left:${box.x * 100}%;top:${box.y * 100}%;transform:translateY(-110%);`}
                        >
                          {helperFieldLabels[key] || key}
                        </div>
                      {/each}
                    {:else if scanResult?.boxes}
                      {#each Object.entries(scanResult.boxes) as [key, box]}
                        <div
                          class="absolute border-2 border-sky-500/70"
                          style={`left:${box.x * 100}%;top:${box.y * 100}%;width:${box.w * 100}%;height:${box.h * 100}%;`}
                          title={key}
                        ></div>
                      {/each}
                    {/if}
                    {#if helperBox}
                      <div
                        class="absolute border-2 border-rose-500/80 bg-rose-200/10"
                        style={`left:${helperBox.x * 100}%;top:${helperBox.y * 100}%;width:${helperBox.w * 100}%;height:${helperBox.h * 100}%;`}
                      ></div>
                      <div
                        class="absolute rounded-full bg-rose-600 px-2 py-0.5 text-[10px] font-semibold text-white shadow"
                        style={`left:${helperBox.x * 100}%;top:${helperBox.y * 100}%;transform:translateY(-110%);`}
                      >
                        {helperFieldLabels[helperAttribute] || helperAttribute}
                      </div>
                    {/if}
                  {/if}
                </div>
                <p class="mt-2 text-xs text-earth-text/60">
                  Green: saved mappings. Rose: active selection.
                </p>
                {#if helperBox}
                  <div class="mt-2 rounded-xl border border-earth-border bg-white/80 p-3 text-xs text-earth-text/70">
                    <div>
                      Normalized: x {helperBox.x.toFixed(3)}, y {helperBox.y.toFixed(3)}, w
                      {helperBox.w.toFixed(3)}, h {helperBox.h.toFixed(3)}
                    </div>
                    <div class="mt-1">
                      Percent: x {(helperBox.x * 100).toFixed(1)}%, y
                      {(helperBox.y * 100).toFixed(1)}%, w
                      {(helperBox.w * 100).toFixed(1)}%, h {(helperBox.h * 100).toFixed(1)}%
                    </div>
                    <div class="mt-2 flex flex-wrap items-center gap-2">
                      <button
                        class="rounded-full border border-earth-border px-3 py-1 text-xs font-semibold"
                        type="button"
                        on:click={copyHelperBox}
                      >
                        {helperCopied ? 'Copied' : 'Copy JSON'}
                      </button>
                      <button
                        class="rounded-full border border-earth-border px-3 py-1 text-xs font-semibold"
                        type="button"
                        on:click={clearHelperBox}
                      >
                        Clear selection
                      </button>
                      <span class="text-[11px] text-earth-text/50">{"{ \"x\": 0.123, \"y\": 0.456, \"w\": 0.222, \"h\": 0.333 }"}</span>
                    </div>
                  </div>
                {/if}
              </div>
              <div>
                <label class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">Name</label>
                <input
                  class="mt-2 w-full rounded-2xl border border-earth-border bg-white px-3 py-2 text-sm"
                  bind:value={name}
                />
              </div>
              <div>
                <label class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">Variety</label>
                <input
                  class="mt-2 w-full rounded-2xl border border-earth-border bg-white px-3 py-2 text-sm"
                  bind:value={variety}
                />
              </div>
              <div>
                <label class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">Season</label>
                <input
                  class="mt-2 w-full rounded-2xl border border-earth-border bg-white px-3 py-2 text-sm"
                  bind:value={season}
                />
              </div>
              <div>
                <label class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">Spacing</label>
                <input
                  class="mt-2 w-full rounded-2xl border border-earth-border bg-white px-3 py-2 text-sm"
                  bind:value={spacing}
                />
              </div>
              <div>
                <label class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">
                  Days to maturity
                </label>
                <input
                  class="mt-2 w-full rounded-2xl border border-earth-border bg-white px-3 py-2 text-sm"
                  bind:value={daysToMaturity}
                />
              </div>
              <div>
                <label class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">
                  Sow depth
                </label>
                <input
                  class="mt-2 w-full rounded-2xl border border-earth-border bg-white px-3 py-2 text-sm"
                  bind:value={sowDepth}
                />
              </div>
              <div>
                <label class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">
                  Days to emerge
                </label>
                <input
                  class="mt-2 w-full rounded-2xl border border-earth-border bg-white px-3 py-2 text-sm"
                  bind:value={daysToEmerge}
                />
              </div>
              <div>
                <label class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">
                  Seed spacing
                </label>
                <input
                  class="mt-2 w-full rounded-2xl border border-earth-border bg-white px-3 py-2 text-sm"
                  bind:value={seedSpacing}
                />
              </div>
              <div class="md:col-span-2">
                <label class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">
                  Row spacing
                </label>
                <input
                  class="mt-2 w-full rounded-2xl border border-earth-border bg-white px-3 py-2 text-sm"
                  bind:value={rowSpacing}
                />
              </div>
              <div class="md:col-span-2">
                <label class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">
                  Thinning
                </label>
                <input
                  class="mt-2 w-full rounded-2xl border border-earth-border bg-white px-3 py-2 text-sm"
                  bind:value={thinning}
                />
              </div>
              <div class="md:col-span-2">
                <label class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">
                  When to sow
                </label>
                <textarea
                  class="mt-2 h-28 w-full rounded-2xl border border-earth-border bg-white p-3 text-sm"
                  bind:value={whenToSow}
                ></textarea>
              </div>
              <div class="md:col-span-2">
                <label class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">Notes</label>
                <textarea
                  class="mt-2 h-24 w-full rounded-2xl border border-earth-border bg-white p-3 text-sm"
                  bind:value={notes}
                ></textarea>
              </div>
            </div>
            <div class="mt-4 flex flex-wrap gap-3">
              <button
                class="rounded-full border border-earth-border px-5 py-2.5 text-sm font-semibold"
                on:click={() => {
                  scanResult = null;
                  scanFile = null;
                  scanPreviewUrl = '';
                }}
              >
                Clear
              </button>
              <button
                class="rounded-full border border-earth-border px-5 py-2.5 text-sm font-semibold"
                on:click={handleScan}
                disabled={scanLoading}
              >
                Retry scan
              </button>
              <button
                class="rounded-full bg-earth-forest px-5 py-2.5 text-sm font-semibold text-earth-bg"
                on:click={async () => {
                  await submitPlant();
                  showScanModal = false;
                  scanResult = null;
                  scanFile = null;
                  scanPreviewUrl = '';
                }}
                disabled={isSaving}
              >
                {isSaving ? 'Saving…' : 'Create plant'}
              </button>
            </div>
          </div>
        {/if}
      </div>
    </div>
  {/if}

  {#if confirmDeletePlantId}
    <div class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 px-4">
      <div class="w-full max-w-lg rounded-3xl border border-earth-border bg-earth-surface p-6">
        <div class="flex items-center justify-between">
          <h2 class="text-lg font-semibold">Delete plant?</h2>
          <button
            class="rounded-full border border-earth-border px-3 py-1 text-xs"
            on:click={() => (confirmDeletePlantId = null)}
          >
            Close
          </button>
        </div>
        <div class="mt-4 rounded-2xl border border-earth-border bg-white/80 p-4 text-sm">
          <p class="text-earth-text/70">This will permanently remove the plant from the registry. This action cannot be undone.</p>
        </div>
        <div class="mt-6 flex flex-wrap gap-3">
          <button
            class="rounded-full border border-earth-border px-5 py-2.5 text-sm font-semibold"
            on:click={() => (confirmDeletePlantId = null)}
          >
            Cancel
          </button>
          <button
            class="rounded-full bg-rose-600 px-5 py-2.5 text-sm font-semibold text-white"
            on:click={async () => {
              await handleDeletePlant(confirmDeletePlantId);
              confirmDeletePlantId = null;
            }}
          >
            Delete plant
          </button>
        </div>
      </div>
    </div>
  {/if}
</section>
