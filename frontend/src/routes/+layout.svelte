<script>
  import '../app.css';
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { gardens, selectedGardenId, loadGardens } from '$lib/stores/garden';

  const links = [
    { href: '/', label: 'Dashboard' },
    { href: '/beds', label: 'Garden Manager' },
    { href: '/planner', label: 'Planner' },
    { href: '/plants', label: 'Plant Registry' },
    { href: '/journal', label: 'Journal' }
  ];

  const isActive = (href) =>
    $page.url.pathname === href || (href !== '/' && $page.url.pathname.startsWith(href));

  onMount(() => {
    loadGardens();
  });
</script>

<div class="min-h-screen bg-earth-bg text-earth-text">
  <header class="border-b border-earth-border bg-earth-surface/80">
    <div class="mx-auto flex w-full max-w-6xl flex-col gap-4 px-4 py-4 md:flex-row md:items-center md:justify-between md:px-6">
      <div>
        <p class="text-xs uppercase tracking-[0.4em] text-earth-terracotta">Garden Planner</p>
        <p class="text-lg font-semibold">Local-first workspace</p>
      </div>

      <div class="w-full md:w-auto md:min-w-[220px]">
        <label class="text-[10px] uppercase tracking-[0.3em] text-earth-terracotta">
          Active garden
        </label>
        <select
          class="mt-2 w-full rounded-full border border-earth-border bg-white px-4 py-2 text-sm"
          bind:value={$selectedGardenId}
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

      <nav class="flex w-full flex-wrap gap-2 text-sm md:w-auto">
        {#each links as link}
          <a
            class={`rounded-full px-4 py-2 font-semibold ${
              isActive(link.href)
                ? 'bg-earth-forest text-earth-bg shadow-sm'
                : 'border border-earth-border bg-earth-bg text-earth-text/80'
            }`}
            href={link.href}
          >
            {link.label}
          </a>
        {/each}
      </nav>
    </div>
  </header>

  <main>
    <slot />
  </main>
</div>
