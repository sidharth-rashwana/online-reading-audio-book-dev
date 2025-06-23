package routes

import (
	"sync"

	"github.com/julienschmidt/httprouter"
	"github.com/sidharth-rashwana/book/internal/controller"
)

func AuthorRoutes(r *httprouter.Router, uc *controller.DBConnection, wg *sync.WaitGroup) {
	r.GET("/author/:id", uc.GetAuthorById)
	r.GET("/authors", uc.GetAuthors)
	r.POST("/author", uc.CreateAuthor)
	r.DELETE("/author/:id", uc.DeleteAuthorById)
	wg.Done()
}
