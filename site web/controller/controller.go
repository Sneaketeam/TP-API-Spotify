package controller

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

type SpotifyData struct {
	Items []Album `json:"items"`
}

type Album struct {
	Name        string  `json:"name"`
	ReleaseDate string  `json:"release_date"`
	TotalTracks int     `json:"total_tracks"`
	Images      []Image `json:"images"`
}

type Image struct {
	Url string `json:"url"`
}

func IndexPage(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("template/index.html")
	t.Execute(w, nil)
}

func AlbumDamsoHandler(w http.ResponseWriter, r *http.Request) {
	token := "BQB7ZhTm98iSt8H133ndPhJSHe9Eo4LKdtb9dBesLnXqeCR3kc93NaSB9JiBqe61CYYsA-YvWdUErC7IKDI0GalF-hy11xP4rkmIomUFrHW7VG8cuSW-KBMsrNe4vzBbU7RcYOKw6yo"
	damsoID := "2RS1iy08e1nNHQA3aqyxgv"
	url := "https://api.spotify.com/v1/artists/" + damsoID + "/albums"

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Fprint(w, "Erreur internet ou api")
		return
	}
	defer resp.Body.Close()

	var result SpotifyData
	json.NewDecoder(resp.Body).Decode(&result)

	t, _ := template.ParseFiles("template/damso.html")
	t.Execute(w, result.Items)
}

func TrackLaylowHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Page Laylow à faire")
}