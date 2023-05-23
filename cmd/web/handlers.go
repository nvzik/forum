package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gofrs/uuid"
	"github.com/tlsh0/internal/models"
	"github.com/tlsh0/internal/validator"
)

type snippetCreateForm struct {
	Title     string
	Content   string
	Backend   string
	Frontend  string
	Fullstack string
	validator.Validator
}

type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

// homepage
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	status := r.Context().Value("data").(*contextData)
	if r.URL.Path != "/" {
		// app.notFound(w)
		app.errorHandler(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	data := app.newTemplateData(r)
	switch r.Method {
	case http.MethodGet:
		snippets, err := app.snippets.Latest()
		if err != nil {
			app.errorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		data.IsAuthenticated = status.IsAuthenticated
		data.Snippets = snippets
		app.render(w, http.StatusOK, "home.html", data)
	case http.MethodPost:
		result := []*models.Snippet{}
		backend, err := strconv.Atoi(r.FormValue("backend"))
		if err != nil {
			backend = 0
		}
		frontend, err := strconv.Atoi(r.FormValue("frontend"))
		if err != nil {
			frontend = 0
		}
		fullstack, err := strconv.Atoi(r.FormValue("fullstack"))
		if err != nil {
			fullstack = 0
		}
		snippets, err := app.snippets.Latest()
		if err != nil {
			app.errorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		if backend == 0 && frontend == 0 && fullstack == 0 {
			data.Snippets = snippets
		} else {
			for i := range snippets {
				switch {
				case snippets[i].Backend == backend && backend != 0:
					result = append(result, snippets[i])
				case snippets[i].Frontend == frontend && frontend != 0:
					result = append(result, snippets[i])
				case snippets[i].Fullstack == fullstack && fullstack != 0:
					result = append(result, snippets[i])
				}
			}
			data.Snippets = result
		}
		data.IsAuthenticated = status.IsAuthenticated
		app.render(w, http.StatusOK, "home.html", data)
	default:
		app.errorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		data.IsAuthenticated = status.IsAuthenticated
		return
	}
}

func (app *application) errorHandler(w http.ResponseWriter, r *http.Request, errorNum int, errorDescript string) {
	status := r.Context().Value("data").(*contextData)
	data := app.newTemplateData(r)
	data.CodeError = errorNum
	data.ErrorDescription = errorDescript
	data.IsAuthenticated = status.IsAuthenticated
	app.render(w, errorNum, "error.html", data)
}

// onepost
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	status := r.Context().Value("data").(*contextData)
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.errorHandler(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	data := app.newTemplateData(r)

	switch r.Method {
	case http.MethodGet:
		snippet, err := app.snippets.Get(id)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.errorHandler(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
				return
			} else {
				app.errorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			}
			return
		}
		data.IsAuthenticated = status.IsAuthenticated
		data.Snippet = snippet
		data.LikedPost = app.reactions.GetLikeByUserID(id, status.UserID)
		data.DislikedPost = app.reactions.GetDislikeByUserID(id, status.UserID)
		comments, err := app.comments.CommentsByPostID(snippet.ID, status.UserID)
		if err != nil && !errors.Is(err, models.ErrNoRecord) {
			app.errorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		data.Comments = comments
		app.render(w, http.StatusOK, "view.html", data)
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			app.errorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		comment := &models.Comment{
			SnippetID: id,
			ParentID:  id,
			UserID:    status.UserID,
			Content:   r.PostForm.Get("content"),
		}

		app.comments.Insert(comment)
		http.Redirect(w, r, fmt.Sprintf("/post/view?id=%d", id), http.StatusSeeOther)
	default:
		app.errorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
}

// posts
func (app *application) likePost(w http.ResponseWriter, r *http.Request) {
	status := r.Context().Value("data").(*contextData)
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.errorHandler(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	_ = app.reactions.LikePost(status.UserID, id)
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

func (app *application) dislikePost(w http.ResponseWriter, r *http.Request) {
	status := r.Context().Value("data").(*contextData)
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.errorHandler(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	_ = app.reactions.DislikePost(status.UserID, id)
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

// comments
func (app *application) likeComment(w http.ResponseWriter, r *http.Request) {
	status := r.Context().Value("data").(*contextData)
	commentID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || commentID < 1 {
		app.errorHandler(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	_ = app.commentReactions.LikeComment(status.UserID, commentID)
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

func (app *application) dislikeComment(w http.ResponseWriter, r *http.Request) {
	status := r.Context().Value("data").(*contextData)
	commentID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || commentID < 1 {
		app.errorHandler(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	_ = app.commentReactions.DislikeComment(status.UserID, commentID)
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

// posts by creator
func (app *application) mysnippets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		app.errorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	status := r.Context().Value("data").(*contextData)
	if r.URL.Path != "/post/myposts" {
		app.errorHandler(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	userSnippets, err := app.snippets.GetSnippetsByUserID(status.UserID)
	if err != nil {
		app.errorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	data := app.newTemplateData(r)
	data.IsAuthenticated = status.IsAuthenticated
	data.Snippets = userSnippets
	app.render(w, http.StatusOK, "mysnippets.html", data)
}

// likedsnippets
func (app *application) likedsnippets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		app.errorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	status := r.Context().Value("data").(*contextData)
	if r.URL.Path != "/post/likedposts" {
		app.errorHandler(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	userSnippets, err := app.snippets.GetLikedSnippetsByUserID(status.UserID)
	if err != nil {
		app.errorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	data := app.newTemplateData(r)
	data.IsAuthenticated = status.IsAuthenticated
	data.Snippets = userSnippets
	app.render(w, http.StatusOK, "likedsnippets.html", data)
}

// createsnippet
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	status := r.Context().Value("data").(*contextData)
	data := app.newTemplateData(r)
	switch r.Method {
	case http.MethodGet:
		data.IsAuthenticated = status.IsAuthenticated
		data.Form = snippetCreateForm{}
		app.render(w, http.StatusOK, "create.html", data)
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			app.errorHandler(w, r, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}
		form := snippetCreateForm{
			Title:     r.PostForm.Get("title"),
			Content:   r.PostForm.Get("content"),
			Backend:   r.PostForm.Get("backend"),
			Frontend:  r.PostForm.Get("frontend"),
			Fullstack: r.PostForm.Get("fullstack"),
		}
		form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
		form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characteres long")
		form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
		form.CheckField(validator.AtLeastOne(form.Backend, form.Frontend, form.Fullstack), "category", "At least one category must be checked")

		if !form.Valid() {
			data = app.newTemplateData(r)
			data.Form = form
			data.IsAuthenticated = status.IsAuthenticated
			app.render(w, http.StatusUnprocessableEntity, "create.html", data)
			return
		}

		backend, err := strconv.Atoi(form.Backend)
		if err != nil {
			backend = 0
		}
		frontend, err := strconv.Atoi(form.Frontend)
		if err != nil {
			frontend = 0
		}
		fullstack, err := strconv.Atoi(form.Fullstack)
		if err != nil {
			fullstack = 0
		}

		snippet := &models.Snippet{
			UserID:    status.UserID,
			Title:     form.Title,
			Content:   form.Content,
			Created:   time.Now(),
			Backend:   backend,
			Frontend:  frontend,
			Fullstack: fullstack,
		}

		id, err := app.snippets.Insert(snippet)
		if err != nil {
			app.errorHandler(w, r, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}

		route := fmt.Sprintf("/post/view?id=%d", id)
		data.IsAuthenticated = status.IsAuthenticated
		data.Form = form

		http.Redirect(w, r, route, http.StatusSeeOther)
	default:
		app.errorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		data.IsAuthenticated = status.IsAuthenticated
		return
	}
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		data := app.newTemplateData(r)
		data.Form = userSignupForm{}
		app.render(w, http.StatusOK, "signup.html", data)
	case http.MethodPost:
		var form userSignupForm

		err := r.ParseForm()
		if err != nil {
			app.errorHandler(w, r, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}

		form = userSignupForm{
			Name:     r.PostForm.Get("name"),
			Email:    r.PostForm.Get("email"),
			Password: r.PostForm.Get("password"),
		}

		form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
		form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
		form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
		form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
		form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")
		form.CheckField(validator.SpecialChars(form.Password), "password", "The password must contain at least of these characters: !@#$%^&*()_+-=")

		if !form.Valid() {
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signup.html", data)
			return
		}

		err = app.users.Insert(form.Name, form.Email, form.Password)
		if err != nil {
			if errors.Is(err, models.ErrDuplicateEmail) {
				form.AddFieldError("email", "Email address is already in use")

				data := app.newTemplateData(r)
				data.Form = form
				app.render(w, http.StatusUnprocessableEntity, "signup.html", data)
			} else {
				app.errorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			}

			return
		}

		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
	default:
		app.errorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		data := app.newTemplateData(r)
		data.Form = userLoginForm{}
		app.render(w, http.StatusOK, "login.html", data)
	case http.MethodPost:

		err := r.ParseForm()
		if err != nil {
			app.errorHandler(w, r, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}

		form := userSignupForm{
			Email:    r.PostForm.Get("email"),
			Password: r.PostForm.Get("password"),
		}

		form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
		form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
		form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

		if !form.Valid() {
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "login.html", data)
			return
		}

		id, err := app.users.Authenticate(form.Email, form.Password)
		if err != nil {
			if errors.Is(err, models.ErrInvalidCredentials) {
				form.AddNonFieldError("Email or password is incorrect")

				data := app.newTemplateData(r)
				data.Form = form
				app.render(w, http.StatusUnprocessableEntity, "login.html", data)
			} else {
				app.errorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			}
			return
		}
		
		token := uuid.Must(uuid.NewV4()).String()
		newCookie := &http.Cookie{
			Name:  "token",
			Value: token,
			Path:  "/",
		}
		http.SetCookie(w, newCookie)
		err = app.sessions.StartSession(id, token)
		if err != nil {
			app.errorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)

	default:
		app.errorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
}

func (app *application) userLogout(w http.ResponseWriter, r *http.Request) {
	status := r.Context().Value("data").(*contextData)
	err := app.sessions.StopSession(status.Token)
	if err != nil {
		app.errorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
