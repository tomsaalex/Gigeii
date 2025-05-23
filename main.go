package main

import (
	"fmt"
	"log"
	"net/http"

	"example.com/db"
	"example.com/handler"
	"example.com/repository"
	"example.com/service"
)

func main() {
	// connect to the db
	pool := db.NewConnectionPool()
	queries := db.New(pool)

	userRepository := repository.NewDBUserRepository(queries)

	argonHelper := service.StandardArgon2idHash()
	userService := service.NewUserService(userRepository, argonHelper)

	jwtHelper := service.NewJwtUtil()
	userHandler := handler.NewUserHandler(userService, jwtHelper)

	pageHandler := handler.NewPageHandler()

	r := handler.SetupRoutes(handler.RouteDependencies{
		UserHandler: userHandler,
		PageHandler: pageHandler,
		JwtHelper:   jwtHelper,
	})

	fmt.Println("Server is listening on :8008")
	log.Fatal(http.ListenAndServe(":8008", r))
}
