package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	_ "sync/atomic"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

/* when you want to hardcode DB connection string values use below snip

const (
	DB_USER     = "db_username"
	DB_PASSWORD = "db_password"
	DB_NAME     = "db_name"
	DB_HOST     = "db_servername"
	DB_PORT     = db_port_number
)

*/

var (
	DB_USER     = os.Getenv("DB_USER")
	DB_PASSWORD = os.Getenv("DB_PASSWORD")
	DB_NAME     = os.Getenv("DB_NAME")
	DB_HOST     = os.Getenv("DB_HOST")
	DB_PORT     = os.Getenv("DB_PORT")
)

func setupDB() *sql.DB {
	dbinfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected to database!")

	//checkErr(err)
	return db

}

type Movie struct {
	MovieID   string `json:"movieid"`
	MovieName string `json:"moviename"`
}

type JsonResponse struct {
	Type    string  `json:"type"`
	Data    []Movie `json:"data"`
	Message string  `json:"message"`
}

func main() {
	router := mux.NewRouter()

	isReady := &atomic.Value{}
	isReady.Store(false)
	go func() {
		log.Printf("Readyz probe is negative by default...")
		time.Sleep(10 * time.Second)
		isReady.Store(true)
		log.Printf("Readyz probe is positive.")
	}()

	router.HandleFunc("/movies/", GetMovies).Methods("GET")

	router.HandleFunc("/movies/", CreateMovie).Methods("POST")

	router.HandleFunc("/movies/{movieid}", DeleteMovie).Methods("DELETE")

	router.HandleFunc("/movies/", DeleteMovies).Methods("DELETE")

	router.HandleFunc("/movies/{movieid}", UpdateMovies).Methods("PUT")

	router.HandleFunc("/healthz", healthz)

	router.HandleFunc("/readyz", readyz(isReady))

	fmt.Println("Server Listening at 8001")

	log.Fatal(http.ListenAndServe(":8001", router))

}

func printMessage(message string) {
	fmt.Println("")
	fmt.Println(message)
	fmt.Println("")
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// fetch all the records for GET call

func GetMovies(w http.ResponseWriter, r *http.Request) {
	db := setupDB()

	printMessage("Getting movies...")

	// Get all movies from movies table that don't have movieID = "1"
	rows, err := db.Query("SELECT * FROM movies")

	// check errors
	checkErr(err)

	// var response []JsonResponse
	var movies []Movie

	// Foreach movie
	for rows.Next() {
		var id int
		var movieID string
		var movieName string

		err = rows.Scan(&id, &movieID, &movieName)

		// check errors
		checkErr(err)

		movies = append(movies, Movie{MovieID: movieID, MovieName: movieName})
	}

	var response = JsonResponse{Type: "success", Data: movies}

	json.NewEncoder(w).Encode(response)
	defer db.Close()

}

// Create a movie

// response and request handlers
func CreateMovie(w http.ResponseWriter, r *http.Request) {
	movieID := r.FormValue("movieid")
	movieName := r.FormValue("moviename")
	//db := setupDB()

	var response = JsonResponse{}

	if movieID == "" || movieName == "" {
		response = JsonResponse{Type: "error", Message: "You are missing movieID or movieName parameter."}
	} else {
		db := setupDB()

		printMessage("Inserting movie into DB")

		fmt.Println("Inserting new movie with ID: " + movieID + " and name: " + movieName)

		var lastInsertID int
		err := db.QueryRow("INSERT INTO movies(movieID, movieName) VALUES($1, $2) returning id;", movieID, movieName).Scan(&lastInsertID)

		// check errors
		checkErr(err)

		response = JsonResponse{Type: "success", Message: "The movie has been inserted successfully!"}
		defer db.Close()

	}

	json.NewEncoder(w).Encode(response)

}

// Delete a movie

// response and request handlers
func DeleteMovie(w http.ResponseWriter, r *http.Request) {
	//db := setupDB()
	params := mux.Vars(r)

	movieID := params["movieid"]

	var response = JsonResponse{}

	if movieID == "" {
		response = JsonResponse{Type: "error", Message: "You are missing movieID parameter."}
	} else {
		db := setupDB()

		printMessage("Deleting movie from DB")

		_, err := db.Exec("DELETE FROM movies where movieID = $1", movieID)

		// check errors
		checkErr(err)

		response = JsonResponse{Type: "success", Message: "The movie has been deleted successfully!"}
		defer db.Close()

	}

	json.NewEncoder(w).Encode(response)

}

// Delete all movies

// response and request handlers
func DeleteMovies(w http.ResponseWriter, r *http.Request) {
	db := setupDB()

	printMessage("Deleting all movies...")

	_, err := db.Exec("DELETE FROM movies")

	// check errors
	checkErr(err)

	printMessage("All movies have been deleted successfully!")

	var response = JsonResponse{Type: "success", Message: "All movies have been deleted successfully!"}

	json.NewEncoder(w).Encode(response)
	defer db.Close()
}

// Update movies using id of a movie

func UpdateMovies(w http.ResponseWriter, r *http.Request) {
	//db := setupDB()
	params := mux.Vars(r)

	movieID := params["movieid"]
	movieName := r.FormValue("moviename")

	var response = JsonResponse{}

	if movieID == "" {
		response = JsonResponse{Type: "error", Message: "You are missing movieID parameter to update a movie."}
	} else {
		db := setupDB()

		printMessage("Updating movie from DB")

		_, err := db.Exec("UPDATE movies SET movieName=$2  where movieID = $1", movieID, movieName)

		// check errors
		checkErr(err)

		response = JsonResponse{Type: "success", Message: "The moviename has been updated successfully!"}
		defer db.Close()

	}

	json.NewEncoder(w).Encode(response)

}

// healthz handler serves as a liveness probe

func healthz(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// readyz handler serves as a readiness probe

func readyz(isReady *atomic.Value) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		if isReady == nil || !isReady.Load().(bool) {
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
