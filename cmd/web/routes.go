package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/post/view", app.snippetView)
	mux.HandleFunc("/post/create", app.requireAuth(app.snippetCreate))
	mux.HandleFunc("/post/myposts", app.requireAuth(app.mysnippets))
	mux.HandleFunc("/post/likedposts", app.requireAuth(app.likedsnippets))
	mux.HandleFunc("/user/signup", app.userSignup)
	mux.HandleFunc("/user/login", app.userLogin)
	mux.HandleFunc("/user/logout", app.requireAuth(app.userLogout))
	// mux.HandleFunc("/filter", app.filterPosts)
	mux.HandleFunc("/likePost", app.requireAuth(app.likePost))
	mux.HandleFunc("/dislikePost", app.requireAuth(app.dislikePost))
	mux.HandleFunc("/likeComment", app.requireAuth(app.likeComment))
	mux.HandleFunc("/dislikeComment", app.requireAuth(app.dislikeComment))
	mux.HandleFunc("/error", app.userLogin)
	return app.checkSession(app.recoverPanic(app.logRequest(secureHeaders(mux))))
}
