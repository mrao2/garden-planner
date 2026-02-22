package httpapi

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handlers struct {
	db         *sql.DB
	companions *CompanionIndex
}

type errorEnvelope struct {
	Error errorPayload `json:"error"`
}

type errorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type dataEnvelope[T any] struct {
	Data T `json:"data"`
}

type healthData struct {
	Status string    `json:"status"`
	Time   time.Time `json:"time"`
}

func (h *Handlers) health(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, dataEnvelope[healthData]{
		Data: healthData{
			Status: "ok",
			Time:   time.Now().UTC(),
		},
	})
}

type garden struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type bed struct {
	ID        string    `json:"id"`
	GardenID  string    `json:"garden_id"`
	Name      string    `json:"name"`
	Notes     string    `json:"notes"`
	Type      string    `json:"type"`
	Shape     string    `json:"shape"`
	SunLevel  string    `json:"sun_level"`
	Watering  string    `json:"watering_type"`
	Width     float64   `json:"width"`
	Length    *float64  `json:"length,omitempty"`
	Height    float64   `json:"height"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type planter struct {
	ID          string    `json:"id"`
	GardenID    string    `json:"garden_id"`
	Name        string    `json:"name"`
	Notes       string    `json:"notes"`
	Type        string    `json:"type"`
	PlanterType string    `json:"planter_type"`
	Shape       string    `json:"shape"`
	SunLevel    string    `json:"sun_level"`
	Watering    string    `json:"watering_type"`
	Width       float64   `json:"width"`
	Length      *float64  `json:"length,omitempty"`
	Height      float64   `json:"height"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type plant struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Variety         string    `json:"variety"`
	Season          string    `json:"season"`
	Spacing         string    `json:"spacing"`
	Notes           string    `json:"notes"`
	DaysToMaturity  *int      `json:"days_to_maturity,omitempty"`
	SowDepth        string    `json:"sow_depth"`
	GerminationTime string    `json:"germination_time"`
	WhenToSow       string    `json:"when_to_sow"`
	DaysToEmerge    string    `json:"days_to_emerge"`
	SeedSpacing     string    `json:"seed_spacing"`
	RowSpacing      string    `json:"row_spacing"`
	Thinning        string    `json:"thinning"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type journalPlantRef struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Variety string `json:"variety"`
}

type journalEntry struct {
	ID        string            `json:"id"`
	GardenID  string            `json:"garden_id"`
	EntryDate time.Time         `json:"entry_date"`
	Text      string            `json:"text"`
	Plants    []journalPlantRef `json:"plants"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

type bedPlanting struct {
	ID        string     `json:"id"`
	BedID     string     `json:"bed_id"`
	PlantID   string     `json:"plant_id"`
	PlantName string     `json:"plant_name"`
	Variety   string     `json:"variety"`
	StartDate time.Time  `json:"start_date"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	Notes     string     `json:"notes"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type createGardenRequest struct {
	Name  string `json:"name"`
	Notes string `json:"notes"`
}

type updateGardenRequest struct {
	Name  *string `json:"name"`
	Notes *string `json:"notes"`
}

type createBedRequest struct {
	Name     string   `json:"name"`
	Notes    string   `json:"notes"`
	Type     string   `json:"type"`
	Shape    string   `json:"shape"`
	SunLevel string   `json:"sun_level"`
	Watering string   `json:"watering_type"`
	Width    float64  `json:"width"`
	Length   *float64 `json:"length"`
	Height   float64  `json:"height"`
}

type updateBedRequest struct {
	Name     *string  `json:"name"`
	Notes    *string  `json:"notes"`
	Type     *string  `json:"type"`
	Shape    *string  `json:"shape"`
	SunLevel *string  `json:"sun_level"`
	Watering *string  `json:"watering_type"`
	Width    *float64 `json:"width"`
	Length   *float64 `json:"length"`
	Height   *float64 `json:"height"`
}

type createPlantRequest struct {
	Name            string `json:"name"`
	Variety         string `json:"variety"`
	Season          string `json:"season"`
	Spacing         string `json:"spacing"`
	Notes           string `json:"notes"`
	DaysToMaturity  *int   `json:"days_to_maturity"`
	SowDepth        string `json:"sow_depth"`
	GerminationTime string `json:"germination_time"`
	WhenToSow       string `json:"when_to_sow"`
	DaysToEmerge    string `json:"days_to_emerge"`
	SeedSpacing     string `json:"seed_spacing"`
	RowSpacing      string `json:"row_spacing"`
	Thinning        string `json:"thinning"`
}

type updatePlantRequest struct {
	Name            *string `json:"name"`
	Variety         *string `json:"variety"`
	Season          *string `json:"season"`
	Spacing         *string `json:"spacing"`
	Notes           *string `json:"notes"`
	DaysToMaturity  *int    `json:"days_to_maturity"`
	SowDepth        *string `json:"sow_depth"`
	GerminationTime *string `json:"germination_time"`
	WhenToSow       *string `json:"when_to_sow"`
	DaysToEmerge    *string `json:"days_to_emerge"`
	SeedSpacing     *string `json:"seed_spacing"`
	RowSpacing      *string `json:"row_spacing"`
	Thinning        *string `json:"thinning"`
}

type createJournalEntryRequest struct {
	Text      string   `json:"text"`
	EntryDate string   `json:"entry_date"`
	PlantIDs  []string `json:"plant_ids"`
	GardenID  string   `json:"garden_id"`
}

type updateJournalEntryRequest struct {
	Text      *string `json:"text"`
	EntryDate *string `json:"entry_date"`
	GardenID  *string `json:"garden_id"`
}

type createBedPlantingRequest struct {
	PlantID   string `json:"plant_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Notes     string `json:"notes"`
}

type updateBedPlantingRequest struct {
	PlantID   *string `json:"plant_id"`
	StartDate *string `json:"start_date"`
	EndDate   *string `json:"end_date"`
	Notes     *string `json:"notes"`
}

type createPlanterRequest struct {
	Name        string   `json:"name"`
	Notes       string   `json:"notes"`
	Type        string   `json:"type"`
	PlanterType string   `json:"planter_type"`
	Shape       string   `json:"shape"`
	SunLevel    string   `json:"sun_level"`
	Watering    string   `json:"watering_type"`
	Width       float64  `json:"width"`
	Length      *float64 `json:"length"`
	Height      float64  `json:"height"`
}

type updatePlanterRequest struct {
	Name        *string  `json:"name"`
	Notes       *string  `json:"notes"`
	Type        *string  `json:"type"`
	PlanterType *string  `json:"planter_type"`
	Shape       *string  `json:"shape"`
	SunLevel    *string  `json:"sun_level"`
	Watering    *string  `json:"watering_type"`
	Width       *float64 `json:"width"`
	Length      *float64 `json:"length"`
	Height      *float64 `json:"height"`
}

func (h *Handlers) listGardens(w http.ResponseWriter, r *http.Request) {
	limit := clampQueryInt(r, "limit", 25, 1, 100)
	offset := clampQueryInt(r, "offset", 0, 0, 10000)
	queryText := strings.TrimSpace(r.URL.Query().Get("q"))

	var rows *sql.Rows
	var err error
	if queryText != "" {
		rows, err = h.db.Query(`
			SELECT id, name, notes, created_at, updated_at
			FROM gardens
			WHERE name ILIKE '%' || $1 || '%'
			ORDER BY created_at DESC
			LIMIT $2 OFFSET $3
		`, queryText, limit, offset)
	} else {
		rows, err = h.db.Query(`
			SELECT id, name, notes, created_at, updated_at
			FROM gardens
			ORDER BY created_at DESC
			LIMIT $1 OFFSET $2
		`, limit, offset)
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "db_query_failed", "Unable to load gardens")
		return
	}
	defer rows.Close()

	gardens := make([]garden, 0)
	for rows.Next() {
		var g garden
		if err := rows.Scan(&g.ID, &g.Name, &g.Notes, &g.CreatedAt, &g.UpdatedAt); err != nil {
			respondError(w, http.StatusInternalServerError, "db_scan_failed", "Unable to read gardens")
			return
		}
		gardens = append(gardens, g)
	}

	respondJSON(w, http.StatusOK, dataEnvelope[[]garden]{Data: gardens})
}

func (h *Handlers) createGarden(w http.ResponseWriter, r *http.Request) {
	var req createGardenRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid_json", "Invalid JSON payload")
		return
	}
	if strings.TrimSpace(req.Name) == "" {
		respondError(w, http.StatusBadRequest, "missing_name", "Garden name is required")
		return
	}
	if req.Notes == "" {
		req.Notes = ""
	}

	var g garden
	err := h.db.QueryRow(`
		INSERT INTO gardens (name, notes)
		VALUES ($1, $2)
		RETURNING id, name, notes, created_at, updated_at
	`, req.Name, req.Notes).Scan(&g.ID, &g.Name, &g.Notes, &g.CreatedAt, &g.UpdatedAt)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "db_insert_failed", "Unable to create garden")
		return
	}

	respondJSON(w, http.StatusCreated, dataEnvelope[garden]{Data: g})
}

func (h *Handlers) getGarden(w http.ResponseWriter, r *http.Request) {
	gardenID, err := parseUUIDParam(r, "gardenID")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid_id", "Garden ID is invalid")
		return
	}

	var g garden
	row := h.db.QueryRow(`
		SELECT id, name, notes, created_at, updated_at
		FROM gardens
		WHERE id = $1
	`, gardenID)
	if err := row.Scan(&g.ID, &g.Name, &g.Notes, &g.CreatedAt, &g.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondError(w, http.StatusNotFound, "not_found", "Garden not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "db_query_failed", "Unable to load garden")
		return
	}

	respondJSON(w, http.StatusOK, dataEnvelope[garden]{Data: g})
}

func (h *Handlers) updateGarden(w http.ResponseWriter, r *http.Request) {
	gardenID, err := parseUUIDParam(r, "gardenID")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid_id", "Garden ID is invalid")
		return
	}

	var req updateGardenRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid_json", "Invalid JSON payload")
		return
	}

	updates := make([]string, 0, 2)
	args := make([]any, 0, 3)
	argID := 1

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			respondError(w, http.StatusBadRequest, "missing_name", "Garden name cannot be empty")
			return
		}
		updates = append(updates, "name = $"+itoa(argID))
		args = append(args, name)
		argID++
	}
	if req.Notes != nil {
		updates = append(updates, "notes = $"+itoa(argID))
		args = append(args, *req.Notes)
		argID++
	}

	if len(updates) == 0 {
		respondError(w, http.StatusBadRequest, "no_updates", "No fields provided for update")
		return
	}

	updates = append(updates, "updated_at = NOW()")
	args = append(args, gardenID)

	query := "UPDATE gardens SET " + strings.Join(updates, ", ") + " WHERE id = $" + itoa(argID) + " RETURNING id, name, notes, created_at, updated_at"

	var g garden
	if err := h.db.QueryRow(query, args...).Scan(&g.ID, &g.Name, &g.Notes, &g.CreatedAt, &g.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondError(w, http.StatusNotFound, "not_found", "Garden not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "db_update_failed", "Unable to update garden")
		return
	}

	respondJSON(w, http.StatusOK, dataEnvelope[garden]{Data: g})
}

func (h *Handlers) deleteGarden(w http.ResponseWriter, r *http.Request) {
	gardenID, err := parseUUIDParam(r, "gardenID")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid_id", "Garden ID is invalid")
		return
	}

	res, err := h.db.Exec("DELETE FROM gardens WHERE id = $1", gardenID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "db_delete_failed", "Unable to delete garden")
		return
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		respondError(w, http.StatusNotFound, "not_found", "Garden not found")
		return
	}

	respondJSON(w, http.StatusOK, dataEnvelope[map[string]string]{Data: map[string]string{"status": "deleted"}})
}

func (h *Handlers) listBedsByGarden(w http.ResponseWriter, r *http.Request) {
	gardenID, err := parseUUIDParam(r, "gardenID")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid_id", "Garden ID is invalid")
		return
	}

	limit := clampQueryInt(r, "limit", 50, 1, 200)
	offset := clampQueryInt(r, "offset", 0, 0, 10000)
	queryText := strings.TrimSpace(r.URL.Query().Get("q"))

	var rows *sql.Rows
	if queryText != "" {
		rows, err = h.db.Query(`
			SELECT id, garden_id, name, notes, type, shape, sun_level, watering_type, width, length, height, created_at, updated_at
			FROM beds
			WHERE garden_id = $1 AND name ILIKE '%' || $2 || '%'
			ORDER BY created_at DESC
			LIMIT $3 OFFSET $4
		`, gardenID, queryText, limit, offset)
	} else {
		rows, err = h.db.Query(`
			SELECT id, garden_id, name, notes, type, shape, sun_level, watering_type, width, length, height, created_at, updated_at
			FROM beds
			WHERE garden_id = $1
			ORDER BY created_at DESC
			LIMIT $2 OFFSET $3
		`, gardenID, limit, offset)
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "db_query_failed", "Unable to load beds")
		return
	}
	defer rows.Close()

	beds := make([]bed, 0)
	for rows.Next() {
		var b bed
		var length sql.NullFloat64
		if err := rows.Scan(&b.ID, &b.GardenID, &b.Name, &b.Notes, &b.Type, &b.Shape, &b.SunLevel, &b.Watering, &b.Width, &length, &b.Height, &b.CreatedAt, &b.UpdatedAt); err != nil {
			respondError(w, http.StatusInternalServerError, "db_scan_failed", "Unable to read beds")
			return
		}
		if length.Valid {
			b.Length = &length.Float64
		}
		beds = append(beds, b)
	}

	respondJSON(w, http.StatusOK, dataEnvelope[[]bed]{Data: beds})
}

func (h *Handlers) createBed(w http.ResponseWriter, r *http.Request) {
	gardenID, err := parseUUIDParam(r, "gardenID")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid_id", "Garden ID is invalid")
		return
	}

	var req createBedRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid_json", "Invalid JSON payload")
		return
	}
	if strings.TrimSpace(req.Name) == "" {
		respondError(w, http.StatusBadRequest, "missing_name", "Bed name is required")
		return
	}
	shape := normalizeEnum(req.Shape)
	if shape != "rectangle" && shape != "circle" {
		respondError(w, http.StatusBadRequest, "invalid_shape", "Shape must be rectangle or circle")
		return
	}
	if req.Width <= 0 || req.Height <= 0 {
		respondError(w, http.StatusBadRequest, "invalid_size", "Bed width and height must be positive")
		return
	}
	if shape == "rectangle" {
		if req.Length == nil || *req.Length <= 0 {
			respondError(w, http.StatusBadRequest, "invalid_size", "Bed length must be positive for rectangles")
			return
		}
	} else if req.Length != nil {
		respondError(w, http.StatusBadRequest, "invalid_size", "Circular beds do not take a length")
		return
	}

	var b bed
	var length sql.NullFloat64
	if req.Length != nil {
		length = sql.NullFloat64{Float64: *req.Length, Valid: true}
	}
	row := h.db.QueryRow(`
		INSERT INTO beds (garden_id, name, notes, type, shape, sun_level, watering_type, width, length, height)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, garden_id, name, notes, type, shape, sun_level, watering_type, width, length, height, created_at, updated_at
	`, gardenID, req.Name, req.Notes, req.Type, shape, req.SunLevel, req.Watering, req.Width, length, req.Height)
	var lengthOut sql.NullFloat64
	if err := row.Scan(&b.ID, &b.GardenID, &b.Name, &b.Notes, &b.Type, &b.Shape, &b.SunLevel, &b.Watering, &b.Width, &lengthOut, &b.Height, &b.CreatedAt, &b.UpdatedAt); err != nil {
		respondError(w, http.StatusInternalServerError, "db_insert_failed", "Unable to create bed")
		return
	}
	if lengthOut.Valid {
		b.Length = &lengthOut.Float64
	}

	respondJSON(w, http.StatusCreated, dataEnvelope[bed]{Data: b})
}

func (h *Handlers) getBed(w http.ResponseWriter, r *http.Request) {
	bedID, err := parseUUIDParam(r, "bedID")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid_id", "Bed ID is invalid")
		return
	}

	var b bed
	row := h.db.QueryRow(`
		SELECT id, garden_id, name, notes, type, shape, sun_level, watering_type, width, length, height, created_at, updated_at
		FROM beds
		WHERE id = $1
	`, bedID)
	var length sql.NullFloat64
	if err := row.Scan(&b.ID, &b.GardenID, &b.Name, &b.Notes, &b.Type, &b.Shape, &b.SunLevel, &b.Watering, &b.Width, &length, &b.Height, &b.CreatedAt, &b.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondError(w, http.StatusNotFound, "not_found", "Bed not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "db_query_failed", "Unable to load bed")
		return
	}
	if length.Valid {
		b.Length = &length.Float64
	}

	respondJSON(w, http.StatusOK, dataEnvelope[bed]{Data: b})
}

func (h *Handlers) updateBed(w http.ResponseWriter, r *http.Request) {
	bedID, err := parseUUIDParam(r, "bedID")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid_id", "Bed ID is invalid")
		return
	}

	var req updateBedRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid_json", "Invalid JSON payload")
		return
	}

	updates := make([]string, 0, 8)
	args := make([]any, 0, 9)
	argID := 1

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			respondError(w, http.StatusBadRequest, "missing_name", "Bed name cannot be empty")
			return
		}
		updates = append(updates, "name = $"+itoa(argID))
		args = append(args, name)
		argID++
	}
	if req.Notes != nil {
		updates = append(updates, "notes = $"+itoa(argID))
		args = append(args, *req.Notes)
		argID++
	}
	if req.Type != nil {
		updates = append(updates, "type = $"+itoa(argID))
		args = append(args, strings.TrimSpace(*req.Type))
		argID++
	}
	if req.Shape != nil {
		shape := normalizeEnum(*req.Shape)
		if shape != "rectangle" && shape != "circle" {
			respondError(w, http.StatusBadRequest, "invalid_shape", "Shape must be rectangle or circle")
			return
		}
		updates = append(updates, "shape = $"+itoa(argID))
		args = append(args, shape)
		argID++
		if shape == "circle" {
			updates = append(updates, "length = NULL")
		}
		if shape == "rectangle" && req.Length == nil {
			respondError(w, http.StatusBadRequest, "invalid_size", "Length is required for rectangles")
			return
		}
	}
	if req.SunLevel != nil {
		updates = append(updates, "sun_level = $"+itoa(argID))
		args = append(args, strings.TrimSpace(*req.SunLevel))
		argID++
	}
	if req.Watering != nil {
		updates = append(updates, "watering_type = $"+itoa(argID))
		args = append(args, strings.TrimSpace(*req.Watering))
		argID++
	}
	if req.Width != nil {
		if *req.Width <= 0 {
			respondError(w, http.StatusBadRequest, "invalid_size", "Bed width must be positive")
			return
		}
		updates = append(updates, "width = $"+itoa(argID))
		args = append(args, *req.Width)
		argID++
	}
	if req.Length != nil {
		if *req.Length <= 0 {
			respondError(w, http.StatusBadRequest, "invalid_size", "Bed length must be positive")
			return
		}
		updates = append(updates, "length = $"+itoa(argID))
		args = append(args, *req.Length)
		argID++
	}
	if req.Height != nil {
		if *req.Height <= 0 {
			respondError(w, http.StatusBadRequest, "invalid_size", "Bed height must be positive")
			return
		}
		updates = append(updates, "height = $"+itoa(argID))
		args = append(args, *req.Height)
		argID++
	}

	if len(updates) == 0 {
		respondError(w, http.StatusBadRequest, "no_updates", "No fields provided for update")
		return
	}

	updates = append(updates, "updated_at = NOW()")
	args = append(args, bedID)

	query := "UPDATE beds SET " + strings.Join(updates, ", ") + " WHERE id = $" + itoa(argID) + " RETURNING id, garden_id, name, notes, type, shape, sun_level, watering_type, width, length, height, created_at, updated_at"

	var b bed
	var length sql.NullFloat64
	if err := h.db.QueryRow(query, args...).Scan(&b.ID, &b.GardenID, &b.Name, &b.Notes, &b.Type, &b.Shape, &b.SunLevel, &b.Watering, &b.Width, &length, &b.Height, &b.CreatedAt, &b.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondError(w, http.StatusNotFound, "not_found", "Bed not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "db_update_failed", "Unable to update bed")
		return
	}
	if length.Valid {
		b.Length = &length.Float64
	}

	respondJSON(w, http.StatusOK, dataEnvelope[bed]{Data: b})
}

func (h *Handlers) deleteBed(w http.ResponseWriter, r *http.Request) {
	bedID, err := parseUUIDParam(r, "bedID")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid_id", "Bed ID is invalid")
		return
	}

	res, err := h.db.Exec("DELETE FROM beds WHERE id = $1", bedID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "db_delete_failed", "Unable to delete bed")
		return
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		respondError(w, http.StatusNotFound, "not_found", "Bed not found")
		return
	}

	respondJSON(w, http.StatusOK, dataEnvelope[map[string]string]{Data: map[string]string{"status": "deleted"}})
}

func (h *Handlers) listPlantersByGarden(w http.ResponseWriter, r *http.Request) {
	gardenID, err := parseUUIDParam(r, "gardenID")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid_id", "Garden ID is invalid")
		return
	}

	limit := clampQueryInt(r, "limit", 50, 1, 200)
	offset := clampQueryInt(r, "offset", 0, 0, 10000)
	queryText := strings.TrimSpace(r.URL.Query().Get("q"))

	var rows *sql.Rows
	if queryText != "" {
		rows, err = h.db.Query(`
			SELECT id, garden_id, name, notes, type, planter_type, shape, sun_level, watering_type, width, length, height, created_at, updated_at
			FROM planters
			WHERE garden_id = $1 AND name ILIKE '%' || $2 || '%'
			ORDER BY created_at DESC
			LIMIT $3 OFFSET $4
		`, gardenID, queryText, limit, offset)
	} else {
		rows, err = h.db.Query(`
			SELECT id, garden_id, name, notes, type, planter_type, shape, sun_level, watering_type, width, length, height, created_at, updated_at
			FROM planters
			WHERE garden_id = $1
			ORDER BY created_at DESC
			LIMIT $2 OFFSET $3
		`, gardenID, limit, offset)
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "db_query_failed", "Unable to load planters")
		return
	}
	defer rows.Close()

	planters := make([]planter, 0)
	for rows.Next() {
		var p planter
		var length sql.NullFloat64
		if err := rows.Scan(&p.ID, &p.GardenID, &p.Name, &p.Notes, &p.Type, &p.PlanterType, &p.Shape, &p.SunLevel, &p.Watering, &p.Width, &length, &p.Height, &p.CreatedAt, &p.UpdatedAt); err != nil {
			respondError(w, http.StatusInternalServerError, "db_scan_failed", "Unable to read planters")
			return
		}
		if length.Valid {
			p.Length = &length.Float64
		}
		planters = append(planters, p)
	}

	respondJSON(w, http.StatusOK, dataEnvelope[[]planter]{Data: planters})
}

func (h *Handlers) createPlanter(w http.ResponseWriter, r *http.Request) {
	gardenID, err := parseUUIDParam(r, "gardenID")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid_id", "Garden ID is invalid")
		return
	}

	var req createPlanterRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid_json", "Invalid JSON payload")
		return
	}
	if strings.TrimSpace(req.Name) == "" {
		respondError(w, http.StatusBadRequest, "missing_name", "Planter name is required")
		return
	}
	shape := normalizeEnum(req.Shape)
	if shape != "rectangle" && shape != "circle" {
		respondError(w, http.StatusBadRequest, "invalid_shape", "Shape must be rectangle or circle")
		return
	}
	if req.Width <= 0 || req.Height <= 0 {
		respondError(w, http.StatusBadRequest, "invalid_size", "Planter width and height must be positive")
		return
	}
	if shape == "rectangle" {
		if req.Length == nil || *req.Length <= 0 {
			respondError(w, http.StatusBadRequest, "invalid_size", "Planter length must be positive for rectangles")
			return
		}
	} else if req.Length != nil {
		respondError(w, http.StatusBadRequest, "invalid_size", "Circular planters do not take a length")
		return
	}

	var p planter
	var length sql.NullFloat64
	if req.Length != nil {
		length = sql.NullFloat64{Float64: *req.Length, Valid: true}
	}
	row := h.db.QueryRow(`
		INSERT INTO planters (garden_id, name, notes, type, planter_type, shape, sun_level, watering_type, width, length, height)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, garden_id, name, notes, type, planter_type, shape, sun_level, watering_type, width, length, height, created_at, updated_at
	`, gardenID, req.Name, req.Notes, req.Type, req.PlanterType, shape, req.SunLevel, req.Watering, req.Width, length, req.Height)
	var lengthOut sql.NullFloat64
	if err := row.Scan(&p.ID, &p.GardenID, &p.Name, &p.Notes, &p.Type, &p.PlanterType, &p.Shape, &p.SunLevel, &p.Watering, &p.Width, &lengthOut, &p.Height, &p.CreatedAt, &p.UpdatedAt); err != nil {
		respondError(w, http.StatusInternalServerError, "db_insert_failed", "Unable to create planter")
		return
	}
	if lengthOut.Valid {
		p.Length = &lengthOut.Float64
	}

	respondJSON(w, http.StatusCreated, dataEnvelope[planter]{Data: p})
}

func (h *Handlers) getPlanter(w http.ResponseWriter, r *http.Request) {
	planterID, err := parseUUIDParam(r, "planterID")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid_id", "Planter ID is invalid")
		return
	}

	var p planter
	row := h.db.QueryRow(`
		SELECT id, garden_id, name, notes, type, planter_type, shape, sun_level, watering_type, width, length, height, created_at, updated_at
		FROM planters
		WHERE id = $1
	`, planterID)
	var length sql.NullFloat64
	if err := row.Scan(&p.ID, &p.GardenID, &p.Name, &p.Notes, &p.Type, &p.PlanterType, &p.Shape, &p.SunLevel, &p.Watering, &p.Width, &length, &p.Height, &p.CreatedAt, &p.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondError(w, http.StatusNotFound, "not_found", "Planter not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "db_query_failed", "Unable to load planter")
		return
	}
	if length.Valid {
		p.Length = &length.Float64
	}

	respondJSON(w, http.StatusOK, dataEnvelope[planter]{Data: p})
}

func (h *Handlers) updatePlanter(w http.ResponseWriter, r *http.Request) {
	planterID, err := parseUUIDParam(r, "planterID")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid_id", "Planter ID is invalid")
		return
	}

	var req updatePlanterRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid_json", "Invalid JSON payload")
		return
	}

	updates := make([]string, 0, 9)
	args := make([]any, 0, 10)
	argID := 1

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			respondError(w, http.StatusBadRequest, "missing_name", "Planter name cannot be empty")
			return
		}
		updates = append(updates, "name = $"+itoa(argID))
		args = append(args, name)
		argID++
	}
	if req.Notes != nil {
		updates = append(updates, "notes = $"+itoa(argID))
		args = append(args, *req.Notes)
		argID++
	}
	if req.Type != nil {
		updates = append(updates, "type = $"+itoa(argID))
		args = append(args, strings.TrimSpace(*req.Type))
		argID++
	}
	if req.PlanterType != nil {
		updates = append(updates, "planter_type = $"+itoa(argID))
		args = append(args, strings.TrimSpace(*req.PlanterType))
		argID++
	}
	if req.Shape != nil {
		shape := normalizeEnum(*req.Shape)
		if shape != "rectangle" && shape != "circle" {
			respondError(w, http.StatusBadRequest, "invalid_shape", "Shape must be rectangle or circle")
			return
		}
		updates = append(updates, "shape = $"+itoa(argID))
		args = append(args, shape)
		argID++
		if shape == "circle" {
			updates = append(updates, "length = NULL")
		}
		if shape == "rectangle" && req.Length == nil {
			respondError(w, http.StatusBadRequest, "invalid_size", "Length is required for rectangles")
			return
		}
	}
	if req.SunLevel != nil {
		updates = append(updates, "sun_level = $"+itoa(argID))
		args = append(args, strings.TrimSpace(*req.SunLevel))
		argID++
	}
	if req.Watering != nil {
		updates = append(updates, "watering_type = $"+itoa(argID))
		args = append(args, strings.TrimSpace(*req.Watering))
		argID++
	}
	if req.Width != nil {
		if *req.Width <= 0 {
			respondError(w, http.StatusBadRequest, "invalid_size", "Planter width must be positive")
			return
		}
		updates = append(updates, "width = $"+itoa(argID))
		args = append(args, *req.Width)
		argID++
	}
	if req.Length != nil {
		if *req.Length <= 0 {
			respondError(w, http.StatusBadRequest, "invalid_size", "Planter length must be positive")
			return
		}
		updates = append(updates, "length = $"+itoa(argID))
		args = append(args, *req.Length)
		argID++
	}
	if req.Height != nil {
		if *req.Height <= 0 {
			respondError(w, http.StatusBadRequest, "invalid_size", "Planter height must be positive")
			return
		}
		updates = append(updates, "height = $"+itoa(argID))
		args = append(args, *req.Height)
		argID++
	}

	if len(updates) == 0 {
		respondError(w, http.StatusBadRequest, "no_updates", "No fields provided for update")
		return
	}

	updates = append(updates, "updated_at = NOW()")
	args = append(args, planterID)

	query := "UPDATE planters SET " + strings.Join(updates, ", ") + " WHERE id = $" + itoa(argID) +
		" RETURNING id, garden_id, name, notes, type, planter_type, shape, sun_level, watering_type, width, length, height, created_at, updated_at"

	var p planter
	var length sql.NullFloat64
	if err := h.db.QueryRow(query, args...).Scan(&p.ID, &p.GardenID, &p.Name, &p.Notes, &p.Type, &p.PlanterType, &p.Shape, &p.SunLevel, &p.Watering, &p.Width, &length, &p.Height, &p.CreatedAt, &p.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondError(w, http.StatusNotFound, "not_found", "Planter not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "db_update_failed", "Unable to update planter")
		return
	}
	if length.Valid {
		p.Length = &length.Float64
	}

	respondJSON(w, http.StatusOK, dataEnvelope[planter]{Data: p})
}

func (h *Handlers) deletePlanter(w http.ResponseWriter, r *http.Request) {
	planterID, err := parseUUIDParam(r, "planterID")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid_id", "Planter ID is invalid")
		return
	}

	res, err := h.db.Exec("DELETE FROM planters WHERE id = $1", planterID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "db_delete_failed", "Unable to delete planter")
		return
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		respondError(w, http.StatusNotFound, "not_found", "Planter not found")
		return
	}

	respondJSON(w, http.StatusOK, dataEnvelope[map[string]string]{Data: map[string]string{"status": "deleted"}})
}
func (h *Handlers) listPlants(w http.ResponseWriter, r *http.Request) {
	limit := clampQueryInt(r, "limit", 50, 1, 200)
	offset := clampQueryInt(r, "offset", 0, 0, 10000)
	queryText := strings.TrimSpace(r.URL.Query().Get("q"))

	var rows *sql.Rows
	var err error
	if queryText != "" {
		rows, err = h.db.Query(`
			SELECT id, name, variety, season, spacing, notes, days_to_maturity, sow_depth, germination_time, when_to_sow, days_to_emerge, seed_spacing, row_spacing, thinning, created_at, updated_at
			FROM plants
			WHERE name ILIKE '%' || $1 || '%' OR variety ILIKE '%' || $1 || '%'
			ORDER BY created_at DESC
			LIMIT $2 OFFSET $3
		`, queryText, limit, offset)
	} else {
		rows, err = h.db.Query(`
			SELECT id, name, variety, season, spacing, notes, days_to_maturity, sow_depth, germination_time, when_to_sow, days_to_emerge, seed_spacing, row_spacing, thinning, created_at, updated_at
			FROM plants
			ORDER BY created_at DESC
			LIMIT $1 OFFSET $2
		`, limit, offset)
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "db_query_failed", "Unable to load plants")
		return
	}
	defer rows.Close()

	plants := make([]plant, 0)
	for rows.Next() {
		var p plant
		var days sql.NullInt64
		if err := rows.Scan(&p.ID, &p.Name, &p.Variety, &p.Season, &p.Spacing, &p.Notes, &days, &p.SowDepth, &p.GerminationTime, &p.WhenToSow, &p.DaysToEmerge, &p.SeedSpacing, &p.RowSpacing, &p.Thinning, &p.CreatedAt, &p.UpdatedAt); err != nil {
			respondError(w, http.StatusInternalServerError, "db_scan_failed", "Unable to read plants")
			return
		}
		if days.Valid {
			value := int(days.Int64)
			p.DaysToMaturity = &value
		}
		plants = append(plants, p)
	}

	respondJSON(w, http.StatusOK, dataEnvelope[[]plant]{Data: plants})
}

func (h *Handlers) createPlant(w http.ResponseWriter, r *http.Request) {
	var req createPlantRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid_json", "Invalid JSON payload")
		return
	}
	if strings.TrimSpace(req.Name) == "" {
		respondError(w, http.StatusBadRequest, "missing_name", "Plant name is required")
		return
	}

	var p plant
	var days sql.NullInt64
	row := h.db.QueryRow(`
		INSERT INTO plants (name, variety, season, spacing, notes, days_to_maturity, sow_depth, germination_time, when_to_sow, days_to_emerge, seed_spacing, row_spacing, thinning)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, name, variety, season, spacing, notes, days_to_maturity, sow_depth, germination_time, when_to_sow, days_to_emerge, seed_spacing, row_spacing, thinning, created_at, updated_at
	`, req.Name, req.Variety, req.Season, req.Spacing, req.Notes, req.DaysToMaturity, req.SowDepth, req.GerminationTime, req.WhenToSow, req.DaysToEmerge, req.SeedSpacing, req.RowSpacing, req.Thinning)
	if err := row.Scan(&p.ID, &p.Name, &p.Variety, &p.Season, &p.Spacing, &p.Notes, &days, &p.SowDepth, &p.GerminationTime, &p.WhenToSow, &p.DaysToEmerge, &p.SeedSpacing, &p.RowSpacing, &p.Thinning, &p.CreatedAt, &p.UpdatedAt); err != nil {
		respondError(w, http.StatusInternalServerError, "db_insert_failed", "Unable to create plant")
		return
	}
	if days.Valid {
		value := int(days.Int64)
		p.DaysToMaturity = &value
	}

	respondJSON(w, http.StatusCreated, dataEnvelope[plant]{Data: p})
}

func (h *Handlers) getPlant(w http.ResponseWriter, r *http.Request) {
	plantID, err := parseUUIDParam(r, "plantID")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid_id", "Plant ID is invalid")
		return
	}

	var p plant
	var days sql.NullInt64
	row := h.db.QueryRow(`
		SELECT id, name, variety, season, spacing, notes, days_to_maturity, sow_depth, germination_time, when_to_sow, days_to_emerge, seed_spacing, row_spacing, thinning, created_at, updated_at
		FROM plants
		WHERE id = $1
	`, plantID)
	if err := row.Scan(&p.ID, &p.Name, &p.Variety, &p.Season, &p.Spacing, &p.Notes, &days, &p.SowDepth, &p.GerminationTime, &p.WhenToSow, &p.DaysToEmerge, &p.SeedSpacing, &p.RowSpacing, &p.Thinning, &p.CreatedAt, &p.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondError(w, http.StatusNotFound, "not_found", "Plant not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "db_query_failed", "Unable to load plant")
		return
	}
	if days.Valid {
		value := int(days.Int64)
		p.DaysToMaturity = &value
	}

	respondJSON(w, http.StatusOK, dataEnvelope[plant]{Data: p})
}

func (h *Handlers) updatePlant(w http.ResponseWriter, r *http.Request) {
	plantID, err := parseUUIDParam(r, "plantID")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid_id", "Plant ID is invalid")
		return
	}

	var req updatePlantRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid_json", "Invalid JSON payload")
		return
	}

	updates := make([]string, 0, 13)
	args := make([]any, 0, 14)
	argID := 1

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			respondError(w, http.StatusBadRequest, "missing_name", "Plant name cannot be empty")
			return
		}
		updates = append(updates, "name = $"+itoa(argID))
		args = append(args, name)
		argID++
	}
	if req.Variety != nil {
		updates = append(updates, "variety = $"+itoa(argID))
		args = append(args, *req.Variety)
		argID++
	}
	if req.Season != nil {
		updates = append(updates, "season = $"+itoa(argID))
		args = append(args, *req.Season)
		argID++
	}
	if req.Spacing != nil {
		updates = append(updates, "spacing = $"+itoa(argID))
		args = append(args, *req.Spacing)
		argID++
	}
	if req.Notes != nil {
		updates = append(updates, "notes = $"+itoa(argID))
		args = append(args, *req.Notes)
		argID++
	}
	if req.DaysToMaturity != nil {
		updates = append(updates, "days_to_maturity = $"+itoa(argID))
		args = append(args, *req.DaysToMaturity)
		argID++
	}
	if req.SowDepth != nil {
		updates = append(updates, "sow_depth = $"+itoa(argID))
		args = append(args, *req.SowDepth)
		argID++
	}
	if req.GerminationTime != nil {
		updates = append(updates, "germination_time = $"+itoa(argID))
		args = append(args, *req.GerminationTime)
		argID++
	}
	if req.WhenToSow != nil {
		updates = append(updates, "when_to_sow = $"+itoa(argID))
		args = append(args, *req.WhenToSow)
		argID++
	}
	if req.DaysToEmerge != nil {
		updates = append(updates, "days_to_emerge = $"+itoa(argID))
		args = append(args, *req.DaysToEmerge)
		argID++
	}
	if req.SeedSpacing != nil {
		updates = append(updates, "seed_spacing = $"+itoa(argID))
		args = append(args, *req.SeedSpacing)
		argID++
	}
	if req.RowSpacing != nil {
		updates = append(updates, "row_spacing = $"+itoa(argID))
		args = append(args, *req.RowSpacing)
		argID++
	}
	if req.Thinning != nil {
		updates = append(updates, "thinning = $"+itoa(argID))
		args = append(args, *req.Thinning)
		argID++
	}

	if len(updates) == 0 {
		respondError(w, http.StatusBadRequest, "no_updates", "No fields provided for update")
		return
	}

	updates = append(updates, "updated_at = NOW()")
	args = append(args, plantID)

	query := "UPDATE plants SET " + strings.Join(updates, ", ") + " WHERE id = $" + itoa(argID) +
		" RETURNING id, name, variety, season, spacing, notes, days_to_maturity, sow_depth, germination_time, when_to_sow, days_to_emerge, seed_spacing, row_spacing, thinning, created_at, updated_at"

	var p plant
	var days sql.NullInt64
	if err := h.db.QueryRow(query, args...).Scan(&p.ID, &p.Name, &p.Variety, &p.Season, &p.Spacing, &p.Notes, &days, &p.SowDepth, &p.GerminationTime, &p.WhenToSow, &p.DaysToEmerge, &p.SeedSpacing, &p.RowSpacing, &p.Thinning, &p.CreatedAt, &p.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondError(w, http.StatusNotFound, "not_found", "Plant not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "db_update_failed", "Unable to update plant")
		return
	}
	if days.Valid {
		value := int(days.Int64)
		p.DaysToMaturity = &value
	}

	respondJSON(w, http.StatusOK, dataEnvelope[plant]{Data: p})
}

func (h *Handlers) deletePlant(w http.ResponseWriter, r *http.Request) {
	plantID, err := parseUUIDParam(r, "plantID")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid_id", "Plant ID is invalid")
		return
	}

	res, err := h.db.Exec("DELETE FROM plants WHERE id = $1", plantID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "db_delete_failed", "Unable to delete plant")
		return
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		respondError(w, http.StatusNotFound, "not_found", "Plant not found")
		return
	}

	respondJSON(w, http.StatusOK, dataEnvelope[map[string]string]{Data: map[string]string{"status": "deleted"}})
}

func (h *Handlers) listJournalEntries(w http.ResponseWriter, r *http.Request) {
	limit := clampQueryInt(r, "limit", 50, 1, 200)
	offset := clampQueryInt(r, "offset", 0, 0, 10000)

	rows, err := h.db.Query(`
		SELECT id, garden_id, entry_date, text, created_at, updated_at
		FROM journal_entries
		ORDER BY entry_date DESC, created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "db_query_failed", "Unable to load journal entries")
		return
	}
	defer rows.Close()

	entries := make([]journalEntry, 0)
	entryIDs := make([]string, 0)
	for rows.Next() {
		var entry journalEntry
		if err := rows.Scan(&entry.ID, &entry.GardenID, &entry.EntryDate, &entry.Text, &entry.CreatedAt, &entry.UpdatedAt); err != nil {
			respondError(w, http.StatusInternalServerError, "db_scan_failed", "Unable to read journal entries")
			return
		}
		entries = append(entries, entry)
		entryIDs = append(entryIDs, entry.ID)
	}

	plantsByEntry := make(map[string][]journalPlantRef)
	if len(entryIDs) > 0 {
		query, args := buildInQuery("SELECT jep.entry_id, p.id, p.name, p.variety FROM journal_entry_plants jep JOIN plants p ON p.id = jep.plant_id WHERE jep.entry_id IN ", entryIDs)
		rows, err := h.db.Query(query, args...)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "db_query_failed", "Unable to load journal entry plants")
			return
		}
		defer rows.Close()
		for rows.Next() {
			var entryID string
			var plant journalPlantRef
			if err := rows.Scan(&entryID, &plant.ID, &plant.Name, &plant.Variety); err != nil {
				respondError(w, http.StatusInternalServerError, "db_scan_failed", "Unable to read journal entry plants")
				return
			}
			plantsByEntry[entryID] = append(plantsByEntry[entryID], plant)
		}
	}

	for i := range entries {
		entries[i].Plants = plantsByEntry[entries[i].ID]
	}

	respondJSON(w, http.StatusOK, dataEnvelope[[]journalEntry]{Data: entries})
}

