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
type Artists []Artist

func main() {
	fmt.Printf("Starting server at port 8080\n")
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	global := callAPI()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		selectedArtist := selectArtist(global, "Kendrick Lamar")
		selectedArtist2 := selectArtist(global, "XXXTentacion")
		selectedArtist3 := selectArtist(global, "Rihanna")
		selectedArtist4 := selectArtist(global, "Katy Perry")
		selectedArtist5 := selectArtist(global, "Imagine Dragons")
		fmt.Println(selectedArtist3)

		tmpl := template.Must(template.ParseFiles("./static/index.html"))

		tmpl.Execute(w, struct {
			Artists         Artists
			SelectedArtist  Artist
			SelectedArtist2 Artist
			SelectedArtist3 Artist
			SelectedArtist4 Artist
			SelectedArtist5 Artist
		}{
			global,
			selectedArtist,
			selectedArtist2,
			selectedArtist3,
			selectedArtist4,
			selectedArtist5,
		})
	})
	http.HandleFunc("/artists", func(w http.ResponseWriter, r *http.Request) {
		tmpl2 := template.Must(template.ParseFiles("./static/artists.html"))
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
	// fmt.Println(string(body))
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(body, &artist)
	// fmt.Println(artist)
	return artist
}
