package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	model "github.com/mfc-creations/mongoapi/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const connectionString = ""  // Use ur URI
const dbName = "netflix"
const colName = "watchlist"

var collection *mongo.Collection

func init()  {
	clientOption:=options.Client().ApplyURI(connectionString);
	client,err:= mongo.Connect(context.TODO(),clientOption)
	if err!=nil{
		log.Fatal(err)
	}
	fmt.Println("DB connected successfully!!")
	
	collection=client.Database(dbName).Collection(colName)
	fmt.Println("Collection instance ready")
}


func GetAllMovies(w http.ResponseWriter,r *http.Request)  {
	w.Header().Set("Content-Type","application/x-www-form-urlencode")
	allMovies:=getAllMovies()
	json.NewEncoder(w).Encode(allMovies)
}

func CreateMovie(w http.ResponseWriter,r *http.Request) {
	w.Header().Set("Content-Type","application/x-www-form-urlencode")
	w.Header().Set("Allow-control-Allow-Methods","POST")

	var movie model.Netflix
	_=json.NewDecoder(r.Body).Decode(&movie)
	insertOneMovie(movie)
	json.NewEncoder(w).Encode(movie)
}

func MarkAsWatched(w http.ResponseWriter,r *http.Request)  {
	w.Header().Set("Content-Type","application/x-www-form-urlencode")
	w.Header().Set("Allow-control-Allow-Methods","PUT")

	params:=mux.Vars(r);
	updateOneMovie(params["id"]);
	json.NewEncoder(w).Encode(params["id"])
}

func DeleteAMovie(w http.ResponseWriter,r *http.Request)  {
	w.Header().Set("Content-Type","application/x-www-form-urlencode")
	w.Header().Set("Allow-control-Allow-Methods","DELETE")

	params:=mux.Vars(r);
	deleteOneMovie(params["id"])
	json.NewEncoder(w).Encode(params["id"])

}

func DeleteAllMovie(w http.ResponseWriter,r *http.Request)  {
	w.Header().Set("Content-Type","application/x-www-form-urlencode")
	w.Header().Set("Allow-control-Allow-Methods","DELETE")
	count:= deleteAllMovie()
	json.NewEncoder(w).Encode(count)

}

//########### DB helpers ####################
func insertOneMovie(movie model.Netflix)  {
	inserted,err:=collection.InsertOne(context.Background(),movie);
	if err!=nil{
		log.Fatal(err)
	}
	fmt.Println("Inserted one movie with id: ",inserted.InsertedID)
}

func updateOneMovie(movieId string)  {
	id,_:=primitive.ObjectIDFromHex(movieId);
	filter:=bson.M{"_id":id};
	update:=bson.M{"$set":bson.M{"watched":true}}

	result,err:=collection.UpdateOne(context.Background(),filter,update);
	if err!=nil {
		log.Fatal(err)
	}
	fmt.Println("Modified count",result.ModifiedCount)
}

func deleteOneMovie(movieId string)  {
	id,_:=primitive.ObjectIDFromHex(movieId);
	filter:=bson.M{"_id":id};
	delete,err:=collection.DeleteOne(context.Background(),filter)
	if err!=nil {
		log.Fatal(err)
	}
	fmt.Println("Modified count",delete)
}

func deleteAllMovie() int64  {
	delete,err:=collection.DeleteMany(context.Background(),bson.D{{}},nil)
		if err!=nil {
		log.Fatal(err)
	}
	fmt.Println("Modified count",delete.DeletedCount)
	return delete.DeletedCount
}

func getAllMovies() []primitive.M{
	cur,err:=collection.Find(context.Background(),bson.D{{}});
	if err!=nil {
		log.Fatal(err)
	}
	var movies  []primitive.M
	for cur.Next(context.Background()){
		var movie bson.M
		err:=cur.Decode(&movie);
		if err!=nil {
		log.Fatal(err)
	}
	movies = append(movies, movie)
	}
	defer cur.Close(context.Background());
	return movies
}



