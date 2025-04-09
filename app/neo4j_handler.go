package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
	"os"
)

type Neo4jHandler struct {
	Driver neo4j.DriverWithContext
	Ctx    context.Context
}

func NewNeo4jHandler() (*Neo4jHandler, error) {
	//dotenvErr := godotenv.Load(".env")
	//if dotenvErr != nil {
	//	return nil, fmt.Errorf("erro carregando .env: %w", dotenvErr)
	//}

	ctx := context.Background()
	dbUri := os.Getenv("NEO4J_URI")
	dbUser := os.Getenv("NEO4J_USER")
	dbPassword := os.Getenv("NEO4J_PASSWORD")

	driver, err := neo4j.NewDriverWithContext(dbUri, neo4j.BasicAuth(dbUser, dbPassword, ""))
	if err != nil {
		return nil, fmt.Errorf("erro na criação do driver Neo4j: %w", err)
	}

	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar conectividade com Neo4j: %w", err)
	}

	fmt.Println("Conexão com Neo4j foi realizada com sucesso!")
	return &Neo4jHandler{Driver: driver, Ctx: ctx}, nil
}

func (handler *Neo4jHandler) Close() {
	if err := handler.Driver.Close(handler.Ctx); err != nil {
		log.Printf("Erro ao fechar o driver do Neo4j: %s", err)
	}
}
