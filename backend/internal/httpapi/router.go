package httpapi

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Router(db *sql.DB, companions *CompanionIndex) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	h := &Handlers{db: db, companions: companions}

	r.Get("/api/health", h.health)

	r.Route("/api/gardens", func(r chi.Router) {
		r.Get("/", h.listGardens)
		r.Post("/", h.createGarden)
		r.Route("/{gardenID}", func(r chi.Router) {
			r.Get("/", h.getGarden)
			r.Patch("/", h.updateGarden)
			r.Delete("/", h.deleteGarden)
			r.Route("/beds", func(r chi.Router) {
				r.Get("/", h.listBedsByGarden)
				r.Post("/", h.createBed)
			})
			r.Route("/planters", func(r chi.Router) {
				r.Get("/", h.listPlantersByGarden)
				r.Post("/", h.createPlanter)
			})
		})
	})

	r.Route("/api/beds", func(r chi.Router) {
		r.Get("/{bedID}", h.getBed)
		r.Patch("/{bedID}", h.updateBed)
		r.Delete("/{bedID}", h.deleteBed)
		r.Route("/{bedID}/plantings", func(r chi.Router) {
			r.Get("/", h.listBedPlantings)
			r.Post("/", h.createBedPlanting)
		})
	})

	r.Route("/api/planters", func(r chi.Router) {
		r.Get("/{planterID}", h.getPlanter)
		r.Patch("/{planterID}", h.updatePlanter)
		r.Delete("/{planterID}", h.deletePlanter)
	})

	r.Route("/api/plants", func(r chi.Router) {
		r.Get("/", h.listPlants)
		r.Post("/", h.createPlant)
		r.Get("/{plantID}", h.getPlant)
		r.Patch("/{plantID}", h.updatePlant)
		r.Delete("/{plantID}", h.deletePlant)
	})

	r.Route("/api/journal", func(r chi.Router) {
		r.Get("/", h.listJournalEntries)
		r.Post("/", h.createJournalEntry)
		r.Route("/{entryID}", func(r chi.Router) {
			r.Patch("/", h.updateJournalEntry)
			r.Delete("/", h.deleteJournalEntry)
		})
	})

	r.Post("/api/seed-packets/scan", h.scanSeedPacket)

	r.Get("/api/companions", h.getCompanionData)

	r.Route("/api/plantings", func(r chi.Router) {
		r.Patch("/{plantingID}", h.updateBedPlanting)
		r.Delete("/{plantingID}", h.deleteBedPlanting)
	})

	return r
}