func (h *Handlers) createJournalEntry(w http.ResponseWriter, r *http.Request) {
	var req createJournalEntryRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid_json", "Invalid JSON payload")
		return
	}
	if strings.TrimSpace(req.Text) == "" {
		respondError(w, http.StatusBadRequest, "missing_text", "Journal text is required")
		return
	}
	if _, err := uuid.Parse(req.GardenID); err != nil {
		respondError(w, http.StatusBadRequest, "invalid_id", "Garden ID is invalid")
		return
	}
	if !h.idExists("gardens", req.GardenID) {
		respondError(w, http.StatusNotFound, "not_found", "Garden not found")
		return
	}

	entryDate := time.Now().UTC()
	if strings.TrimSpace(req.EntryDate) != "" {
		parsed, err := parseDate(req.EntryDate)
		if err != nil {
			respondError(w, http.StatusBadRequest, "invalid_date", "Entry date must be YYYY-MM-DD")
			return
		}
		entryDate = parsed
	}

	plants, err := h.loadPlants()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "db_query_failed", "Unable to load plants")
		return
	}

	explicitIDs := uniqueValidPlantIDs(req.PlantIDs)
	for _, id := range explicitIDs {
		if !h.idExists("plants", id) {
			respondError(w, http.StatusNotFound, "not_found", "Plant not found")
			return
		}
	}

	matchedIDs := matchPlantIDs(req.Text, plants)
	plantIDs := unionIDs(explicitIDs, matchedIDs)

	tx, err := h.db.Begin()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "db_tx_failed", "Unable to create journal entry")
		return
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	var entry journalEntry
	err = tx.QueryRow(`
		INSERT INTO journal_entries (garden_id, entry_date, text)
		VALUES ($1, $2, $3)
		RETURNING id, garden_id, entry_date, text, created_at, updated_at
	`, req.GardenID, entryDate, strings.TrimSpace(req.Text)).Scan(&entry.ID, &entry.GardenID, &entry.EntryDate, &entry.Text, &entry.CreatedAt, &entry.UpdatedAt)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "db_insert_failed", "Unable to create journal entry")
		return
	}

	if len(plantIDs) > 0 {
		for _, plantID := range plantIDs {
			if _, err := tx.Exec(`INSERT INTO journal_entry_plants (entry_id, plant_id) VALUES ($1, $2)`, entry.ID, plantID); err != nil {
				respondError(w, http.StatusInternalServerError, "db_insert_failed", "Unable to link journal entry to plants")
				return
			}
		}
	}

	if err := tx.Commit(); err != nil {
		respondError(w, http.StatusInternalServerError, "db_tx_failed", "Unable to create journal entry")
		return
	}

	entry.Plants, err = h.loadPlantRefs(plantIDs)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "db_query_failed", "Unable to load linked plants")
		return
	}

	respondJSON(w, http.StatusCreated, dataEnvelope[journalEntry]{Data: entry})
}

