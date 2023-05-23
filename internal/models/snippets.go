package models

import (
	"database/sql"
	"errors"
	"time"
)

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(s *Snippet) (int, error) {
	stmt := `INSERT INTO snippets (userID, title, content, created, backend, frontend, fullstack) VALUES(?, ?, ?, ?, ?, ?, ?)`

	result, err := m.DB.Exec(stmt, s.UserID, s.Title, s.Content, s.Created, s.Backend, s.Frontend, s.Fullstack)
	if err != nil {
		return 0, nil
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	stmt := `SELECT snippetID, (SELECT name FROM users WHERE users.userID = snippets.userID), title, content, created, backend, frontend, fullstack FROM snippets
	WHERE snippetID = ?`

	row := m.DB.QueryRow(stmt, id)

	s := &Snippet{}
	err := row.Scan(&s.ID, &s.Username, &s.Title, &s.Content, &s.Created, &s.Backend, &s.Frontend, &s.Fullstack)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	stmt = `SELECT (reactions) FROM reactions WHERE snippetID = ?`

	rows, err := m.DB.Query(stmt, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var a int
		err = rows.Scan(&a)
		if err != nil {
			return nil, err
		}

		s.Reactions += a
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return s, nil
}

func (m *SnippetModel) Latest() ([]*Snippet, error) {
	stmt := `SELECT snippetID,  (SELECT name FROM users WHERE users.userID = snippets.userID), title, content, created, backend, frontend, fullstack, (SELECT COUNT(*) FROM comments where comments.snippetID=snippets.snippetID) FROM snippets
	WHERE created < ? ORDER BY snippetID DESC`

	rows, err := m.DB.Query(stmt, time.Now().Local())
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	snippets := []*Snippet{}

	for rows.Next() {
		s := &Snippet{}
		err = rows.Scan(&s.ID, &s.Username, &s.Title, &s.Content, &s.Created, &s.Backend, &s.Frontend, &s.Fullstack, &s.SumComments)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return snippets, nil
}

func (m *SnippetModel) GetSnippetsByUserID(userID int) ([]*Snippet, error) {
	stmt := `SELECT snippetID, (SELECT name FROM users WHERE users.userID = snippets.userID), title, content, created, backend, frontend, fullstack, (SELECT COUNT(*) FROM comments where comments.snippetID=snippets.snippetID) FROM snippets WHERE userID = ? ORDER BY snippetID DESC`

	rows, err := m.DB.Query(stmt, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	snippets := []*Snippet{}

	for rows.Next() {
		s := &Snippet{}
		err = rows.Scan(&s.ID, &s.Username, &s.Title, &s.Content, &s.Created, &s.Backend, &s.Frontend, &s.Fullstack, &s.SumComments)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return snippets, nil
}

func (m *SnippetModel) GetLikedSnippetsByUserID(userID int) ([]*Snippet, error) {
	stmt := `SELECT *,(SELECT COUNT(*) FROM comments where comments.snippetID=snippets.SnippetID) FROM snippets INNER JOIN reactions ON reactions.snippetID=snippets.snippetID WHERE reactions.userID = ? AND reactions=1  ORDER BY snippetID DESC`
	rows, err := m.DB.Query(stmt, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	snippets := []*Snippet{}

	for rows.Next() {
		s := &Snippet{}
		err = rows.Scan(&s.ID, &s.Username, &s.Title, &s.Content, &s.Created, &s.Backend, &s.Frontend, &s.Fullstack, &s.Reactions, &s.snippetID, &s.userID, &s.SumComments)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// for i := range snippets {
	// 	stmt = `SELECT COUNT(*) FROM comments where comments.snippetID=?`
	// 	rows, err = m.DB.Query(stmt, snippets[i].snippetID)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	defer rows.Close()

	// 	for rows.Next() {
	// 		err = rows.Scan(&snippets[i].SumComments)
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 	}
	// }

	return snippets, nil
}
