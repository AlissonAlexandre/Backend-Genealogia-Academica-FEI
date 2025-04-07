package main

import (
	"app/routes"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func main() {
	neo4jHandler, err := NewNeo4jHandler()
	if err != nil {
		log.Fatalf("Erro ao inicializar o handler Neo4j: %s", err)
	}
	defer neo4jHandler.Close()

	r := gin.Default()

	r.GET("/grafos", routes.GetGrafosHandler(&neo4jHandler.Driver))
	r.GET("/metagrafos", routes.GetMetaGrafosHandler(&neo4jHandler.Driver))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Servidor rodando na porta %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %s", err)
	}
}
