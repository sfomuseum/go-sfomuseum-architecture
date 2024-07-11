package campus

import (
	"context"
	"database/sql"
	"fmt"
)

// SFO Terminal Complex (1954~ to 1963~)
// https://millsfield.sfomuseum.org/buildings/1159396329/
const FIRST_SFO int64 = 1159396329

// MostRecentComplexWithIterator will return a `Complex` instance representing the most recent relationships of the SFO terminal complex
// and its descendants derived from data defined in 'iterator_uri' and 'paths'.
func MostRecentComplexWithIterator(ctx context.Context, iterator_uri string, paths ...string) (*Complex, error) {

	dsn := ":memory:"

	db, err := newWhosOnFirstDatabaseFromIterator(ctx, dsn, iterator_uri, paths...)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive new WOF database, %w", err)
	}

	sql_db, err := db.Conn()

	if err != nil {
		return nil, fmt.Errorf("Failed to create database connection, %w", err)
	}

	return FindMostRecentComplexWithDatabase(ctx, sql_db)
}

func FindMostRecentComplexWithDatabase(ctx context.Context, db *sql.DB) (*Complex, error) {

	sfo_id, err := findMostRecentComplexID(ctx, db, FIRST_SFO)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive most recent complex ID, %w", err)
	}

	return FindComplex(ctx, db, sfo_id)
}

func findMostRecentComplexID(ctx context.Context, db *sql.DB, id int64) (int64, error) {

	q := "SELECT DISTINCT(superseded_by_id) FROM supersedes WHERE superseded_id = ?"

	rows, err := db.QueryContext(ctx, q, id)

	if err != nil {

		if err == sql.ErrNoRows {
			return id, nil
		}

		return -1, err
	}

	defer rows.Close()
	possible := make([]int64, 0)

	for rows.Next() {

		var superseded_by int64
		err := rows.Scan(&superseded_by)

		if err != nil {
			return -1, fmt.Errorf("Failed to scan row, %w", err)
		}

		possible = append(possible, superseded_by)
	}

	err = rows.Close()

	if err != nil {
		return -1, err
	}

	err = rows.Err()

	if err != nil {
		return -1, err
	}

	switch len(possible) {
	case 0:
		return id, nil
	case 1:
		return findMostRecentComplexID(ctx, db, possible[0])
	default:
		return -1, fmt.Errorf("Multiple results for '%d', not implemented", id)
	}
}
