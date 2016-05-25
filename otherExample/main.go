/*
  This exmaple is using "raw" go, no gin. This requires a bit more work but not much else is different
*/

package main

import (
	"html/template"
	"net/http"
)

func myquery1(val string) string {
	// test data
	return val
}

func myquery2() []string {
	return []string{"Ned", "Caetlyn", "Rob", "Ygritte", "Osha", "Hodor"}
}

func form(w http.ResponseWriter, r *http.Request) {
	// get the template html file
	form, _ := template.ParseFiles("form.html")

	// make sure that the method sent here is a POST
	switch r.Method {
	case "POST":
		r.ParseForm()
		// r.FormValue(..) gets the value that was sent in with the post request
		result := myquery1(r.FormValue("zip"))
		// put in the value from result in the place of {{ . }}
		form.Execute(w, result)
	default:
		form.Execute(w, "NO POST DATA")
	}
}

func list(w http.ResponseWriter, r *http.Request) {
	page, _ := template.ParseFiles("list.html")

	betrayals := myquery2()

	page.Execute(w, betrayals)
}

func main() {
	server := http.Server{
		Addr: "127.0.0.1:8001",
	}
	// handle each of the requests to either
	http.HandleFunc("/form", form)
	http.HandleFunc("/list", list)
	server.ListenAndServe()
}
