package main

import (
    "context"
    "fmt"
	"log"
//	"reflect"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type TODO struct{
	Title   string    `json:"title"`
	Note    string    `json:"note"`
	DueDate time.Time `json:"due_date"`
}

//CreateConnection creates a mongodb connection
func CreateConnection() *mongo.Client {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
	return client
}

//InsertData inserts data to collection
 func InsertData(client *mongo.Client){
		// Set collection name
		collection := client.Database("TodoApp").Collection("todolist")

		// Some data to insert
		yoga := TODO{"Do Yoga", "After Brushing your tooth perform yoga.",time.Now()}
		shower := TODO{"Take Shower", "After yoga Take Shower.",time.Now()}
		getready := TODO{"Get Ready", "After Shower, Take Breakfast and get ready.",time.Now() }
	
		// Insert a single document
		insertResult, err := collection.InsertOne(context.TODO(), yoga)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Created a single todo: ", insertResult.InsertedID)
	
		// Insert multiple documents
		todolist := []interface{}{shower, getready}
	
		insertManyResult, err := collection.InsertMany(context.TODO(), todolist)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Created multiple Todo: ", insertManyResult.InsertedIDs)
 }

//UpdateData updates data
func UpdateData(client *mongo.Client){
		// Set collection name
		collection := client.Database("TodoApp").Collection("todolist")

		// Where condition or what data to update
		where := bson.D{{"title", "Do Yoga"}}

		update := bson.M{"$set": bson.M{"note": "With yoga also perform microexcercise."}}

		updateResult, err := collection.UpdateOne(context.TODO(), where, update)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
}

//FindData finds data
func FindData(client *mongo.Client){
		// Set collection name
		collection := client.Database("TodoApp").Collection("todolist")

		// Find a single document
		var result TODO

		// Where condition for search
		where := bson.D{{"title", "Do Yoga"}}

		err := collection.FindOne(context.TODO(), where).Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
	
		fmt.Printf("Found a single document: %+v\n", result)
	
		findOptions := options.Find()
		findOptions.SetLimit(2)
	
		var results []*TODO
	
		// Finding multiple documents returns a cursor
		cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
		if err != nil {
			log.Fatal(err)
		}
	
		// Iterate through the cursor
		for cur.Next(context.TODO()) {
			var elem TODO
			err := cur.Decode(&elem)
			if err != nil {
				log.Fatal(err)
			}
	
			results = append(results, &elem)
		}
	
		if err := cur.Err(); err != nil {
			log.Fatal(err)
		}
	
		// Close the cursor once finished
		cur.Close(context.TODO())
	
		fmt.Printf("Found multiple documents (array of pointers): %+v\n", results)
}

//DeleteData Deletes collection data
func DeleteData(client *mongo.Client){
		// Set collection name
		collection := client.Database("TodoApp").Collection("todolist")

		// Delete all the documents in the collection
		deleteResult, err := collection.DeleteMany(context.TODO(), bson.D{{}})
		if err != nil {
			log.Fatal(err)
		}
	
		fmt.Printf("Deleted %v documents in the todolist collection\n", deleteResult.DeletedCount)
	
		// Close the connection once no longer needed
		err = client.Disconnect(context.TODO())
	
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Println("Connection to MongoDB closed.")
		}

}

func main() {
	client := CreateConnection()
	InsertData(client)
	UpdateData(client)
	FindData(client)
	DeleteData(client)
}