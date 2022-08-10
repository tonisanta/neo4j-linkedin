package main

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"neo4j-linkedin/config"
	"neo4j-linkedin/handlers"
	"neo4j-linkedin/repositories"
)

func main() {
	settings, err := config.ReadConfig("config.json")
	if err != nil {
		panic("Error: couldn't load the config")
	}
	router := gin.Default()

	driver, err := config.NewDriver(settings)
	neo4jRepo := repositories.NewNeo4jRepository(driver)
	handler := handlers.RequestsHandler{Neo4jRepository: neo4jRepo}

	router.GET("/company/:name", handler.GetCompanyHandler)
	router.GET("/profile-views/:name", handler.GetProfileViewersHandler)

	router.POST("/person", handler.NewPersonHandler)
	router.POST("/company", handler.NewCompanyHandler)
	router.POST("/relationship-p2p", handler.NewRelationshipPersonToPersonHandler)
	router.POST("/relationship-p2c", handler.NewRelationshipPersonToCompanyHandler)

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.Run() // listen and serve on 0.0.0.0:8080
}
