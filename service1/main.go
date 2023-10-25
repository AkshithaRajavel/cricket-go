package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"service1/models"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MONGO *mongo.Client
var Coll *mongo.Collection

func Mongo_connect() {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI("mongodb+srv://akshitha:akshitha@cluster0.4iviwmv.mongodb.net/?retryWrites=true&w=majority").SetServerAPIOptions(serverAPI)
	var err error
	MONGO, err = mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	Coll = MONGO.Database("sports").Collection("players")
	fmt.Println("mongo connection success")
}
func main() {
	Mongo_connect()
	router := mux.NewRouter().PathPrefix("/team").Subrouter()
	router.HandleFunc("/getTeams", func(w http.ResponseWriter, r *http.Request) {
		teams := []models.Team{}
		pipeline := []bson.M{bson.M{"$match": bson.M{"country": bson.M{"$exists": true}}}, bson.M{"$group": bson.M{"_id": "$country"}}}
		cursor, err := Coll.Aggregate(context.TODO(), pipeline)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "internal server error"})
			return
		}
		type result struct {
			Id string `json:"_id" bson:"_id"`
		}
		res := []result{}
		err = cursor.All(context.TODO(), &res)
		if err != nil {
			fmt.Println(err)
		}
		for _, t := range res {
			var team = models.Team{}
			team.Country = t.Id
			cursor, err = Coll.Find(context.TODO(), bson.M{"country": team.Country})
			team.Players = []models.Player{}
			err = cursor.All(context.TODO(), &team.Players)
			err = Coll.FindOne(context.TODO(), bson.M{"country": team.Country, "captain": true}).Decode(&team.Captain)
			teams = append(teams, team)
		}
		json.NewEncoder(w).Encode(teams)
	})
	router.HandleFunc("/getPlayers", func(w http.ResponseWriter, r *http.Request) {
		players := []models.Player{}
		cursor, err := Coll.Find(context.TODO(), bson.D{})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "internal server error"})
			return
		}
		err = cursor.All(context.TODO(), &players)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "internal server error"})
			return
		}
		json.NewEncoder(w).Encode(players)

	}).Methods("GET")
	router.HandleFunc("/createPlayer", func(w http.ResponseWriter, r *http.Request) {
		player := models.Player{}
		e := json.NewDecoder(r.Body).Decode(&player)
		if e != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "bad request"})
			return
		}
		player.ID = primitive.NewObjectID()
		fmt.Println(player.ID)
		insertRes, err := Coll.InsertOne(context.Background(), player)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "internal server error"})
			return
		}
		json.NewEncoder(w).Encode(insertRes.InsertedID)
	}).Methods("POST")
	router.HandleFunc("/getPlayer/{id}", func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, _ := primitive.ObjectIDFromHex(params["id"])
		var player = models.Player{}
		Coll.FindOne(context.TODO(), map[string]primitive.ObjectID{"_id": id}).Decode(&player)
		json.NewEncoder(w).Encode(&player)
	}).Methods("GET")
	router.HandleFunc("/updatePlayer/{id}", func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, _ := primitive.ObjectIDFromHex(params["id"])
		var player = models.Player{}
		json.NewDecoder(r.Body).Decode(&player)
		updateRes, err := Coll.UpdateByID(context.TODO(), id, bson.D{
			{"$set", player}, // Replace with your update criteria
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "internal server error"})
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"modified_players": updateRes.ModifiedCount})
	}).Methods("POST")
	router.HandleFunc("/deletePlayer/{id}", func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, _ := primitive.ObjectIDFromHex(params["id"])
		deleteRes, err := Coll.DeleteOne(context.TODO(), map[string]primitive.ObjectID{"_id": id})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "internal server error"})
			return
		}
		json.NewEncoder(w).Encode(map[string]int64{"deleted": deleteRes.DeletedCount})
	}).Methods("GET")
	http.ListenAndServe(":8081", router)
}
