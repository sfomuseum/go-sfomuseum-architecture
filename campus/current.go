package campus

// We could (might still) use go-whosonfirst-travel to determine the most recent SFO
// but since we have to (want to) use SQLite for querying relationships its just as
// easy to use that same database to follow superseded_by breadcrumbs to the "present"
// (20210513/thisisaaronland)

/*

> ./bin/wof-travel-id -superseded-by -source fs:///usr/local/data/sfomuseum-data-architecture/data 1159396329
[1] 1159396329 SFO Terminal Complex [1954~] [1963~]
[2] 1159396325 SFO Terminal Complex [1963~] [1974~]
[3] 1159396339 SFO Terminal Complex [1974~] [1979~]
[4] 1159396331 SFO Terminal Complex [1979~] [1981~]
[5] 1159396327 SFO Terminal Complex [1981~] [1983~]
[6] 1159554801 SFO Terminal Complex [1983~] [1988~]
[7] 1159554803 SFO Terminal Complex [1988~] [2000~]
[8] 1159396319 SFO Terminal Complex [2000~] [2006~]
[9] 1159396337 SFO Terminal Complex [2006~] [2011~]
[10] 1159396333 SFO Terminal Complex [2011~] [2014~]
[11] 1159396321 SFO Terminal Complex [2014~] [2017~]
[12] 1159157271 SFO Terminal Complex [2017~] [2019-07-23]
[13] 1477855605 SFO Terminal Complex [2019-07-23] [2020-~05]
[14] 1729792387 SFO Terminal Complex [2020-~05] [2021-05-25]
[15] 1745882083 SFO Terminal Complex [2021-05-25] [..]

*/

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/aaronland/go-sqlite"
	aa_database "github.com/aaronland/go-sqlite/database"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-whosonfirst-sqlite-features-index"
	"github.com/whosonfirst/go-whosonfirst-sqlite-features/tables"
	sql_index "github.com/whosonfirst/go-whosonfirst-sqlite-index/v3"
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

	terminals, err := findTerminals(ctx, db, sfo_id)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive terminals for complex %d, %w", sfo_id, err)
	}

	c := &Complex{
		WhosOnFirstId: sfo_id,
		SFOId:         "SFO",
		Terminals:     terminals,
	}

	return c, nil
}

func findTerminals(ctx context.Context, db *sql.DB, sfo_id int64) ([]*Terminal, error) {

	slog.Info("Find terminals", "parent id", sfo_id)

	terminal_ids, err := findChildIDs(ctx, db, sfo_id, "terminal")

	if err != nil {
		return nil, fmt.Errorf("Failed to find any child records (terminals) for %d, %v", sfo_id, err)
	}

	terminals := make([]*Terminal, len(terminal_ids))

	for idx, t_id := range terminal_ids {

		commonareas, err := findCommonAreas(ctx, db, t_id)

		if err != nil {
			return nil, fmt.Errorf("Failed to derive common areas for %d, %w", t_id, err)
		}

		boardingareas, err := findBoardingAreas(ctx, db, t_id)

		if err != nil {
			return nil, fmt.Errorf("Failed to derive boarding areas for %d, %w", t_id, err)
		}

		t_body, err := loadFeatureWithDBAndChecks(ctx, db, t_id)

		if err != nil {
			return nil, fmt.Errorf("Failed to load feature for %d, %w", t_id, err)
		}

		var sfoid string

		rsp := gjson.GetBytes(t_body, "properties.sfo:id")

		if rsp.Exists() {

			sfoid = rsp.String()

		} else {

			rsp := gjson.GetBytes(t_body, "properties.sfomuseum:terminal_id")

			if !rsp.Exists() {
				return nil, fmt.Errorf("Missing properties.sfomuseum:terminal_id property for terminal %d", t_id)
			}

			switch rsp.String() {
			case "ITB":
				sfoid = "100" // gis.BUILDING_ITB
			case "T1":
				sfoid = "200" // gis.BUILDING_T1
			case "T2":
				sfoid = "300" // gis.BUILDING_T2
			case "T3":
				sfoid = "400" // gis.BUILDING_T3
			default:
				return nil, fmt.Errorf("Unrecognized terminal_id '%s' for %d", rsp.String(), t_id)
			}
		}

		slog.Info("Add terminal", "parent id", sfoid, "id", t_id)

		terminal := &Terminal{
			WhosOnFirstId: t_id,
			SFOId:         sfoid,
		}

		if len(commonareas) > 0 {
			terminal.CommonAreas = commonareas
		}

		if len(boardingareas) > 0 {
			terminal.BoardingAreas = boardingareas
		}

		terminals[idx] = terminal
	}

	return terminals, nil
}

