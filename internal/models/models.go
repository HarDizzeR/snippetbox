package models

import (
	"database/sql"
	"log"
)

// InitDB initializes the database schema if it doesn't exist
func InitDB(db *sql.DB) error {
	// Create snippets table if it doesn't exist
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS snippets (
            id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
            title VARCHAR(100) NOT NULL,
            content TEXT NOT NULL,
            created DATETIME NOT NULL,
            expires DATETIME NOT NULL
        )
    `)
	if err != nil {
		return err
	}

	// Check if index exists before creating it
	var indexExists int
	err = db.QueryRow("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema=DATABASE() AND table_name='snippets' AND index_name='idx_snippets_created'").Scan(&indexExists)
	if err != nil {
		return err
	}

	if indexExists == 0 {
		_, err = db.Exec("CREATE INDEX idx_snippets_created ON snippets(created)")
		if err != nil {
			return err
		}
		log.Println("Created index on snippets.created")
	}

	// Check if snippets table is empty before inserting sample data
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM snippets").Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		// Insert sample data only if table is empty - each insert as a separate statement
		_, err = db.Exec(`
            INSERT INTO snippets (title, content, created, expires) VALUES (
                'An old silent pond',
                'An old silent pond...\nA frog jumps into the pond,\nsplash! Silence again.\n\n– Matsuo Bashō',
                UTC_TIMESTAMP(),
                DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY)
            )
        `)
		if err != nil {
			return err
		}

		_, err = db.Exec(`
            INSERT INTO snippets (title, content, created, expires) VALUES (
                'Over the wintry forest',
                'Over the wintry\nforest, winds howl in rage\nwith no leaves to blow.\n\n– Natsume Soseki',
                UTC_TIMESTAMP(),
                DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY)
            )
        `)
		if err != nil {
			return err
		}

		_, err = db.Exec(`
            INSERT INTO snippets (title, content, created, expires) VALUES (
                'First autumn morning',
                'First autumn morning\nthe mirror I stare into\nshows my father''s face.\n\n– Murakami Kijo',
                UTC_TIMESTAMP(),
                DATE_ADD(UTC_TIMESTAMP(), INTERVAL 7 DAY)
            )
        `)
		if err != nil {
			return err
		}
		log.Println("Inserted sample data into snippets table")
	}

	// Create sessions table if it doesn't exist
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS sessions (
            token CHAR(43) PRIMARY KEY,
            data BLOB NOT NULL,
            expiry TIMESTAMP(6) NOT NULL
        )
    `)
	if err != nil {
		return err
	}

	// Check if sessions index exists before creating it
	err = db.QueryRow("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema=DATABASE() AND table_name='sessions' AND index_name='sessions_expiry_idx'").Scan(&indexExists)
	if err != nil {
		return err
	}

	if indexExists == 0 {
		_, err = db.Exec("CREATE INDEX sessions_expiry_idx ON sessions(expiry)")
		if err != nil {
			return err
		}
		log.Println("Created index on sessions.expiry")
	}

	// Create users table if it doesn't exist
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
            name VARCHAR(255) NOT NULL,
            email VARCHAR(255) NOT NULL,
            hashed_password CHAR(60) NOT NULL,
            created DATETIME NOT NULL,
            CONSTRAINT users_uc_email UNIQUE (email)
        )
    `)
	if err != nil {
		return err
	}

	log.Println("Database schema initialized successfully")
	return nil
}
