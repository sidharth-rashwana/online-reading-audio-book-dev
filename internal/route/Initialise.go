package routes

import (
	"sync"

	"github.com/julienschmidt/httprouter"
	controllers "github.com/sidharth-rashwana/book/internal/controller"
	"go.mongodb.org/mongo-driver/mongo"
)

func InitalizeRoutes(db *mongo.Database) *httprouter.Router {
	var wg sync.WaitGroup

	r := httprouter.New()

	uc := controllers.DBConnector(db)

	wg.Add(3)

	go AuthorRoutes(r, uc, &wg)
	go BookRoutes(r, uc, &wg)
	go GenreRoutes(r, uc, &wg)

	wg.Wait()

	return r
}
