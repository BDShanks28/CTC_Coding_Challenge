package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB
var err error

// Attempt at error handling
type Response struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

func sendJSONResponse(w http.ResponseWriter, statusCode int, message string, status string) {
	response := Response{
		Message: message,
		Status:  status,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	if db == nil {
		sendJSONResponse(w, http.StatusInternalServerError, "Database not initialized", "error")
		return
	}
	if r.Method != http.MethodPost {
		sendJSONResponse(w, http.StatusMethodNotAllowed, "Invalid request method", "error")
		return
	}

	var user struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		sendJSONResponse(w, http.StatusBadRequest, "Invalid request payload", "error")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		sendJSONResponse(w, http.StatusInternalServerError, "Error hashing password", "error")
		return
	}

	_, err = db.Exec("INSERT INTO users (email, password_hash) VALUES ($1, $2)", user.Email, string(hashedPassword))
	if err != nil {
		log.Printf("Error saving user: %v", err)
		sendJSONResponse(w, http.StatusInternalServerError, "Error saving user", "error")
		return
	}

	sendJSONResponse(w, http.StatusCreated, "Sign up successful!", "success")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var user struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var storedHash string
	err := db.QueryRow("SELECT password_hash FROM users WHERE email = $1", user.Email).Scan(&storedHash)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Error checking credentials", http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(user.Password))
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// On successful login, sets a session cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "some-token",
		Expires: time.Now().Add(24 * time.Hour),
		Path:    "/",
	})

	// Respond with a redirect to the congrats page
	http.Redirect(w, r, "/congrats", http.StatusSeeOther)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	// Clear the session cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now().Add(-1 * time.Hour), // Expire the cookie immediately
		Path:    "/",
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func ErrorHandler(w http.ResponseWriter, errMsq string, statusCode int) {
	response := Response{
		Message: errMsq,
		Status:  "error",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func main() {
	godotenv.Load("info.env")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	fmt.Println("DB_HOST:", dbHost)
	fmt.Println("DB_PORT:", dbPort)
	fmt.Println("DB_USER:", dbUser)
	fmt.Println("DB_PASSWORD:", dbPassword)
	fmt.Println("DB_NAME:", dbName)

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging database: ", err)
	}

	router := mux.NewRouter()

	// Serve static files
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Serve the main HTML file
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	}).Methods("GET")

	// Serve the congrats page
	router.HandleFunc("/congrats", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/congrats.html")
	}).Methods("GET")

	router.HandleFunc("/signup", signupHandler).Methods("POST")
	router.HandleFunc("/login", loginHandler).Methods("POST")
	router.HandleFunc("/logout", logoutHandler).Methods("POST")

	log.Println("Server started at http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
