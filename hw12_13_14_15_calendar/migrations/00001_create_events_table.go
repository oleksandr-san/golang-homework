package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(Up00001, Down00001)
}

func Up00001(tx *sql.Tx) error {
	_, err := tx.Exec(`CREATE TABLE events (
		id UUID PRIMARY KEY,
		owner UUID,
		title TEXT,
		descr TEXT,
		start_date DATE NOT NULL,
		start_time TIME,
		end_date DATE NOT NULL,
		end_time TIME
	);
	CREATE INDEX owner_idx ON events (owner);
	CREATE INDEX start_idx ON events USING BTREE (start_date, start_time);
`)
	if err != nil {
		return err
	}
	return nil
}

func Down00001(tx *sql.Tx) error {
	_, err := tx.Exec(`
	DROP INDEX IF EXISTS start_idx;
	DROP INDEX IF EXISTS owner_idx;
	DROP TABLE IF EXISTS events;`)
	if err != nil {
		return err
	}
	return nil
}
