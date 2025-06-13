package main

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"slices"
	"time"

	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Defines a "model" that we can use to communicate with the
// frontend or the database
// More on these "tags" like `bson:"_id,omitempty"`: https://go.dev/wiki/Well-known-struct-tags
type BookStore struct {
	MongoID     primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	ID          string             `bson:"id" json:"id"`
	BookName    string             `bson:"title" json:"title"`
	BookAuthor  string             `bson:"author" json:"author"`
	BookEdition string             `bson:"edition,omitempty" json:"edition"`
	BookPages   string             `bson:"pages,omitempty" json:"pages"`
	BookYear    string             `bson:"year,omitempty" json:"year"`
}

// Wraps the "Template" struct to associate a necessary method
// to determine the rendering procedure
type Template struct {
	tmpl *template.Template
}

// Preload the available templates for the view folder.
// This builds a local "database" of all available "blocks"
// to render upon request, i.e., replace the respective
// variable or expression.
// For more on templating, visit https://jinja.palletsprojects.com/en/3.0.x/templates/
// to get to know more about templating
// You can also read Golang's documentation on their templating
// https://pkg.go.dev/text/template
func loadTemplates() *Template {
	return &Template{
		tmpl: template.Must(template.ParseGlob("views/*.html")),
	}
}

// Method definition of the required "Render" to be passed for the Rendering
// engine.
// Contraire to method declaration, such syntax defines methods for a given
// struct. "Interfaces" and "structs" can have methods associated with it.
// The difference lies that interfaces declare methods whether struct only
// implement them, i.e., only define them. Such differentiation is important
// for a compiler to ensure types provide implementations of such methods.
func (t *Template) Render(w io.Writer, name string, data interface{}, ctx echo.Context) error {
	return t.tmpl.ExecuteTemplate(w, name, data)
}

// Here we make sure the connection to the database is correct and initial
// configurations exists. Otherwise, we create the proper database and collection
// we will store the data.
// To ensure correct management of the collection, we create a return a
// reference to the collection to always be used. Make sure if you create other
// files, that you pass the proper value to ensure communication with the
// database
// More on what bson means: https://www.mongodb.com/docs/drivers/go/current/fundamentals/bson/
func prepareDatabase(client *mongo.Client, dbName string, collecName string) (*mongo.Collection, error) {
	db := client.Database(dbName)

	names, err := db.ListCollectionNames(context.TODO(), bson.D{{}})
	if err != nil {
		return nil, err
	}
	if !slices.Contains(names, collecName) {
		cmd := bson.D{{"create", collecName}}
		var result bson.M
		if err = db.RunCommand(context.TODO(), cmd).Decode(&result); err != nil {
			log.Fatal(err)
			return nil, err
		}
	}

	coll := db.Collection(collecName)
	return coll, nil
}

