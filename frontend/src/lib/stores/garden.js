import { writable } from 'svelte/store';
import { createGarden, listGardens, updateGarden, deleteGarden } from '$lib/api';

export const gardens = writable([]);
export const selectedGardenId = writable('');
export const gardensLoading = writable(false);
export const gardensError = writable(null);

export async function loadGardens() {
  gardensLoading.set(true);
  gardensError.set(null);
  try {
    const data = await listGardens({ limit: 200 });
    gardens.set(data);
    selectedGardenId.update((current) => {
      if (current) return current;
      return data[0]?.id ?? '';
    });
  } catch (err) {
    gardensError.set(err instanceof Error ? err.message : 'Unable to load gardens.');
  } finally {
    gardensLoading.set(false);
  }
}

export async function addGarden(payload) {
  const created = await createGarden(payload);
  gardens.update((items) => [created, ...items]);
  selectedGardenId.set(created.id);
  return created;
}

export async function editGarden(id, payload) {
  const updated = await updateGarden(id, payload);
  gardens.update((items) => items.map((g) => (g.id === id ? updated : g)));
  return updated;
}

export async function removeGarden(id) {
  await deleteGarden(id);
  let nextId = '';
  gardens.update((items) => {
    const remaining = items.filter((g) => g.id !== id);
    nextId = remaining[0]?.id ?? '';
    return remaining;
  });
  selectedGardenId.update((current) => {
    if (current !== id) return current;
    return nextId;
  });
}
