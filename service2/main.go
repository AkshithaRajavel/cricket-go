package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"service2/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/gorilla/mux"
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
	var matches []models.Match
	var countries []string
	router := mux.NewRouter().PathPrefix("/match").Subrouter()
	router.HandleFunc("/schedule", func(w http.ResponseWriter, r *http.Request) {
		matches = []models.Match{}
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
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "internal server error"})
			return
		}
		countries = []string{}
		for _, t := range res {
			countries = append(countries, t.Id)
		}
		N := len(countries)
		n := 1
		for i := 0; i < N; i++ {
			for j := i + 1; j < N; j++ {
				var match = models.Match{}
				match.Name = fmt.Sprintf("match_%d", n)
				match.Team1 = countries[i]
				match.Team2 = countries[j]
				match.Date = time.Now().AddDate(0, 0, n)
				matches = append(matches, match)
				n += 1
			}
		}
		matches = append(matches, models.Match{
			Name:  "final",
			Team1: "first-in-table",
			Team2: "second-in-table",
			Date:  time.Now().AddDate(0, 0, n),
		})
		json.NewEncoder(w).Encode(matches)
	})
	router.HandleFunc("/play", func(w http.ResponseWriter, r *http.Request) {
		var PointsTable = map[string]int{}
		for _, t := range countries {
			PointsTable[t] = 0
		}
		N := len(matches)
		if N == 0 {
			json.NewEncoder(w).Encode("Schedule matches before playing")
			return
		}
		type lead struct {
			team   string
			points int
		}
		first := lead{team: "", points: 0}
		second := lead{team: "", points: 0}
		for i := 0; i < N-1; i++ {
			matchRes := rand.Intn(2)
			if matchRes == 0 {
				matches[i].Winner = matches[i].Team1
			} else {
				matches[i].Winner = matches[i].Team2
			}
			PointsTable[matches[i].Winner] += 2
			if PointsTable[matches[i].Winner] >= first.points {
				second.team = first.team
				second.points = first.points
				first.team = matches[i].Winner
				first.points = PointsTable[matches[i].Winner]
			} else if PointsTable[matches[i].Winner] >= second.points {
				second.team = matches[i].Winner
				second.points = PointsTable[matches[i].Winner]
			}
		}
		matches[N-1].Team1 = first.team
		matches[N-1].Team2 = second.team
		matchRes := rand.Intn(1)
		if matchRes == 0 {
			matches[N-1].Winner = matches[N-1].Team1
		} else {
			matches[N-1].Winner = matches[N-1].Team2
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"PointsTable": PointsTable,
			"Matches":     matches,
			"Winner":      matches[N-1].Winner})
	})
	http.ListenAndServe(":8082", router)
}
