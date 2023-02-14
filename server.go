package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
)

type Artist struct {
	Membre string
}
type Artists []struct {
	Id      int
	Image   string
	Name    string
	Members []string
	// Dates   []string
}

//	type Planning struct {
//		Dates []Dates
//	}
type Dates struct {
	Index []struct {
		ID    int      `json:"id"`
		Dates []string `json:"dates"`
	} `json:"index"`
}

func main() {

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	global := callAPI()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("./static/index.html"))
		tmpl.Execute(w, global)
	})

	http.HandleFunc("/artists", func(w http.ResponseWriter, r *http.Request) {
		tmpl2 := template.Must(template.ParseFiles("./static/artists.html"))
		tmpl2.Execute(w, global)
	})

	http.HandleFunc("/planning", func(w http.ResponseWriter, r *http.Request) {
		tmpl3 := template.Must(template.ParseFiles("./static/planning.html"))
		dates := callDate()
		fmt.Println(dates.Index)
		tmpl3.Execute(w, dates)
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
	// fmt.Println(string(body))
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(body, &artist)
	// fmt.Println(artist)
	return artist
}

func callDate() Dates {
	response, err := http.Get("https://groupietrackers.herokuapp.com/api/dates")
	var dates Dates

	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	// fmt.Println(string(body))

	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(body, &dates)
	// fmt.Println(dates)

	return dates
}
