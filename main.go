package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


type Todo struct{
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Completed bool `json:"completed"`
	Body string `json:"body"`
}

var collection *mongo.Collection

func main(){
	fmt.Println("Hello")

	err:=godotenv.Load(".env")

	if err!=nil{
		log.Fatal("Error loading .env file",err)
	}

	MONGODB_URI:=os.Getenv("MONGODB_URI")
	clientOptions:=options.Client().ApplyURI(MONGODB_URI)

	client,err:=mongo.Connect(context.Background(),clientOptions)

	if err!=nil{
		log.Fatal(err)
	}

	err=client.Ping(context.Background(),nil)

	defer client.Disconnect(context.Background())

	if err!=nil{
		log.Fatal(err)
	}

	fmt.Println("Conntected to MONGODB ATLAS")

	collection=client.Database("golang_db").Collection("todos")

	app:=fiber.New()
	
	app.Get("/api/todos",getTodos)
	app.Post("/api/todos",createTodo)
	app.Patch("/api/todos/:id",updateTodo)
	app.Delete("/api/todos/:id",deleteTodo)

	port:=os.Getenv("PORT")

	if port==""{
		port="5000"
	}
	log.Fatal(app.Listen("0.0.0.0:"+port))
}

func getTodos(c* fiber.Ctx)error {
var todos[] Todo

cursor,err:=collection.Find(context.Background(),bson.M{})

if err!=nil{
	return err
}

defer cursor.Close(context.Background())

fmt.Println("Cursor",cursor)
fmt.Println("Context.Background()",context.Background())
 
for cursor.Next(context.Background()){
	var todo Todo
	if err:=cursor.Decode(&todo); err!=nil{
		return err
	}
// fmt.Println("Todo",todo)

	todos=append(todos, todo)
}
return c.JSON(todos)
}

func createTodo(c*fiber.Ctx)error{
	//both are same
	todo:=new(Todo)
	// todo :=&Todo{}
	if err:=c.BodyParser(todo)
	err!=nil{
		return err
	}

	if todo.Body==""{
		return c.Status(400).JSON(fiber.Map{"error":"Todo body cannot be empty"})
	}

	insertResult,err:=collection.InsertOne(context.Background(),todo)
	if err!=nil{
		return err
	}
fmt.Println("todo.ID=",todo)

	todo.ID=insertResult.InsertedID.(primitive.ObjectID)
fmt.Println("todo.ID=",todo)
	return c.Status(201).JSON(todo)

}
func updateTodo(c*fiber.Ctx)error{
id:=c.Params("id")
objectId,err:=primitive.ObjectIDFromHex(id)

if err!=nil{
	return c.Status(400).JSON(fiber.Map{"error":"Invalid todo Id"})
}
filter:=bson.M{"_id":objectId}
update:=bson.M{"$set":bson.M{"completed":true}}

sto,err:=collection.UpdateOne(context.Background(),filter,update)

fmt.Printf("Value of sto: %+v\n", sto)

if err!=nil{
	return c.Status(400).JSON(fiber.Map{"error":"Invalid todo Id"})
}

return c.Status(200).JSON(fiber.Map{"success":true})
}

func deleteTodo(c*fiber.Ctx)error{
id:=c.Params("id")

objectId,err:=primitive.ObjectIDFromHex(id)

if err!=nil{
	return c.Status(400).JSON(fiber.Map{"error":"invalid"})
}
filter:=bson.M{"_id":objectId}
_,err=collection.DeleteOne(context.Background(),filter)
if err!=nil{
	return err
}

return c.Status(200).JSON(fiber.Map{"success":true})
}