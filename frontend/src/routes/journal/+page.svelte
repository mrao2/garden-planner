<script>
  import { onMount } from 'svelte';
  import {
    createJournalEntry,
    listJournalEntries,
    listPlants,
    listBeds,
    listPlanters,
    updateJournalEntry,
    deleteJournalEntry
  } from '$lib/api';
  import { gardens, selectedGardenId } from '$lib/stores/garden';

  let entryTextRaw = '';
  let entryGardenId = '';
  let entryDate = '';
  let showModal = false;
  let entries = [];
  let plants = [];
  let beds = [];
  let planters = [];
  let detectedPlants = [];
  let entryInput;
  let showPlantPicker = false;
  let plantQuery = '';
  let plantMatches = [];
  let slashCommandStart = null;
  let slashCommandEnd = null;
  let slashCommandType = null;
  let savedSelectionRange = null;
  let loading = true;
  let error = null;
  let formError = null;
  let isSaving = false;
  let editingEntryId = null;
  let confirmDeleteId = null;

  async function loadEntries() {
    loading = true;
    error = null;
    try {
      entries = await listJournalEntries({
        limit: 200,
        garden_id: $selectedGardenId || undefined
      });
    } catch (err) {
      error = err instanceof Error ? err.message : 'Unable to load journal entries.';
    } finally {
      loading = false;
    }
  }

  async function loadPlantsForDetection() {
    try {
      plants = await listPlants({ limit: 500 });
    } catch (err) {
      // Keep detection best-effort; journal entries can still be saved.
    }
  }

  async function addEntry() {
    formError = null;
    const plainText = stripMarkers(entryTextRaw).trim();
    if (!plainText) {
      formError = 'Journal entry text is required.';
      return;
    }
    if (!entryGardenId) {
      formError = 'Select a garden for this entry.';
      return;
    }
    isSaving = true;
    try {
      if (editingEntryId) {
        const updated = await updateJournalEntry(editingEntryId, {
          text: plainText,
          garden_id: entryGardenId,
          entry_date: entryDate || undefined
        });
        entries = entries.map((entry) => (entry.id === updated.id ? updated : entry));
      } else {
        const created = await createJournalEntry({
          text: plainText,
          garden_id: entryGardenId,
          entry_date: entryDate || undefined
        });
        entries = [created, ...entries];
      }
      entryTextRaw = '';
      entryDate = '';
      editingEntryId = null;
      if (entryInput) {
        entryInput.innerHTML = '';
      }
    } catch (err) {
      formError = err instanceof Error ? err.message : 'Unable to save entry.';
    } finally {
      isSaving = false;
    }
  }

  onMount(loadEntries);
  onMount(loadPlantsForDetection);
  onMount(() => {
    if ($selectedGardenId && !entryGardenId) {
      entryGardenId = $selectedGardenId;
    }
    loadScopedOptions();
  });

  const formatDate = (value) => {
    if (!value) return '';
    const date = new Date(value);
    if (Number.isNaN(date.getTime())) return value;
    return date.toISOString().slice(0, 10);
  };

  const normalizeForMatch = (value) => {
    const lowered = value.toLowerCase();
    let out = '';
    let lastWasSpace = false;
    for (const ch of lowered) {
      const code = ch.charCodeAt(0);
      const isAlphaNum =
        (code >= 97 && code <= 122) || (code >= 48 && code <= 57);
      if (isAlphaNum) {
        out += ch;
        lastWasSpace = false;
      } else if (!lastWasSpace) {
        out += ' ';
        lastWasSpace = true;
      }
    }
    return out.trim();
  };

  const detectPlants = (text, list) => {
    if (!text || !list.length) return [];
    const normalizedText = ` ${normalizeForMatch(text)} `;
    let varietyMatched = false;
    const matches = [];
    const seen = new Set();
    list.forEach((plant) => {
      const nameKey = normalizeForMatch(plant.name || '');
      const varietyKey = normalizeForMatch(plant.variety || '');
      if (varietyKey && normalizedText.includes(` ${varietyKey} `)) {
        if (!varietyMatched) {
          varietyMatched = true;
          matches.length = 0;
          seen.clear();
        }
        if (!seen.has(plant.id)) {
          seen.add(plant.id);
          matches.push(plant);
        }
        return;
      }
      if (varietyMatched) return;
      if (nameKey && normalizedText.includes(` ${nameKey} `)) {
        if (!seen.has(plant.id)) {
          seen.add(plant.id);
          matches.push(plant);
        }
      }
    });
    return matches;
  };

  $: detectedPlants = detectPlants(stripMarkers(entryTextRaw), plants);

  const stripMarkers = (value) =>
    value.replace(/\[\[(plant|bed|planter):[^\]|]+\|([^\]]+)\]\]/g, '$2');

  const renderEntryHtml = (value) => {
    if (!value) return '';
    const escaped = escapeHtml(value);
    return escaped.replace(
      /\[\[(plant|bed|planter):([^\]|]+)\|([^\]]+)\]\]/g,
      (_, type, id, label) =>
        `<span class="inline-flex items-center gap-2 rounded-full border border-earth-border bg-earth-bg/70 px-3 py-1 text-xs font-semibold text-earth-text" data-token="${type}" data-id="${id}" contenteditable="false" title="${escapeHtml(label)}">${escapeHtml(label)}</span>`
    );
  };

  const escapeHtml = (value) =>
    value
      .replaceAll('&', '&amp;')
      .replaceAll('<', '&lt;')
      .replaceAll('>', '&gt;')
      .replaceAll('"', '&quot;')
      .replaceAll("'", '&#39;');

  function serializeEntryContent() {
    if (!entryInput) return '';
    const parts = [];
    entryInput.childNodes.forEach((node) => {
      if (node.nodeType === Node.TEXT_NODE) {
        parts.push(node.textContent || '');
        return;
      }
      if (node.nodeType === Node.ELEMENT_NODE) {
        const element = node;
        const token = element.dataset?.token;
        const id = element.dataset?.id;
        if (token && id) {
          const label = element.textContent || '';
          parts.push(`[[${token}:${id}|${label}]]`);
          return;
        }
        parts.push(element.textContent || '');
      }
    });
    return parts.join('').replace(/\u200B/g, '');
  }

  function updateSlashCommand() {
    if (!entryInput) return;
    const before = getTextBeforeCaret(entryInput);
    const candidates = [
      { type: 'plant', index: before.lastIndexOf('/plant') },
      { type: 'bed', index: before.lastIndexOf('/bed') },
      { type: 'planter', index: before.lastIndexOf('/container') }
    ];
    const active = candidates.reduce((best, current) => {
      if (current.index > best.index) return current;
      return best;
    }, { type: null, index: -1 });
    if (active.index === -1) {
      resetPlantPicker();
      return;
    }
    if (active.index > 0 && !/\s/.test(before[active.index - 1])) {
      resetPlantPicker();
      return;
    }
    const commandSlice = before.slice(active.index);
    if (commandSlice.includes('\n')) {
      resetPlantPicker();
      return;
    }
    const query = commandSlice.replace(/^\/(plant|bed|container)\s*/i, '');
    slashCommandStart = active.index;
    slashCommandEnd = before.length;
    slashCommandType = active.type;
    plantQuery = query;
    plantMatches = filterMatches(active.type, query);
    showPlantPicker = true;
  }

  function resetPlantPicker() {
    showPlantPicker = false;
    plantQuery = '';
    plantMatches = [];
    slashCommandStart = null;
    slashCommandEnd = null;
    slashCommandType = null;
  }

  function filterMatches(type, query) {
    const trimmed = normalizeForMatch(query || '');
    if (type === 'plant') {
      if (!plants.length) return [];
      if (!trimmed) return plants.slice(0, 50);
      const byName = plants.filter((plant) =>
        normalizeForMatch(plant.name || '').includes(trimmed)
      );
      if (byName.length > 0) return byName.slice(0, 50);
      return plants
        .filter((plant) => normalizeForMatch(plant.variety || '').includes(trimmed))
        .slice(0, 50);
    }
    if (type === 'bed') {
      if (!beds.length) return [];
      if (!trimmed) return beds.slice(0, 50);
      return beds
        .filter((bed) => normalizeForMatch(bed.name || '').includes(trimmed))
        .slice(0, 50);
    }
    if (type === 'planter') {
      if (!planters.length) return [];
      if (!trimmed) return planters.slice(0, 50);
      return planters
        .filter((planter) => normalizeForMatch(planter.name || '').includes(trimmed))
        .slice(0, 50);
    }
    return [];
  }

  $: plantMatches = showPlantPicker ? filterMatches(slashCommandType, plantQuery) : [];

  function insertSlashMatch(match) {
    if (slashCommandStart === null || slashCommandEnd === null || !slashCommandType) return;
    if (slashCommandType === 'plant') {
      const label = `${match.name}${match.variety ? ` ${match.variety}` : ''}`.trim();
      insertToken(label, 'plant', match.id, slashCommandStart, slashCommandEnd);
      return;
    }
    const label = `${match.name}`.trim();
    insertToken(label, slashCommandType, match.id, slashCommandStart, slashCommandEnd);
  }

  function saveContentEditableSelection() {
    if (!entryInput) return;
    const sel = window.getSelection();
    if (sel && sel.rangeCount > 0 && entryInput.contains(sel.anchorNode)) {
      savedSelectionRange = sel.getRangeAt(0).cloneRange();
    }
  }

  function insertToken(label, type, id, startIndex, endIndex) {
    if (!entryInput) return;
    const range = getRangeFromIndices(entryInput, startIndex, endIndex);
    if (!range) return;
    range.deleteContents();
    const spacer = startIndex > 0 ? ' ' : '';
    if (spacer) {
      range.insertNode(document.createTextNode(spacer));
      range.collapse(false);
    }
    const pill = document.createElement('span');
    pill.className =
      'inline-flex items-center gap-2 rounded-full border border-earth-border bg-earth-bg/70 px-3 py-1 text-xs font-semibold text-earth-text';
    pill.dataset.token = type;
    pill.dataset.id = id;
    pill.setAttribute('contenteditable', 'false');
    pill.textContent = label;
    range.insertNode(pill);
    range.collapse(false);
    const spaceNode = document.createTextNode('\u200B ');
    range.insertNode(spaceNode);
    resetPlantPicker();
    entryTextRaw = serializeEntryContent();
    updatePlantAttributes();
    // Use setTimeout to let mobile browsers settle before restoring focus/caret
    setTimeout(() => {
      entryInput.focus();
      const selection = window.getSelection();
      if (selection) {
        selection.removeAllRanges();
        const newRange = document.createRange();
        newRange.setStartAfter(spaceNode);
        newRange.collapse(true);
        selection.addRange(newRange);
      }
    }, 10);
  }

  function getTextBeforeCaret(root) {
    const selection = window.getSelection();
    if (!selection || selection.rangeCount === 0) return '';
    const range = selection.getRangeAt(0).cloneRange();
    range.selectNodeContents(root);
    range.setEnd(selection.getRangeAt(0).endContainer, selection.getRangeAt(0).endOffset);
    return range.toString();
  }

  function getRangeFromIndices(root, start, end) {
    const range = document.createRange();
    let currentIndex = 0;
    let startSet = false;

    const walker = document.createTreeWalker(root, NodeFilter.SHOW_TEXT, null);
    while (walker.nextNode()) {
      const node = walker.currentNode;
      const textLength = node.textContent?.length ?? 0;
      if (!startSet && currentIndex + textLength >= start) {
        range.setStart(node, Math.max(0, start - currentIndex));
        startSet = true;
      }
      if (startSet && currentIndex + textLength >= end) {
        range.setEnd(node, Math.max(0, end - currentIndex));
        return range;
      }
      currentIndex += textLength;
    }
    return null;
  }

  function updatePlantAttributes() {
    if (!entryInput) return;
    entryInput.querySelectorAll('[data-token="plant"]').forEach((node) => {
      const plantId = node.dataset.id;
      const plant = plants.find((item) => item.id === plantId);
      if (!plant) return;
      const details = [
        plant.name,
        plant.variety,
        plant.season,
        plant.spacing
      ]
        .filter(Boolean)
        .join(' · ');
      node.title = details;
    });
  }

  async function loadScopedOptions() {
    if (!entryGardenId) {
      beds = [];
      planters = [];
      return;
    }
    try {
      beds = await listBeds(entryGardenId, { limit: 200 });
    } catch (err) {
      beds = [];
    }
    try {
      planters = await listPlanters(entryGardenId, { limit: 200 });
    } catch (err) {
      planters = [];
    }
  }

  $: if ($selectedGardenId && !entryGardenId) {
    entryGardenId = $selectedGardenId;
  }

  $: if ($selectedGardenId !== undefined) {
    loadEntries();
  }

  $: if (entryGardenId) {
    loadScopedOptions();
  }
 
  async function syncEntryInput() {
    if (!entryInput) return;
    entryInput.innerHTML = renderEntryHtml(entryTextRaw);
    updatePlantAttributes();
  }

  function handleEntryInput() {
    entryTextRaw = serializeEntryContent();
    updateSlashCommand();
  }

  async function handleDeleteEntry(entryId) {
    try {
      await deleteJournalEntry(entryId);
      entries = entries.filter((entry) => entry.id !== entryId);
    } catch (err) {
      error = err instanceof Error ? err.message : 'Unable to delete entry.';
    }
  }
