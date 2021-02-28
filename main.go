package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

var tpl *template.Template

var cats categories
var artists artistsArr

func init() {
	go initData()
	tpl = template.Must(template.ParseGlob("templates/*html"))

}

func initData() {
	err := cats.dataFromURL(baseURL)
	if err != nil {
		log.Fatal(err)
	}
	err = artists.dataFromURL(cats.Artists)
	if err != nil {
		log.Fatal(err)
	}
}

const (
	baseURL = "https://groupietrackers.herokuapp.com/api"
)

type categories struct {
	Artists   string `json:"artists"`
	Locations string `json:"locations"`
	Dates     string `json:"dates"`
	Relations string `json:"relation"`
}

type artistID struct {
	ID int `json:"id"`
}

type artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

type datesLocations struct {
	ID             int `json:"id"`
	DatesLocations map[string][]string
}

type artistsArr []artist

func (c *categories) dataFromURL(url string) (err error) {
	pageContent, err := http.Get(url)
	if err != nil {
		return
	}
	data, err := ioutil.ReadAll(pageContent.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, c)
	return
}

func (a *artistsArr) dataFromURL(url string) (err error) {
	pageContent, err := http.Get(url)
	if err != nil {
		return
	}
	data, err := ioutil.ReadAll(pageContent.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, a)
	return
}

func (d *datesLocations) dataFromURL(url string) (err error) {
	pageContent, err := http.Get(url)
	if err != nil {
		return
	}
	data, err := ioutil.ReadAll(pageContent.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, d)
	return
}

// func dataFromURL(url string) ([]byte, error) {
// 	return nil
// }
type errType struct {
	Status  int
	Message string
}

func errHandler(w http.ResponseWriter, r *http.Request, er errType) {
	w.WriteHeader(er.Status)
	tpl.ExecuteTemplate(w, "error", er)
}

func index(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "layout", artists)
}

func artistInfo(w http.ResponseWriter, r *http.Request) {
	artistID := r.URL.Query().Get("id")
	id, err := strconv.Atoi(artistID)
	fmt.Println(r.Header)
	if err != nil {
		errHandler(w, r, errType{Status: http.StatusInternalServerError, Message: "Something Went Wrong"})
		return
	}
	fmt.Println(id, len(artists))
	if id > len(artists) {
		errHandler(w, r, errType{Status: http.StatusNotFound, Message: "Artist Not Found"})
		return
	}
	var dl datesLocations
	err = dl.dataFromURL(artists[id-1].Relations)
	d := struct {
		ArtistInfo     artist
		DatesLocations datesLocations
	}{
		ArtistInfo:     artists[id-1],
		DatesLocations: dl,
	}
	tpl.ExecuteTemplate(w, "artist.html", d)
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("public"))))
	mux.HandleFunc("/", index)
	mux.HandleFunc("/artist", artistInfo)
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	fmt.Println("Listening on port ", server.Addr)
	server.ListenAndServe()
}
