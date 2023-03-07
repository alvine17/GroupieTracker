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
		tmpl2 := template.Must(template.ParseFiles("./static/artists.html"))
		tmpl2.Execute(w, global)
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

		// Lire le corps de la réponse
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		// Désérialiser la réponse JSON dans une structure Artist
		var artist Artist
		err = json.Unmarshal(body, &artist)
		if err != nil {
			log.Fatal(err)
		}

		// Charger le template HTML pour la description de l'artiste
		tmpl3 := template.Must(template.ParseFiles("./static/description.html"))

		// Exécuter le template en envoyant les informations sur l'artiste à la page HTML
		err = tmpl3.Execute(w, artist)
		if err != nil {
			log.Fatal(err)
		}
	})

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("query")
		fmt.Println("Recherche de l'artiste : ", query)
	})

	http.HandleFunc("/contact", contactHandler)

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
	fmt.Fprintf(w, "<h1>Thank you for contacting us!</h1>")

}
