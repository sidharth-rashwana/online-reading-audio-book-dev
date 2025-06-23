package controller

import (
	"context"
	"encoding/json"
	"errors"
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

func (uc *DBConnection) CreateBook(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var book models.CreateBook
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var author models.GetAuthor
	err = uc.db.Collection(database.AUTHORS).FindOne(context.Background(), bson.M{"_id": book.AuthorID}).Decode(&author)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			http.Error(w, "author does not exist", http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var genre models.GetGenre
	err = uc.db.Collection(database.GENRES).FindOne(context.Background(), bson.M{"_id": book.GenreID}).Decode(&genre)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			http.Error(w, "genre does not exist", http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var existingBook models.CreateBook
	book.Name = strings.ToLower(book.Name)
	err = uc.db.Collection(database.BOOKS).FindOne(context.Background(), bson.M{"name": book.Name}).Decode(&existingBook)
	if err == nil {
		http.Error(w, "book already exists", http.StatusBadRequest)
		return
	} else if err != mongo.ErrNoDocuments {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	book.Name = strings.ToLower(book.Name)
	_, err = uc.db.Collection(database.BOOKS).InsertOne(context.Background(), book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := map[string]interface{}{
		"status": http.StatusCreated,
		"data":   "book created successfully",
	}
	json.NewEncoder(w).Encode(response)
}

func (uc *DBConnection) GetBooks(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	authorMap := make(map[string]string)
	genreMap := make(map[string]string)

	var books []map[string]interface{}

	cursor, err := uc.db.Collection(database.BOOKS).Find(context.Background(), bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var book models.GetAllBooks
		if err := cursor.Decode(&book); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		authorID := book.AuthorID.Hex()
		genreID := book.GenreID.Hex()

		authorName, ok := authorMap[authorID]
		if !ok {
			var author models.GetAuthor
			err := uc.db.Collection(database.AUTHORS).FindOne(context.Background(), bson.M{"_id": book.AuthorID}).Decode(&author)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			authorName = author.Name
			authorMap[authorID] = authorName
		}

		genreName, ok := genreMap[genreID]
		if !ok {
			var genre models.GetGenre
			err := uc.db.Collection(database.GENRES).FindOne(context.Background(), bson.M{"_id": book.GenreID}).Decode(&genre)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			genreName = genre.Name
			genreMap[genreID] = genreName
		}

		bookData := map[string]interface{}{
			"_id":          book.Id,
			"name":         book.Name,
			"author":       authorName,
			"bookPhotoURL": book.BookPhotoURL,
			"audioBookURL": book.AudioBookURL,
			"genreName":    genreName,
		}
		books = append(books, bookData)
	}

	if err := cursor.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	response := map[string]interface{}{
		"status": http.StatusOK,
		"data":   books,
	}
	json.NewEncoder(w).Encode(response)
}

func (uc *DBConnection) GetBookByAuthorID(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id") // authorId

	authorID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "invalid author ID", http.StatusNotFound)
		return
	}

	var author models.GetAuthor
	err = uc.db.Collection(database.AUTHORS).FindOne(context.Background(), bson.M{"_id": authorID}).Decode(&author)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			http.Error(w, "author does not exist", http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	genreMap := make(map[string]string)

	var books []map[string]interface{}
	cursor, err := uc.db.Collection(database.BOOKS).Find(context.Background(), bson.M{"authorId": authorID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var book models.GetAllBooksByGenreId
		if err := cursor.Decode(&book); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		genreID := book.GenreID.Hex()

		genreName, ok := genreMap[genreID]
		if !ok {
			var genre models.GetGenre
			err := uc.db.Collection(database.GENRES).FindOne(context.Background(), bson.M{"_id": book.GenreID}).Decode(&genre)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			genreName = genre.Name
			genreMap[genreID] = genreName
		}

		bookData := map[string]interface{}{
			"name":         book.Name,
			"author":       author.Name,
			"bookPhotoURL": book.BookPhotoURL,
			"audioBookURL": book.AudioBookURL,
			"genreName":    genreName,
		}
		books = append(books, bookData)
	}

	if err := cursor.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	response := map[string]interface{}{
		"status": http.StatusOK,
		"data":   books,
	}
	json.NewEncoder(w).Encode(response)
}

func (uc *DBConnection) GetBookByGenreID(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id") // genreId

	genreID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "invalid genre ID", http.StatusNotFound)
		return
	}

	var genre models.GetGenre
	err = uc.db.Collection(database.GENRES).FindOne(context.Background(), bson.M{"_id": genreID}).Decode(&genre)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			http.Error(w, "genre does not exist", http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	authorMap := make(map[string]string)

	var books []map[string]interface{}
	cursor, err := uc.db.Collection(database.BOOKS).Find(context.Background(), bson.M{"genreId": genreID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var book models.GetAllBooksByAuthorId
		if err := cursor.Decode(&book); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		authorID := book.AuthorID.Hex()

		authorName, ok := authorMap[authorID]
		if !ok {
			var author models.GetAuthor
			err := uc.db.Collection(database.AUTHORS).FindOne(context.Background(), bson.M{"_id": book.AuthorID}).Decode(&author)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			authorName = author.Name
			authorMap[authorID] = authorName
		}

		bookData := map[string]interface{}{
			"name":         book.Name,
			"author":       authorName,
			"bookPhotoURL": book.BookPhotoURL,
			"audioBookURL": book.AudioBookURL,
			"genreName":    genre.Name,
		}
		books = append(books, bookData)
	}

	if err := cursor.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	response := map[string]interface{}{
		"status": http.StatusOK,
		"data":   books,
	}
	json.NewEncoder(w).Encode(response)
}

func (uc *DBConnection) GetBookById(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var book models.GetAllBooks
	var authorName, genreName string

	err = uc.db.Collection(database.BOOKS).FindOne(context.Background(), bson.M{"_id": objID}).Decode(&book)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var author models.GetAuthor
	err = uc.db.Collection(database.AUTHORS).FindOne(context.Background(), bson.M{"_id": book.AuthorID}).Decode(&author)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	authorName = author.Name

	var genre models.GetGenre
	err = uc.db.Collection(database.GENRES).FindOne(context.Background(), bson.M{"_id": book.GenreID}).Decode(&genre)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	genreName = genre.Name

	w.Header().Set("Content-Type", "application/json")

	response := map[string]interface{}{
		"status": http.StatusOK,
		"data": map[string]interface{}{
			"book": map[string]interface{}{
				"_id":          book.Id,
				"name":         book.Name,
				"authorName":   authorName,
				"bookPhotoURL": book.BookPhotoURL,
				"audioBookURL": book.AudioBookURL,
				"genreName":    genreName,
			},
		},
	}
	json.NewEncoder(w).Encode(response)
}

func (uc *DBConnection) DeleteBookById(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	res, err := uc.db.Collection(database.BOOKS).DeleteOne(context.Background(), bson.M{"_id": objID})
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
		"data":   fmt.Sprintf("book deleted successfully with id: %s", id),
	}
	json.NewEncoder(w).Encode(response)

}
