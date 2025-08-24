package main

import (
	"net/http"

	"github.com/justinas/alice"
	"github.com/bmizerany/pat"
)


// updated. instead of *http.ServeMux now we use http.Handler
func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureMiddleware)
	
	// creating new middleware chain containing the middleware to
	// our dinamic application routes.
	dynamicMiddleware := alice.New(app.session.Enable, noSurf)

	// changing after pat
	// mux := http.NewServeMux()
	mux := pat.New()
	// if we working without justinas/alice package then we should write like this:
	// mux.Get("/", app.session.Enable(http.HandlerFunc(app.home)))
	// mux.HandleFunc("/", app.home)
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	// mux.HandleFunc("/snippet/create", app.createSnippet)
	mux.Get("/snippet/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createSnippetForm))
	mux.Post("/snippet/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createSnippet))
	// mux.HandleFunc("/snippet", app.showSnippet)
	mux.Get("/snippet/:id", dynamicMiddleware.ThenFunc(app.showSnippet))

	// user sign routes
	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Post("/user/logout", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.logoutUser))

	
	// создаем сервер файл котороый ищет файлы в указанной директории
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	// Без StripPrefix запрос /static/css/main.css попытается искать файл в:
	// text ./ui/static/static/css/main.css | Лишний /static
	// mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))
	

	// updated. instead of mux now we use secureHeaders(mux)
	return standardMiddleware.Then(mux)
}