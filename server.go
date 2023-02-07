package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
)

type Artist struct {
	Image  string
	Membre string
}
type Artists []struct {
	Id      int
	Image   string
	Name    string
	Members []string
}

func main() {

	fs := http.FileServer(http.Dir("./static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	global := callAPI()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("./static/index.html"))
		tmpl.Execute(w, global)
	})

	http.HandleFunc("/rony", func(w http.ResponseWriter, r *http.Request) {
		tmpl2 := template.Must(template.ParseFiles("./static/test.html"))
		tmpl2.Execute(w, global)
	})

	http.ListenAndServe(":8080", nil)

}

func callAPI() Artists {
	response, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	var artist Artists

	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(body, &artist)

	return artist
}
