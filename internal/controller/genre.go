package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/sidharth-rashwana/book/internal/database"
	models "github.com/sidharth-rashwana/book/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (uc *DBConnection) CreateGenre(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var genre models.CreateGenre
	err := json.NewDecoder(r.Body).Decode(&genre)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	genre.Name = strings.ToLower(genre.Name)
	var existingAuthor models.CreateGenre
	err = uc.db.Collection(database.GENRES).FindOne(context.Background(), bson.M{"name": genre.Name}).Decode(&existingAuthor)
	if err == nil {
		http.Error(w, "genre already exists", http.StatusBadRequest)
		return
	} else if err != mongo.ErrNoDocuments {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = uc.db.Collection(database.GENRES).InsertOne(context.Background(), genre)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	//Response indicating successful creation, without the ID
	//Response := map[string]string{"message": "genre created successfully"}
	response := map[string]interface{}{
		"status": http.StatusCreated,
		"data":   "genre created successfully",
	}
	json.NewEncoder(w).Encode(response)
}

func (uc *DBConnection) GetGenres(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var genres []models.GetAllGenres

	cursor, err := uc.db.Collection(database.GENRES).Find(context.Background(), bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	if err := cursor.All(context.Background(), &genres); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	response := map[string]interface{}{
		"status": http.StatusOK,
		"data":   genres,
	}
	json.NewEncoder(w).Encode(response)
}

func (uc *DBConnection) GetGenreById(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var genre models.GetGenre

	err = uc.db.Collection(database.GENRES).FindOne(context.Background(), bson.M{"_id": objID}).Decode(&genre)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	response := map[string]interface{}{
		"status": http.StatusOK,
		"data":   genre,
	}
	json.NewEncoder(w).Encode(response)
}

func (uc *DBConnection) DeleteGenreById(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	res, err := uc.db.Collection(database.GENRES).DeleteOne(context.Background(), bson.M{"_id": objID})
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
		"data":   fmt.Sprintf("genre deleted successfully with id: %s", id),
	}
	json.NewEncoder(w).Encode(response)

}