func findObservationDecks(ctx context.Context, db *sql.DB, t_id int64) ([]*ObservationDeck, error) {

	slog.Info("Find observation decks", "parent id", t_id)

	deck_ids, err := findChildIDs(ctx, db, t_id, "observationdeck")

	if err != nil {
		return nil, fmt.Errorf("Failed to find any child records (observation decks) for %d, %v", t_id, err)
	}

	decks := make([]*ObservationDeck, len(deck_ids))

	for idx, d_id := range deck_ids {

		galleries, err := findGalleries(ctx, db, d_id)

		if err != nil {
			return nil, fmt.Errorf("Failed to derive galleries for observation deck %d, %w", d_id, err)
		}

		publicart, err := findPublicArt(ctx, db, d_id)

		if err != nil {
			return nil, fmt.Errorf("Failed to derive public art for observation deck %d, %w", d_id, err)
		}

		d_body, err := loadFeatureWithDBAndChecks(ctx, db, d_id)

		if err != nil {
			return nil, fmt.Errorf("Failed to load feature for observation deck %d, %w", d_id, err)
		}

		var sfoid string

		rsp := gjson.GetBytes(d_body, "properties.sfo:id")

		if !rsp.Exists() {
			return nil, fmt.Errorf("Unable to find sfo:id for WOF record, %d", d_id)
		}

		sfoid = rsp.String()

		slog.Info("Add observation deck", "sfo id", sfoid, "parent id", t_id, "id", d_id)

		deck := &ObservationDeck{
			WhosOnFirstId: d_id,
			SFOId:         sfoid,
		}

		if len(galleries) > 0 {
			deck.Galleries = galleries
		}

		if len(publicart) > 0 {
			deck.PublicArt = publicart
		}

		decks[idx] = deck
	}

	return decks, nil
}

