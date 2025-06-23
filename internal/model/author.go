package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateAuthor struct {
	Name           string `json:"name" bson:"name" validate:"required"`
	AuthorPhotoURL string `json:"authorPhotoURL" bson:"authorPhotoURL" validate:"required"`
}

type GetAllAuthors struct {
	ID             primitive.ObjectID `json:"id" bson:"_id" validate:"required"`
	Name           string             `json:"name" bson:"name" validate:"required"`
	AuthorPhotoURL string             `json:"authorPhotoURL" bson:"authorPhotoURL" validate:"required"`
}

type GetAuthor struct  {
	Name           string `json:"name" bson:"name" validate:"required"`
	AuthorPhotoURL string `json:"authorPhotoURL" bson:"authorPhotoURL" validate:"required"`
}
