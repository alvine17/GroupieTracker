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
type Artists []Artist

func main() {
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	global := callAPI()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		selectedArtist := selectArtist(global, "Queen")
		tmpl := template.Must(template.ParseFiles("./static/index.html"))

		tmpl.Execute(w, struct {
			Artists        Artists
			SelectedArtist Artist
		}{
			global,
			selectedArtist,
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

	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(body, &artist)

	return artist
}

func selectArtist(artists Artists, name string) Artist {
	var selectedArtist Artist
	for _, artist := range artists {
		if artist.Name == name {
			selectedArtist = artist
			break
		}
	}
	return selectedArtist
}
