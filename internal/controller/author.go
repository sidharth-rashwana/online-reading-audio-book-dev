package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	database "github.com/sidharth-rashwana/book/internal/database"
	models "github.com/sidharth-rashwana/book/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (uc *DBConnection) CreateAuthor(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var author models.CreateAuthor
	err := json.NewDecoder(r.Body).Decode(&author)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	author.Name = strings.ToLower(author.Name)

	var existingAuthor models.CreateAuthor
	err = uc.db.Collection(database.AUTHORS).FindOne(context.Background(), bson.M{"name": author.Name}).Decode(&existingAuthor)
	if err == nil {
		http.Error(w, "author already exists", http.StatusBadRequest)
		return
	} else if err != mongo.ErrNoDocuments {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = uc.db.Collection(database.AUTHORS).InsertOne(context.Background(), author)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := map[string]interface{}{
		"status": http.StatusCreated,
		"data":   "author created successfully",
	}
	json.NewEncoder(w).Encode(response)
}

func (uc *DBConnection) GetAuthors(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var authors []models.GetAllAuthors

	cursor, err := uc.db.Collection(database.AUTHORS).Find(context.Background(), bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	if err := cursor.All(context.Background(), &authors); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	response := map[string]interface{}{
		"status": http.StatusOK,
		"data":   authors,
	}
	json.NewEncoder(w).Encode(response)
}

func (uc *DBConnection) GetAuthorById(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var author models.GetAuthor

	err = uc.db.Collection(database.AUTHORS).FindOne(context.Background(), bson.M{"_id": objID}).Decode(&author)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	response := map[string]interface{}{
		"status": http.StatusOK,
		"data":   author,
	}
	json.NewEncoder(w).Encode(response)
}

func (uc *DBConnection) DeleteAuthorById(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	res, err := uc.db.Collection(database.AUTHORS).DeleteOne(context.Background(), bson.M{"_id": objID})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if res.DeletedCount == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"status": http.StatusOK,
		"data":   fmt.Sprintf("author deleted successfully with id: %s", id),
	}
	json.NewEncoder(w).Encode(response)
}
