package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateGenre struct {
	Name          string `json:"name" bson:"name" validate:"required"`
	GenrePhotoURL string `json:"genrePhotoURL" bson:"genrePhotoURL" validate:"required"`
}

type GetAllGenres struct {
	ID            primitive.ObjectID `json:"id" bson:"_id" validate:"required"`
	Name          string             `json:"name" bson:"name" validate:"required"`
	GenrePhotoURL string             `json:"genrePhotoURL" bson:"genrePhotoURL" validate:"required"`
}

type GetGenre struct {
	Name          string `json:"name" bson:"name" validate:"required"`
	GenrePhotoURL string `json:"genrePhotoURL" bson:"genrePhotoURL" validate:"required"`
}