func (h *Handlers) updateJournalEntry(w http.ResponseWriter, r *http.Request) {
	entryID, err := parseUUIDParam(r, "entryID")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid_id", "Entry ID is invalid")
		return
	}

	var req updateJournalEntryRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid_json", "Invalid JSON payload")
		return
	}

	updates := make([]string, 0, 3)
	args := make([]any, 0, 4)
	argID := 1

	if req.Text != nil {
		text := strings.TrimSpace(*req.Text)
		if text == "" {
			respondError(w, http.StatusBadRequest, "missing_text", "Journal text is required")
			return
		}
		updates = append(updates, "text = $"+itoa(argID))
		args = append(args, text)
		argID++
	}
	if req.EntryDate != nil {
		parsed, err := parseDate(*req.EntryDate)
		if err != nil {
			respondError(w, http.StatusBadRequest, "invalid_date", "Entry date must be YYYY-MM-DD")
			return
		}
		updates = append(updates, "entry_date = $"+itoa(argID))
		args = append(args, parsed)
		argID++
	}
	if req.GardenID != nil {
		if _, err := uuid.Parse(*req.GardenID); err != nil {
			respondError(w, http.StatusBadRequest, "invalid_id", "Garden ID is invalid")
			return
		}
		if !h.idExists("gardens", *req.GardenID) {
			respondError(w, http.StatusNotFound, "not_found", "Garden not found")
			return
		}
		updates = append(updates, "garden_id = $"+itoa(argID))
		args = append(args, *req.GardenID)
		argID++
	}

	if len(updates) == 0 {
		respondError(w, http.StatusBadRequest, "no_updates", "No fields provided for update")
		return
	}

	updates = append(updates, "updated_at = NOW()")
	args = append(args, entryID)

	tx, err := h.db.Begin()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "db_tx_failed", "Unable to update journal entry")
		return
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	query := "UPDATE journal_entries SET " + strings.Join(updates, ", ") +
		" WHERE id = $" + itoa(argID) +
		" RETURNING id, garden_id, entry_date, text, created_at, updated_at"

	var entry journalEntry
	if err := tx.QueryRow(query, args...).Scan(&entry.ID, &entry.GardenID, &entry.EntryDate, &entry.Text, &entry.CreatedAt, &entry.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondError(w, http.StatusNotFound, "not_found", "Journal entry not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "db_update_failed", "Unable to update journal entry")
		return
	}

	if req.Text != nil {
		plants, err := h.loadPlants()
		if err != nil {
			respondError(w, http.StatusInternalServerError, "db_query_failed", "Unable to load plants")
			return
		}
		matchedIDs := matchPlantIDs(entry.Text, plants)

		if _, err := tx.Exec(`DELETE FROM journal_entry_plants WHERE entry_id = $1`, entry.ID); err != nil {
			respondError(w, http.StatusInternalServerError, "db_delete_failed", "Unable to update entry plants")
			return
		}
		for _, plantID := range matchedIDs {
			if _, err := tx.Exec(`INSERT INTO journal_entry_plants (entry_id, plant_id) VALUES ($1, $2)`, entry.ID, plantID); err != nil {
				respondError(w, http.StatusInternalServerError, "db_insert_failed", "Unable to update entry plants")
				return
			}
		}
		entry.Plants, err = h.loadPlantRefs(matchedIDs)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "db_query_failed", "Unable to load linked plants")
			return
		}
	} else {
		entry.Plants, err = h.loadEntryPlants(entry.ID)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "db_query_failed", "Unable to load linked plants")
			return
		}
	}

	if err := tx.Commit(); err != nil {
		respondError(w, http.StatusInternalServerError, "db_tx_failed", "Unable to update journal entry")
		return
	}

	respondJSON(w, http.StatusOK, dataEnvelope[journalEntry]{Data: entry})
}

