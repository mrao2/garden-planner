<script>
  import { gardens, selectedGardenId } from '$lib/stores/garden';

  const tasks = [
    {
      id: 't1',
      date: '2026-03-01',
      type: 'Start Indoors',
      plant: 'Tomato',
      variety: 'San Marzano',
      notes: '8 weeks before last frost',
      gardenId: ''
    },
    {
      id: 't2',
      date: '2026-03-15',
      type: 'Start Indoors',
      plant: 'Pepper',
      variety: 'Sweet Bell',
      notes: '10 weeks before last frost',
      gardenId: ''
    },
    {
      id: 't3',
      date: '2026-04-10',
      type: 'Direct Sow',
      plant: 'Carrot',
      variety: 'Nantes',
      notes: 'Cool soil OK',
      gardenId: ''
    },
    {
      id: 't4',
      date: '2026-04-25',
      type: 'Transplant',
      plant: 'Lettuce',
      variety: 'Butterhead',
      notes: 'Harden off first',
      gardenId: ''
    },
    {
      id: 't5',
      date: '2026-05-15',
      type: 'Transplant',
      plant: 'Tomato',
      variety: 'San Marzano',
      notes: 'After last frost',
      gardenId: ''
    },
    {
      id: 't6',
      date: '2026-05-18',
      type: 'Direct Sow',
      plant: 'Basil',
      variety: 'Genovese',
      notes: 'Warm soil',
      gardenId: ''
    }
  ];

  const formatter = new Intl.DateTimeFormat('en-US', {
    month: 'short',
    day: 'numeric'
  });

  function groupByWeek(items) {
    const groups = new Map();
    items.forEach((task) => {
      const date = new Date(task.date + 'T00:00:00');
      const weekStart = new Date(date);
      weekStart.setDate(date.getDate() - date.getDay());
      const key = weekStart.toISOString().slice(0, 10);
      if (!groups.has(key)) groups.set(key, []);
      groups.get(key).push(task);
    });
    return Array.from(groups.entries()).sort((a, b) => (a[0] > b[0] ? 1 : -1));
  }

  $: activeGardenName = $gardens.find((g) => g.id === $selectedGardenId)?.name || 'All gardens';
  $: filteredTasks =
    $selectedGardenId ? tasks.filter((t) => t.gardenId === $selectedGardenId) : tasks;
  $: grouped = groupByWeek(filteredTasks);
</script>

<section class="mx-auto w-full max-w-6xl px-6 py-12">
  <div class="flex flex-wrap items-end justify-between gap-4">
    <div>
      <p class="text-xs uppercase tracking-[0.4em] text-earth-terracotta">Planner</p>
      <h1 class="mt-3 text-3xl font-semibold">Planting Timeline</h1>
      <p class="mt-2 text-sm text-earth-text/70">
        Tasks for {activeGardenName}. Dates are based on a USDA 6b last frost assumption.
      </p>
    </div>
  </div>

  <div class="mt-8 space-y-6">
    {#if grouped.length === 0}
      <div class="rounded-2xl border border-earth-border bg-earth-bg/70 p-4">
        <p class="text-sm text-earth-text/60">No planner tasks yet.</p>
      </div>
    {:else}
      {#each grouped as [weekStart, items]}
        <div class="rounded-3xl border border-earth-border bg-earth-surface p-6">
          <div class="flex flex-wrap items-center justify-between gap-4">
            <div>
              <p class="text-xs uppercase tracking-[0.3em] text-earth-terracotta">Week of</p>
              <h2 class="mt-2 text-xl font-semibold">
                {formatter.format(new Date(weekStart + 'T00:00:00'))}
              </h2>
            </div>
            <span class="rounded-full border border-earth-border px-3 py-1 text-xs">
              {items.length} tasks
            </span>
          </div>
          <div class="mt-4 space-y-3">
            {#each items as task}
              <div class="rounded-2xl border border-earth-border bg-white p-4">
                <div class="flex flex-wrap items-center justify-between gap-3">
                  <div>
                    <p class="text-sm font-semibold">
                      {task.type}: {task.plant}
                      {task.variety ? `· ${task.variety}` : ''}
                    </p>
                    <p class="text-xs text-earth-text/60">{task.notes}</p>
                  </div>
                  <span class="rounded-full bg-earth-forest/10 px-3 py-1 text-xs text-earth-forest">
                    {formatter.format(new Date(task.date + 'T00:00:00'))}
                  </span>
                </div>
              </div>
            {/each}
          </div>
        </div>
      {/each}
    {/if}
  </div>
</section>
