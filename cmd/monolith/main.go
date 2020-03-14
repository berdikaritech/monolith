package main

import (
	"log"
	"net/http"

	"github.com/berdikaritech/monolith/internal/pkg/middleware"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/rs/cors"
)

func main() {
	router := httprouter.New()
	//router.Handler("GET", "/")

	corsOptions := cors.Options{
		AllowedHeaders:  []string{"Authorization", "Content-Type"},
		AllowedMethods:  []string{"GET", "OPTIONS", "POST"},
		AllowOriginFunc: func(origin string) bool { return true },
	}

	chain := alice.New(
		cors.New(corsOptions).Handler,
		//middleware.RequestIDHandler,
		middleware.LogHandler,
		//middleware.AuthorizationHandler(res.Resolvers.(*resolver.Resolver)),
	).Then(router)

	log.Fatal(http.ListenAndServe(":8080", chain))
}