func (h *Handlers) deleteJournalEntry(w http.ResponseWriter, r *http.Request) {
	entryID, err := parseUUIDParam(r, "entryID")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid_id", "Entry ID is invalid")
		return
	}

	res, err := h.db.Exec("DELETE FROM journal_entries WHERE id = $1", entryID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "db_delete_failed", "Unable to delete journal entry")
		return
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		respondError(w, http.StatusNotFound, "not_found", "Journal entry not found")
		return
	}

	respondJSON(w, http.StatusOK, dataEnvelope[map[string]string]{Data: map[string]string{"status": "deleted"}})
}

func (h *Handlers) getCompanionData(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.URL.Query().Get("plant"))
	if query == "" {
		respondError(w, http.StatusBadRequest, "missing_query", "Plant name is required")
		return
	}
	if h.companions == nil {
		respondError(w, http.StatusServiceUnavailable, "not_ready", "Companion dataset not available")
		return
	}
	record, ok := h.companions.FindByName(query)
	if !ok {
		respondError(w, http.StatusNotFound, "not_found", "No companion data found")
		return
	}

	respondJSON(w, http.StatusOK, dataEnvelope[CompanionRecord]{Data: record})
}

func (h *Handlers) listBedPlantings(w http.ResponseWriter, r *http.Request) {
	bedID, err := parseUUIDParam(r, "bedID")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid_id", "Bed ID is invalid")
		return
	}

	rows, err := h.db.Query(`
		SELECT bp.id, bp.bed_id, bp.plant_id, p.name, p.variety, bp.start_date, bp.end_date,
		       bp.notes, bp.created_at, bp.updated_at
		FROM bed_plantings bp
		JOIN plants p ON p.id = bp.plant_id
		WHERE bp.bed_id = $1
		ORDER BY bp.start_date DESC, bp.created_at DESC
	`, bedID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "db_query_failed", "Unable to load plantings")
		return
	}
	defer rows.Close()

	plantings := make([]bedPlanting, 0)
	for rows.Next() {
		var p bedPlanting
		var endDate sql.NullTime
		if err := rows.Scan(&p.ID, &p.BedID, &p.PlantID, &p.PlantName, &p.Variety, &p.StartDate, &endDate, &p.Notes, &p.CreatedAt, &p.UpdatedAt); err != nil {
			respondError(w, http.StatusInternalServerError, "db_scan_failed", "Unable to read plantings")
			return
		}
		if endDate.Valid {
			p.EndDate = &endDate.Time
		}
		plantings = append(plantings, p)
	}

	respondJSON(w, http.StatusOK, dataEnvelope[[]bedPlanting]{Data: plantings})
}

