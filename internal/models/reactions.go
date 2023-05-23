package models

import (
	"database/sql"
	"errors"
)

type ReactionModel struct {
	DB *sql.DB
}

func (r *ReactionModel) LikePost(userID, postID int) error {
	stmt := `SELECT (reactions) FROM reactions WHERE snippetID = ? AND userID = ?`

	var reaction int
	err := r.DB.QueryRow(stmt, postID, userID).Scan(&reaction)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_, err = r.DB.Exec(`INSERT INTO reactions (snippetID, userID, reactions) VALUES (?, ?, ?)`, postID, userID, 1)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	if reaction == -1 {
		_, err := r.DB.Exec(`UPDATE reactions SET reactions = ? WHERE snippetID = ? AND userID = ?`, 1, postID, userID)
		if err != nil {
			return err
		}
	} else if reaction == 1 {
		_, err := r.DB.Exec(`DELETE FROM reactions WHERE snippetID = ? AND userID = ?`, postID, userID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *ReactionModel) DislikePost(userID, postID int) error {
	stmt := `SELECT (reactions) FROM reactions WHERE snippetID = ? AND userID = ?`

	var reaction int
	err := r.DB.QueryRow(stmt, postID, userID).Scan(&reaction)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_, err = r.DB.Exec(`INSERT INTO reactions (snippetID, userID, reactions) VALUES (?, ?, ?)`, postID, userID, -1)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	if reaction == 1 {
		_, err := r.DB.Exec(`UPDATE reactions SET reactions = ? WHERE snippetID = ? AND userID = ?`, -1, postID, userID)
		if err != nil {
			return err
		}
	} else if reaction == -1 {
		_, err := r.DB.Exec(`DELETE FROM reactions WHERE snippetID = ? AND userID = ?`, postID, userID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *ReactionModel) GetLikeByUserID(postID, userID int) bool {
	stmt := `SELECT (reactions) FROM reactions WHERE snippetID = ? AND userID = ?`

	var reaction int
	err := r.DB.QueryRow(stmt, postID, userID).Scan(&reaction)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false
		} else {
			return false
		}
	}

	if reaction == 1 {
		return true
	}
	return false
}

func (r *ReactionModel) GetDislikeByUserID(postID, userID int) bool {
	stmt := `SELECT (reactions) FROM reactions WHERE snippetID = ? AND userID = ?`

	var reaction int
	err := r.DB.QueryRow(stmt, postID, userID).Scan(&reaction)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false
		} else {
			return false
		}
	}

	if reaction == -1 {
		return true
	}
	return false
}
