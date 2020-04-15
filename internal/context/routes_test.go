package context

import (
	"github.com/gin-gonic/gin"
	"testing"
)

func TestRouter(t *testing.T) {
	router := gin.Default()
	RegisterRoute(router)
	router.Run(":8021")
}
