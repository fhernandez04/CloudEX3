package books

import (
	"context"
	"fmt"
	"log"
	"slices"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func PrepareDatabase(client *mongo.Client, dbName, collecName string) (*mongo.Collection, error) {
	db := client.Database(dbName)
	names, err := db.ListCollectionNames(context.TODO(), bson.D{{}})
	if err != nil {
		return nil, err
	}
	if !slices.Contains(names, collecName) {
		cmd := bson.D{{Key: "create", Value: collecName}}
		var result bson.M
		if err = db.RunCommand(context.TODO(), cmd).Decode(&result); err != nil {
			log.Fatal(err)
			return nil, err
		}
	}
	return db.Collection(collecName), nil
}

func PrepareData(coll *mongo.Collection) {
	startData := []BookStore{
		{
			ID:          "example1",
			BookName:    "The Vortex",
			BookAuthor:  "José Eustasio Rivera",
			BookEdition: "958-30-0804-4",
			BookPages:   "292",
			BookYear:    "1924",
		},
		{
			ID:          "example2",
			BookName:    "Frankenstein",
			BookAuthor:  "Mary Shelley",
			BookEdition: "978-3-649-64609-9",
			BookPages:   "280",
			BookYear:    "1818",
		},
		{
			ID:          "example3",
			BookName:    "The Black Cat",
			BookAuthor:  "Edgar Allan Poe",
			BookEdition: "978-3-99168-238-7",
			BookPages:   "280",
			BookYear:    "1843",
		},
	}

	for _, book := range startData {
		cursor, err := coll.Find(context.TODO(), book)
		if err != nil {
			panic(err)
		}
		var results []BookStore
		if err = cursor.All(context.TODO(), &results); err != nil {
			panic(err)
		}
		if len(results) > 1 {
			log.Fatal("more records were found")
		} else if len(results) == 0 {
			result, err := coll.InsertOne(context.TODO(), book)
			if err != nil {
				panic(err)
			} else {
				fmt.Printf("%+v\n", result)
			}
		} else {
			for _, res := range results {
				cursor.Decode(&res)
				fmt.Printf("%+v\n", res)
			}
		}
	}
}
