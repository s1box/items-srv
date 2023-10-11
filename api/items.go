package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/s1box/dc-go-srv/db"
)

type ItemApi interface {
	GetItems(c *gin.Context)
	GetItemById(c *gin.Context)
	PostItem(c *gin.Context)
	DeleteItemById(c *gin.Context)
}

type itemsApiImpl struct {
	dbClient db.ItemsDbClient
}

var _ ItemApi = (*itemsApiImpl)(nil)

func NewItemsRestApi() AddToApiToRouter {
	dbConnConf := db.ReadDatabaseConnectionConfig()
	return &itemsApiImpl{
		dbClient: db.NewItemsDbClient(dbConnConf),
	}
}

func (ia *itemsApiImpl) AddApiToRouter(router *gin.Engine) {
	router.GET("/items", ia.GetItems)
	router.GET("/items/:id", ia.GetItemById)
	router.GET("/items/random", ia.GetRandomItem)
	router.POST("/items", ia.PostItem)
	router.DELETE("/items/:id", ia.DeleteItemById)
	router.OPTIONS("/items/:id", ia.OptionsItemById)
}

func (ia *itemsApiImpl) GetItems(c *gin.Context) {
	ia.addCommonHeaders(c)

	items, err := ia.dbClient.SelectAllItems()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, items)
}

func (ia *itemsApiImpl) GetItemById(c *gin.Context) {
	ia.addCommonHeaders(c)

	id := c.Param("id")
	item, err := ia.dbClient.GetItemById(id)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	if item == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("item with ID %s not found", id)})
		return
	}
	c.IndentedJSON(http.StatusOK, item)
}

func (ia *itemsApiImpl) GetRandomItem(c *gin.Context) {
	ia.addCommonHeaders(c)

	item := &db.Item{Id: -1, Name: "unexisting", Num: -1.1}
	var err error
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	if item == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "no items in the database"})
		return
	}
	c.IndentedJSON(http.StatusOK, item)
}

func (ia *itemsApiImpl) PostItem(c *gin.Context) {
	ia.addCommonHeaders(c)

	var newItem db.Item
	if err := c.BindJSON(&newItem); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if newItem.Name == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Item name cannot be empty"})
		return
	}
	id, err := ia.dbClient.InsertItem(newItem)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusCreated, gin.H{"id": id})
}

func (ia *itemsApiImpl) DeleteItemById(c *gin.Context) {
	ia.addCommonHeaders(c)

	id := c.Param("id")
	deleted, err := ia.dbClient.DeleteItemById(id)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	if !deleted {
		c.IndentedJSON(http.StatusNotFound, nil)
		return
	}
	c.IndentedJSON(http.StatusNoContent, nil)
}

func (ia *itemsApiImpl) OptionsItemById(c *gin.Context) {
	ia.addCommonHeaders(c)
	c.Header("Access-Control-Allow-Methods", "DELETE")
	c.IndentedJSON(http.StatusOK, nil)
}

func (ia *itemsApiImpl) addCommonHeaders(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
}