// Here we prepare some fictional data and we insert it into the database
// the first time we connect to it. Otherwise, we check if it already exists.
func prepareData(client *mongo.Client, coll *mongo.Collection) {
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

	// This syntax helps us iterate over arrays. It behaves similar to Python
	// However, range always returns a tuple: (idx, elem). You can ignore the idx
	// by using _.
	// In the topic of function returns: sadly, there is no standard on return types from function. Most functions
	// return a tuple with (res, err), but this is not granted. Some functions
	// might return a ret value that includes res and the err, others might have
	// an out parameter.
	for _, book := range startData {
		cursor, err := coll.Find(context.TODO(), book)
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

// Generic method to perform "SELECT * FROM BOOKS" (if this was SQL, which
// it is not :D ), and then we convert it into an array of map. In Golang, you
// define a map by writing map[<key type>]<value type>{<key>:<value>}.
// interface{} is a special type in Golang, basically a wildcard...
func findAllBooks(coll *mongo.Collection) []BookStore {
	cursor, err := coll.Find(context.TODO(), bson.D{{}})
	var results []BookStore
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	return results
}

func main() {
	// Connect to the database. Such defer keywords are used once the local
	// context returns; for this case, the local context is the main function
	// By user defer function, we make sure we don't leave connections
	// dangling despite the program crashing. Isn't this nice? :D
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// TODO: make sure to pass the proper username, password, and port
	// client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	// client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://mongodb:testmongo@localhost:27017"))
	databaseUri := os.Getenv("DATABASE_URI")
	if databaseUri == "" {
		databaseUri = "mongodb://localhost:27017"
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(databaseUri))

	// This is another way to specify the call of a function. You can define inline
	// functions (or anonymous functions, similar to the behavior in Python)
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	// You can use such name for the database and collection, or come up with
	// one by yourself!
	coll, err := prepareDatabase(client, "exercise-1", "information")

	prepareData(client, coll)

	// Here we prepare the server
	e := echo.New()

	// Define our custom renderer
	e.Renderer = loadTemplates()

	// Log the requests. Please have a look at echo's documentation on more
	// middleware
	e.Use(middleware.Logger())

	e.Static("/css", "css")

	// Endpoint definition. Here, we divided into two groups: top-level routes
	// starting with /, which usually serve webpages. For our RESTful endpoints,
	// we prefix the route with /api to indicate more information or resources
	// are available under such route.
	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index", nil)
	})

	e.GET("/books", func(c echo.Context) error {
		books := findAllBooks(coll)
		return c.Render(200, "book-table", books)
	})

	e.GET("/authors", func(c echo.Context) error {
		// Search for all books in the collection
		cursor, err := coll.Find(context.TODO(), bson.M{})
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to fetch books")
		}
		defer cursor.Close(context.TODO())

		// Use a map to store unique authors
		authorsMap := make(map[string]bool)
		for cursor.Next(context.TODO()) {
			var book BookStore
			if err := cursor.Decode(&book); err == nil && book.BookAuthor != "" {
				authorsMap[book.BookAuthor] = true
			}
		}

		// Convert the map keys to a slice
		authors := make([]string, 0, len(authorsMap))
		for author := range authorsMap {
			authors = append(authors, author)
		}

		return c.Render(http.StatusOK, "authors-table", authors)
	})

	e.GET("/years", func(c echo.Context) error {
		// Search for all books in the collection
		cursor, err := coll.Find(context.TODO(), bson.M{})
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to fetch books")
		}
		defer cursor.Close(context.TODO())

		// Use a map to store unique years
		yearsMap := make(map[string]bool)
		for cursor.Next(context.TODO()) {
			var book BookStore
			if err := cursor.Decode(&book); err == nil && book.BookYear != "" {
				yearsMap[book.BookYear] = true
			}
		}

		// Convert the map keys to a slice
		years := make([]string, 0, len(yearsMap))
		for year := range yearsMap {
			years = append(years, year)
		}

		return c.Render(http.StatusOK, "years-table", years)
	})

	e.GET("/search", func(c echo.Context) error {
		return c.Render(200, "search-bar", nil)
	})

	e.GET("/create", func(c echo.Context) error {
		return c.NoContent(http.StatusNoContent)
	})

	// You will have to expand on the allowed methods for the path
	// `/api/route`, following the common standard.
	// A very good documentation is found here:
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Methods
	// It specifies the expected returned codes for each type of request
	// method.
	e.GET("/api/books", func(c echo.Context) error {
		books := findAllBooks(coll)
		return c.JSON(http.StatusOK, books)
	})

	e.POST("/api/books", func(c echo.Context) error {
		var book BookStore
		if err := c.Bind(&book); err != nil { // bind the request body to the book struct
			return c.String(http.StatusBadRequest, "Invalid input")
		}

		if book.ID == "" || book.BookName == "" || book.BookAuthor == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "id, title and author are required"})
		}

		// Check if the book already exists
		filter := bson.M{"id": book.ID}
		count, err := coll.CountDocuments(context.TODO(), filter)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to check for duplicates")
		}
		if count > 0 {
			return c.JSON(http.StatusConflict, map[string]string{"error": "Book with the same ID already exists"})
		}

		// Insert the new book
		result, err := coll.InsertOne(context.TODO(), book)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to insert book")
		}

		book.MongoID = result.InsertedID.(primitive.ObjectID)
		return c.JSON(http.StatusCreated, book)
	})

	e.PUT("/api/books/:id", func(c echo.Context) error {
		id := c.Param("id") // get the id from the URL parameter

		var book BookStore
		if err := c.Bind(&book); err != nil { // bind the request body to the book struct
			return c.String(http.StatusBadRequest, "Invalid input")
		}

		if book.BookName == "" || book.BookAuthor == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "title and author are required"})
		}

		book.ID = id // ensure the ID is set to the one from the URL

		filter := bson.M{"id": id}     // filter to find the book by ID
		update := bson.M{"$set": book} // update the book with the new data

		result, err := coll.UpdateOne(context.TODO(), filter, update) // perform the update
		// Check if the update was successful
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to update book")
		}

		if result.MatchedCount == 0 { // no book found with the given ID
			return c.String(http.StatusNotFound, "Book not found")
		}

		return c.JSON(http.StatusOK, book)
	})

	e.DELETE("/api/books/:id", func(c echo.Context) error {
		id := c.Param("id")                                   // get the id from the URL parameter
		filter := bson.M{"id": id}                            // filter to find the book by ID
		result, err := coll.DeleteOne(context.TODO(), filter) // perform the delete
		// Check if the delete was successful
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to delete book")
		}
		if result.DeletedCount == 0 { // no book found with the given ID
			return c.String(http.StatusNotFound, "Book not found")
		}

		// return status 200 success
		return c.NoContent(http.StatusOK)
	})

	// We start the server and bind it to port 3030. For future references, this
	// is the application's port and not the external one. For this first exercise,
	// they could be the same if you use a Cloud Provider. If you use ngrok or similar,
	// they might differ.
	// In the submission website for this exercise, you will have to provide the internet-reachable
	// endpoint: http://<host>:<external-port>
	e.Logger.Fatal(e.Start(":3030"))
}
