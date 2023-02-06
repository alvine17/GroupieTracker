package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
)

type Artists []struct {
	ID      int      `json:"id"`
	Image   string   `json:"image"`
	Name    string   `json:"name"`
	Members []string `json:"members"`
}

func main() {

	tmpl := template.Must(template.ParseFiles("index.html"))

	// global := callAPI()
	global := callAPI()

	for _, i := range global {
		fmt.Println(i.Name)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, global)
	})

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.ListenAndServe(":8080", nil)

}

func callAPI() Artists {
	response, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")

	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	// fmt.Println(string(body))
	if err != nil {
		log.Fatal(err)
	}
	array := Artists{}
	json.Unmarshal(body, &array)

	return array

}
