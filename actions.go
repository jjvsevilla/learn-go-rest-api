package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func getSession() *mgo.Session {
	var mongoURI = "mongodb://juanjo:juanjodb@godemo-shard-00-00-fslfq.mongodb.net:27017,godemo-shard-00-01-fslfq.mongodb.net:27017,godemo-shard-00-02-fslfq.mongodb.net:27017/test?replicaSet=GoDemo-shard-0&authSource=admin"

	dialInfo, dialInfoErr := mgo.ParseURL(mongoURI)
	if dialInfoErr != nil {
		panic(dialInfoErr)
	}

	tlsConfig := &tls.Config{}
	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
		return conn, err
	}

	session, sessionErr := mgo.DialWithInfo(dialInfo)
	if sessionErr != nil {
		panic(sessionErr)
	}

	return session
}

func responseMovie(w http.ResponseWriter, status int, results Movie) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(results)
}

func responseMovies(w http.ResponseWriter, status int, results []Movie) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(results)
}

var collection = getSession().DB("GoDemo").C("Movies")

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "home")
}

func MovieList(w http.ResponseWriter, r *http.Request) {
	var results []Movie
	err := collection.Find(nil).Sort("-year").All(&results)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Results", results)
	}

	responseMovies(w, 200, results)
}

func MovieDetail(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	movieID := params["movieId"]

	if !bson.IsObjectIdHex(movieID) {
		w.WriteHeader(404)
		return
	}

	oID := bson.ObjectIdHex(movieID)

	results := Movie{}
	err := collection.FindId(oID).One(&results)

	if err != nil {
		w.WriteHeader(404)
		return
	}

	responseMovie(w, 200, results)
}

func MovieAdd(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var movieData Movie
	err := decoder.Decode(&movieData)
	if err != nil {
		panic(err)
	}

	defer r.Body.Close()

	errMongo := collection.Insert(movieData)
	if errMongo != nil {
		w.WriteHeader(500)
		return
	}

	responseMovie(w, 200, movieData)
}

func MovieUpdate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	movieID := params["movieId"]

	if !bson.IsObjectIdHex(movieID) {
		w.WriteHeader(404)
		return
	}

	oid := bson.ObjectIdHex(movieID)

	decoder := json.NewDecoder(r.Body)

	var movie_data Movie
	err := decoder.Decode(&movie_data)

	if err != nil {
		panic(err)
		w.WriteHeader(500)
		return
	}

	defer r.Body.Close()

	document := bson.M{"_id": oid}
	change := bson.M{"$set": movie_data}
	err = collection.Update(document, change)

	if err != nil {
		w.WriteHeader(404)
		return
	}

	responseMovie(w, 200, movie_data)
}

type Message struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (this *Message) setStatus(data string) {
	this.Status = data
}

func (this *Message) setMessage(data string) {
	this.Message = data
}

func MovieRemove(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	movieID := params["movieId"]

	if !bson.IsObjectIdHex(movieID) {
		w.WriteHeader(404)
		return
	}

	oid := bson.ObjectIdHex(movieID)

	err := collection.RemoveId(oid)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	message := new(Message)
	message.setStatus("success")
	message.setMessage("La pelicula con ID " + movieID + " ha sido borrada correctamente")

	results := message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(results)
}
