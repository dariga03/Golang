package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Response struct {
	AnimeList []Anime `json:"animeList"`
}

type Anime struct {
	Title      string   `json:"title"`
	Characters []string `json:"characters"`
	ReleaseDate string   `json:"release_date"`
}

func main() {
	log.Println("starting API server")
	//create a new router
	router := mux.NewRouter()
	log.Println("creating routes")
	//specify endpoints
	router.HandleFunc("/health-check", HealthCheck).Methods("GET")
	router.HandleFunc("/anime", GetAnimeList).Methods("GET")
	router.HandleFunc("/anime/{title}", GetAnime).Methods("GET")
	http.Handle("/", router)

	//start and listen to requests
	http.ListenAndServe(":8080", router)

}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	log.Println("entering health check end point")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "API is up and running")
}

func GetAnimeList(w http.ResponseWriter, r *http.Request) {
	log.Println("entering anime list end point")
	var response Response
	animeList := prepareAnimeList()

	response.AnimeList = animeList

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return
	}

	w.Write(jsonResponse)
}

func GetAnime(w http.ResponseWriter, r *http.Request) {
	log.Println("entering anime end point")
	vars := mux.Vars(r)
	title := vars["title"]

	animeList := prepareAnimeList()

	for _, anime := range animeList {
		if anime.Title == title {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			jsonResponse, err := json.Marshal(anime)
			if err != nil {
				return
			}
			w.Write(jsonResponse)
			return
		}
	}

	http.NotFound(w, r)
}


func prepareAnimeList() []Anime {
	var animeList []Anime

	anime := Anime{
		Title:      "Attack on Titan",
		Characters: []string{"Eren Yeager", "Mikasa Ackerman", "Armin Arlert"},
		ReleaseDate: "2013-04-06",
	}
	animeList = append(animeList, anime)

	// Добавьте другие записи аниме по мере необходимости

	return animeList
}

//localhost:8080/anime
//http://localhost:8080/health-check
//http://localhost:8080/anime/Attack%20on%20Titan