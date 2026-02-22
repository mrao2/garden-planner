<script>
  import { onMount } from 'svelte';
  import { listBeds } from '$lib/api';
  import {
    gardens,
    selectedGardenId,
    gardensLoading,
    gardensError,
    loadGardens,
    addGarden,
    editGarden,
    removeGarden
  } from '$lib/stores/garden';

  let beds = [];
  let loadingBeds = true;
  let bedsError = null;
  let displayedGardens = [];

  let showGardenModal = false;
  let gardenFormError = null;
  let gardenIsSaving = false;
  let gardenName = '';
  let gardenNotes = '';
  let editingGardenId = null;

  async function loadBedsForGarden() {
    if (!$selectedGardenId) {
      beds = [];
      loadingBeds = false;
      return;
    }
    loadingBeds = true;
    bedsError = null;
    try {
      beds = await listBeds($selectedGardenId, { limit: 6 });
    } catch (err) {
      bedsError = err instanceof Error ? err.message : 'Unable to load beds.';
    } finally {
      loadingBeds = false;
    }
  }

  async function submitGarden() {
    gardenFormError = null;
    if (!gardenName.trim()) {
      gardenFormError = 'Garden name is required.';
      return;
    }
    gardenIsSaving = true;
    try {
      if (editingGardenId) {
        await editGarden(editingGardenId, {
          name: gardenName.trim(),
          notes: gardenNotes.trim()
        });
      } else {
        await addGarden({ name: gardenName.trim(), notes: gardenNotes.trim() });
      }
      gardenName = '';
      gardenNotes = '';
      editingGardenId = null;
      showGardenModal = false;
    } catch (err) {
      gardenFormError = err instanceof Error ? err.message : 'Unable to save garden.';
    } finally {
      gardenIsSaving = false;
    }
  }

  async function handleDeleteGarden(id) {
    if (!confirm('Delete this garden and its beds?')) return;
    try {
      await removeGarden(id);
    } catch (err) {
      gardenFormError = err instanceof Error ? err.message : 'Unable to delete garden.';
    }
  }

  onMount(async () => {
    await loadGardens();
    await loadBedsForGarden();
  });

  $: if ($selectedGardenId) {
    loadBedsForGarden();
  }

  $: displayedGardens =
    $selectedGardenId && $gardens.length
      ? $gardens.filter((garden) => garden.id === $selectedGardenId)
      : $gardens;
</script>

