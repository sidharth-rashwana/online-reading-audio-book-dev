package routes

import (
	"sync"

	"github.com/julienschmidt/httprouter"
	"github.com/sidharth-rashwana/book/internal/controller"
)

func BookRoutes(r *httprouter.Router, uc *controller.DBConnection, wg *sync.WaitGroup) {

	r.GET("/book/:id", uc.GetBookById)
	r.GET("/books", uc.GetBooks)
	r.GET("/books/author/:id", uc.GetBookByAuthorID)
	r.GET("/books/genre/:id", uc.GetBookByGenreID)
	r.POST("/book", uc.CreateBook)
	r.DELETE("/book/:id", uc.DeleteBookById)
	wg.Done()
}
