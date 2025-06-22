package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func main() {
	router := gin.Default()

	// Status endpoint
	router.GET("/match/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Match Service is running ✅"})
	})

	// Create match
	router.POST("/match/create", func(c *gin.Context) {
		var body struct {
			TeamA string `json:"teamA"`
			TeamB string `json:"teamB"`
		}
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		uri := os.Getenv("NEO4J_URI")
		user := os.Getenv("NEO4J_USER")
		pass := os.Getenv("NEO4J_PASSWORD")

		driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(user, pass, ""))
		if err != nil {
			log.Fatal(err)
		}
		defer driver.Close()

		session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
		defer session.Close()

		_, err = session.Run(
			`MERGE (a:Team {name: $teamA})
             MERGE (b:Team {name: $teamB})
             CREATE (a)-[:MATCHED_WITH]->(b)`,
			map[string]interface{}{"teamA": body.TeamA, "teamB": body.TeamB},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Neo4j error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Match created"})
	})

	port := "8083"
	fmt.Println("Match Service running on port", port)
	router.Run(":" + port)
}