func findCommonAreas(ctx context.Context, db *sql.DB, parent_id int64) ([]*CommonArea, error) {

	slog.Info("Find common areas", "parent", parent_id)

	commonarea_ids, err := findChildIDs(ctx, db, parent_id, "commonarea")

	if err != nil {
		return nil, fmt.Errorf("Failed to find any child records (common areas) for %d, %v", parent_id, err)
	}

	commonareas := make([]*CommonArea, len(commonarea_ids))

	for idx, c_id := range commonarea_ids {

		gates, err := findGates(ctx, db, c_id)

		if err != nil {
			return nil, fmt.Errorf("Failed to derive gates for common area %d, %w", c_id, err)
		}

		checkpoints, err := findCheckpoints(ctx, db, c_id)

		if err != nil {
			return nil, fmt.Errorf("Failed to derive gates for check points %d, %w", c_id, err)
		}

		galleries, err := findGalleries(ctx, db, c_id)

		if err != nil {
			return nil, fmt.Errorf("Failed to derive gates for galleries %d, %w", c_id, err)
		}

		observation_decks, err := findObservationDecks(ctx, db, c_id)

		if err != nil {
			return nil, fmt.Errorf("Failed to derive observation decks for galleries %d, %w", c_id, err)
		}

		publicart, err := findPublicArt(ctx, db, c_id)

		if err != nil {
			return nil, fmt.Errorf("Failed to derive public art for galleries %d, %w", c_id, err)
		}

		c_body, err := loadFeatureWithDBAndChecks(ctx, db, c_id)

		if err != nil {
			return nil, fmt.Errorf("Failed to load feature %d, %w", c_id, err)
		}

		var sfoid string

		rsp := gjson.GetBytes(c_body, "properties.sfo:id")

		if rsp.Exists() {

			sfoid = rsp.String()

		} else {

			rsp := gjson.GetBytes(c_body, "properties.sfo:building_id")

			if !rsp.Exists() {
				return nil, fmt.Errorf("Unable to find sfo:building_id for %d", c_id)
			}

			switch rsp.String() {
			case "ITB", "100":
				sfoid = "100CAD" // gis.COMMONAREA_ITB_DEPARTURES
			case "T1", "200":
				sfoid = "200CAD" // gis.COMMONAREA_T1_DEPARTURES
			case "T2", "300":
				sfoid = "300CAD" // gis.COMMONAREA_T2_DEPARTURES
			case "T3", "400":
				sfoid = "400CAD" // gis.COMMONAREA_T3_DEPARTURES
			default:
				return nil, fmt.Errorf("Unrecognized sfo:id '%s' for %d", rsp.String(), c_id)
			}
		}

		area := &CommonArea{
			WhosOnFirstId: c_id,
			SFOId:         sfoid,
		}

		if len(gates) > 0 {
			area.Gates = gates
		}

		if len(checkpoints) > 0 {
			area.Checkpoints = checkpoints
		}

		if len(galleries) > 0 {
			area.Galleries = galleries
		}

		if len(publicart) > 0 {
			area.PublicArt = publicart
		}

		if len(observation_decks) > 0 {
			area.ObservationDecks = observation_decks
		}

		commonareas[idx] = area
	}

	return commonareas, nil
}

func findBoardingAreas(ctx context.Context, db *sql.DB, id int64) ([]*BoardingArea, error) {

	slog.Info("Find boarding areas", "parent", id)

	boardingarea_ids, err := findChildIDs(ctx, db, id, "boardingarea")

	if err != nil {
		return nil, fmt.Errorf("Failed to find any child records (boarding areas areas) for %d, %v", id, err)
	}

	boardingareas := make([]*BoardingArea, len(boardingarea_ids))

	for idx, b_id := range boardingarea_ids {

		gates, err := findGates(ctx, db, b_id)

		if err != nil {
			return nil, fmt.Errorf("Failed to derive gates for boarding area %d, %w", b_id, err)
		}

		checkpoints, err := findCheckpoints(ctx, db, b_id)

		if err != nil {
			return nil, fmt.Errorf("Failed to derive check points for boarding area %d, %w", b_id, err)
		}

		galleries, err := findGalleries(ctx, db, b_id)

		if err != nil {
			return nil, fmt.Errorf("Failed to derive galleries for boarding area %d, %w", b_id, err)
		}

		publicart, err := findPublicArt(ctx, db, b_id)

		if err != nil {
			return nil, fmt.Errorf("Failed to derive public art for boarding area %d, %w", b_id, err)
		}

		observation_decks, err := findObservationDecks(ctx, db, b_id)

		if err != nil {
			return nil, fmt.Errorf("Failed to derive observation decks for boarding area %d, %w", b_id, err)
		}

		b_body, err := loadFeatureWithDBAndChecks(ctx, db, b_id)

		if err != nil {
			return nil, fmt.Errorf("Failed to load feature for %d, %w", b_id, err)
		}

		var sfoid string

		rsp := gjson.GetBytes(b_body, "properties.sfo:id")

		if rsp.Exists() {

			sfoid = rsp.String()

		} else {

			rsp := gjson.GetBytes(b_body, "properties.sfo:building_id")

			if !rsp.Exists() {
				return nil, fmt.Errorf("Missing sfo:building_id for boarding area %d", b_id)
			}

			sfoid = rsp.String()
		}

		area := &BoardingArea{
			WhosOnFirstId: b_id,
			SFOId:         sfoid,
		}

		if len(gates) > 0 {
			area.Gates = gates
		}

		if len(checkpoints) > 0 {
			area.Checkpoints = checkpoints
		}

		if len(galleries) > 0 {
			area.Galleries = galleries
		}

		if len(publicart) > 0 {
			area.PublicArt = publicart
		}

		if len(observation_decks) > 0 {
			area.ObservationDecks = observation_decks
		}

		boardingareas[idx] = area
	}

	return boardingareas, nil

}

