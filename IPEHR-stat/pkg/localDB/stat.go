package localDB

import (
	"fmt"
	"time"
)

const (
	TableNamePatients  = "stat_patients"
	TableNameDocuments = "stat_documents"
)

func (db *DB) StatPatientsCountGet(start, end int64) (uint64, error) {
	query := `SELECT SUM(count) 
			  FROM ` + TableNamePatients + `
			  WHERE timestamp_day >= ? AND timestamp_day < ?`

	row := db.db.QueryRow(query, start, end)

	var count uint64
	err := row.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("row.Scan error: %w", err)
	}

	return count, nil
}

func (db *DB) StatPatientsCountIncrement(timestamp time.Time) error {
	db.Lock()
	defer db.Unlock()

	timestamp = timestamp.Truncate(time.Hour * 24)

	query := `INSERT INTO ` + TableNamePatients + ` (timestamp_day, count) VALUES (?, 1)
			  ON CONFLICT (timestamp_day) DO UPDATE SET 
			  count = count + 1`

	_, err := db.db.Exec(query, timestamp.Unix())
	if err != nil {
		return fmt.Errorf("StatPatientsCountIncrement error: %w query: %s timestamp: %d", err, query, timestamp.Unix())
	}

	return nil
}

func (db *DB) StatDocumentsCountGet(start, end int64) (uint64, error) {
	query := `SELECT SUM(count) 
			  FROM ` + TableNameDocuments + ` 
			  WHERE timestamp_day >= ? AND timestamp_day < ?`

	row := db.db.QueryRow(query, start, end)

	var count uint64
	err := row.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("row.Scan error: %w", err)
	}

	return count, nil
}

func (db *DB) StatDocumentsCountIncrement(timestamp time.Time) error {
	db.Lock()
	defer db.Unlock()

	timestamp = timestamp.Truncate(time.Hour * 24)

	query := `INSERT INTO ` + TableNameDocuments + ` (timestamp_day, count) VALUES (?, 1)
			  ON CONFLICT (timestamp_day) DO UPDATE SET 
			  count = count + 1`

	_, err := db.db.Exec(query, timestamp.Unix())
	if err != nil {
		return fmt.Errorf("StatPatientsCountIncrement error: %w query: %s timestamp: %d", err, query, timestamp.Unix())
	}

	return nil
}
