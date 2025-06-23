package environment

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found")
	}
}

func InitalizeEnv() (string, string, string, string) {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("unable to fetch URI")
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		log.Fatal("unable to fetch DB Name")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("unable to fetch port")
	}

	host := os.Getenv("HOST")
	if host == "" {
		log.Fatal("unable to fetch host")
	}
	return uri, dbName, port, host
}
