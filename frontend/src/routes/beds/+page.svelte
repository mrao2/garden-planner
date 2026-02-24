<script>
  import { onMount } from 'svelte';
  import {
    listBeds,
    updateBed,
    deleteBed,
    listPlants,
    listBedPlantings,
    createBedPlanting,
    deleteBedPlanting,
    getCompanionData,
    listPlanters,
    createPlanter,
    deletePlanter
  } from '$lib/api';
  import { gardens, selectedGardenId, loadGardens } from '$lib/stores/garden';

  let beds = [];
  let planters = [];
  let plants = [];
  let plantingsByBed = {};
  let loading = true;
  let error = null;
  let formError = null;
  let isSaving = false;
  let showModal = false;
  let editingId = null;
  let plantingError = null;
  let plantingIsSaving = false;
  let showPlantingModal = false;
  let plantingBedId = '';
  let plantingPlantId = '';
  let plantingStartDate = '';
  let plantingEndDate = '';
  let plantingNotes = '';
  let plantingCompanion = null;
  let plantingCompanionLoading = false;
  let plantingCompanionError = null;

  let bedName = '';
  let bedNotes = '';
  let bedType = '';
  let bedShape = 'rectangle';
  let bedSun = 'full-sun';
  let bedWatering = 'hand-watered';
  let bedWidth = '';
  let bedLength = '';
  let bedHeight = '';

  let showPlanterModal = false;
  let planterFormError = null;
  let planterIsSaving = false;
  let planterName = '';
  let planterNotes = '';
  let planterType = '';
  let planterKind = 'on-ground';
  let planterShape = 'rectangle';
  let planterSun = 'full-sun';
  let planterWatering = 'hand-watered';
  let planterWidth = '';
  let planterLength = '';
  let planterHeight = '';

  async function loadBeds() {
    if (!$selectedGardenId) {
      beds = [];
      plantingsByBed = {};
      loading = false;
      return;
    }
    loading = true;
    error = null;
    try {
      beds = await listBeds($selectedGardenId, { limit: 200 });
      await loadPlanters();
      await loadPlantings();
    } catch (err) {
      error = err instanceof Error ? err.message : 'Unable to load beds.';
    } finally {
      loading = false;
    }
  }

  async function loadPlanters() {
    if (!$selectedGardenId) {
      planters = [];
      return;
    }
    try {
      planters = await listPlanters($selectedGardenId, { limit: 200 });
    } catch (err) {
      error = err instanceof Error ? err.message : 'Unable to load planters.';
    }
  }

  async function loadPlants() {
    try {
      plants = await listPlants({ limit: 500 });
    } catch (err) {
      plantingError = err instanceof Error ? err.message : 'Unable to load plants.';
    }
  }

  async function loadPlantings() {
    if (!beds.length) {
      plantingsByBed = {};
      return;
    }
    plantingError = null;
    try {
      const entries = await Promise.all(
        beds.map(async (bed) => [bed.id, await listBedPlantings(bed.id)])
      );
      plantingsByBed = Object.fromEntries(entries);
    } catch (err) {
      plantingError = err instanceof Error ? err.message : 'Unable to load plantings.';
    }
  }

  async function submitBed() {
    formError = null;
    if (!$selectedGardenId) {
      formError = 'Select a garden first.';
      return;
    }
    if (!bedName.trim()) {
      formError = 'Bed name is required.';
      return;
    }

    const width = Number(bedWidth);
    const length = bedShape === 'rectangle' ? Number(bedLength) : null;
    const height = Number(bedHeight);

    if (!Number.isFinite(width) || !Number.isFinite(height)) {
      formError = 'Width and height must be numbers.';
      return;
    }
    if (bedShape === 'rectangle' && !Number.isFinite(length)) {
      formError = 'Length must be a number for rectangles.';
      return;
    }
    if (width <= 0 || height <= 0 || (bedShape === 'rectangle' && length !== null && length <= 0)) {
      formError = 'Size values must be positive.';
      return;
    }

    isSaving = true;
    try {
      if (editingId) {
        const updated = await updateBed(editingId, {
          name: bedName.trim(),
          notes: bedNotes.trim(),
          type: bedType.trim(),
          shape: bedShape,
          sun_level: bedSun,
          watering_type: bedWatering,
          width,
          length: bedShape === 'rectangle' ? length : null,
          height
        });
        beds = beds.map((bed) => (bed.id === updated.id ? updated : bed));
      } else {
        const { createBed } = await import('$lib/api');
        const created = await createBed($selectedGardenId, {
          name: bedName.trim(),
          notes: bedNotes.trim(),
          type: bedType.trim(),
          shape: bedShape,
          sun_level: bedSun,
          watering_type: bedWatering,
          width,
          length: bedShape === 'rectangle' ? length : null,
          height
        });
        beds = [created, ...beds];
      }
      bedName = '';
      bedNotes = '';
      bedType = '';
      bedShape = 'rectangle';
      bedSun = 'full-sun';
      bedWatering = 'hand-watered';
      bedWidth = '';
      bedLength = '';
      bedHeight = '';
      editingId = null;
      showModal = false;
    } catch (err) {
      formError = err instanceof Error ? err.message : 'Unable to save bed.';
    } finally {
      isSaving = false;
    }
  }

  async function submitPlanter() {
    planterFormError = null;
    if (!$selectedGardenId) {
      planterFormError = 'Select a garden first.';
      return;
    }
    if (!planterName.trim()) {
      planterFormError = 'Planter name is required.';
      return;
    }

    const width = Number(planterWidth);
    const length = planterShape === 'rectangle' ? Number(planterLength) : null;
    const height = Number(planterHeight);

    if (!Number.isFinite(width) || !Number.isFinite(height)) {
      planterFormError = 'Width and height must be numbers.';
      return;
    }
    if (planterShape === 'rectangle' && !Number.isFinite(length)) {
      planterFormError = 'Length must be a number for rectangles.';
      return;
    }
    if (width <= 0 || height <= 0 || (planterShape === 'rectangle' && length !== null && length <= 0)) {
      planterFormError = 'Size values must be positive.';
      return;
    }

    planterIsSaving = true;
    try {
      const created = await createPlanter($selectedGardenId, {
        name: planterName.trim(),
        notes: planterNotes.trim(),
        type: planterType.trim(),
        planter_type: planterKind,
        shape: planterShape,
        sun_level: planterSun,
        watering_type: planterWatering,
        width,
        length: planterShape === 'rectangle' ? length : null,
        height
      });
      planters = [created, ...planters];
      planterName = '';
      planterNotes = '';
      planterType = '';
      planterKind = 'on-ground';
      planterShape = 'rectangle';
      planterSun = 'full-sun';
      planterWatering = 'hand-watered';
      planterWidth = '';
      planterLength = '';
      planterHeight = '';
      showPlanterModal = false;
    } catch (err) {
      planterFormError = err instanceof Error ? err.message : 'Unable to save planter.';
    } finally {
      planterIsSaving = false;
    }
  }

  function openPlantingModal(bedId) {
    plantingError = null;
    plantingCompanion = null;
    plantingCompanionError = null;
    plantingBedId = bedId;
    plantingPlantId = '';
    plantingStartDate = '';
    plantingEndDate = '';
    plantingNotes = '';
    showPlantingModal = true;
  }

  async function submitPlanting() {
    plantingError = null;
    if (!plantingBedId) {
      plantingError = 'Select a bed.';
      return;
    }
    if (!plantingPlantId) {
      plantingError = 'Select a plant.';
      return;
    }
    if (!plantingStartDate) {
      plantingError = 'Start date is required.';
      return;
    }
    plantingIsSaving = true;
    try {
      const created = await createBedPlanting(plantingBedId, {
        plant_id: plantingPlantId,
        start_date: plantingStartDate,
        end_date: plantingEndDate,
        notes: plantingNotes
      });
      const current = plantingsByBed[plantingBedId] ?? [];
      plantingsByBed = {
        ...plantingsByBed,
        [plantingBedId]: [created, ...current]
      };
      showPlantingModal = false;
    } catch (err) {
      plantingError = err instanceof Error ? err.message : 'Unable to assign plant.';
    } finally {
      plantingIsSaving = false;
    }
  }

  async function handleDelete(id) {
    if (!confirm('Delete this bed?')) return;
    try {
      await deleteBed(id);
      beds = beds.filter((bed) => bed.id !== id);
      const next = { ...plantingsByBed };
      delete next[id];
      plantingsByBed = next;
    } catch (err) {
      error = err instanceof Error ? err.message : 'Unable to delete bed.';
    }
  }

  async function handleDeletePlanter(id) {
    if (!confirm('Delete this planter?')) return;
    try {
      await deletePlanter(id);
      planters = planters.filter((planter) => planter.id !== id);
    } catch (err) {
      error = err instanceof Error ? err.message : 'Unable to delete planter.';
    }
  }

  async function handleDeletePlanting(plantingId, bedId) {
    if (!confirm('Remove this planting?')) return;
    try {
      await deleteBedPlanting(plantingId);
      plantingsByBed = {
        ...plantingsByBed,
        [bedId]: (plantingsByBed[bedId] ?? []).filter((item) => item.id !== plantingId)
      };
    } catch (err) {
      plantingError = err instanceof Error ? err.message : 'Unable to delete planting.';
    }
  }

  onMount(async () => {
    await loadGardens();
    await loadPlants();
    await loadBeds();
  });

  $: if ($selectedGardenId) {
    loadBeds();
  }

  const formatDate = (value) => {
    if (!value) return '—';
    const date = new Date(value);
    if (Number.isNaN(date.getTime())) return value;
    return date.toISOString().slice(0, 10);
  };

  $: if (showPlantingModal && plantingPlantId) {
    const selected = plants.find((plant) => plant.id === plantingPlantId);
    if (selected) {
      loadCompanion(selected.name);
    }
  }

  async function loadCompanion(name) {
    plantingCompanionLoading = true;
    plantingCompanionError = null;
    try {
      plantingCompanion = await getCompanionData(name);
    } catch (err) {
      plantingCompanion = null;
      plantingCompanionError = err instanceof Error ? err.message : 'Unable to load companion data.';
    } finally {
      plantingCompanionLoading = false;
    }
  }
