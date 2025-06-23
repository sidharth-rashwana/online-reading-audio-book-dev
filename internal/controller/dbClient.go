package controller

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type DBConnection struct {
	db *mongo.Database
}

func DBConnector(db *mongo.Database) *DBConnection {
	return &DBConnection{db}
}
