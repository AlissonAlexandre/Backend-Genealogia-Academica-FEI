package Routes

import (
	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func GetMetagrafosHandler(driver neo4j.DriverWithContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx := c.Request.Context()
		instituicao := c.Query("instituicaoDoutorado")

		session := driver.NewSession(ctx, neo4j.SessionConfig{
			AccessMode: neo4j.AccessModeRead,
		})
		defer session.Close(ctx)

		query := `
			MATCH (orientador:Pesquisador)-[:ORIENTOU]->(orientado:Pesquisador)
			WHERE orientador.instituicaoCorrespondente IS NOT NULL
			  AND orientado.instituicaoCorrespondente IS NOT NULL
			  AND (
			    $instituicao IS NULL OR
			    orientador.instituicaoCorrespondente = $instituicao OR
			    orientado.instituicaoCorrespondente = $instituicao
			  )
			WITH 
			  orientador.instituicaoCorrespondente AS sourceInst, 
			  orientado.instituicaoCorrespondente AS targetInst, 
			  count(*) AS weight
			RETURN sourceInst, targetInst, weight
			ORDER BY weight DESC
		`

		var instParam any = nil
		if instituicao != "" && instituicao != "null" {
			instParam = instituicao
		}

		params := map[string]any{
			"instituicao": instParam,
		}

		result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			records, err := tx.Run(ctx, query, params)
			if err != nil {
				return nil, err
			}

			var data []map[string]any
			for records.Next(ctx) {
				record := records.Record()
				data = append(data, map[string]any{
					"source": record.Values[0],
					"target": record.Values[1],
					"weight": record.Values[2],
				})
			}

			return data, records.Err()
		})

		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, result)
	}
}
