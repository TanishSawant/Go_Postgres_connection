package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"go_postgres.com/Models"
)

var db *gorm.DB

var err error

var (
	drivers = []Models.Driver{
		{Name: "Jimmy Johnson", License: "ABC123"},

		{Name: "Howard Hills", License: "XYZ789"},

		{Name: "Craig Colbin", License: "DEF333"},
	}

	cars = []Models.Car{

		{Year: 2000, Make: "Toyota", ModelName: "Tundra", DriverID: 1},

		{Year: 2001, Make: "Honda", ModelName: "Accord", DriverID: 1},

		{Year: 2002, Make: "Nissan", ModelName: "Sentra", DriverID: 2},

		{Year: 2003, Make: "Ford", ModelName: "F-150", DriverID: 3},
	}
)

func GetCars(w http.ResponseWriter, r *http.Request) {

	var cars []Models.Car

	db.Find(&cars)

	json.NewEncoder(w).Encode(&cars)

}

func GetCar(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	var car Models.Car

	db.First(&car, params["id"])

	json.NewEncoder(w).Encode(&car)

}

func GetDriver(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	var driver Models.Driver

	var cars []Models.Car

	db.First(&driver, params["id"])

	db.Model(&driver).Related(&cars)

	driver.Cars = cars

	json.NewEncoder(w).Encode(&driver)

}

func DeleteCar(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	var car Models.Car

	db.First(&car, params["id"])

	db.Delete(&car)

	var cars []Models.Car

	db.Find(&cars)

	json.NewEncoder(w).Encode(&cars)

}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	dbname   = "postgres"
	password = "pablo"
)

func main() {

	fmt.Println("Creating Router....")

	router := mux.NewRouter()

	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=postgres sslmode=disable password=pablo")

	if err != nil {
		panic(err)
	}

	defer db.Close()

	db.AutoMigrate(&Models.Driver{})

	db.AutoMigrate(&Models.Car{})

	for index := range cars {
		db.Create(&cars[index])
	}

	for index := range drivers {
		db.Create(&drivers[index])
	}

	router.HandleFunc("/cars", GetCars).Methods("GET")

	router.HandleFunc("/cars/{id}", GetCar).Methods("GET")

	router.HandleFunc("/drivers/{id}", GetDriver).Methods("GET")

	router.HandleFunc("/cars/{id}", DeleteCar).Methods("DELETE")

	handler := cors.Default().Handler(router)

	log.Fatal(http.ListenAndServe(":8080", handler))
}
