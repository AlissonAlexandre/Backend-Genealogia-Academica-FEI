package main

import (
	"backend_genealogia_academica/Routes"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func main() {
	neo4jHandler, err := NewNeo4jHandler()
	if err != nil {
		log.Fatalf("Erro ao inicializar o handler do Neo4j: %s", err)
	}
	defer neo4jHandler.Close()

	r := gin.Default()

	r.GET("/grafos", Routes.GetGrafosHandler(neo4jHandler.Driver))
	r.GET("/metagrafos", Routes.GetMetagrafosHandler(neo4jHandler.Driver))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8093"
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Erro ao iniciar o servidor:  %s", err)
	}
	log.Printf("Servidor rodando na porta %s", port)

}