</script>

<section class="mx-auto w-full max-w-6xl px-6 py-12">
  <div class="flex flex-wrap items-end justify-between gap-4">
    <div>
      <p class="text-xs uppercase tracking-[0.4em] text-earth-terracotta">Journal</p>
      <h1 class="mt-3 text-3xl font-semibold">Garden Journal</h1>
      <p class="mt-2 text-sm text-earth-text/70">
        Capture observations, weather notes, and planting decisions.
      </p>
    </div>
  </div>

  <div class="mt-8 grid gap-6 lg:grid-cols-[1fr]">
    <div class="rounded-3xl border border-earth-border bg-white p-6">
      <div class="flex flex-wrap items-center justify-between gap-4">
        <div>
          <h2 class="text-lg font-semibold">Recent entries</h2>
          <p class="mt-1 text-sm text-earth-text/70">Capture observations and decisions.</p>
        </div>
        <button
          class="rounded-full bg-earth-forest px-5 py-2.5 text-sm font-semibold text-earth-bg"
          on:click={() => {
            if ($selectedGardenId) {
              entryGardenId = $selectedGardenId;
            }
            entryTextRaw = '';
            entryDate = '';
            editingEntryId = null;
            queueMicrotask(syncEntryInput);
            showModal = true;
          }}
        >
          Add entry
        </button>
      </div>
      <div class="mt-4 space-y-4">
        {#if loading}
          <p class="text-sm text-earth-text/60">Loading journal entries…</p>
        {:else if error}
          <div class="rounded-2xl border border-rose-200 bg-rose-50 p-4 text-sm text-rose-700">
            {error}
          </div>
        {:else if entries.length === 0}
          <p class="text-sm text-earth-text/60">No journal entries yet.</p>
        {:else}
          {#each entries as entry}
            <div class="rounded-2xl border border-earth-border bg-earth-bg/70 p-4">
              <div class="flex flex-wrap items-center justify-between gap-2">
                <p class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">
                  {formatDate(entry.entry_date)}
                </p>
                <div class="flex flex-wrap gap-2">
                  <span class="rounded-full border border-earth-border px-3 py-1 text-xs">
                    {$gardens.find((g) => g.id === entry.garden_id)?.name ?? 'Garden'}
                  </span>
                  {#if entry.plants?.length}
                    {#each entry.plants as plant}
                      <span class="rounded-full bg-earth-forest/10 px-3 py-1 text-xs text-earth-forest">
                        {plant.name}{plant.variety ? ` · ${plant.variety}` : ''}
                      </span>
                    {/each}
                  {/if}
                </div>
              </div>
              <p class="mt-2 text-sm text-earth-text/80">{entry.text}</p>
              <div class="mt-4 flex flex-wrap gap-2">
                <button
                  class="rounded-full border border-earth-border px-3 py-1 text-xs"
                  on:click={() => {
                    editingEntryId = entry.id;
                    entryGardenId = entry.garden_id;
                    entryDate = formatDate(entry.entry_date);
                    entryTextRaw = entry.text;
                    queueMicrotask(syncEntryInput);
                    showModal = true;
                  }}
                >
                  Edit
                </button>
                <button
                  class="rounded-full border border-earth-border px-3 py-1 text-xs text-rose-700"
                  on:click={() => (confirmDeleteId = entry.id)}
                >
                  Delete
                </button>
              </div>
            </div>
          {/each}
        {/if}
      </div>
    </div>
  </div>

  {#if showModal}
    <div class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 px-4 py-6">
      <div class="w-full max-w-xl max-h-[90vh] overflow-y-auto rounded-3xl border border-earth-border bg-earth-surface p-6">
        <div class="flex items-center justify-between">
          <h2 class="text-lg font-semibold">New journal entry</h2>
          <button
            class="rounded-full border border-earth-border px-3 py-1 text-xs"
            on:click={() => (showModal = false)}
          >
            Close
          </button>
        </div>
        <div class="mt-4">
          <label class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">
            Garden
          </label>
          <select
            class="mt-2 w-full rounded-2xl border border-earth-border bg-white px-3 py-2 text-sm"
            bind:value={entryGardenId}
          >
            <option value="">Select garden</option>
            {#each $gardens as garden}
              <option value={garden.id}>{garden.name}</option>
            {/each}
          </select>
        </div>
        <div class="mt-4">
          <label class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">
            Entry date
          </label>
          <input
            class="mt-2 w-full rounded-2xl border border-earth-border bg-white px-3 py-2 text-sm"
            type="date"
            bind:value={entryDate}
          />
        </div>
        <div class="relative mt-4">
          <div
            class="min-h-[160px] max-h-[40vh] overflow-y-auto w-full rounded-2xl border border-earth-border bg-white p-3 text-sm"
            contenteditable
            role="textbox"
            aria-multiline="true"
            data-placeholder="Today I observed…"
            bind:this={entryInput}
            on:input={handleEntryInput}
            on:click={updateSlashCommand}
            on:keyup={updateSlashCommand}
            on:blur={saveContentEditableSelection}
          ></div>
          {#if showPlantPicker}
            <div class="absolute left-3 right-3 top-full z-10 mt-1 rounded-2xl border border-earth-border bg-white p-3 shadow-lg">
              <p class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">
                /{slashCommandType === 'planter' ? 'container' : slashCommandType} matches
              </p>
              <input
                type="text"
                placeholder="Search…"
                class="mt-2 w-full rounded-xl border border-earth-border bg-earth-bg px-3 py-1.5 text-sm outline-none focus:border-earth-terracotta"
                bind:value={plantQuery}
              />
              {#if plantMatches.length === 0}
                <p class="mt-2 text-sm text-earth-text/60">No matches.</p>
              {:else}
                <div class="mt-2 flex max-h-64 flex-col gap-2 overflow-y-auto">
                  {#each plantMatches as match}
                    <button
                      type="button"
                      class="rounded-xl border border-earth-border px-3 py-2 text-left text-sm hover:bg-earth-bg/70 active:bg-earth-bg"
                      on:mousedown|preventDefault={() => {}}
                      on:click={() => insertSlashMatch(match)}
                    >
                      <span class="font-semibold">{match.name}</span>
                      {match.variety ? ` · ${match.variety}` : ''}
                    </button>
                  {/each}
                </div>
              {/if}
            </div>
          {/if}
        </div>
        <div class="mt-4 rounded-2xl border border-earth-border bg-white p-3">
          <p class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">
            Detected plants
          </p>
          {#if detectedPlants.length === 0}
            <p class="mt-2 text-sm text-earth-text/60">No plant matches yet.</p>
          {:else}
            <div class="mt-2 flex flex-wrap gap-2">
              {#each detectedPlants as plant}
                <span class="rounded-full bg-earth-forest/10 px-3 py-1 text-xs text-earth-forest">
                  {plant.name}{plant.variety ? ` · ${plant.variety}` : ''}
                </span>
              {/each}
            </div>
          {/if}
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
            on:click={async () => {
              await addEntry();
              if (!formError) {
                showModal = false;
              }
            }}
            disabled={isSaving || !entryGardenId}
          >
            {isSaving ? 'Saving…' : 'Save entry'}
          </button>
        </div>
      </div>
    </div>
  {/if}

  {#if confirmDeleteId}
    <div class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 px-4">
      <div class="w-full max-w-lg rounded-3xl border border-earth-border bg-earth-surface p-6">
        <div class="flex items-center justify-between">
          <h2 class="text-lg font-semibold">Delete journal entry?</h2>
          <button
            class="rounded-full border border-earth-border px-3 py-1 text-xs"
            on:click={() => (confirmDeleteId = null)}
          >
            Close
          </button>
        </div>
        <div class="mt-4 rounded-2xl border border-earth-border bg-white/80 p-4 text-sm">
          <p class="text-earth-text/70">This action cannot be undone.</p>
        </div>
        <div class="mt-6 flex flex-wrap gap-3">
          <button
            class="rounded-full border border-earth-border px-5 py-2.5 text-sm font-semibold"
            on:click={() => (confirmDeleteId = null)}
          >
            Cancel
          </button>
          <button
            class="rounded-full bg-rose-600 px-5 py-2.5 text-sm font-semibold text-white"
            on:click={async () => {
              await handleDeleteEntry(confirmDeleteId);
              confirmDeleteId = null;
            }}
          >
            Delete entry
          </button>
        </div>
      </div>
    </div>
  {/if}
</section>
