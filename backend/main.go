package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"context"
	"encoding/json"

	"github.com/gorilla/mux"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB Database Name
var dbName string = "TO_DO_LIST"

// MongoDB Collection Name
var colName string = "LIST"

// Task struct
type taskType struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title     string             `json:"title,omitempty"`
	Completed bool               `json:"completed,omitempty"`
	Details   string             `json:"details,omitempty"`
	Priority  string             `json:"priority,omitempty"`
}

// Returns router to handle all http requests in their own handlers
func Router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/task", GetAllTasks).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/task", CreateTask).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/task/{id}", TaskComplete).Methods("PUT", "OPTIONS")
	r.HandleFunc("/api/undoTask/{id}", UndoTask).Methods("PUT", "OPTIONS")
	r.HandleFunc("/api/deleteTask/{id}", DeleteTask).Methods("DELETE", "OPTIONS")
	return r
}

// collection object/instance
var collection *mongo.Collection

// empty array of Item structs
var taskTypes []taskType

func getDB() (*mongo.Database, error) {
	//Connects to database
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI("mongodb://localhost:27017").
		SetServerAPIOptions(serverAPIOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	//Gets database
	db := client.Database(dbName)
	return db, nil
}

// connect to database and return collection pointer
func getCollection(db *mongo.Database, colName string) *mongo.Collection {
	return db.Collection(colName)
}

func updateDocuments() {
	//gets all documents in collection
	cursor, err := collection.Aggregate(context.TODO(), mongo.Pipeline{})
	if err != nil {
		log.Fatal(err) // if error getting the collection, log it and exit program
	}
	defer cursor.Close(context.TODO())
	taskTypes = nil
	//goes through all documents in collection and adds them to the items to the webpage
	for cursor.Next(context.TODO()) {
		var result bson.M
		err := cursor.Decode(&result)
		if err != nil {
			log.Fatal(err) // if error decoding the documents, log it and exit program
		}
		taskTypes = append(taskTypes,
			taskType{
				Title:     result["title"].(string),
				ID:        result["_id"].(primitive.ObjectID),
				Completed: result["completed"].(bool),
				Details:   result["details"].(string),
				Priority:  result["priority"].(string),
			})
	}
}

// GetAllTasks handler
func GetAllTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	updateDocuments()
	json.NewEncoder(w).Encode(taskTypes)
}

// CreateTask handler
func CreateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	var newTask taskType
	_ = json.NewDecoder(r.Body).Decode(&newTask)
	insertOneTask(newTask)
	json.NewEncoder(w).Encode(newTask)
}

// TaskComplete handler
func TaskComplete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	options := mux.Vars(r)
	taskComplete(options["id"])
	json.NewEncoder(w).Encode(options["id"])
}

// UndoTask handler
func UndoTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	options := mux.Vars(r)
	undoTask(options["id"])
	json.NewEncoder(w).Encode(options["id"])
}

// DeleteTask handler
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	options := mux.Vars(r)
	deleteOneTask(options["id"])
	json.NewEncoder(w).Encode(options["id"])
}

// insert one task into mongo server
func insertOneTask(task taskType) {
	_, err := collection.InsertOne(context.Background(), task)
	if err != nil {
		log.Fatal(err)
	}
}

// completes one task from mongo server using the tasks' id
func taskComplete(task string) {
	id, _ := primitive.ObjectIDFromHex(task)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"completed": true}}
	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
}

// undo completion from one task from mongo server using the tasks' id
func undoTask(task string) {
	id, _ := primitive.ObjectIDFromHex(task)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"completed": false}}
	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
}

// delete one task from mongo server using the tasks' id
func deleteOneTask(task string) {
	id, _ := primitive.ObjectIDFromHex(task)
	filter := bson.M{"_id": id}
	_, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	//request handler
	r := Router()

	//mongodb connection
	db, err := getDB()
	if err != nil {
		log.Fatal(err)
	}

	//initializes the collection
	collection = getCollection(db, colName)

	// Starts server
	fmt.Println("Server started on port 9000")
	log.Fatal(http.ListenAndServe(":9000", r))
}
