package campus

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/aaronland/go-sqlite"
	aa_database "github.com/aaronland/go-sqlite/database"
	"github.com/sfomuseum/go-edtf"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-whosonfirst-sqlite-features-index"
	"github.com/whosonfirst/go-whosonfirst-sqlite-features/tables"
	sql_index "github.com/whosonfirst/go-whosonfirst-sqlite-index/v3"
)

var WARN_IS_CURRENT = true

func NewDatabaseWithIterator(ctx context.Context, dsn string, iterator_uri string, paths ...string) (*aa_database.SQLiteDatabase, error) {

	driver := "sqlite3"

	db, err := aa_database.NewDBWithDriver(ctx, driver, dsn)

	if err != nil {
		return nil, err
	}

	err = db.LiveHardDieFast()

	if err != nil {
		return nil, err
	}

	to_index := make([]sqlite.Table, 0)

	geojson_opts, err := tables.DefaultGeoJSONTableOptions()

	if err != nil {
		return nil, err
	}

	geojson_opts.IndexAltFiles = false

	geojson_table, err := tables.NewGeoJSONTableWithDatabaseAndOptions(ctx, db, geojson_opts)

	if err != nil {
		return nil, err
	}

	to_index = append(to_index, geojson_table)

	supersedes_table, err := tables.NewSupersedesTableWithDatabase(ctx, db)

	if err != nil {
		return nil, err
	}

	to_index = append(to_index, supersedes_table)

	spr_opts, err := tables.DefaultSPRTableOptions()

	if err != nil {
		return nil, err
	}

	spr_table, err := tables.NewSPRTableWithDatabaseAndOptions(ctx, db, spr_opts)

	if err != nil {
		return nil, err
	}

	to_index = append(to_index, spr_table)

	record_opts := &index.SQLiteFeaturesLoadRecordFuncOptions{
		StrictAltFiles: false,
	}

	record_func := index.SQLiteFeaturesLoadRecordFunc(record_opts)

	idx_opts := &sql_index.SQLiteIndexerOptions{
		DB:             db,
		Tables:         to_index,
		LoadRecordFunc: record_func,
	}

	idx, err := sql_index.NewSQLiteIndexer(idx_opts)

	if err != nil {
		return nil, err
	}

	err = idx.IndexPaths(ctx, iterator_uri, paths)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func findChildIDs(ctx context.Context, db *sql.DB, parent_id int64, placetype string) ([]int64, error) {

	q := `SELECT s.id FROM spr s, geojson g WHERE s.id=g.id AND s.parent_id=? AND JSON_EXTRACT(g.body, '$.properties."sfomuseum:placetype"')=?`

	slog.Debug(q, "parent_id", parent_id, "placetype", placetype)

	rows, err := db.QueryContext(ctx, q, parent_id, placetype)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	children := make([]int64, 0)

	for rows.Next() {

		var superseded_by int64
		err := rows.Scan(&superseded_by)

		if err != nil {
			return nil, err
		}

		children = append(children, superseded_by)
	}

	err = rows.Close()

	if err != nil {
		return nil, err
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	if placetype == "publicart" {

		slog.Info("WTF", "parent id", parent_id, "children", children)
	}

	return children, nil
}

func loadFeatureWithDBAndChecks(ctx context.Context, db *sql.DB, id int64) ([]byte, error) {

	body, err := loadFeatureWithDB(ctx, db, id)

	if err != nil {
		return nil, fmt.Errorf("Failed to load feature for record %d, %w", id, err)
	}

	name_rsp := gjson.GetBytes(body, "properties.wof:name")
	inception_rsp := gjson.GetBytes(body, "properties.edtf:inception")
	cessation_rsp := gjson.GetBytes(body, "properties.edtf:cessation")

	deprecated_rsp := gjson.GetBytes(body, "properties.edtf:deprecated")

	if deprecated_rsp.Exists() && deprecated_rsp.String() != "" {
		return nil, nil
	}

	current_rsp := gjson.GetBytes(body, "properties.mz:is_current")

	if !current_rsp.Exists() {
		return nil, fmt.Errorf("Missing properties.mz:is_current property for record %d", id)
	}

	if current_rsp.Int() != 1 && WARN_IS_CURRENT {

		cessation_str := cessation_rsp.String()

		if cessation_str == "" || cessation_str == edtf.OPEN {
			slog.Warn("Unexpected mz:is_current property", "id", id, "mz:is_current", current_rsp.Int(), "name", name_rsp.String(), "inception", inception_rsp.String(), "cessation", cessation_rsp.String())
		}

		// return nil, nil
	}

	return body, nil
}

func loadFeatureWithDB(ctx context.Context, db *sql.DB, id int64) ([]byte, error) {

	q := "SELECT body FROM geojson WHERE id = ?"

	row := db.QueryRowContext(ctx, q, id)

	var body string
	err := row.Scan(&body)

	if err != nil {
		return nil, err
	}

	return []byte(body), nil
}