</script>

<section class="mx-auto w-full max-w-6xl px-6 py-12">
  <div class="flex flex-wrap items-end justify-between gap-4">
    <div>
      <p class="text-xs uppercase tracking-[0.4em] text-earth-terracotta">Gardens</p>
      <h1 class="mt-3 text-3xl font-semibold">Garden Manager</h1>
      <p class="mt-2 text-sm text-earth-text/70">
        Track beds and containers with their care details.
      </p>
    </div>
  </div>

  <div class="mt-8 rounded-3xl border border-earth-border bg-earth-surface p-6">
    <div class="flex flex-wrap items-center gap-4">
      <div class="w-full sm:w-auto sm:min-w-[220px]">
        <label class="text-xs uppercase tracking-[0.3em] text-earth-terracotta">Garden</label>
        <select
          class="mt-2 w-full rounded-2xl border border-earth-border bg-white px-3 py-2 text-sm"
          bind:value={$selectedGardenId}
          on:change={loadBeds}
        >
          {#if $gardens.length === 0}
            <option value="">No gardens yet</option>
          {:else}
            {#each $gardens as garden}
              <option value={garden.id}>{garden.name}</option>
            {/each}
          {/if}
        </select>
      </div>
      <button
        class="rounded-full bg-earth-forest px-5 py-2.5 text-sm font-semibold text-earth-bg"
        on:click={() => {
          formError = null;
          editingId = null;
          showModal = true;
        }}
      >
        New Bed
      </button>
      <button
        class="rounded-full border border-earth-border px-5 py-2.5 text-sm font-semibold"
        on:click={() => {
          planterFormError = null;
          showPlanterModal = true;
        }}
      >
        New Planter
      </button>
    </div>

    {#if error}
      <div class="mt-6 rounded-2xl border border-rose-200 bg-rose-50 p-4 text-sm text-rose-700">
        {error}
      </div>
    {/if}

    <div class="mt-6 hidden overflow-x-auto rounded-2xl border border-earth-border bg-white sm:block">
      <table class="w-full text-left text-sm">
        <thead class="bg-earth-bg/80 text-xs uppercase tracking-[0.2em] text-earth-terracotta">
          <tr>
            <th class="px-4 py-3">Bed</th>
            <th class="px-4 py-3">Notes</th>
            <th class="px-4 py-3">Plantings</th>
            <th class="px-4 py-3">Type</th>
            <th class="px-4 py-3">Shape</th>
            <th class="px-4 py-3">Sun</th>
            <th class="px-4 py-3">Watering</th>
            <th class="px-4 py-3">Size</th>
            <th class="px-4 py-3">Actions</th>
          </tr>
        </thead>
        <tbody>
          {#if loading}
            <tr>
              <td class="px-4 py-4 text-earth-text/60" colspan="9">Loading beds…</td>
            </tr>
          {:else if beds.length === 0}
            <tr>
              <td class="px-4 py-4 text-earth-text/60" colspan="9">
                No beds yet for this garden.
              </td>
            </tr>
          {:else}
            {#each beds as bed}
              <tr class="border-t border-earth-border">
                <td class="px-4 py-3 font-semibold">{bed.name}</td>
                <td class="px-4 py-3 text-earth-text/70">{bed.notes || '—'}</td>
                <td class="px-4 py-3">
                  <div class="space-y-2">
                    {(plantingsByBed[bed.id] ?? []).length === 0 ? '—' : ''}
                    {#each plantingsByBed[bed.id] ?? [] as planting}
                      <div class="flex flex-wrap items-center gap-2">
                        <span class="rounded-full bg-earth-bg/70 px-2 py-1 text-xs">
                          {planting.plant_name}
                          {planting.variety ? ` · ${planting.variety}` : ''} ·
                          {formatDate(planting.start_date)}
                          {planting.end_date ? ` → ${formatDate(planting.end_date)}` : ''}
                        </span>
                        <button
                          class="rounded-full border border-earth-border px-2 py-0.5 text-[10px] text-rose-700"
                          on:click={() => handleDeletePlanting(planting.id, bed.id)}
                        >
                          Remove
                        </button>
                      </div>
                    {/each}
                    <button
                      class="rounded-full border border-earth-border px-2 py-1 text-xs"
                      on:click={() => openPlantingModal(bed.id)}
                    >
                      Assign plant
                    </button>
                  </div>
                </td>
                <td class="px-4 py-3">{bed.type || '—'}</td>
                <td class="px-4 py-3">{bed.shape}</td>
                <td class="px-4 py-3">{bed.sun_level || '—'}</td>
                <td class="px-4 py-3">{bed.watering_type || '—'}</td>
                <td class="px-4 py-3">
                  {Math.round(bed.width)}
                  {bed.shape === 'rectangle' && bed.length ? ` × ${Math.round(bed.length)}` : ''}
                  {` × ${Math.round(bed.height)}`}
                </td>
                <td class="px-4 py-3">
                  <div class="flex gap-2">
                    <button
                      class="rounded-full border border-earth-border px-3 py-1 text-xs"
                      on:click={() => {
                        formError = null;
                        editingId = bed.id;
                        bedName = bed.name;
                        bedNotes = bed.notes ?? '';
                        bedType = bed.type ?? '';
                        bedShape = bed.shape ?? 'rectangle';
                        bedSun = bed.sun_level ?? 'full-sun';
                        bedWatering = bed.watering_type ?? 'hand-watered';
                        bedWidth = String(bed.width);
                        bedLength = bed.length ? String(bed.length) : '';
                        bedHeight = String(bed.height);
                        showModal = true;
                      }}
                    >
                      Edit
                    </button>
                    <button
                      class="rounded-full border border-earth-border px-3 py-1 text-xs text-rose-700"
                      on:click={() => handleDelete(bed.id)}
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

    <div class="mt-6 space-y-4 sm:hidden">
      {#if loading}
        <div class="rounded-2xl border border-earth-border bg-white p-4 text-sm text-earth-text/60">Loading beds…</div>
      {:else if beds.length === 0}
        <div class="rounded-2xl border border-earth-border bg-white p-4 text-sm text-earth-text/60">No beds yet for this garden.</div>
      {:else}
        {#each beds as bed}
          <div class="rounded-2xl border border-earth-border bg-white p-4">
            <div class="flex items-start justify-between gap-2">
              <p class="font-semibold">{bed.name}</p>
              <div class="flex shrink-0 gap-2">
                <button
                  class="rounded-full border border-earth-border px-3 py-1 text-xs"
                  on:click={() => {
                    formError = null;
                    editingId = bed.id;
                    bedName = bed.name;
                    bedNotes = bed.notes ?? '';
                    bedType = bed.type ?? '';
                    bedShape = bed.shape ?? 'rectangle';
                    bedSun = bed.sun_level ?? 'full-sun';
                    bedWatering = bed.watering_type ?? 'hand-watered';
                    bedWidth = String(bed.width);
                    bedLength = bed.length ? String(bed.length) : '';
                    bedHeight = String(bed.height);
                    showModal = true;
                  }}
                >Edit</button>
                <button
                  class="rounded-full border border-earth-border px-3 py-1 text-xs text-rose-700"
                  on:click={() => handleDelete(bed.id)}
                >Delete</button>
              </div>
            </div>
            <div class="mt-2 grid grid-cols-2 gap-x-4 gap-y-1 text-xs text-earth-text/70">
              <div>Type: {bed.type || '—'}</div>
              <div>Shape: {bed.shape}</div>
              <div>Sun: {bed.sun_level || '—'}</div>
              <div>Watering: {bed.watering_type || '—'}</div>
              <div class="col-span-2">Size: {Math.round(bed.width)}{bed.shape === 'rectangle' && bed.length ? ` × ${Math.round(bed.length)}` : ''} × {Math.round(bed.height)}</div>
            </div>
            {#if bed.notes}
              <p class="mt-2 text-xs text-earth-text/60">{bed.notes}</p>
            {/if}
            <div class="mt-3 space-y-2">
              <p class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">Plantings</p>
              {#if (plantingsByBed[bed.id] ?? []).length === 0}
                <p class="text-xs text-earth-text/60">None yet.</p>
              {:else}
                {#each plantingsByBed[bed.id] ?? [] as planting}
                  <div class="flex flex-wrap items-center gap-2">
                    <span class="rounded-full bg-earth-bg/70 px-2 py-1 text-xs">
                      {planting.plant_name}{planting.variety ? ` · ${planting.variety}` : ''} · {formatDate(planting.start_date)}{planting.end_date ? ` → ${formatDate(planting.end_date)}` : ''}
                    </span>
                    <button
                      class="rounded-full border border-earth-border px-2 py-0.5 text-[10px] text-rose-700"
                      on:click={() => handleDeletePlanting(planting.id, bed.id)}
                    >Remove</button>
                  </div>
                {/each}
              {/if}
              <button
                class="rounded-full border border-earth-border px-2 py-1 text-xs"
                on:click={() => openPlantingModal(bed.id)}
              >Assign plant</button>
            </div>
          </div>
        {/each}
      {/if}
    </div>
  </div>

  <div class="mt-8 rounded-3xl border border-earth-border bg-earth-surface p-6">
    <div class="flex flex-wrap items-center justify-between gap-4">
      <div>
        <h2 class="text-xl font-semibold">Planters</h2>
        <p class="mt-1 text-sm text-earth-text/70">Containers for patios and porches.</p>
      </div>
      <button
        class="rounded-full border border-earth-border px-4 py-2 text-xs font-semibold"
        on:click={() => {
          planterFormError = null;
          showPlanterModal = true;
        }}
      >
        Add planter
      </button>
    </div>
    <div class="mt-6 hidden overflow-x-auto rounded-2xl border border-earth-border bg-white sm:block">
      <table class="w-full text-left text-sm">
        <thead class="bg-earth-bg/80 text-xs uppercase tracking-[0.2em] text-earth-terracotta">
          <tr>
            <th class="px-4 py-3">Planter</th>
            <th class="px-4 py-3">Type</th>
            <th class="px-4 py-3">Kind</th>
            <th class="px-4 py-3">Sun</th>
            <th class="px-4 py-3">Watering</th>
            <th class="px-4 py-3">Size</th>
            <th class="px-4 py-3">Actions</th>
          </tr>
        </thead>
        <tbody>
          {#if loading}
            <tr>
              <td class="px-4 py-4 text-earth-text/60" colspan="7">Loading planters…</td>
            </tr>
          {:else if planters.length === 0}
            <tr>
              <td class="px-4 py-4 text-earth-text/60" colspan="7">No planters yet.</td>
            </tr>
          {:else}
            {#each planters as planter}
              <tr class="border-t border-earth-border">
                <td class="px-4 py-3 font-semibold">{planter.name}</td>
                <td class="px-4 py-3">{planter.type || '—'}</td>
                <td class="px-4 py-3">{planter.planter_type || '—'}</td>
                <td class="px-4 py-3">{planter.sun_level || '—'}</td>
                <td class="px-4 py-3">{planter.watering_type || '—'}</td>
                <td class="px-4 py-3">
                  {Math.round(planter.width)}
                  {planter.shape === 'rectangle' && planter.length ? ` × ${Math.round(planter.length)}` : ''}
                  {` × ${Math.round(planter.height)}`} · {planter.shape}
                </td>
                <td class="px-4 py-3">
                  <button
                    class="rounded-full border border-earth-border px-3 py-1 text-xs text-rose-700"
                    on:click={() => handleDeletePlanter(planter.id)}
                  >
                    Delete
                  </button>
                </td>
              </tr>
            {/each}
          {/if}
        </tbody>
      </table>
    </div>

    <div class="mt-6 space-y-4 sm:hidden">
      {#if loading}
        <div class="rounded-2xl border border-earth-border bg-white p-4 text-sm text-earth-text/60">Loading planters…</div>
      {:else if planters.length === 0}
        <div class="rounded-2xl border border-earth-border bg-white p-4 text-sm text-earth-text/60">No planters yet.</div>
      {:else}
        {#each planters as planter}
          <div class="rounded-2xl border border-earth-border bg-white p-4">
            <div class="flex items-start justify-between gap-2">
              <p class="font-semibold">{planter.name}</p>
              <button
                class="rounded-full border border-earth-border px-3 py-1 text-xs text-rose-700"
                on:click={() => handleDeletePlanter(planter.id)}
              >Delete</button>
            </div>
            <div class="mt-2 grid grid-cols-2 gap-x-4 gap-y-1 text-xs text-earth-text/70">
              <div>Type: {planter.type || '—'}</div>
              <div>Kind: {planter.planter_type || '—'}</div>
              <div>Sun: {planter.sun_level || '—'}</div>
              <div>Watering: {planter.watering_type || '—'}</div>
              <div class="col-span-2">Size: {Math.round(planter.width)}{planter.shape === 'rectangle' && planter.length ? ` × ${Math.round(planter.length)}` : ''} × {Math.round(planter.height)} · {planter.shape}</div>
            </div>
          </div>
        {/each}
      {/if}
    </div>
  </div>

  {#if plantingError}
    <div class="mt-6 rounded-2xl border border-rose-200 bg-rose-50 p-4 text-sm text-rose-700">
      {plantingError}
    </div>
  {/if}

  {#if showModal}
    <div class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 px-4 py-6">
      <div class="w-full max-w-xl max-h-[90vh] overflow-y-auto rounded-3xl border border-earth-border bg-earth-surface p-6">
        <div class="flex items-center justify-between">
          <h2 class="text-lg font-semibold">{editingId ? 'Edit bed' : 'Add bed'}</h2>
          <button
            class="rounded-full border border-earth-border px-3 py-1 text-xs"
            on:click={() => (showModal = false)}
          >
            Close
          </button>
        </div>

        <div class="mt-4 grid gap-4 md:grid-cols-2">
          <div class="md:col-span-2">
            <label class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">Bed name</label>
            <input
              class="mt-2 w-full rounded-2xl border border-earth-border bg-white px-3 py-2 text-sm"
              placeholder="Bed A"
              bind:value={bedName}
            />
          </div>
          <div>
            <label class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">Type</label>
            <input
              class="mt-2 w-full rounded-2xl border border-earth-border bg-white px-3 py-2 text-sm"
              placeholder="Raised"
              bind:value={bedType}
            />
          </div>
          <div>
            <label class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">Shape</label>
            <select
              class="mt-2 w-full rounded-2xl border border-earth-border bg-white px-3 py-2 text-sm"
              bind:value={bedShape}
            >
              <option value="rectangle">Rectangle</option>
              <option value="circle">Circle</option>
            </select>
          </div>
          <div>
            <label class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">Sun level</label>
            <select
              class="mt-2 w-full rounded-2xl border border-earth-border bg-white px-3 py-2 text-sm"
              bind:value={bedSun}
            >
              <option value="full-sun">Full sun</option>
              <option value="partial-sun">Partial sun</option>
              <option value="partial-shade">Partial shade</option>
              <option value="full-shade">Full shade</option>
            </select>
          </div>
          <div>
            <label class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">
              Watering
            </label>
            <select
              class="mt-2 w-full rounded-2xl border border-earth-border bg-white px-3 py-2 text-sm"
              bind:value={bedWatering}
            >
              <option value="hand-watered">Hand watered</option>
              <option value="irrigated">Irrigated</option>
            </select>
          </div>
          <div>
            <label class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">
              {bedShape === 'circle' ? 'Diameter' : 'Width'}
            </label>
            <input
              class="mt-2 w-full rounded-2xl border border-earth-border bg-white px-3 py-2 text-sm"
              placeholder="200"
              bind:value={bedWidth}
            />
          </div>
          {#if bedShape === 'rectangle'}
            <div>
              <label class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">Length</label>
              <input
                class="mt-2 w-full rounded-2xl border border-earth-border bg-white px-3 py-2 text-sm"
                placeholder="100"
                bind:value={bedLength}
              />
            </div>
            <div>
              <label class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">Height</label>
              <input
                class="mt-2 w-full rounded-2xl border border-earth-border bg-white px-3 py-2 text-sm"
                placeholder="18"
                bind:value={bedHeight}
              />
            </div>
          {:else}
            <div>
              <label class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">Height</label>
              <input
                class="mt-2 w-full rounded-2xl border border-earth-border bg-white px-3 py-2 text-sm"
                placeholder="18"
                bind:value={bedHeight}
              />
            </div>
          {/if}
          <div class="md:col-span-2">
            <label class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">Notes</label>
            <textarea
              class="mt-2 h-24 w-full rounded-2xl border border-earth-border bg-white p-3 text-sm"
              placeholder="Full sun, raised."
              bind:value={bedNotes}
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
            on:click={() => (showModal = false)}
          >
            Cancel
          </button>
          <button
            class="rounded-full bg-earth-forest px-5 py-2.5 text-sm font-semibold text-earth-bg"
            on:click={submitBed}
            disabled={isSaving}
          >
            {isSaving ? 'Saving…' : editingId ? 'Update bed' : 'Save bed'}
          </button>
        </div>
      </div>
    </div>
  {/if}

  {#if showPlantingModal}
    <div class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 px-4 py-6">
      <div class="w-full max-w-xl max-h-[90vh] overflow-y-auto rounded-3xl border border-earth-border bg-earth-surface p-6">
        <div class="flex items-center justify-between">
          <h2 class="text-lg font-semibold">Assign plant to bed</h2>
          <button
            class="rounded-full border border-earth-border px-3 py-1 text-xs"
            on:click={() => (showPlantingModal = false)}
          >
            Close
          </button>
        </div>
        <div class="mt-4 grid gap-4 md:grid-cols-2">
          <div class="md:col-span-2">
            <label class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">Plant</label>
            <select
              class="mt-2 w-full rounded-2xl border border-earth-border bg-white px-3 py-2 text-sm"
              bind:value={plantingPlantId}
            >
              <option value="">Select plant</option>
              {#each plants as plant}
                <option value={plant.id}>
                  {plant.name}{plant.variety ? ` · ${plant.variety}` : ''}
                </option>
              {/each}
            </select>
          </div>
          <div>
            <label class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">
              Start date
            </label>
            <input
              class="mt-2 w-full rounded-2xl border border-earth-border bg-white px-3 py-2 text-sm"
              type="date"
              bind:value={plantingStartDate}
            />
          </div>
          <div>
            <label class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">
              End date (optional)
            </label>
            <input
              class="mt-2 w-full rounded-2xl border border-earth-border bg-white px-3 py-2 text-sm"
              type="date"
              bind:value={plantingEndDate}
            />
          </div>
          <div class="md:col-span-2">
            <label class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">Notes</label>
            <textarea
              class="mt-2 h-24 w-full rounded-2xl border border-earth-border bg-white p-3 text-sm"
              placeholder="Succession planting, shade cloth, etc."
              bind:value={plantingNotes}
            ></textarea>
          </div>
        </div>

        <div class="mt-4 rounded-2xl border border-earth-border bg-white p-4">
          <p class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">
            Companion guidance
          </p>
          {#if plantingCompanionLoading}
            <p class="mt-2 text-sm text-earth-text/60">Loading companion data…</p>
          {:else if plantingCompanionError}
            <p class="mt-2 text-sm text-rose-700">{plantingCompanionError}</p>
          {:else if plantingCompanion}
            <div class="mt-2 grid gap-2 text-sm text-earth-text/80">
              <div><strong>Helps:</strong> {plantingCompanion.helps || '—'}</div>
              <div><strong>Helped by:</strong> {plantingCompanion.helped_by || '—'}</div>
              <div><strong>Avoid:</strong> {plantingCompanion.avoid || '—'}</div>
            </div>
          {:else}
            <p class="mt-2 text-sm text-earth-text/60">Select a plant to view guidance.</p>
          {/if}
        </div>

        {#if plantingError}
          <div class="mt-4 rounded-2xl border border-rose-200 bg-rose-50 p-3 text-sm text-rose-700">
            {plantingError}
          </div>
        {/if}

        <div class="mt-4 flex flex-wrap gap-3">
          <button
            class="rounded-full border border-earth-border px-5 py-2.5 text-sm font-semibold"
            on:click={() => (showPlantingModal = false)}
          >
            Cancel
          </button>
          <button
            class="rounded-full bg-earth-forest px-5 py-2.5 text-sm font-semibold text-earth-bg"
            on:click={submitPlanting}
            disabled={plantingIsSaving}
          >
            {plantingIsSaving ? 'Saving…' : 'Assign plant'}
          </button>
        </div>
      </div>
    </div>
  {/if}

  {#if showPlanterModal}
    <div class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 px-4 py-6">
      <div class="w-full max-w-xl max-h-[90vh] overflow-y-auto rounded-3xl border border-earth-border bg-earth-surface p-6">
        <div class="flex items-center justify-between">
          <h2 class="text-lg font-semibold">Add planter</h2>
          <button
            class="rounded-full border border-earth-border px-3 py-1 text-xs"
            on:click={() => (showPlanterModal = false)}
          >
            Close
          </button>
        </div>
        <div class="mt-4 grid gap-4 md:grid-cols-2">
          <div class="md:col-span-2">
            <label class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">Name</label>
            <input
              class="mt-2 w-full rounded-2xl border border-earth-border bg-white px-3 py-2 text-sm"
              placeholder="Patio pot 1"
              bind:value={planterName}
            />
          </div>
          <div>
            <label class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">Type</label>
            <input
              class="mt-2 w-full rounded-2xl border border-earth-border bg-white px-3 py-2 text-sm"
              placeholder="Terracotta"
              bind:value={planterType}
            />
          </div>
          <div>
            <label class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">
              Planter type
            </label>
            <select
              class="mt-2 w-full rounded-2xl border border-earth-border bg-white px-3 py-2 text-sm"
              bind:value={planterKind}
            >
              <option value="hanging">Hanging</option>
              <option value="grow-bag">Grow bag</option>
              <option value="on-ground">On ground</option>
            </select>
          </div>
          <div>
            <label class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">Shape</label>
            <select
              class="mt-2 w-full rounded-2xl border border-earth-border bg-white px-3 py-2 text-sm"
              bind:value={planterShape}
            >
              <option value="rectangle">Rectangle</option>
              <option value="circle">Circle</option>
            </select>
          </div>
          <div>
            <label class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">Sun level</label>
            <select
              class="mt-2 w-full rounded-2xl border border-earth-border bg-white px-3 py-2 text-sm"
              bind:value={planterSun}
            >
              <option value="full-sun">Full sun</option>
              <option value="partial-sun">Partial sun</option>
              <option value="partial-shade">Partial shade</option>
              <option value="full-shade">Full shade</option>
            </select>
          </div>
          <div>
            <label class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">
              Watering
            </label>
            <select
              class="mt-2 w-full rounded-2xl border border-earth-border bg-white px-3 py-2 text-sm"
              bind:value={planterWatering}
            >
              <option value="hand-watered">Hand watered</option>
              <option value="irrigated">Irrigated</option>
            </select>
          </div>
          <div>
            <label class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">
              {planterShape === 'circle' ? 'Diameter' : 'Width'}
            </label>
            <input
              class="mt-2 w-full rounded-2xl border border-earth-border bg-white px-3 py-2 text-sm"
              placeholder="18"
              bind:value={planterWidth}
            />
          </div>
          {#if planterShape === 'rectangle'}
            <div>
              <label class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">Length</label>
              <input
                class="mt-2 w-full rounded-2xl border border-earth-border bg-white px-3 py-2 text-sm"
                placeholder="24"
                bind:value={planterLength}
              />
            </div>
            <div>
              <label class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">Height</label>
              <input
                class="mt-2 w-full rounded-2xl border border-earth-border bg-white px-3 py-2 text-sm"
                placeholder="12"
                bind:value={planterHeight}
              />
            </div>
          {:else}
            <div>
              <label class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">Height</label>
              <input
                class="mt-2 w-full rounded-2xl border border-earth-border bg-white px-3 py-2 text-sm"
                placeholder="12"
                bind:value={planterHeight}
              />
            </div>
          {/if}
          <div class="md:col-span-2">
            <label class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">Notes</label>
            <textarea
              class="mt-2 h-24 w-full rounded-2xl border border-earth-border bg-white p-3 text-sm"
              placeholder="Needs extra drainage."
              bind:value={planterNotes}
            ></textarea>
          </div>
        </div>

        {#if planterFormError}
          <div class="mt-4 rounded-2xl border border-rose-200 bg-rose-50 p-3 text-sm text-rose-700">
            {planterFormError}
          </div>
        {/if}

        <div class="mt-4 flex flex-wrap gap-3">
          <button
            class="rounded-full border border-earth-border px-5 py-2.5 text-sm font-semibold"
            on:click={() => (showPlanterModal = false)}
          >
            Cancel
          </button>
          <button
            class="rounded-full bg-earth-forest px-5 py-2.5 text-sm font-semibold text-earth-bg"
            on:click={submitPlanter}
            disabled={planterIsSaving}
          >
            {planterIsSaving ? 'Saving…' : 'Save planter'}
          </button>
        </div>
      </div>
    </div>
  {/if}
</section>
