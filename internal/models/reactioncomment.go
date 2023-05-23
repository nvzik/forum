package models

import (
	"database/sql"
	"errors"
)

type ReactionCommentModel struct {
	DB *sql.DB
}

func (r *ReactionCommentModel) LikeComment(userID, commentID int) error {
	stmt := `SELECT (reactionsComment) FROM reactComments WHERE commentID = ? AND userID = ?`

	var reaction int
	err := r.DB.QueryRow(stmt, commentID, userID).Scan(&reaction)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_, err = r.DB.Exec(`INSERT INTO reactComments (commentID, userID, reactionsComment) VALUES (?, ?, ?)`, commentID, userID, 1)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	if reaction == -1 {
		_, err := r.DB.Exec(`UPDATE reactComments SET reactionsComment = ? WHERE commentID = ? AND userID = ?`, 1, commentID, userID)
		if err != nil {
			return err
		}
	} else if reaction == 1 {
		_, err := r.DB.Exec(`DELETE FROM reactComments WHERE commentID = ? AND userID = ?`, commentID, userID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *ReactionCommentModel) DislikeComment(userID, commentID int) error {
	stmt := `SELECT (reactionsComment) FROM reactComments WHERE commentID = ? AND userID = ?`

	var reaction int
	err := r.DB.QueryRow(stmt, commentID, userID).Scan(&reaction)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_, err = r.DB.Exec(`INSERT INTO reactComments (commentID, userID, reactionsComment) VALUES (?, ?, ?)`, commentID, userID, -1)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	if reaction == 1 {
		_, err := r.DB.Exec(`UPDATE reactComments SET reactionsComment = ? WHERE commentID = ? AND userID = ?`, -1, commentID, userID)
		if err != nil {
			return err
		}
	} else if reaction == -1 {
		_, err := r.DB.Exec(`DELETE FROM reactComments WHERE commentID = ? AND userID = ?`, commentID, userID)
		if err != nil {
			return err
		}
	}
	return nil
}
