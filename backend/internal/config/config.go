package config

import "os"

type Config struct {
	Port string
	DatabaseURL   string
	RunMigrations bool
	MigrationsDir string
}

func FromEnv() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbURL := os.Getenv("DATABASE_URL")
	runMigrations := os.Getenv("RUN_MIGRATIONS") == "true"
	migrationsDir := os.Getenv("MIGRATIONS_DIR")
	if migrationsDir == "" {
		migrationsDir = "migrations"
	}

	return Config{
		Port:          port,
		DatabaseURL:   dbURL,
		RunMigrations: runMigrations,
		MigrationsDir: migrationsDir,
	}
}
