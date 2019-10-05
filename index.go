package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	mgo "github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/gorilla/mux"
)

type db struct {
	collection *mgo.Collection
}

type user struct {
	UserID int    `bson:"userid"`
	Name   string `bson:"name"`
}

func (db *db) createUsers(users ...*user) {
	for _, user := range users {
		err := db.collection.Insert(user)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (db *db) getByUserID(userID int) user {
	var user user
	db.collection.Find(bson.M{"userid": userID}).One(&user)
	return user
}

// Ping echoes a Pong message
func Ping(w http.ResponseWriter, r *http.Request) {
	response, _ := json.Marshal("Pong")
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// GetByUserID responds with a JSON of the user with the given id
func (db *db) GetByUserID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDInt, _ := strconv.Atoi(vars["userid"])
	user := db.getByUserID(userIDInt)
	response, _ := json.Marshal(user)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func main() {
	session, err := mgo.Dial(os.Getenv("DB_HOST"))
	if err != nil {
		panic(err)
	}
	defer session.Close()

	c := session.DB("usersdb").C("user")
	db := &db{collection: c}

	// If there are no users, create some
	docCount, err := db.collection.Count()
	if docCount == 0 {
		db.createUsers(&user{
			UserID: 1,
			Name:   "John",
		}, &user{
			UserID: 2,
			Name:   "George",
		})
	}

	r := mux.NewRouter()

	r.HandleFunc(
		"/ping",
		Ping).Methods("GET")

	r.HandleFunc(
		"/user/{userid:[0-9]+}",
		db.GetByUserID).Methods("GET")

	srv := &http.Server{
		Handler: r,
		Addr:    "0.0.0.0:8000",
	}
	log.Fatal(srv.ListenAndServe())

}
