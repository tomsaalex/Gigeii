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

	//trebuie sa conv asa
	// date:= "2025-06-06T00:00:00+02:00"
	// data,_:=time.Parse("2006-01-02T15:04:05-07:00", date)
	// fmt.Println(data.Format(time.DateOnly))
	// fmt.Println(data.UTC().Format(time.DateOnly))
	// loc:=time.FixedZone("Europe/Bucharest", 2*60*60)
	// fmt.Println(data.UTC().In(loc).Format(time.DateOnly))


	



	userRepository := repository.NewDBUserRepository(queries)
	availabilityRepository := repository.NewDBAvailabilityRepository(pool, queries)

	argonHelper := service.StandardArgon2idHash()
	userService := service.NewUserService(userRepository, argonHelper)
	availabilityService := service.NewAvailabilityService(availabilityRepository, pool, queries)

	jwtHelper := service.NewJwtUtil()
	userHandler := handler.NewUserHandler(userService, jwtHelper)
	availabilityHandler := handler.NewAvailabilityHandler(userService, availabilityService)
	pageHandler := handler.NewPageHandler()


	resellerRepo:= repository.NewDbResellerRepository(queries)
	resellerService := service.NewResellerService(resellerRepo)
	resellerHandler := handler.NewResellerHandler(resellerService, availabilityHandler)

	r := handler.SetupRoutes(handler.RouteDependencies{
		UserHandler:         userHandler,
		ResellerHandler:     resellerHandler,
		AvailabilityHandler: availabilityHandler,
		PageHandler:         pageHandler,
		JwtHelper:           jwtHelper,
	})

	fmt.Println("Server is listening on http://localhost:8008")
	log.Fatal(http.ListenAndServe(":8008", r))
}
