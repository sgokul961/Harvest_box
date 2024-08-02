package main

import (
	"HarvestBox/config"
	"HarvestBox/routes"
	"fmt"
	"log"
	"net/http"

	"github.com/rs/cors"
	"github.com/urfave/negroni"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	if err := config.MigrateTables(); err != nil {
		log.Fatalln("Error performing migration:", err)
	}

	myCors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		ExposedHeaders:   []string{"Authorization"},
	})

	n := negroni.Classic()
	n.Use(myCors)
	n.UseHandler(routes.Routes())

	server := http.Server{
		Addr:    cfg.Port,
		Handler: n,
	}
	fmt.Printf("Server listening on port %s...\n", cfg.Port)
	err = server.ListenAndServe()
	if err != nil {
		fmt.Println("Error:", err)
	}
}