func findGates(ctx context.Context, db *sql.DB, parent_id int64) ([]*Gate, error) {

	slog.Info("Find gates", "parent", parent_id)

	gate_ids, err := findChildIDs(ctx, db, parent_id, "gate")

	if err != nil {
		return nil, fmt.Errorf("Failed to find any child records (gates) for %d, %w", parent_id, err)
	}

	gates := make([]*Gate, len(gate_ids))

	for idx, g_id := range gate_ids {

		g_body, err := loadFeatureWithDBAndChecks(ctx, db, g_id)

		if err != nil {
			return nil, fmt.Errorf("Failed to load feature for date %d, %w", g_id, err)
		}

		var sfoid string

		rsp := gjson.GetBytes(g_body, "properties.sfo:id")

		if rsp.Exists() {

			sfoid = rsp.String()

		} else {

			rsp := gjson.GetBytes(g_body, "properties.wof:name")

			if !rsp.Exists() {
				return nil, fmt.Errorf("Missing wof:name for %d", g_id)
			}

			sfoid = rsp.String()
		}

		slog.Info("Add gate", "sfo id", sfoid, "parent_id", parent_id, "id", g_id)

		g := &Gate{
			WhosOnFirstId: g_id,
			SFOId:         sfoid,
		}

		gates[idx] = g
	}

	return gates, nil
}

func findCheckpoints(ctx context.Context, db *sql.DB, parent_id int64) ([]*Checkpoint, error) {

	slog.Info("Find check points", "parent id", parent_id)

	checkpoint_ids, err := findChildIDs(ctx, db, parent_id, "checkpoint")

	if err != nil {
		return nil, fmt.Errorf("Failed to find any child records (checkpoints) for %d, %w", parent_id, err)
	}

	checkpoints := make([]*Checkpoint, len(checkpoint_ids))

	for idx, cp_id := range checkpoint_ids {

		cp_body, err := loadFeatureWithDBAndChecks(ctx, db, cp_id)

		if err != nil {
			return nil, fmt.Errorf("Failed to load feature for %d, %w", cp_id, err)
		}

		var sfoid string

		rsp := gjson.GetBytes(cp_body, "properties.sfo:id")

		if !rsp.Exists() {
			return nil, fmt.Errorf("Missing sfo:id for %d", cp_id)
		}

		sfoid = rsp.String()

		slog.Info("Add checkpoint", "sfo id", sfoid, "parent id", parent_id, "id", cp_id)

		cp := &Checkpoint{
			WhosOnFirstId: cp_id,
			SFOId:         sfoid,
		}

		checkpoints[idx] = cp
	}

	return checkpoints, nil
}

func findGalleries(ctx context.Context, db *sql.DB, parent_id int64) ([]*Gallery, error) {

	slog.Info("Find galleries", "parent id", parent_id)

	gallery_ids, err := findChildIDs(ctx, db, parent_id, "gallery")

	if err != nil {
		return nil, fmt.Errorf("Failed to find any child records (galleries) for %d, %w", parent_id, err)
	}

	galleries := make([]*Gallery, len(gallery_ids))

	for idx, g_id := range gallery_ids {

		g_body, err := loadFeatureWithDBAndChecks(ctx, db, g_id)

		if err != nil {
			return nil, fmt.Errorf("Failed load feature for gallery %d, %w", g_id, err)
		}

		var sfomid string

		rsp := gjson.GetBytes(g_body, "properties.sfomuseum:map_id")

		if rsp.Exists() {

			sfomid = rsp.String()

		} else {

			rsp := gjson.GetBytes(g_body, "properties.sfomuseum:gallery_id")

			if !rsp.Exists() {
				return nil, fmt.Errorf("Missing sfomuseum:gallery_id property for gallery %d", g_id)
			}

			sfomid = rsp.String()
		}

		slog.Info("Add gallery", "sfo id", sfomid, "parent id", parent_id, "id", g_id)

		g := &Gallery{
			WhosOnFirstId: g_id,
			SFOId:         sfomid,
		}

		galleries[idx] = g
	}

	return galleries, nil
}