func (h *Handlers) createBedPlanting(w http.ResponseWriter, r *http.Request) {
	bedID, err := parseUUIDParam(r, "bedID")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid_id", "Bed ID is invalid")
		return
	}

	var req createBedPlantingRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid_json", "Invalid JSON payload")
		return
	}
	if _, err := uuid.Parse(req.PlantID); err != nil {
		respondError(w, http.StatusBadRequest, "invalid_id", "Plant ID is invalid")
		return
	}
	startDate, err := parseDate(req.StartDate)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid_date", "Start date must be YYYY-MM-DD")
		return
	}

	var endDate sql.NullTime
	if strings.TrimSpace(req.EndDate) != "" {
		parsed, err := parseDate(req.EndDate)
		if err != nil {
			respondError(w, http.StatusBadRequest, "invalid_date", "End date must be YYYY-MM-DD")
			return
		}
		endDate = sql.NullTime{Time: parsed, Valid: true}
	}

	if !h.idExists("beds", bedID) {
		respondError(w, http.StatusNotFound, "not_found", "Bed not found")
		return
	}
	if !h.idExists("plants", req.PlantID) {
		respondError(w, http.StatusNotFound, "not_found", "Plant not found")
		return
	}

	var planting bedPlanting
	var endDateOut sql.NullTime
	row := h.db.QueryRow(`
		INSERT INTO bed_plantings (bed_id, plant_id, start_date, end_date, notes)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, bed_id, plant_id, start_date, end_date, notes, created_at, updated_at
	`, bedID, req.PlantID, startDate, endDate, req.Notes)
	if err := row.Scan(&planting.ID, &planting.BedID, &planting.PlantID, &planting.StartDate, &endDateOut, &planting.Notes, &planting.CreatedAt, &planting.UpdatedAt); err != nil {
		respondError(w, http.StatusInternalServerError, "db_insert_failed", "Unable to create planting")
		return
	}
	if endDateOut.Valid {
		planting.EndDate = &endDateOut.Time
	}

	err = h.db.QueryRow(`SELECT name, variety FROM plants WHERE id = $1`, planting.PlantID).Scan(&planting.PlantName, &planting.Variety)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "db_query_failed", "Unable to load plant")
		return
	}

	respondJSON(w, http.StatusCreated, dataEnvelope[bedPlanting]{Data: planting})
}

