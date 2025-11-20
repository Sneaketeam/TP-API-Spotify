package controller

import (
	"encoding/base64"
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// --- TES IDENTIFIANTS SPOTIFY (À REMPLIR !!!) ---
const (
	ClientID     = "bfd18fe9f8fe4350973a982051bf8366"
	ClientSecret = "799165282057445582d03bd575714af4"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type AlbumResponse struct {
	Items []Album `json:"items"`
}

type SearchResponse struct {
	Tracks TrackList `json:"tracks"`
}

type TrackList struct {
	Items []Track `json:"items"`
}

type Album struct {
	Name        string  `json:"name"`
	ReleaseDate string  `json:"release_date"`
	TotalTracks int     `json:"total_tracks"`
	Images      []Image `json:"images"`
}

type Track struct {
	Name         string       `json:"name"`
	Album        AlbumInfo    `json:"album"`
	Artists      []Artist     `json:"artists"`
	ExternalUrls ExternalUrls `json:"external_urls"`
}

type AlbumInfo struct {
	Name        string  `json:"name"`
	ReleaseDate string  `json:"release_date"`
	Images      []Image `json:"images"`
}

type Artist struct {
	Name string `json:"name"`
}

type ExternalUrls struct {
	Spotify string `json:"spotify"`
}

type Image struct {
	Url string `json:"url"`
}

// --- FONCTION POUR RÉCUPÉRER LE TOKEN AUTOMATIQUEMENT ---
func getAutoToken() (string, error) {
	authURL := "https://accounts.spotify.com/api/token"

	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", authURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	auth := base64.StdEncoding.EncodeToString([]byte(ClientID + ":" + ClientSecret))
	req.Header.Add("Authorization", "Basic "+auth)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return "", &SpotifyError{Status: resp.StatusCode, Body: bodyBytes}
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(bodyBytes, &tokenResp); err != nil {
		return "", err
	}

	if tokenResp.AccessToken == "" {
		return "", &SpotifyError{Status: resp.StatusCode, Body: bodyBytes}
	}

	return tokenResp.AccessToken, nil
}

// --- HANDLERS ---

func IndexPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("template/index.html")
	if err != nil {
		http.Error(w, "Template introuvable", http.StatusInternalServerError)
		return
	}
	_ = t.Execute(w, nil)
}

func AlbumDamsoHandler(w http.ResponseWriter, r *http.Request) {
	token, err := getAutoToken()
	if err != nil {
		http.Error(w, "Impossible d'obtenir le token Spotify", http.StatusInternalServerError)
		return
	}

	client := &http.Client{}
	apiUrl := "https://api.spotify.com/v1/artists/2UwqpfQtNuhBwviIC0f2ie/albums?include_groups=album"

	req, _ := http.NewRequest("GET", apiUrl, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Erreur de connexion à Spotify", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		http.Error(w, "Erreur Spotify", resp.StatusCode)
		return
	}

	var data AlbumResponse
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		http.Error(w, "Erreur décodage Spotify", http.StatusInternalServerError)
		return
	}

	t, err := template.ParseFiles("template/damso.html")
	if err != nil {
		http.Error(w, "Template damso introuvable", http.StatusInternalServerError)
		return
	}
	_ = t.Execute(w, data.Items)
}

func TrackLaylowHandler(w http.ResponseWriter, r *http.Request) {
	token, err := getAutoToken()
	if err != nil {
		http.Error(w, "Impossible d'obtenir le token Spotify", http.StatusInternalServerError)
		return
	}

	client := &http.Client{}
	apiUrl := "https://api.spotify.com/v1/search?q=track:Maladresse%20artist:Laylow&type=track&limit=1"

	req, _ := http.NewRequest("GET", apiUrl, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Erreur de connexion à Spotify", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		http.Error(w, "Erreur Spotify", resp.StatusCode)
		return
	}

	var result SearchResponse
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		http.Error(w, "Erreur décodage Spotify", http.StatusInternalServerError)
		return
	}

	if len(result.Tracks.Items) == 0 {
		http.Error(w, "Aucun morceau trouvé", http.StatusNotFound)
		return
	}

	t, err := template.ParseFiles("template/laylow.html")
	if err != nil {
		http.Error(w, "Template laylow introuvable", http.StatusInternalServerError)
		return
	}
	_ = t.Execute(w, result.Tracks.Items[0])
}

// --- STRUCTURE D'ERREUR SPOTIFY PERSONNALISÉE ---
type SpotifyError struct {
	Status int
	Body   []byte
}

func (e *SpotifyError) Error() string {
	return "Spotify API error"
}