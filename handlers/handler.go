package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"neo4j-linkedin/models"
	"neo4j-linkedin/repositories"
	"net/http"
)

type RequestsHandler struct {
	Neo4jRepository *repositories.Neo4jRepository
}

func newItemHandler[T any](c *gin.Context, save func(item T) error) {
	var item T
	err := c.Bind(&item)
	if err != nil {
		log.Printf("Error while binding the body request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	err = save(item)
	if err != nil {
		log.Printf("Error while inserting item to repository: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, item)
}

func newItemWithValidationHandler[T any](c *gin.Context, isValid func(item T) bool, save func(item T) error) {
	var item T
	err := c.Bind(&item)
	if err != nil {
		log.Printf("Error while binding the body request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if !isValid(item) {
		log.Printf("Error: the item does not meet the requirements")
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	err = save(item)
	if err != nil {
		log.Printf("Error while inserting item to repository: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, item)
}

func getItemHandler[T any, V any](c *gin.Context, getDetails func(item T) (V, error)) {
	var item T
	err := c.BindUri(&item)

	if err != nil {
		log.Printf("Error while binding the Uri: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	itemDetails, err := getDetails(item)
	if err != nil {
		log.Printf("Error while getting item from repository: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, itemDetails)
}

func (handler *RequestsHandler) NewCompanyHandler(c *gin.Context) {
	log.Println("New Company request")
	newItemHandler(c, handler.Neo4jRepository.AddCompany)
}

func (handler *RequestsHandler) NewPersonHandler(c *gin.Context) {
	log.Println("New Person request")
	isValidPerson := func(p models.Person) bool {
		// define the custom logic for validation ...
		return len(p.Name) > 2
	}
	newItemWithValidationHandler(c, isValidPerson, handler.Neo4jRepository.AddPerson)
}

func (handler *RequestsHandler) NewRelationshipPersonToPersonHandler(c *gin.Context) {
	log.Println("New relationship p2p request")
	newItemHandler(c, handler.Neo4jRepository.AddRelationshipPersonToPerson)
}

func (handler *RequestsHandler) NewRelationshipPersonToCompanyHandler(c *gin.Context) {
	log.Println("New relationship p2c request")
	newItemHandler(c, handler.Neo4jRepository.AddRelationshipPersonToCompany)
}

func (handler *RequestsHandler) GetCompanyHandler(c *gin.Context) {
	log.Println("Get Company request")
	getItemHandler(c, handler.Neo4jRepository.GetCompanyInfo)
}

func (handler *RequestsHandler) GetProfileViewersHandler(c *gin.Context) {
	log.Println("Get profile viewers request")
	getItemHandler(c, handler.Neo4jRepository.GetProfileViewers)
}
