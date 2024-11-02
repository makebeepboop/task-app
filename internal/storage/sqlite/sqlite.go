package sqlite

import (
	"database/sql"
	"fmt"
	"log/slog"
)

type Storage struct {
	log *slog.Logger
	db  *sql.DB
}

func New(log *slog.Logger, storagePath string) (*Storage, error) {
	const operation = "storage.sqlite.New"

	log = log.With(
		slog.String("operation", operation),
	)

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	log.Info("database opened")

	return &Storage{
		log: log,
		db:  db,
	}, nil
}

func (s *Storage) Stop() {
	const operation = "storage.sqlite.Stop"

	s.log.With(
		slog.String("operation", operation),
	)

	if err := s.db.Close(); err != nil {
		s.log.Warn("close database failed")
		return
	}

	s.log.Info("database closed")
}
