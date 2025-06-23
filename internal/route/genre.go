package routes

import (
	"sync"

	"github.com/julienschmidt/httprouter"
	"github.com/sidharth-rashwana/book/internal/controller"
)

func GenreRoutes(r *httprouter.Router, uc *controller.DBConnection, wg *sync.WaitGroup) {
	// routes
	r.GET("/genre/:id", uc.GetGenreById)
	r.GET("/genres", uc.GetGenres)
	r.POST("/genre", uc.CreateGenre)
	r.DELETE("/genre/:id", uc.DeleteGenreById)
	wg.Done()
}
