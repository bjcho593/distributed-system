package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

var (
	ctx    = context.Background()
	driver neo4j.DriverWithContext
)

func main() {
	// Obtener variables de entorno
	uri := os.Getenv("NEO4J_URI")
	user := os.Getenv("NEO4J_USER")
	pass := os.Getenv("NEO4J_PASSWORD")

	// Conexión a Neo4j
	var err error
	driver, err = neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(user, pass, ""))
	if err != nil {
		log.Fatal("No se pudo conectar a Neo4j:", err)
	}

	defer driver.Close(ctx)

	router := gin.Default()

	// Endpoint: estado
	router.GET("/match/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Match Service is running ✅"})
	})

	// Endpoint: crear partido
	router.POST("/match/create", func(c *gin.Context) {
		var payload struct {
			TeamA string `json:"teamA"`
			TeamB string `json:"teamB"`
		}
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
			return
		}

		session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
		defer session.Close(ctx)

		_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			query := `
				MERGE (a:Team {name: $teamA})
				MERGE (b:Team {name: $teamB})
				MERGE (a)-[:PLAYS]->(b)
			`
			params := map[string]any{
				"teamA": payload.TeamA,
				"teamB": payload.TeamB,
			}
			_, err := tx.Run(ctx, query, params)
			return nil, err
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Neo4j error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Match created"})
	})

	// Endpoint: obtener todos los partidos
	router.GET("/match/all", func(c *gin.Context) {
		session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
		defer session.Close(ctx)

		result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			query := `
				MATCH (a:Team)-[:PLAYS]->(b:Team)
				RETURN a.name AS teamA, b.name AS teamB
			`
			records, err := tx.Run(ctx, query, nil)
			if err != nil {
				return nil, err
			}

			var matches []map[string]string
			for records.Next(ctx) {
				record := records.Record()
				match := map[string]string{
					"teamA": record.Values[0].(string),
					"teamB": record.Values[1].(string),
				}
				matches = append(matches, match)
			}
			return matches, nil
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Neo4j read error"})
			return
		}

		c.JSON(http.StatusOK, result)
	})

	// Ejecutar servidor
	router.Run(":8083")
}
