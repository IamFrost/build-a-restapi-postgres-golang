package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"github.com/rs/cors"
	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

// const (
// 	hostname = "localhost"
// 	port = 3000
// 	username = "postgres"
// 	password = "postgres"
// 	database = "test1"
// 	)

// Post struct (Model)
type Post struct {
	Title  string `json:"title"`
	Body   string `json:"body"`
	Userid string `json:"userid"`
	ID     string `json:"id"`
}

// Init posts var as a slice Post struct
var posts []Post

func createConnection() *sql.DB {

	// Connect to the DB, panic if failed
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost/test1?sslmode=disable")
	if err != nil {
		fmt.Println(`Could not connect to db`)
		panic(err)
	}

	// check the connection
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	// fmt.Println("Successfully connected!")
	// return the connection
	return db
}

// Get all posts
func getPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// create the postgres db connection
	db := createConnection()
	// close the db connection
	defer db.Close()

	rows, err := db.Query(`SELECT * FROM posts`)
	if err != nil {
		panic(err)
	}

	fmt.Println(rows)

	var col1 string
	var col2 string
	var col3 string
	var col4 string
	posts = nil
	for rows.Next() {
		rows.Scan(&col1, &col2, &col3, &col4)
		// fmt.Println(col1, col2, col3, col4)
		posts = append(posts, Post{Title: col1, Body: col2, Userid: col3, ID: col4})
	}
	json.NewEncoder(w).Encode(posts)
}

// Get single post
func getPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// create the postgres db connection
	db := createConnection()
	// close the db connection
	defer db.Close()

	// get the postid from the request params, key is "id"
	params := mux.Vars(r)

	// convert the id type from string to int
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	rows := db.QueryRow(`SELECT * FROM posts WHERE id=$1`, id)

	var col1 string
	var col2 string
	var col3 string
	var col4 string
	posts = nil

	rows.Scan(&col1, &col2, &col3, &col4)
	// fmt.Println(col1, col2, col3, col4)
	posts = append(posts, Post{Title: col1, Body: col2, Userid: col3, ID: col4})

	json.NewEncoder(w).Encode(posts)

}

// Add new post
func createPost(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")

	var p Post
	err := json.NewDecoder(r.Body).Decode(&p)
	// fmt.Println("this is post p : ", p)
	// fmt.Println("this is error: ", err)
	useridConv, err := strconv.Atoi(p.Userid)
	if err != nil {
		panic(err)
	}
	idConv, err := strconv.Atoi(p.ID)
	if err != nil {
		panic(err)
	}
	fmt.Println("here is id1: ", useridConv)
	fmt.Println("here is id2: ", idConv)

	// create the postgres db connection
	db := createConnection()
	// close the db connection
	defer db.Close()

	// fmt.Println(`INSERT INTO posts (title, body, userid, id) VALUES ($1, $2, $3, $4)`,p.Title,p.Body,p.Userid,p.ID)
	row, err := db.Exec(`INSERT INTO posts (title, body, userid, id) VALUES ($1, $2, $3, $4)`, p.Title, p.Body, p.Userid, p.ID)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}
	fmt.Println(row)
	// fmt.Printf("Inserted a single record %v", p.ID)
}

// Update post
func updatePost(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "application/json")
	// w.Header().Set("Access-Control-Allow-Origin", "*")

	var p Post
	err := json.NewDecoder(r.Body).Decode(&p)
	fmt.Println("this is post p : ", p)
	fmt.Println("this is error: ", err)
	// useridConv, err := strconv.Atoi(p.Userid)
	// if err != nil {
	// 	panic(err)
	// }
	// idConv, err := strconv.Atoi(p.ID)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("here is id1: ", useridConv)
	// fmt.Println("here is id2: ", idConv)

	//get the postid from the request params, key is "id"
	params := mux.Vars(r)

	//convert the id type from string to int
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}
	fmt.Println("this is id: ", id)
	// fmt.Println()
	//create the postgres db connection
	db := createConnection()
	//close the db connection
	defer db.Close()

	fmt.Println(`UPDATE posts SET title=$1, body=$2, userid=$3, id=$4 WHERE id=$5`, p.Title, p.Body, p.Userid, p.ID, id)
	row, err := db.Exec(`UPDATE posts SET title=$1, body=$2, userid=$3, id=$4 WHERE id=$5`, p.Title, p.Body, p.Userid, p.ID, id)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}
	fmt.Println(row)
	fmt.Printf("Inserted a single record %v", p.ID)
}

// Delete post
func deletePost(w http.ResponseWriter, r *http.Request) {

	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	// w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	// w.Header().Add("Content-Type", "application/json")

	// create the postgres db connection
	db := createConnection()
	// close the db connection
	defer db.Close()

	// get the postid from the request params, key is "id"
	params := mux.Vars(r)

	// convert the id type from string to int
	id, err := strconv.Atoi(params["id"])
	fmt.Println(id)

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	rows, err := db.Exec(`DELETE FROM posts WHERE id=$1`, id)
	if err != nil {
		panic(err)
	}
	fmt.Println(rows)
}

// Main function
func main() {
	// pgConString := fmt.Sprintf("host=%s port=%d user=%s "+
	// 							"password=%s dbname=%s sslmode=disable",
	// 							hostname, port, username, password, database)
	// Init router
	r := mux.NewRouter()
	// fmt.Println(r)

	// Hardcoded data - @todo: add database
	// posts = append(posts, Post{Title: "hi", Body: "card", Userid: 200, ID: 66})
	// posts = append(posts, Post{Title: "hello", Body: "hi there", Userid: 2400, ID: 656})

	// Route handles & endpoints
	r.HandleFunc("/posts", getPosts).Methods("GET")
	r.HandleFunc("/posts/{id}", getPost).Methods("GET")
	r.HandleFunc("/posts", createPost).Methods("POST")
	r.HandleFunc("/posts/{id}", updatePost).Methods("PUT")
	r.HandleFunc("/posts/{id}", deletePost).Methods("DELETE")

	// Start server
	 handler := cors.AllowAll().Handler(r)
	 // handler := cors.Default().Handler(r)
	log.Fatal(http.ListenAndServe(":3000", handler))
}
