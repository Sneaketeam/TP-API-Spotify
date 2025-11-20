package router

import (
	"TpSpotify/controller"
	"fmt"
	"net/http"
)

func InitServer() {

	fileServer := http.FileServer(http.Dir("./stylecss"))
	http.Handle("/stylecss/", http.StripPrefix("/stylecss/", fileServer))

	imgServer := http.FileServer(http.Dir("./image"))
	http.Handle("/image/", http.StripPrefix("/image/", imgServer))


	http.HandleFunc("/", controller.IndexPage)
	http.HandleFunc("/album/damso", controller.AlbumDamsoHandler)
	http.HandleFunc("/track/laylow", controller.TrackLaylowHandler)

	fmt.Println("Le serveur est lancé sur http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}