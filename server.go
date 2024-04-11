package main

import (
	"encoding/json"
	"errors"
	"fmt" // formatting and printing values to the console.
	"io"
	"log"      // logging messages to the console.
	"net/http" // Used for build HTTP servers and clients.
)

// Port we listen on.
const portNum string = ":8080"

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var users []User

// Handler functions.
func Home(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "Homepage")
	if r.Method == "GET" {
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(`hello world`))
		return
	}

	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("wrong method"))

}

func HandleUser(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		usersJson, err := json.Marshal(users)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("server is broke"))
		}

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(usersJson))

	}

	if r.Method == "POST" {
		var user User
		body, err := io.ReadAll(r.Body)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("server is broke"))
		}

		err = json.Unmarshal(body, &user)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("server is broke"))
		}

		users = append(users, user)

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(body))
		fmt.Println(users)

	}

}

func DisplayMethod(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			fmt.Println("POST USERS")
		} else if r.Method == "GET" {
			fmt.Println("GET USERS")
		}

		next.ServeHTTP(w, r)

	})
}

func DisplayHello(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("hello console")
		next.ServeHTTP(w, r)

	})
}

type MiddleWare struct {
	handler http.Handler
}

func (M *MiddleWare) DisplayHello() {
	fmt.Println("hello console")
}

func (M *MiddleWare) HandleBasicAuth(w http.ResponseWriter, r *http.Request) error {
	username, password, ok := r.BasicAuth()
	if ok {
		if username == "admin" && password == "admin" {
			return nil
		}

	}
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
	return errors.New("Unauthorized")
}

func (M *MiddleWare) DisplayMethod(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		fmt.Println("POST USERS")
	} else if r.Method == "GET" {
		fmt.Println("GET USERS")
	}
}

func (M *MiddleWare) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := M.HandleBasicAuth(w, r)
	if err != nil {
		return
	}
	M.DisplayHello()
	M.DisplayMethod(w, r)
	M.handler.ServeHTTP(w, r)
}

func main() {
	log.Println("Starting our simple http server.")

	mux := http.NewServeMux()

	// Registering our handler functions, and creating paths.
	mux.HandleFunc("/", Home)

	mux.HandleFunc("/users", HandleUser)

	log.Println("Started on port", portNum)
	fmt.Println("To close connection CTRL+C :-)")

	middleWare := &MiddleWare{handler: mux}
	srv := http.Server{
		Addr:    portNum,
		Handler: middleWare,
	}

	// Spinning up the server.
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
