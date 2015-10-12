package guestbook

import (
        "html/template"
        "net/http"
        "time"
	"strconv"

        "appengine"
        "appengine/datastore"
        "appengine/user"
)

var templates = template.Must(template.ParseFiles(
	"guestbookform.html",
	"view.html",
	"edit.html",
))

type Greeting struct {
	Author   string
        FirstName  string
	LastName  string
	Class	string
	Grad	string
        Content string
        Date    time.Time
}

type Return struct {
        Key     *datastore.Key
	ID	int64
	Data	interface{}
}

type Post struct {
    Title string
    Content interface{}
}

func init() {
    http.HandleFunc("/", rootHandler)
    http.HandleFunc("/sign", sign)
    http.HandleFunc("/edit/", editHandler)
    http.HandleFunc("/view", view)
}

func renderTemplate(w http.ResponseWriter, tmpl string, r *Return) {
    err := templates.ExecuteTemplate(w, tmpl, r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func renderSliceTemplate(w http.ResponseWriter, tmpl string, r []Return) {
    err := templates.ExecuteTemplate(w, tmpl, r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
        c := appengine.NewContext(r)
        q := datastore.NewQuery("Greeting").Order("-Date").Limit(10)
        greetings := make([]Greeting, 0, 10)
        keys, err := q.GetAll(c, &greetings)
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }

        res := make([]Return, 0, 10)
        for i, r := range greetings {
            k := keys[i]
            y := Return{
                Key: k,
                ID: k.IntID(),
                Data: r,
            }
            res = append(res, y)
        }
        renderSliceTemplate(w, "guestbookform.html", res)

        if err != nil {
            panic(err)
        }
}

func sign(w http.ResponseWriter, r *http.Request) {
        c := appengine.NewContext(r)
        r.ParseForm()
        g := Greeting{
		FirstName: r.FormValue("fname"),
		LastName: r.FormValue("lname"),
                Content: r.FormValue("content"),
		Class: r.FormValue("class"),
		Grad: r.FormValue("grad"),
                Date:    time.Now(),
        }
	if r.FormValue("grad")!= ""{
                g.Grad = "Yes"
        }else{
		g.Grad="No"
	}
        if u := user.Current(c); u != nil {
                g.Author = u.String()
        }

        key := datastore.NewKey(c, "Greeting", "", 0, nil)
        _, err := datastore.Put(c, key, &g)
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }
        http.Redirect(w, r, "/", http.StatusFound)
}

func view(w http.ResponseWriter, r *http.Request) {
        c := appengine.NewContext(r)
	q := datastore.NewQuery("Greeting").Order("-Date").Limit(10)
        greetings := make([]Greeting, 0, 10)

        keys, err := q.GetAll(c, &greetings)
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        } 

	res := make([]Return, 0, 10)
        for i, r := range greetings {
            k := keys[i]
            y := Return{
                Key: k,
                ID: k.IntID(),
                Data: r,
            }
            res = append(res, y)
        }
        renderSliceTemplate(w, "view.html", res)

        if err != nil {
            panic(err)
        }
}

func editHandler(w http.ResponseWriter, r *http.Request) {

    c := appengine.NewContext(r)
    pageID, _ := strconv.ParseInt(r.URL.Path[len("/edit/"):], 10, 64)
    pageKey := datastore.NewKey(c, "Greeting", "", pageID, nil)

    if r.Method == "GET" {
        var result Greeting
        err := datastore.Get(c, pageKey, &result)

            res := Return{
            Key: pageKey,
            ID: pageID,
            Data: result,
            }
    
    //error page for entity that doesn't exist
        if err != nil {
            renderTemplate(w, "edit.html", &res)
            return
        }

        renderTemplate(w, "edit.html", &res)
    }else if r.Method == "POST" {
        //get data form from and build Greeting struct
        r.ParseForm()
        g := Greeting{
                FirstName: r.FormValue("fname"),
                LastName: r.FormValue("lname"),
                Content: r.FormValue("content"),
                Class: r.FormValue("class"),
                Grad: r.FormValue("grad"),
                Date:    time.Now(),
        }
        //Save the data to the DB
        //Since Key is the same, an update will occur
        //Reload page for "GET" request
        _, err := datastore.Put(c, pageKey, &g)
        if err != nil {
        } else {
            http.Redirect(w, r, "/", http.StatusFound)
        }
        }
    }

