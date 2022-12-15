package localDB

import (
	"fmt"
)

func (db *DB) SyncLastBlockGet() (uint64, error) {
	row := db.db.QueryRow("SELECT value FROM sync WHERE key = 'last_synced_block'")

	var count uint64
	err := row.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("row.Scan error: %w", err)
	}

	return count, nil
}

func (db *DB) SyncLastBlockSet(lastSyncedBlock uint64) error {
	db.Lock()
	defer db.Unlock()

	query := `INSERT INTO sync (key, value) VALUES ('last_synced_block', $1)
			  ON CONFLICT (key) DO UPDATE SET 
			  value = $1`

	_, err := db.db.Exec(query, lastSyncedBlock)
	if err != nil {
		return err
	}

	return nil
}