func (h *Handlers) updateBedPlanting(w http.ResponseWriter, r *http.Request) {
	plantingID, err := parseUUIDParam(r, "plantingID")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid_id", "Planting ID is invalid")
		return
	}

	var req updateBedPlantingRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid_json", "Invalid JSON payload")
		return
	}

	updates := make([]string, 0, 4)
	args := make([]any, 0, 5)
	argID := 1

	if req.PlantID != nil {
		if _, err := uuid.Parse(*req.PlantID); err != nil {
			respondError(w, http.StatusBadRequest, "invalid_id", "Plant ID is invalid")
			return
		}
		if !h.idExists("plants", *req.PlantID) {
			respondError(w, http.StatusNotFound, "not_found", "Plant not found")
			return
		}
		updates = append(updates, "plant_id = $"+itoa(argID))
		args = append(args, *req.PlantID)
		argID++
	}
	if req.StartDate != nil {
		startDate, err := parseDate(*req.StartDate)
		if err != nil {
			respondError(w, http.StatusBadRequest, "invalid_date", "Start date must be YYYY-MM-DD")
			return
		}
		updates = append(updates, "start_date = $"+itoa(argID))
		args = append(args, startDate)
		argID++
	}
	if req.EndDate != nil {
		if strings.TrimSpace(*req.EndDate) == "" {
			updates = append(updates, "end_date = NULL")
		} else {
			endDate, err := parseDate(*req.EndDate)
			if err != nil {
				respondError(w, http.StatusBadRequest, "invalid_date", "End date must be YYYY-MM-DD")
				return
			}
			updates = append(updates, "end_date = $"+itoa(argID))
			args = append(args, endDate)
			argID++
		}
	}
	if req.Notes != nil {
		updates = append(updates, "notes = $"+itoa(argID))
		args = append(args, *req.Notes)
		argID++
	}

	if len(updates) == 0 {
		respondError(w, http.StatusBadRequest, "no_updates", "No fields provided for update")
		return
	}

	updates = append(updates, "updated_at = NOW()")
	args = append(args, plantingID)

	query := "UPDATE bed_plantings SET " + strings.Join(updates, ", ") + " WHERE id = $" + itoa(argID) +
		" RETURNING id, bed_id, plant_id, start_date, end_date, notes, created_at, updated_at"

	var planting bedPlanting
	var endDate sql.NullTime
	if err := h.db.QueryRow(query, args...).Scan(&planting.ID, &planting.BedID, &planting.PlantID, &planting.StartDate, &endDate, &planting.Notes, &planting.CreatedAt, &planting.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondError(w, http.StatusNotFound, "not_found", "Planting not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "db_update_failed", "Unable to update planting")
		return
	}
	if endDate.Valid {
		planting.EndDate = &endDate.Time
	}

	err = h.db.QueryRow(`SELECT name, variety FROM plants WHERE id = $1`, planting.PlantID).Scan(&planting.PlantName, &planting.Variety)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "db_query_failed", "Unable to load plant")
		return
	}

	respondJSON(w, http.StatusOK, dataEnvelope[bedPlanting]{Data: planting})
}

