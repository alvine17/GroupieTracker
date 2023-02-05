package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
)

type Artist struct {
	Id      int
	Image   string
	Name    string
	Members []string
}

func main() {

	tmpl := template.Must(template.ParseFiles("index.html"))

	global := callAPI()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, global)
	})

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.ListenAndServe(":8080", nil)

}

func callAPI() Artist {
	response, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	var artist Artist

	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(body, &artist.Members)

	return artist
}
