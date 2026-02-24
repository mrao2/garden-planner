export type Garden = {
  id: string;
  name: string;
  notes: string;
  created_at: string;
  updated_at: string;
};

export type Bed = {
  id: string;
  garden_id: string;
  name: string;
  notes: string;
  type: string;
  shape: string;
  sun_level: string;
  watering_type: string;
  width: number;
  length?: number | null;
  height: number;
  created_at: string;
  updated_at: string;
};

export type Planter = {
  id: string;
  garden_id: string;
  name: string;
  notes: string;
  type: string;
  planter_type: string;
  shape: string;
  sun_level: string;
  watering_type: string;
  width: number;
  length?: number | null;
  height: number;
  created_at: string;
  updated_at: string;
};

export type Plant = {
  id: string;
  name: string;
  variety: string;
  season: string;
  spacing: string;
  notes: string;
  days_to_maturity?: number | null;
  sow_depth: string;
  when_to_sow: string;
  days_to_emerge: string;
  seed_spacing: string;
  row_spacing: string;
  thinning: string;
  is_starred: boolean;
  created_at: string;
  updated_at: string;
};

export type BedPlanting = {
  id: string;
  bed_id: string;
  plant_id: string;
  plant_name: string;
  variety: string;
  start_date: string;
  end_date?: string | null;
  notes: string;
  created_at: string;
  updated_at: string;
};

export type JournalPlantRef = {
  id: string;
  name: string;
  variety: string;
};

export type JournalEntry = {
  id: string;
  garden_id: string;
  entry_date: string;
  text: string;
  plants: JournalPlantRef[];
  created_at: string;
  updated_at: string;
};

export type CompanionRecord = {
  common_name: string;
  scientific_name: string;
  helps: string;
  helped_by: string;
  attracts: string;
  repels: string;
  avoid: string;
  comments: string;
  helps_list: string[];
  helped_by_list: string[];
  attracts_list: string[];
  repels_list: string[];
  avoid_list: string[];
};

type DataEnvelope<T> = {
  data: T;
};

type ListParams = {
  limit?: number;
  offset?: number;
  q?: string;
  garden_id?: string;
};

type CreateGardenInput = {
  name: string;
  notes?: string;
};

type UpdateGardenInput = {
  name?: string;
  notes?: string;
};

type CreateBedInput = {
  name: string;
  notes?: string;
  type: string;
  shape: string;
  sun_level: string;
  watering_type: string;
  width: number;
  length?: number | null;
  height: number;
};

type UpdateBedInput = Partial<CreateBedInput>;

type CreatePlanterInput = {
  name: string;
  notes?: string;
  type: string;
  planter_type: string;
  shape: string;
  sun_level: string;
  watering_type: string;
  width: number;
  length?: number | null;
  height: number;
};

type UpdatePlanterInput = Partial<CreatePlanterInput>;

type CreatePlantInput = {
  name: string;
  variety?: string;
  season?: string;
  spacing?: string;
  notes?: string;
  days_to_maturity?: number | null;
  sow_depth?: string;
  when_to_sow?: string;
  days_to_emerge?: string;
  seed_spacing?: string;
  row_spacing?: string;
  thinning?: string;
};

type UpdatePlantInput = Partial<CreatePlantInput>;

type CreateBedPlantingInput = {
  plant_id: string;
  start_date: string;
  end_date?: string;
  notes?: string;
};

type UpdateBedPlantingInput = Partial<CreateBedPlantingInput>;

type CreateJournalEntryInput = {
  text: string;
  entry_date?: string;
  garden_id: string;
};

type UpdateJournalEntryInput = {
  text?: string;
  entry_date?: string;
  garden_id?: string;
};

export async function apiFetch<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(path, {
    headers: {
      'Content-Type': 'application/json',
      ...(init?.headers ?? {})
    },
    ...init
  });

  if (!response.ok) {
    const message = await response.text();
    throw new Error(message || `Request failed: ${response.status}`);
  }

  return (await response.json()) as T;
}

function toQuery(params?: ListParams): string {
  if (!params) return '';
  const search = new URLSearchParams();
  if (params.limit !== undefined) search.set('limit', String(params.limit));
  if (params.offset !== undefined) search.set('offset', String(params.offset));
  if (params.q) search.set('q', params.q);
  if (params.garden_id) search.set('garden_id', params.garden_id);
  const qs = search.toString();
  return qs ? `?${qs}` : '';
}

export async function listGardens(params?: ListParams): Promise<Garden[]> {
  const res = await apiFetch<DataEnvelope<Garden[]>>(`/api/gardens${toQuery(params)}`);
  return res.data;
}

export async function getGarden(id: string): Promise<Garden> {
  const res = await apiFetch<DataEnvelope<Garden>>(`/api/gardens/${id}`);
  return res.data;
}

