package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

type Repository struct {
	db *sql.DB
}

func completeTransaction(tx *sql.Tx, e *error) {
	if *e != nil {
		_ = tx.Rollback()
		*e = fmt.Errorf("rollback transaction due to: %w", *e)
	} else {
		err := tx.Commit()
		if err != nil {
			*e = fmt.Errorf("commit transaction error: %w", err)
		}
	}
}

func New(storagePath string) (r *Repository, e error) {
	defer func() {
		if e != nil {
			r = nil
		}
	}()
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("open db error: %w", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("start transaction error: %w", err)
	}
	defer completeTransaction(tx, &e)
	stmt, err := tx.Prepare(`
	CREATE TABLE IF NOT EXISTS url(
		id INTEGER PRIMARY KEY,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL);
	CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`)
	if err != nil {
		return nil, fmt.Errorf("prepare statement error: %w", err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("execute statement error: %w", err)
	}
	return &Repository{db: db}, nil
}

func (r *Repository) SaveUrl(urToSave string, alias string) (idx int64, e error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("start transaction error: %w", err)
	}
	defer completeTransaction(tx, &e)
	stmt, err := tx.Prepare(`	INSERT INTO url(alias, url) VALUES (?, ?)`)
	if err != nil {
		return 0, fmt.Errorf("prepare statement error: %w", err)
	}
	res, err := stmt.Exec(alias, urToSave)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, errors.New(fmt.Sprintf("couldn't insert given alias (%s) as it already exists", alias))
		}
		return 0, fmt.Errorf("execute statement error: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("couldn't retrieve last inserted id: %w", err)
	}
	return id, nil
}

func (r *Repository) GetUrl(alias string) (resultUrl string, e error) {
	tx, err := r.db.Begin()
	if err != nil {
		return "", fmt.Errorf("start transaction error: %w", err)
	}
	defer completeTransaction(tx, &e)
	row := tx.QueryRow("SELECT url FROM url WHERE alias=?", alias)
	err = row.Scan(&resultUrl)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("no alias found (%s): %w", alias, err)
		}
		return "", fmt.Errorf("couldn't parse result row: %w", err)
	}
	return
}
