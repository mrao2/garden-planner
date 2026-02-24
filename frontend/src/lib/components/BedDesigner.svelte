<script>
  import { createEventDispatcher } from 'svelte';

  export let beds = [];

  const dispatch = createEventDispatcher();
  const minSize = 24;
  const handles = ['nw', 'ne', 'sw', 'se'];

  let selectedId = beds[0]?.id ?? null;
  let dragState = null;

  function startMove(event, bed) {
    event.preventDefault();
    selectedId = bed.id;
    dragState = {
      id: bed.id,
      mode: 'move',
      handle: 'move',
      startX: event.clientX,
      startY: event.clientY,
      origin: { ...bed }
    };
    event.currentTarget.setPointerCapture(event.pointerId);
  }

  function startResize(event, bed, handle) {
    event.preventDefault();
    event.stopPropagation();
    selectedId = bed.id;
    dragState = {
      id: bed.id,
      mode: 'resize',
      handle,
      startX: event.clientX,
      startY: event.clientY,
      origin: { ...bed }
    };
    event.currentTarget.setPointerCapture(event.pointerId);
  }

  function onPointerMove(event) {
    if (!dragState) return;
    const dx = event.clientX - dragState.startX;
    const dy = event.clientY - dragState.startY;

    beds = beds.map((bed) => {
      if (bed.id !== dragState.id) return bed;
      if (dragState.mode === 'move') {
        return {
          ...bed,
          x: Math.max(0, dragState.origin.x + dx),
          y: Math.max(0, dragState.origin.y + dy)
        };
      }

      let { x, y, width, height } = dragState.origin;
      if (dragState.handle.includes('n')) {
        y = y + dy;
        height = height - dy;
      }
      if (dragState.handle.includes('s')) {
        height = height + dy;
      }
      if (dragState.handle.includes('w')) {
        x = x + dx;
        width = width - dx;
      }
      if (dragState.handle.includes('e')) {
        width = width + dx;
      }

      width = Math.max(minSize, width);
      height = Math.max(minSize, height);
      x = Math.max(0, x);
      y = Math.max(0, y);

      return { ...bed, x, y, width, height };
    });
  }

  function onPointerUp() {
    if (dragState) {
      const updated = beds.find((bed) => bed.id === dragState.id);
      if (updated) {
        dispatch('commit', updated);
      }
    }
    dragState = null;
  }

  function selectBed(bed) {
    selectedId = bed.id;
  }
</script>

<div class="rounded-3xl border border-earth-border bg-earth-surface p-6">
  <div class="flex items-center justify-between">
    <div>
      <h2 class="text-xl font-semibold">Bed Designer</h2>
      <p class="mt-1 text-sm text-earth-text/70">Drag rectangles, resize via corners.</p>
    </div>
    <span class="rounded-full border border-earth-border px-3 py-1 text-xs">
      {beds.length} beds
    </span>
  </div>

  <div class="mt-6 rounded-2xl border border-earth-border bg-earth-bg/70 p-4">
    <svg
      class="h-48 w-full rounded-xl border border-dashed border-earth-border bg-white md:h-72"
      viewBox="0 0 520 280"
      role="img"
      aria-label="Bed layout editor"
      on:pointermove={onPointerMove}
      on:pointerup={onPointerUp}
      on:pointerleave={onPointerUp}
    >
      <rect x="0" y="0" width="520" height="280" fill="url(#grid)" />
      <defs>
        <pattern id="grid" width="20" height="20" patternUnits="userSpaceOnUse">
          <path d="M 20 0 L 0 0 0 20" fill="none" stroke="rgba(43,42,39,0.08)" stroke-width="1" />
        </pattern>
      </defs>

      {#each beds as bed}
        <g on:pointerdown={(event) => startMove(event, bed)} on:click={() => selectBed(bed)}>
          <rect
            x={bed.x}
            y={bed.y}
            width={bed.width}
            height={bed.height}
            rx="10"
            fill="rgba(47,93,80,0.15)"
            stroke={bed.id === selectedId ? 'rgba(47,93,80,0.8)' : 'rgba(47,93,80,0.4)'}
            stroke-width={bed.id === selectedId ? 2 : 1}
          />
          <text
            x={bed.x + 12}
            y={bed.y + 24}
            font-size="12"
            fill="rgba(43,42,39,0.8)"
          >
            {bed.name}
          </text>

          {#if bed.id === selectedId}
            {#each handles as handle}
              <rect
                class="cursor-pointer"
                x={bed.x + (handle.includes('e') ? bed.width - 6 : -6)}
                y={bed.y + (handle.includes('s') ? bed.height - 6 : -6)}
                width="12"
                height="12"
                rx="3"
                fill="rgba(180,90,60,0.9)"
                on:pointerdown={(event) => startResize(event, bed, handle)}
              />
            {/each}
          {/if}
        </g>
      {/each}
    </svg>
  </div>

  <div class="mt-6 grid gap-4 md:grid-cols-2">
    <div class="rounded-2xl border border-earth-border bg-earth-bg/60 p-4">
      <p class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">Selected Bed</p>
      <p class="mt-2 text-sm font-semibold">
        {beds.find((bed) => bed.id === selectedId)?.name ?? 'None'}
      </p>
      <p class="mt-2 text-xs text-earth-text/60">Drag to move, use corners to resize.</p>
    </div>
    <div class="rounded-2xl border border-earth-border bg-earth-bg/60 p-4">
      <p class="text-xs uppercase tracking-[0.2em] text-earth-terracotta">Geometry</p>
      {#if beds.find((bed) => bed.id === selectedId)}
        <p class="mt-2 text-sm">
          {Math.round(beds.find((bed) => bed.id === selectedId).width)} ×
          {Math.round(beds.find((bed) => bed.id === selectedId).height)} at (
          {Math.round(beds.find((bed) => bed.id === selectedId).x)},
          {Math.round(beds.find((bed) => bed.id === selectedId).y)})
        </p>
      {:else}
        <p class="mt-2 text-sm text-earth-text/60">Select a bed to see measurements.</p>
      {/if}
    </div>
  </div>
</div>
