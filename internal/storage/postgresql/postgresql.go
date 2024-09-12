package postgresql

import (
	"database/sql"
	"fmt"
	"url-shortener/internal/storage"

	"github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(connStr string) (*Storage, error) {
	const op = "storage.postresql.New"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	db.Exec(`DROP TABLE "url"`)

	_, err = db.Exec(`
		CREATE TABLE "url" (
	"id" SERIAL NOT NULL UNIQUE,
	"alias" TEXT NOT NULL UNIQUE,
	"url" TEXT NOT NULL,
	PRIMARY KEY("id")
	);
	CREATE INDEX "url_index"
	ON "url" ("alias");
	`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(alias, urlToSave string) (int64, error) {
	const op = "storage.postresql.SaveURL"

	var id int64

	err := s.db.QueryRow("INSERT INTO url (alias, url) VALUES ($1, $2) RETURNING id", alias, urlToSave).Scan(&id)
	if err != nil {
		// Приведение ошибки к *pq.Error для проверки кода ошибки
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // Код ошибки для нарушения уникальности
				return 0, fmt.Errorf("%s: %w", op, storage.ErrExists)
			}
		}
		// Любая другая ошибка
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *Storage) GETUrl(alias string) (string, error) {
	const op = "storage.postresql.GETUrl"

	var url string

	err := s.db.QueryRow(`SELECT url FROM "url" WHERE alias = $1`, alias).Scan(&url)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("%s: %w", op, storage.ErrNotFound)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return url, nil
}

func (s *Storage) DeleteURL(alias string) error {
	const op = "storage.postresql.DeleteURL"

	_, err := s.db.Exec(`DELETE FROM "url" WHERE alias = $1`, alias)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("%s: %w", op, storage.ErrNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
