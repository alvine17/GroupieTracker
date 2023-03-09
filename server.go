package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"text/template"
)

type Artist struct {
	Id           int
	Image        string
	Name         string
	Members      []string
	CreationDate int
	FirstAlbum   string
	ConcertDate  string
	Relation     string
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
	http.HandleFunc("/artists", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		creationDates, ok := q["creation_date"]

		// Filtrage des artistes selon les années de création
		filteredArtists := global
		if ok {
			for _, year := range creationDates {
				startYear, _ := strconv.Atoi(year)
				filteredArtists = FilterArtistsByYear(filteredArtists, startYear)
			}
		}

		// Génération du template HTML avec les artistes filtrés
		tmpl2, err := template.ParseFiles("./static/artists.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl2.Execute(w, filteredArtists)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.HandleFunc("/planning", func(w http.ResponseWriter, r *http.Request) {
			tmpl3 := template.Must(template.ParseFiles("./static/planning.html"))
			Relation := callRelation()
			tmpl3.Execute(w, Relation)

		})

		http.HandleFunc("/description", func(w http.ResponseWriter, r *http.Request) {
			// Récupérer l'ID de l'artiste sélectionné à partir des paramètres de requête
			selectedArtistID := r.URL.Query().Get("id")

			// Appeler l'API pour récupérer les informations sur l'artiste
			response, err := http.Get("https://groupietrackers.herokuapp.com/api/artists/" + selectedArtistID)
			if err != nil {
				log.Fatal(err)
			}
			defer response.Body.Close()

			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Fatal(err)
			}

			var artist Artist
			err = json.Unmarshal(body, &artist)
			if err != nil {
				log.Fatal(err)
			}

			tmpl4 := template.Must(template.ParseFiles("./static/description.html"))

			err = tmpl4.Execute(w, artist)
			if err != nil {
				log.Fatal(err)
			}
		})

		http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
			query := r.URL.Query().Get("query")
			fmt.Println("Recherche de l'artiste : ", query)
		})

		http.HandleFunc("/contact", contactHandler)

	})
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
	}
	// fmt.Println(relations)
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

func FilterArtistsByYear(artists Artists, startYear int) Artists {
	var filteredArtists Artists
	for _, artist := range artists {
		if artist.CreationDate >= startYear && artist.CreationDate < startYear+10 {
			filteredArtists = append(filteredArtists, artist)
		}
	}
	return filteredArtists
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	r.Method = "POST"
	// Getting the data of the form
	name := r.FormValue("name")
	email := r.FormValue("email")
	phone := r.FormValue("phone")
	message := r.FormValue("message")

	// Print data in terminal
	fmt.Println("Name:", name)
	fmt.Println("Email:", email)
	fmt.Println("Phone:", phone)
	fmt.Println("Message:", message)

	// New page, when data are submitted
	//fmt.Fprintf(w, "<h1>Thank you for contacting us!</h1>")
	http.Redirect(w, r, "/static/thanks.html", http.StatusSeeOther)

	tmpl, err := template.ParseFiles("./static/assets/map.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}