export async function createGarden(input: CreateGardenInput): Promise<Garden> {
  const res = await apiFetch<DataEnvelope<Garden>>('/api/gardens', {
    method: 'POST',
    body: JSON.stringify({
      name: input.name,
      notes: input.notes ?? ''
    })
  });
  return res.data;
}

export async function updateGarden(id: string, input: UpdateGardenInput): Promise<Garden> {
  const res = await apiFetch<DataEnvelope<Garden>>(`/api/gardens/${id}`, {
    method: 'PATCH',
    body: JSON.stringify(input)
  });
  return res.data;
}

export async function deleteGarden(id: string): Promise<void> {
  await apiFetch<DataEnvelope<{ status: string }>>(`/api/gardens/${id}`, {
    method: 'DELETE'
  });
}

export async function listBeds(gardenId: string, params?: ListParams): Promise<Bed[]> {
  const res = await apiFetch<DataEnvelope<Bed[]>>(
    `/api/gardens/${gardenId}/beds${toQuery(params)}`
  );
  return res.data;
}

export async function getBed(id: string): Promise<Bed> {
  const res = await apiFetch<DataEnvelope<Bed>>(`/api/beds/${id}`);
  return res.data;
}

export async function createBed(gardenId: string, input: CreateBedInput): Promise<Bed> {
  const res = await apiFetch<DataEnvelope<Bed>>(`/api/gardens/${gardenId}/beds`, {
    method: 'POST',
    body: JSON.stringify({
      name: input.name,
      notes: input.notes ?? '',
      type: input.type,
      shape: input.shape,
      sun_level: input.sun_level,
      watering_type: input.watering_type,
      width: input.width,
      length: input.length ?? null,
      height: input.height
    })
  });
  return res.data;
}

export async function updateBed(id: string, input: UpdateBedInput): Promise<Bed> {
  const res = await apiFetch<DataEnvelope<Bed>>(`/api/beds/${id}`, {
    method: 'PATCH',
    body: JSON.stringify(input)
  });
  return res.data;
}

export async function deleteBed(id: string): Promise<void> {
  await apiFetch<DataEnvelope<{ status: string }>>(`/api/beds/${id}`, {
    method: 'DELETE'
  });
}

export async function listPlanters(gardenId: string, params?: ListParams): Promise<Planter[]> {
  const res = await apiFetch<DataEnvelope<Planter[]>>(
    `/api/gardens/${gardenId}/planters${toQuery(params)}`
  );
  return res.data;
}

export async function getPlanter(id: string): Promise<Planter> {
  const res = await apiFetch<DataEnvelope<Planter>>(`/api/planters/${id}`);
  return res.data;
}

export async function createPlanter(gardenId: string, input: CreatePlanterInput): Promise<Planter> {
  const res = await apiFetch<DataEnvelope<Planter>>(`/api/gardens/${gardenId}/planters`, {
    method: 'POST',
    body: JSON.stringify({
      name: input.name,
      notes: input.notes ?? '',
      type: input.type,
      planter_type: input.planter_type,
      shape: input.shape,
      sun_level: input.sun_level,
      watering_type: input.watering_type,
      width: input.width,
      length: input.length ?? null,
      height: input.height
    })
  });
  return res.data;
}

export async function updatePlanter(id: string, input: UpdatePlanterInput): Promise<Planter> {
  const res = await apiFetch<DataEnvelope<Planter>>(`/api/planters/${id}`, {
    method: 'PATCH',
    body: JSON.stringify(input)
  });
  return res.data;
}

export async function deletePlanter(id: string): Promise<void> {
  await apiFetch<DataEnvelope<{ status: string }>>(`/api/planters/${id}`, {
    method: 'DELETE'
  });
}

export async function listPlants(params?: ListParams): Promise<Plant[]> {
  const res = await apiFetch<DataEnvelope<Plant[]>>(`/api/plants${toQuery(params)}`);
  return res.data;
}

export async function getPlant(id: string): Promise<Plant> {
  const res = await apiFetch<DataEnvelope<Plant>>(`/api/plants/${id}`);
  return res.data;
}

export async function createPlant(input: CreatePlantInput): Promise<Plant> {
  const res = await apiFetch<DataEnvelope<Plant>>('/api/plants', {
    method: 'POST',
    body: JSON.stringify({
      name: input.name,
      variety: input.variety ?? '',
      season: input.season ?? '',
      spacing: input.spacing ?? '',
      notes: input.notes ?? '',
      days_to_maturity: input.days_to_maturity ?? null,
      sow_depth: input.sow_depth ?? '',
      when_to_sow: input.when_to_sow ?? '',
      days_to_emerge: input.days_to_emerge ?? '',
      seed_spacing: input.seed_spacing ?? '',
      row_spacing: input.row_spacing ?? '',
      thinning: input.thinning ?? ''
    })
  });
  return res.data;
}

