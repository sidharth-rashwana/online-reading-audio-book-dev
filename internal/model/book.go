package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateBook struct {
	Name         string             `json:"name" bson:"name" validate:"required"`
	AuthorID     primitive.ObjectID `json:"authorId" bson:"authorId" validate:"required"`
	GenreID      primitive.ObjectID `json:"genreId" bson:"genreId" validate:"required"`
	BookPhotoURL string             `json:"bookPhotoURL" bson:"bookPhotoURL" validate:"required"` // book image url
	AudioBookURL string             `json:"audioBookURL" bson:"audioBookURL" validate:"required"` // audio of book url
}

type GetAllBooks struct {
	Id           string             `json:"_id" bson:"_id" validate:"required"`
	Name         string             `json:"name" bson:"name" validate:"required"`
	GenreID      primitive.ObjectID `json:"genreId" bson:"genreId" validate:"required"`
	AuthorID     primitive.ObjectID `json:"authorId" bson:"authorId" validate:"required"`
	BookPhotoURL string             `json:"bookPhotoURL" bson:"bookPhotoURL" validate:"required"` // book image url
	AudioBookURL string             `json:"audioBookURL" bson:"audioBookURL" validate:"required"`
}

type GetBook struct {
	Name         string `json:"name" bson:"name" validate:"required"`
	Author       string `json:"author" bson:"author" validate:"required"`
	Genre        string `json:"genre" bson:"genre" validate:"required"`
	BookPhotoURL string `json:"bookPhotoURL" bson:"bookPhotoURL" validate:"required"` // book image url
	AudioBookURL string `json:"audioBookURL" bson:"audioBookURL" validate:"required"`
}

type GetAllBooksByAuthorId struct {
	Id           string             `json:"_id" bson:"_id" validate:"required"`
	Name         string             `json:"name" bson:"name" validate:"required"`
	AuthorID     primitive.ObjectID `json:"authorId" bson:"authorId" validate:"required"`
	BookPhotoURL string             `json:"bookPhotoURL" bson:"bookPhotoURL" validate:"required"` // book image url
	AudioBookURL string             `json:"audioBookURL" bson:"audioBookURL" validate:"required"`
}

type GetAllBooksByGenreId struct {
	Id           string             `json:"_id" bson:"_id" validate:"required"`
	Name         string             `json:"name" bson:"name" validate:"required"`
	GenreID      primitive.ObjectID `json:"genreId" bson:"genreId" validate:"required"`
	BookPhotoURL string             `json:"bookPhotoURL" bson:"bookPhotoURL" validate:"required"` // book image url
	AudioBookURL string             `json:"audioBookURL" bson:"audioBookURL" validate:"required"`
}
