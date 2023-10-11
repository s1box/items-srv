package api

import "github.com/gin-gonic/gin"

type AddToApiToRouter interface {
	AddApiToRouter(router *gin.Engine)
}
