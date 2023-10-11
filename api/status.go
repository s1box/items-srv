package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/s1box/dc-go-srv/db"
)

type CheckApi interface {
	Status(c *gin.Context)
	EnsureDb(c *gin.Context)
}

type checkApiImpl struct{}

var _ CheckApi = (*checkApiImpl)(nil)

func NewStatusRestApi() AddToApiToRouter {
	return &checkApiImpl{}
}

func (ca *checkApiImpl) AddApiToRouter(router *gin.Engine) {
	router.GET("/status", ca.Status)
	router.GET("/items/db", ca.EnsureDb)
	router.GET("/", ca.MainPage)
}

func (ca *checkApiImpl) Status(c *gin.Context) {
	dbConnConf := db.ReadDatabaseConnectionConfig()
	if err := db.Ping(dbConnConf); err != nil {
		msg := "failed to reach database, reason: " + err.Error()
		log.Print(msg)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"database": msg})
		return
	}

	c.Header("Access-Control-Allow-Origin", "*")
	c.IndentedJSON(http.StatusOK, gin.H{"database": "OK"})
}

func (ca *checkApiImpl) EnsureDb(c *gin.Context) {
	dbConnConf := db.ReadDatabaseConnectionConfig()

	createDbErr := db.CreateDatabase(dbConnConf)
	tableErr := db.EnsureTable(dbConnConf, db.ItemsTableName)
	if tableErr == nil {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "OK"})
		return
	}

	c.Header("Access-Control-Allow-Origin", "*")
	c.IndentedJSON(http.StatusInternalServerError, gin.H{
		"database-error": createDbErr.Error(),
		"table-error":    tableErr.Error(),
	})
}

func (ca *checkApiImpl) MainPage(c *gin.Context) {
	// c.IndentedJSON(http.StatusOK, gin.H{"message": "use /status path to check the status"})
	c.Redirect(http.StatusFound, "/status")
}
