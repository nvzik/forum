package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tlsh0/internal/models"
)

type application struct {
	errorLog         *log.Logger
	infoLog          *log.Logger
	snippets         *models.SnippetModel
	users            *models.UserModel
	sessions         *models.SessionModel
	reactions        *models.ReactionModel
	commentReactions *models.ReactionCommentModel
	templateCache    map[string]*template.Template
	comments         *models.CommentModel
}

func main() {
	addr := flag.String("addr", ":8080", "HTTP network address")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := sql.Open("sqlite3", "main.db")
	if err != nil {
		errorLog.Println(err)
	}

	if err = db.Ping(); err != nil {
		errorLog.Println(err)
	}

	defer db.Close()

	// deleteing expired sessions
	go func() {
		for {
			err := models.DeleteExpiredSession(db)
			if err != nil {
				errorLog.Println(err)
			}
			time.Sleep(1 * time.Minute)
		}
	}()

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application{
		errorLog:         errorLog,
		infoLog:          infoLog,
		snippets:         &models.SnippetModel{DB: db},
		users:            &models.UserModel{DB: db},
		sessions:         &models.SessionModel{DB: db},
		reactions:        &models.ReactionModel{DB: db},
		commentReactions: &models.ReactionCommentModel{DB: db},
		templateCache:    templateCache,
		comments:         &models.CommentModel{DB: db},
	}

	srv := &http.Server{
		Addr:           *addr,
		MaxHeaderBytes: 524288,
		ErrorLog:       errorLog,
		Handler:        app.routes(),
		IdleTimeout:    time.Minute,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	infoLog.Printf("Starting server on %s", "http://localhost"+*addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
