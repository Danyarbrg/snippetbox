package main

import (
	"errors"
	"fmt"
	//"html/template"
	"net/http"
	"strconv"
	//"unicode/utf8"
	//"strings"

	"github.com/Danyarbrg/snippetbox/pkg/models"
	"github.com/Danyarbrg/snippetbox/pkg/forms"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// in pat the "/" handle only "/" and we can remove this part
	//if r.URL.Path != "/" {
	//	// http.NotFound(w, r) -- до поялвения хелперов helpers
	//	app.notFound(w)
	//	return
	//}

	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "home.page.tmpl", &templateData{
		Snippets: s,
	})
	//data := &templateData{Snippets: s}

	// инициализируем слайс содержащий пути к файлам шаблонов HTML
	// home.page.tmpl должен быть первым в слайсе
	//files := []string{
	//	"./ui/html/home.page.tmpl",
	//	"./ui/html/base.layout.tmpl",
	//	"./ui/html/footer.partial.tmpl",
	//}

	// используем функцию template.ParseFiles() для чтения шаблонного файла 
	// если ошибка - логируем ошибку и шлем код 500 через http.Error()
	//ts, err := template.ParseFiles(files...)
	//if err != nil {
	//	// app.errorLog.Println(err.Error()) -- до поялвения хелперов helpers
	//	// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	//	app.serverError(w, err)
	//	return
	//}

	//err = ts.Execute(w, data)
	//if err != nil {
	//	// app.errorLog.Println(err.Error()) -- до поялвения хелперов helpers
	//	// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	//	app.serverError(w, err)
	//}
}

// with pat we need to change "id" on ":id"
func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		// http.NotFound(w, r) -- до поялвения хелперов helpers
		app.notFound(w)
		return
	}

	s, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
	}

	// use the PopString() method to retrieve the value for the "flash" key.
	// PopString() also deletes the key and value from the session data, so it
	// acts like a one-time fetch. If there is no matching key in the session
	// data this will return the empty string.
	// flash := app.session.PopString(r, "flash")

	app.render(w, r, "show.page.tmpl", &templateData{
		// Flash:		flash,
		Snippet:	s,
	})
}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

// with pat we can remove POST checker, because we showed method in routes.go
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	//if r.Method != http.MethodPost {
	//	w.Header().Set("Allow", http.MethodPost)
	//	// http.Error(w, "Nethod not allowed...", http.StatusMethodNotAllowed) -- до поялвения хелперов helpers
	//	app.clientError(w, http.StatusMethodNotAllowed)
	//	return
	//}

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	/* before template Form created
	title 	:= r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expires := r.PostForm.Get("expires")

	errors := make(map[string]string)

	if strings.TrimSpace(title) == "" {
		errors["title"] = "This field cant be blank"
	} else if utf8.RuneCountInString(title) > 100 {
		errors["title"] = "This field cant contain more than 100 characters"
	}

	if strings.TrimSpace(content) == "" {
		errors["contenct"] = "This field cant be blank"
	}

	if strings.TrimSpace(expires) == "" {
		errors["expires"] = "This field cant be blank"
	} else if expires != "365" && expires != "7" && expires != "1" {
		errors["expires"] = "This field is invalid"
	}

	if len(errors) > 0 {
		app.render(w, r, "create.page.tmpl", &templateData{
			FormErrors:	errors,
			FormData:	r.PostForm,
		})
		return
	}
	*/

	// creating new forms.Form struct containing the POSTed data from the form.
	form := forms.New(r.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")

	// if form are not valid, redisplay the template in te form.Form object as the data
	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{Form: form})
		return
	}

	// толкаем данные в SnippetModel.Insert(), получая ID обратно
	// before template Form created
	// id, err := app.snippets.Insert(title, content, expires)
	id, err := app.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	// use the Put() method to add a string value ("Your snippet was saved successfully")
	// and the corresponding key ("flash") to the session data.
	// if the user have not existing session, thew will automatically be created by the 
	// session middleware.
	app.session.Put(r, "flash", "Snippet successfully created!")

	//Перенаправить пользователя на соответствующую страницу для фрагмента.
	// change redirect to use new semantic URL style of /snippet/:id
	//http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Create new user...")
}

func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Display the user login form")
}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Authenticate and login user...")
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Logout the user...")
}