func (h *Handlers) deleteBedPlanting(w http.ResponseWriter, r *http.Request) {
	plantingID, err := parseUUIDParam(r, "plantingID")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid_id", "Planting ID is invalid")
		return
	}

	res, err := h.db.Exec("DELETE FROM bed_plantings WHERE id = $1", plantingID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "db_delete_failed", "Unable to delete planting")
		return
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		respondError(w, http.StatusNotFound, "not_found", "Planting not found")
		return
	}

	respondJSON(w, http.StatusOK, dataEnvelope[map[string]string]{Data: map[string]string{"status": "deleted"}})
}

func decodeJSON(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(dst)
}

func respondJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func respondError(w http.ResponseWriter, status int, code string, message string) {
	respondJSON(w, status, errorEnvelope{Error: errorPayload{Code: code, Message: message}})
}

func parseUUIDParam(r *http.Request, param string) (string, error) {
	value := chi.URLParam(r, param)
	_, err := uuid.Parse(value)
	if err != nil {
		return "", err
	}
	return value, nil
}

func parseDate(value string) (time.Time, error) {
	return time.Parse("2006-01-02", strings.TrimSpace(value))
}

func (h *Handlers) idExists(table string, id string) bool {
	var exists bool
	err := h.db.QueryRow("SELECT EXISTS (SELECT 1 FROM "+table+" WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

func (h *Handlers) loadPlants() ([]plant, error) {
	rows, err := h.db.Query(`SELECT id, name, variety, season, spacing, notes, days_to_maturity, sow_depth, germination_time, when_to_sow, days_to_emerge, seed_spacing, row_spacing, thinning, created_at, updated_at FROM plants`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	plants := make([]plant, 0)
	for rows.Next() {
		var p plant
		var days sql.NullInt64
		if err := rows.Scan(&p.ID, &p.Name, &p.Variety, &p.Season, &p.Spacing, &p.Notes, &days, &p.SowDepth, &p.GerminationTime, &p.WhenToSow, &p.DaysToEmerge, &p.SeedSpacing, &p.RowSpacing, &p.Thinning, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		if days.Valid {
			value := int(days.Int64)
			p.DaysToMaturity = &value
		}
		plants = append(plants, p)
	}
	return plants, nil
}

func (h *Handlers) loadPlantRefs(ids []string) ([]journalPlantRef, error) {
	if len(ids) == 0 {
		return []journalPlantRef{}, nil
	}
	query, args := buildInQuery("SELECT id, name, variety FROM plants WHERE id IN ", ids)
	rows, err := h.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	refs := make([]journalPlantRef, 0)
	for rows.Next() {
		var ref journalPlantRef
		if err := rows.Scan(&ref.ID, &ref.Name, &ref.Variety); err != nil {
			return nil, err
		}
		refs = append(refs, ref)
	}
	return refs, nil
}

func (h *Handlers) loadEntryPlants(entryID string) ([]journalPlantRef, error) {
	rows, err := h.db.Query(`
		SELECT p.id, p.name, p.variety
		FROM journal_entry_plants jep
		JOIN plants p ON p.id = jep.plant_id
		WHERE jep.entry_id = $1
	`, entryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	refs := make([]journalPlantRef, 0)
	for rows.Next() {
		var ref journalPlantRef
		if err := rows.Scan(&ref.ID, &ref.Name, &ref.Variety); err != nil {
			return nil, err
		}
		refs = append(refs, ref)
	}
	return refs, nil
}

func matchPlantIDs(text string, plants []plant) []string {
	normalizedText := " " + normalizeForMatch(text) + " "
	matches := make([]string, 0)
	seen := make(map[string]struct{})
	varietyMatched := false
	for _, plant := range plants {
		if plant.ID == "" {
			continue
		}
		nameKey := normalizeForMatch(plant.Name)
		varietyKey := normalizeForMatch(plant.Variety)
		if varietyKey != "" && strings.Contains(normalizedText, " "+varietyKey+" ") {
			if !varietyMatched {
				varietyMatched = true
				matches = matches[:0]
				seen = make(map[string]struct{})
			}
			if _, exists := seen[plant.ID]; !exists {
				seen[plant.ID] = struct{}{}
				matches = append(matches, plant.ID)
			}
			continue
		}
		if varietyMatched {
			continue
		}
		if nameKey != "" && strings.Contains(normalizedText, " "+nameKey+" ") {
			if _, exists := seen[plant.ID]; !exists {
				seen[plant.ID] = struct{}{}
				matches = append(matches, plant.ID)
			}
		}
	}
	return matches
}

func uniqueValidPlantIDs(ids []string) []string {
	seen := make(map[string]struct{})
	valid := make([]string, 0)
	for _, id := range ids {
		trimmed := strings.TrimSpace(id)
		if trimmed == "" {
			continue
		}
		if _, err := uuid.Parse(trimmed); err != nil {
			continue
		}
		if _, exists := seen[trimmed]; exists {
			continue
		}
		seen[trimmed] = struct{}{}
		valid = append(valid, trimmed)
	}
	return valid
}

func unionIDs(primary []string, secondary []string) []string {
	seen := make(map[string]struct{})
	union := make([]string, 0, len(primary)+len(secondary))
	for _, id := range primary {
		if _, exists := seen[id]; exists {
			continue
		}
		seen[id] = struct{}{}
		union = append(union, id)
	}
	for _, id := range secondary {
		if _, exists := seen[id]; exists {
			continue
		}
		seen[id] = struct{}{}
		union = append(union, id)
	}
	return union
}

func buildInQuery(prefix string, values []string) (string, []any) {
	placeholders := make([]string, 0, len(values))
	args := make([]any, 0, len(values))
	for i, value := range values {
		placeholders = append(placeholders, "$"+itoa(i+1))
		args = append(args, value)
	}
	return prefix + "(" + strings.Join(placeholders, ", ") + ")", args
}

func normalizeForMatch(value string) string {
	lowered := strings.ToLower(value)
	var builder strings.Builder
	builder.Grow(len(lowered))
	lastWasSpace := false
	for _, r := range lowered {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			builder.WriteRune(r)
			lastWasSpace = false
			continue
		}
		if !lastWasSpace {
			builder.WriteByte(' ')
			lastWasSpace = true
		}
	}
	return strings.TrimSpace(builder.String())
}

func normalizeEnum(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func itoa(value int) string {
	return strconv.Itoa(value)
}

func clampQueryInt(r *http.Request, key string, fallback int, min int, max int) int {
	valueStr := r.URL.Query().Get(key)
	if valueStr == "" {
		return fallback
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return fallback
	}
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
