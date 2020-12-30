package main

//postgres://nggcyuahlnctol:b5e43ee627bdfa04bf8a2b035b2d1294cadd329069d5fa053b7dca275e7ff7bb@ec2-54-175-243-75.compute-1.amazonaws.com:5432/ddhcqh7dh5dbec
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

		{Year: "2000", Make: "Toyota", ModelName: "Tundra", DriverID: "1"},

		{Year: "2001", Make: "Honda", ModelName: "Accord", DriverID: "2"},

		{Year: "2002", Make: "Nissan", ModelName: "Sentra", DriverID: "3"},

		{Year: "2003", Make: "Ford", ModelName: "F-150", DriverID: "4"},
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

func createCar(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	year := vars["year"]
	make := vars["make"]
	modelname := vars["modelname"]
	did := vars["driverid"]
	fmt.Printf("year = %s, make = %s, model = %s", year, make, modelname)
	db.Create(&Models.Car{Year: year, Make: make, ModelName: modelname, DriverID: did})
	fmt.Fprintf(w, "New Car Successfully Created")

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

	db, err = gorm.Open("postgres", "postgres://nggcyuahlnctol:b5e43ee627bdfa04bf8a2b035b2d1294cadd329069d5fa053b7dca275e7ff7bb@ec2-54-175-243-75.compute-1.amazonaws.com:5432/ddhcqh7dh5dbec")

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

	router.HandleFunc("/cars/{year}/{make}/{modelname}/{driverid}", createCar).Methods("POST")

	handler := cors.Default().Handler(router)

	log.Fatal(http.ListenAndServe(":8080", handler))
}
