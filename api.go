package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/rs/cors"
)

type Profile struct {
	Id        int    `gorm:"primary_key" json:"id"`
	Name      string `json:"name"`
	Bio       string `json:"bio"`
	Birthdate string `json:"bdate"`
}

type Id_container struct {
	Id int
}

func main() {
	db, _ := gorm.Open("sqlite3", "db.db")

	defer db.Close()

	db.AutoMigrate(&Profile{})

	r := mux.NewRouter()
	r.HandleFunc("/getProfiles", getProfiles).Methods("GET")
	r.HandleFunc("/getProfile/{id}", getProfile).Methods("GET")
	r.HandleFunc("/addProfile", addProfile).Methods("POST")
	r.HandleFunc("/updateProfile/{id}", updateProfile).Methods("POST")
	http.Handle("/", cors.Default().Handler(r))
	log.Fatal(http.ListenAndServe("localhost:8080", nil))

}

func getProfiles(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("sqlite3", "db.db")
	var profiles []Profile
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	db.Find(&profiles)
	json.NewEncoder(w).Encode(profiles)
}

func getProfile(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("sqlite3", "db.db")
	var profile Profile
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	db.Where("id = ?", params["id"]).First(&profile)

	json.NewEncoder(w).Encode(profile)
}

func addProfile(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("sqlite3", "db.db")

	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	var profile Profile
	_ = json.NewDecoder(r.Body).Decode(&profile)

	db.Create(&profile)
}

func updateProfile(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("sqlite3", "db.db")

	var profile Profile
	_ = json.NewDecoder(r.Body).Decode(&profile)
	if err != nil {
		log.Fatal(err)
	}

	params := mux.Vars(r)

	w.Header().Set("Content-Type", "application/json")
	var finder Profile
	db.Where("id = ?", params["id"]).First(&finder)

	if finder.Name != profile.Name {
		db.Model(&finder).Where("id = ?", params["id"]).Update("name", profile.Name)
	}

	if finder.Bio != profile.Bio {
		db.Model(&finder).Where("id = ?", params["id"]).Update("bio", profile.Bio)
	}

	if finder.Birthdate != profile.Birthdate {
		db.Model(&finder).Where("id = ?", params["id"]).Update("birthdate", profile.Birthdate)
	}
}
