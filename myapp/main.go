package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func signupHandler(w http.ResponseWriter, r *http.Request) {
	var user struct {
		Email    string `json:"email"`
		Password string `json:password`
	}

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("INSERT INTO users (email, password_hash) VALUES ($1, $2)", user.Email, string(hashedPassword))
	if err != nil {
		log.Printf("Error saving user: %v", err)
		http.Error(w, "Error saving user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
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
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(user.Password))
	if err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	// Simplified token generation for user to log in
	token := "dummy-token success"
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	// Invalidate the token on the client-side.
	w.WriteHeader(http.StatusOK)

	// POST information to determine if logout was successful
	response := map[string]string{"message": "User logged out"}
	json.NewEncoder(w).Encode(response)
}

func main() {
	var err error
	connStr := "user=postgres password=johnD03! dbname=user_login sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := mux.NewRouter()
	router.HandleFunc("/signup", signupHandler).Methods("POST")
	router.HandleFunc("/login", loginHandler).Methods("POST")
	router.HandleFunc("/logout", logoutHandler).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", router))
}
