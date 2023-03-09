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
}

type Artists []Artist

type Relations struct {
	Index []Relation
}

type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
	RelArtist      Artist
}

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

	http.HandleFunc("/planning", func(w http.ResponseWriter, r *http.Request) {
		tmpl3 := template.Must(template.ParseFiles("./static/planning.html"))
		Relation := callRelation()
		tmpl3.Execute(w, Relation)

	})

	// http.HandleFunc("/location", func(w http.ResponseWriter, r *http.Request) {

	// 	tmpl4 := template.Must(template.ParseFiles(".static/planning.html"))
	// 	locations := callLocation()
	// 	// fmt.Println(locations.Index)
	// 	tmpl4.Execute(w, locations)
	// })

	http.ListenAndServe(":8080", nil)

}

func callRelation() Relations {
	response, err := http.Get("https://groupietrackers.herokuapp.com/api/relation")
	var relations Relations

	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(body, &relations)

	artist := callAPI()

	for i, relation := range relations.Index {

		relations.Index[i].RelArtist = artist[relation.ID-1]
		//fmt.Println(relation.RelArtist)
	}
	fmt.Println(relations)
	return relations
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

// func callDate() Datess {
// 	response, err := http.Get("https://groupietrackers.herokuapp.com/api/dates")
// 	var dates Datess

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	defer response.Body.Close()

// 	body, err := ioutil.ReadAll(response.Body)

// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	json.Unmarshal(body, &dates)
// 	// fmt.Println(dates)

// 	return dates
// }

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