func findPublicArt(ctx context.Context, db *sql.DB, parent_id int64) ([]*PublicArt, error) {

	slog.Info("Find public art", "parent id", parent_id)

	publicart_ids, err := findChildIDs(ctx, db, parent_id, "publicart")

	if err != nil {
		return nil, fmt.Errorf("Failed to find any child records (public art) for %d, %w", parent_id, err)
	}

	publicarts := make([]*PublicArt, len(publicart_ids))

	for idx, p_id := range publicart_ids {

		p_body, err := loadFeatureWithDBAndChecks(ctx, db, p_id)

		if err != nil {
			return nil, fmt.Errorf("Failed to load feature for public art %d, %w", p_id, err)
		}

		var sfomid string

		rsp := gjson.GetBytes(p_body, "properties.sfomuseum:map_id")

		if rsp.Exists() {

			sfomid = rsp.String()

		} else {

			rsp := gjson.GetBytes(p_body, "properties.sfomuseum:object_id")

			if !rsp.Exists() {
				return nil, fmt.Errorf("Missing sfomuseum:object_id property for public art %d, %w", p_id, err)
			}

			sfomid = rsp.String()
		}

		slog.Info("Add public art", "sfo id", sfomid, "parent id", parent_id, "id", p_id)

		g := &PublicArt{
			WhosOnFirstId: p_id,
			SFOId:         sfomid,
		}

		publicarts[idx] = g
	}

	return publicarts, nil
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

func findChildIDs(ctx context.Context, db *sql.DB, parent_id int64, placetype string) ([]int64, error) {

	q := `SELECT s.id FROM spr s, geojson g WHERE s.id=g.id AND s.parent_id=? AND JSON_EXTRACT(g.body, '$.properties."sfomuseum:placetype"')=?`

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

	return children, nil
}

func loadFeatureWithDBAndChecks(ctx context.Context, db *sql.DB, id int64) ([]byte, error) {

	body, err := loadFeatureWithDB(ctx, db, id)

	if err != nil {
		return nil, fmt.Errorf("Failed to load feature for record %d, %w", id, err)
	}

	current_rsp := gjson.GetBytes(body, "properties.mz:is_current")

	if !current_rsp.Exists() {
		return nil, fmt.Errorf("Missing properties.mz:is_current property for record %d", id)
	}

	if current_rsp.Int() != 1 {
		return nil, fmt.Errorf("Unexpected mz:is_current property (%d) for record %d", current_rsp.Int(), id)
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

func newWhosOnFirstDatabaseFromIterator(ctx context.Context, dsn string, iterator_uri string, paths ...string) (*aa_database.SQLiteDatabase, error) {

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

	/*
		if opts.IndexRTreeTable {

			rtree_opts, err := tables.DefaultRTreeTableOptions()

			if err != nil {
				return nil, err
			}

			rtree_opts.IndexAltFiles = opts.IndexAltFiles

			rtree_table, err := tables.NewRTreeTableWithDatabaseAndOptions(ctx, db, rtree_opts)

			if err != nil {
				return nil, err
			}

			to_index = append(to_index, rtree_table)
		}

		if opts.IndexPropertiesTable || opts.IndexRTreeTable {

			properties_opts, err := tables.DefaultPropertiesTableOptions()

			if err != nil {
				return nil, err
			}

			properties_opts.IndexAltFiles = opts.IndexAltFiles

			properties_table, err := tables.NewPropertiesTableWithDatabaseAndOptions(ctx, db, properties_opts)

			if err != nil {
				return nil, err
			}

			to_index = append(to_index, properties_table)
		}
	*/

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