<section class="min-h-screen bg-earth-bg text-earth-text">
  <div class="mx-auto flex min-h-screen w-full max-w-6xl flex-col gap-10 px-6 py-16">
    <header class="flex flex-wrap items-end justify-between gap-6">
      <div>
        <p class="text-xs uppercase tracking-[0.4em] text-earth-terracotta">Dashboard</p>
        <h1 class="mt-4 text-4xl font-semibold leading-tight md:text-5xl">
          Garden Planner Overview
        </h1>
        <p class="mt-4 max-w-2xl text-base text-earth-text/80">
          Track gardens, beds, and seasonal tasks from one calm workspace.
        </p>
      </div>
      <div class="flex flex-wrap gap-3">
        <button
          class="rounded-full bg-earth-forest px-5 py-2.5 text-sm font-semibold text-earth-bg"
          on:click={() => {
            gardenFormError = null;
            editingGardenId = null;
            gardenName = '';
            gardenNotes = '';
            showGardenModal = true;
          }}
        >
          New Garden
        </button>
      </div>
    </header>

    {#if $gardensError}
      <div class="rounded-2xl border border-rose-200 bg-rose-50 p-4 text-sm text-rose-700">
        {$gardensError}
      </div>
    {/if}

    <div class="grid gap-6 md:grid-cols-3">
      <div class="rounded-3xl border border-earth-border bg-earth-surface p-6">
        <p class="text-xs uppercase tracking-[0.3em] text-earth-terracotta">Gardens</p>
        <p class="mt-4 text-3xl font-semibold">
          {$gardensLoading ? '—' : $gardens.length}
        </p>
        <p class="mt-2 text-sm text-earth-text/70">Active spaces you are planning.</p>
      </div>
      <div class="rounded-3xl border border-earth-border bg-earth-surface p-6">
        <p class="text-xs uppercase tracking-[0.3em] text-earth-terracotta">Beds</p>
        <p class="mt-4 text-3xl font-semibold">{loadingBeds ? '—' : beds.length}</p>
        <p class="mt-2 text-sm text-earth-text/70">Planting zones defined so far.</p>
      </div>
      <div class="rounded-3xl border border-earth-border bg-earth-surface p-6">
        <p class="text-xs uppercase tracking-[0.3em] text-earth-terracotta">Tasks</p>
        <p class="mt-4 text-3xl font-semibold">0</p>
        <p class="mt-2 text-sm text-earth-text/70">Upcoming seasonal actions.</p>
      </div>
    </div>

    <div class="grid gap-6 lg:grid-cols-[2fr_1fr]">
      <section class="rounded-3xl border border-earth-border bg-white/70 p-6 shadow-sm">
        <div class="flex items-center justify-between">
          <div>
            <h2 class="text-xl font-semibold">Gardens</h2>
            <p class="mt-1 text-sm text-earth-text/70">Manage and review your gardens.</p>
          </div>
          <button
            class="rounded-full bg-earth-terracotta px-4 py-2 text-xs font-semibold text-earth-bg"
            on:click={() => {
              gardenFormError = null;
              editingGardenId = null;
              gardenName = '';
              gardenNotes = '';
              showGardenModal = true;
            }}
          >
            Add garden
          </button>
        </div>

        <div class="mt-6 space-y-4">
          {#if $gardensLoading}
            <div class="rounded-2xl border border-earth-border bg-earth-bg/70 p-4">
              <p class="text-sm text-earth-text/60">Loading gardens…</p>
            </div>
          {:else if displayedGardens.length === 0}
            <div class="rounded-2xl border border-earth-border bg-earth-bg/70 p-4">
              <p class="text-sm text-earth-text/60">No gardens yet. Create your first one.</p>
            </div>
          {:else}
            {#each displayedGardens as garden}
              <div class="rounded-2xl border border-earth-border bg-earth-bg/70 p-4">
                <div class="flex items-start justify-between gap-4">
                  <div>
                    <p class="text-sm font-semibold">{garden.name}</p>
                    <p class="text-xs text-earth-text/60">{garden.notes || 'No notes yet'}</p>
                  </div>
                  <div class="flex gap-2">
                    <button
                      class="rounded-full border border-earth-border px-3 py-1 text-xs"
                      on:click={() => {
                        gardenFormError = null;
                        editingGardenId = garden.id;
                        gardenName = garden.name;
                        gardenNotes = garden.notes ?? '';
                        showGardenModal = true;
                      }}
                    >
                      Edit
                    </button>
                    <button
                      class="rounded-full border border-earth-border px-3 py-1 text-xs text-rose-700"
                      on:click={() => handleDeleteGarden(garden.id)}
                    >
                      Delete
                    </button>
                  </div>
                </div>
              </div>
            {/each}
          {/if}
        </div>
      </section>

      <section class="rounded-3xl border border-earth-border bg-earth-surface p-6">
        <h2 class="text-xl font-semibold">Beds Snapshot</h2>
        <p class="mt-1 text-sm text-earth-text/70">Geometry for the selected garden.</p>
        <div class="mt-6 space-y-4">
          {#if bedsError}
            <div class="rounded-2xl border border-rose-200 bg-rose-50 p-4 text-sm text-rose-700">
              {bedsError}
            </div>
          {:else if loadingBeds}
            <div class="rounded-2xl border border-earth-border bg-earth-bg/70 p-4">
              <p class="text-sm text-earth-text/60">Loading beds…</p>
            </div>
          {:else if beds.length === 0}
            <div class="rounded-2xl border border-earth-border bg-earth-bg/70 p-4">
              <p class="text-sm text-earth-text/60">No beds yet for this garden.</p>
            </div>
          {:else}
            {#each beds as bed}
              <div class="rounded-2xl border border-earth-border bg-earth-bg/70 p-4">
                <p class="text-sm font-semibold">{bed.name}</p>
                <p class="text-xs text-earth-text/60">
                  {Math.round(bed.width)}
                  {bed.shape === 'rectangle' && bed.length ? ` × ${Math.round(bed.length)}` : ''}
                  {` × ${Math.round(bed.height)}`} · {bed.shape}
                </p>
              </div>
            {/each}
          {/if}
        </div>
      </section>
    </div>


    <section class="rounded-3xl border border-earth-border bg-earth-forest/10 p-6">
      <div class="flex flex-wrap items-center justify-between gap-6">
        <div>
          <h2 class="text-xl font-semibold">Next steps</h2>
          <p class="mt-2 text-sm text-earth-text/70">
            Use the tabs above to manage beds, plants, and journal entries.
          </p>
        </div>
        <a
          class="rounded-full bg-earth-forest px-5 py-2.5 text-sm font-semibold text-earth-bg"
          href="/api/health"
        >
          Check API Health
        </a>
      </div>
    </section>
  </div>
</section>

{#if showGardenModal}
  <div class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 px-4">
    <div class="w-full max-w-xl rounded-3xl border border-earth-border bg-earth-surface p-6">
      <div class="flex items-center justify-between">
        <h2 class="text-lg font-semibold">{editingGardenId ? 'Edit garden' : 'Add garden'}</h2>
        <button
          class="rounded-full border border-earth-border px-3 py-1 text-xs"
          on:click={() => (showGardenModal = false)}
        >
          Close
        </button>
      </div>
      <div class="mt-4 grid gap-4">
        <div>
          <label class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">Name</label>
          <input
            class="mt-2 w-full rounded-2xl border border-earth-border bg-white px-3 py-2 text-sm"
            placeholder="Backyard"
            bind:value={gardenName}
          />
        </div>
        <div>
          <label class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">Notes</label>
          <textarea
            class="mt-2 h-24 w-full rounded-2xl border border-earth-border bg-white p-3 text-sm"
            placeholder="South-facing, raised beds"
            bind:value={gardenNotes}
          ></textarea>
        </div>
      </div>

      {#if gardenFormError}
        <div class="mt-4 rounded-2xl border border-rose-200 bg-rose-50 p-3 text-sm text-rose-700">
          {gardenFormError}
        </div>
      {/if}

      <div class="mt-4 flex flex-wrap gap-3">
        <button
          class="rounded-full border border-earth-border px-5 py-2.5 text-sm font-semibold"
          on:click={() => (showGardenModal = false)}
        >
          Cancel
        </button>
        <button
          class="rounded-full bg-earth-forest px-5 py-2.5 text-sm font-semibold text-earth-bg"
          on:click={submitGarden}
          disabled={gardenIsSaving}
        >
          {gardenIsSaving ? 'Saving…' : editingGardenId ? 'Update garden' : 'Save garden'}
        </button>
      </div>
    </div>
  </div>
{/if}
