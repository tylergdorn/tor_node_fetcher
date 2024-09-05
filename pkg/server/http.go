package server

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tylergdorn/prophet_exercise/pkg/storage"
)

type TorNodeRetriever interface {
	GetTorNodes() []storage.TorExitNode
}

type AllowListHandler interface {
	InsertAllowIP(entry storage.AllowListIP)
	DeleteAllowIP(IP string)
	GetAllowListIPs() []storage.AllowListIP
}

type Server struct {
	retriever TorNodeRetriever
	allowList AllowListHandler
}

func StartServer(retriever TorNodeRetriever, allowList AllowListHandler, port string) error {
	serv := Server{retriever: retriever, allowList: allowList}
	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.GET("/nodes", serv.nodes)
	allowlistGroup := router.Group("/allowlist")
	serv.allowListRoutes(*allowlistGroup)
	return router.Run(port)
}

func (s *Server) nodes(c *gin.Context) {
	slog.DebugContext(c.Request.Context(), "got request for nodeshandler")
	nodes := s.retriever.GetTorNodes()
	c.JSON(http.StatusOK, nodes)
}

func (s *Server) allowListRoutes(group gin.RouterGroup) {
	group.GET("", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, s.allowList.GetAllowListIPs())
	})
	group.POST("/:ip", func(ctx *gin.Context) {
		entry := storage.AllowListIP{
			IP:        ctx.Param("ip"),
			Note:      ctx.Query("note"),
			TimeAdded: time.Now(),
		}
		s.allowList.InsertAllowIP(entry)
		ctx.Status(http.StatusOK)
	})
	group.DELETE("/:ip", func(ctx *gin.Context) {
		s.allowList.DeleteAllowIP(ctx.Param("ip"))
	})
}