export async function listBedPlantings(bedId: string): Promise<BedPlanting[]> {
  const res = await apiFetch<DataEnvelope<BedPlanting[]>>(`/api/beds/${bedId}/plantings`);
  return res.data;
}

export async function createBedPlanting(
  bedId: string,
  input: CreateBedPlantingInput
): Promise<BedPlanting> {
  const res = await apiFetch<DataEnvelope<BedPlanting>>(`/api/beds/${bedId}/plantings`, {
    method: 'POST',
    body: JSON.stringify({
      plant_id: input.plant_id,
      start_date: input.start_date,
      end_date: input.end_date ?? '',
      notes: input.notes ?? ''
    })
  });
  return res.data;
}

export async function updateBedPlanting(
  plantingId: string,
  input: UpdateBedPlantingInput
): Promise<BedPlanting> {
  const res = await apiFetch<DataEnvelope<BedPlanting>>(`/api/plantings/${plantingId}`, {
    method: 'PATCH',
    body: JSON.stringify(input)
  });
  return res.data;
}

export async function deleteBedPlanting(plantingId: string): Promise<void> {
  await apiFetch<DataEnvelope<{ status: string }>>(`/api/plantings/${plantingId}`, {
    method: 'DELETE'
  });
}

export async function listJournalEntries(params?: ListParams): Promise<JournalEntry[]> {
  const res = await apiFetch<DataEnvelope<JournalEntry[]>>(`/api/journal${toQuery(params)}`);
  return res.data;
}

export async function createJournalEntry(input: CreateJournalEntryInput): Promise<JournalEntry> {
  const res = await apiFetch<DataEnvelope<JournalEntry>>('/api/journal', {
    method: 'POST',
    body: JSON.stringify({
      text: input.text,
      entry_date: input.entry_date ?? '',
      garden_id: input.garden_id
    })
  });
  return res.data;
}

export async function updateJournalEntry(
  id: string,
  input: UpdateJournalEntryInput
): Promise<JournalEntry> {
  const res = await apiFetch<DataEnvelope<JournalEntry>>(`/api/journal/${id}`, {
    method: 'PATCH',
    body: JSON.stringify(input)
  });
  return res.data;
}

export async function deleteJournalEntry(id: string): Promise<void> {
  await apiFetch<DataEnvelope<{ status: string }>>(`/api/journal/${id}`, {
    method: 'DELETE'
  });
}

export async function getCompanionData(plantName: string): Promise<CompanionRecord> {
  const res = await apiFetch<DataEnvelope<CompanionRecord>>(
    `/api/companions?plant=${encodeURIComponent(plantName)}`
  );
  return res.data;
}

export async function updatePlant(id: string, input: UpdatePlantInput): Promise<Plant> {
  const res = await apiFetch<DataEnvelope<Plant>>(`/api/plants/${id}`, {
    method: 'PATCH',
    body: JSON.stringify(input)
  });
  return res.data;
}

export type SeedPacketScan = {
  text: string;
  name: string;
  variety: string;
  season: string;
  spacing: string;
  notes: string;
  days_to_maturity?: number | null;
  sow_depth: string;
  when_to_sow: string;
  days_to_emerge: string;
  seed_spacing: string;
  row_spacing: string;
  thinning: string;
  boxes?: {
    name?: { x: number; y: number; w: number; h: number };
    specs?: { x: number; y: number; w: number; h: number };
    emerge?: { x: number; y: number; w: number; h: number };
    depth?: { x: number; y: number; w: number; h: number };
    seed?: { x: number; y: number; w: number; h: number };
    row?: { x: number; y: number; w: number; h: number };
    thin?: { x: number; y: number; w: number; h: number };
    mature?: { x: number; y: number; w: number; h: number };
    sow?: { x: number; y: number; w: number; h: number };
  };
};

export async function scanSeedPacket(file: File, company: string): Promise<SeedPacketScan> {
  const body = new FormData();
  body.append('image', file);
  if (company) {
    body.append('company', company);
  }
  const response = await fetch('/api/seed-packets/scan', {
    method: 'POST',
    body
  });
  if (!response.ok) {
    const message = await response.text();
    throw new Error(message || `Request failed: ${response.status}`);
  }
  const payload = (await response.json()) as DataEnvelope<SeedPacketScan>;
  return payload.data;
}

export async function deletePlant(id: string): Promise<void> {
  await apiFetch<DataEnvelope<{ status: string }>>(`/api/plants/${id}`, {
    method: 'DELETE'
  });
}
