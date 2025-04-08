package Routes

import (
	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func GetGrafosHandler(driver neo4j.DriverWithContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx := c.Request.Context()

		instituicao := c.Query("instituicaoDoutorado")
		area := c.Query("areaDoutorado")
		nome := c.Query("nome")

		if instituicao == "" {
			instituicao = "null"
		}
		if area == "" {
			area = "null"
		}
		if nome == "" {
			nome = "null"
		}

		session := driver.NewSession(ctx, neo4j.SessionConfig{
			AccessMode: neo4j.AccessModeRead,
		})
		defer session.Close(ctx)

		query := `
				CALL {
					  WITH $nome AS nomeFiltro
					  OPTIONAL MATCH (n:Pesquisador)
					  WHERE nomeFiltro IS NULL
					  RETURN COLLECT(DISTINCT n) AS resultados
					
					  UNION
					
					  WITH $nome AS nomeFiltro
					  MATCH (root:Pesquisador)
					  WHERE nomeFiltro IS NOT NULL AND toLower(root.nome) CONTAINS toLower(nomeFiltro)
					  OPTIONAL MATCH (root)-[:ORIENTOU*]-(n:Pesquisador)
					  WITH COLLECT(DISTINCT n) + root AS allResults
					  RETURN allResults AS resultados
				}
				WITH DISTINCT resultados AS nodesList
				UNWIND nodesList AS n
				WITH n
				WHERE ($instituicao IS NULL OR n.instituicaoCorrespondente = $instituicao)
				  AND ($area IS NULL OR n.areaDoutorado = $area)
			
				OPTIONAL MATCH (n)-[r:ORIENTOU]->(m:Pesquisador)
			
				WITH COLLECT(DISTINCT n.instituicaoCorrespondente) as instituicaoCorrespondente,
					 COLLECT(DISTINCT n.areaDoutorado) as areas,
					 COLLECT(DISTINCT {
					   id: n.idLattes,
					   label: n.nome,
					   instituicaoCorrespondente: n.instituicaoCorrespondente,
					   areaDoutorado: n.areaDoutorado,
					   indicador_semente: n.indicador_semente
					 }) as nodes,
					 COLLECT(DISTINCT {
					   source: n.idLattes,
					   target: m.idLattes
					 }) as relationships
			
				RETURN {
				  instituicaoCorrespondente: instituicaoCorrespondente,
				  areas: areas,
				  nodes: nodes,
				  edges: [rel IN relationships WHERE rel.source IS NOT NULL AND rel.target IS NOT NULL]
				} as result
			`

		params := map[string]any{
			"instituicao": nullIfString(instituicao),
			"area":        nullIfString(area),
			"nome":        nullIfString(nome),
		}

		result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			records, err := tx.Run(ctx, query, params)
			if err != nil {
				return nil, err
			}
			if records.Next(ctx) {
				return records.Record().Values[0], nil
			}
			return nil, records.Err()
		})

		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, result)
	}
}

func nullIfString(s string) any {
	if s == "null" {
		return nil
	}
	return s
}
