package books

import "go.mongodb.org/mongo-driver/bson/primitive"

type BookStore struct {
	MongoID     primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	ID          string             `bson:"id" json:"id"`
	BookName    string             `bson:"title" json:"title"`
	BookAuthor  string             `bson:"author" json:"author"`
	BookEdition string             `bson:"edition,omitempty" json:"edition"`
	BookPages   string             `bson:"pages,omitempty" json:"pages"`
	BookYear    string             `bson:"year,omitempty" json:"year"`
}
