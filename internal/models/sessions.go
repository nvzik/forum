package models

import (
	"database/sql"
	"errors"
	"time"
)

type SessionModel struct {
	DB *sql.DB
}

func (m *SessionModel) StartSession(userID int, token string) error {
	stmt := `INSERT OR REPLACE INTO sessions (userID, token, expiry) VALUES (?, ?, ?)`

	_, err := m.DB.Exec(stmt, userID, token, time.Now().Add(time.Minute*30))
	if err != nil {
		return err
	}
	return nil
}

func (m *SessionModel) FindSession(token string) (*Session, error) {
	stmt := `SELECT sessionID, userID, token, expiry FROM sessions WHERE token = ?`

	s := &Session{}

	err := m.DB.QueryRow(stmt, token).Scan(&s.ID, &s.UserID, &s.Token, &s.Expiry)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

func (m *SessionModel) FindUserID(userID int) (int, error) {
	stmt := `SELECT COUNT(*) FROM sessions where userID = ?`
	rows, err := m.DB.Query(stmt, userID)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	count := 0
	for rows.Next() {
		count = 0
		err = rows.Scan(&count)
		if err != nil {
			return 0, err
		}
		count++
	}
	return count, nil
}

func (m *SessionModel) StopSession(token string) error {
	stmt := `DELETE FROM sessions WHERE token = ?`

	_, err := m.DB.Exec(stmt, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNoRecord
		} else {
			return err
		}
	}
	return nil
}

func DeleteExpiredSession(db *sql.DB) error {
	if _, err := db.Exec("DELETE FROM sessions WHERE expiry < DATETIME('now')"); err != nil {
		return err
	}
	return nil
}
