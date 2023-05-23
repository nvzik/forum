package models

import (
	"database/sql"
	"errors"
)

type CommentModel struct {
	DB *sql.DB
}

func (r *CommentModel) Insert(c *Comment) error {
	stmt := `INSERT INTO comments (snippetID, userID, parentID, content)
	VALUES(?, ?, ?, ?)`

	_, err := r.DB.Exec(stmt, c.SnippetID, c.UserID, c.ParentID, c.Content)
	if err != nil {
		return err
	}
	return nil
}

func (r *CommentModel) CommentsByPostID(postID, userID int) ([]*Comment, error) {
	stmt := `SELECT commentID, snippetID, userID, content, (SELECT name FROM users WHERE comments.userID = users.userID) FROM comments
	WHERE snippetID = ?`
	comments := []*Comment{}

	rows, err := r.DB.Query(stmt, postID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		c := &Comment{}
		err = rows.Scan(&c.ID, &c.SnippetID, &c.UserID, &c.Content, &c.Username)
		if err != nil {
			return nil, err
		}
		c.Liked = r.GetLikeComByUserID(c.ID, userID)
		c.Disliked = r.GetDislikeComByUserID(c.ID, userID)

		stmt = `SELECT (reactionsComment) FROM reactComments WHERE commentID = ?`
		rowsCom, err := r.DB.Query(stmt, c.ID)
		if err != nil {
			return nil, err
		}
		defer rowsCom.Close()

		for rowsCom.Next() {
			var a int
			err = rowsCom.Scan(&a)
			if err != nil {
				return nil, err
			}

			c.Reactions += a
		}

		if err = rowsCom.Err(); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	return comments, nil
}

func (r *CommentModel) GetLikeComByUserID(commentID, userID int) bool {
	stmt := `SELECT (reactionComment) FROM reactComments WHERE commentID = ? AND userID = ?`

	var reaction int
	err := r.DB.QueryRow(stmt, commentID, userID).Scan(&reaction)
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

func (r *CommentModel) GetDislikeComByUserID(commentID, userID int) bool {
	stmt := `SELECT (reactionComment) FROM reactComments WHERE snippetID = ? AND userID = ?`

	var reaction int
	err := r.DB.QueryRow(stmt, commentID, userID).Scan(&reaction)
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
