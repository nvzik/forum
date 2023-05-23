package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "main.db")
	if err != nil {
		log.Println(err)
		return
	}

	err = db.Ping()
	if err != nil {
		log.Println(err)
		return
	}

	// stmt := `DELETE FROM sqlite_sequence`

	// stmt := `ALTER TABLE sessions TO `

	// stmt := `CREATE TABLE reactComments (
	// 	commentID INTEGER,
	// 	userID INTEGER,
	// 	reactionsComment INTEGER,
	// 	FOREIGN KEY (commentID) REFERENCES comments (commentID),
	// 	FOREIGN KEY (userID) REFERENCES users (userID)
	// );`

	// stmt := `
	// CREATE TABLE users (
	// 	userID INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	// 	name VARCHAR(255) NOT NULL,
	// 	email VARCHAR(255) NOT NULL UNIQUE	,
	// 	hashed_password CHAR(60) NOT NULL
	// 	);`

	stmt := `DROP TABLE sessions;
			CREATE TABLE sessions (
			sessionID INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			userID INTEGER NOT NULL UNIQUE,
			token VARCHAR(60) NOT NULL,
			expiry DATETIME,
			FOREIGN KEY (userID) REFERENCES users (userID)
			);
	`
	// stmt := `	CREATE TABLE snippets (
	// 	snippetID INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	// 	userID INTEGER NOT NULL,
	// 	title VARCHAR(100) NOT NULL,
	// 	content TEXT NOT NULL,
	// 	created DATETIME DEFAULT CURRENT_TIMESTAMP,
	// 	backend BOOL,
	// 	frontend BOOL,
	// 	fullstack BOOL,
	// 	FOREIGN KEY (userID) REFERENCES users (userID)
	// 	);
	// 	`

	// stmt := `CREATE TABLE comments (
	// 	commentID INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	// 	snippetID INTEGER NOT NULL,
	// 	userID INTEGER NOT NULL,
	// 	parentID INTEGER NOT NULL,
	// 	Content TEXT NOT NULL,
	// 	FOREIGN KEY (snippetID) REFERENCES snippets (snippetID),
	// 	FOREIGN KEY (userID) REFERENCES snippets (userID),
	// 	FOREIGN KEY (parentID) REFERENCES snippets (parentID)
	// 	);`

	_, err = db.Exec(stmt)

	if err != nil {
		log.Println(err)
		return
	}
}
