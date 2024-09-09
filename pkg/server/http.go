package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tylergdorn/prophet_exercise/pkg/storage"
)

type TorNodeRetriever interface {
	GetTorNodes(context.Context) ([]storage.TorExitNode, error)
	GetTorNodesPaginated(context.Context, int, int) ([]storage.TorExitNode, error)
}

type AllowListHandler interface {
	InsertAllowIP(ctx context.Context, entry storage.AllowListIP) error
	DeleteAllowIP(ctx context.Context, IP string) error
	GetAllowListIPs(context.Context) ([]storage.AllowListIP, error)
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

func (s *Server) nodes(ctx *gin.Context) {
	slog.DebugContext(ctx.Request.Context(), "got request for nodeshandler")
	var nodes []storage.TorExitNode
	var err error
	limitQuery, limitExists := ctx.GetQuery("limit")
	offsetQuery, offsetExists := ctx.GetQuery("offset")
	if limitExists && offsetExists {
		limit, err := strconv.Atoi(limitQuery)
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("limit not integer"))
		}
		offset, err := strconv.Atoi(offsetQuery)
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("offset not integer"))
		}
		nodes, err = s.retriever.GetTorNodesPaginated(ctx, limit, offset)
	} else {
		nodes, err = s.retriever.GetTorNodes(ctx)
	}
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
	if len(ctx.Errors) > 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{"errors": ctx.Errors})
	} else {
		ctx.JSON(http.StatusOK, nodes)
	}
}

func (s *Server) allowListRoutes(group gin.RouterGroup) {
	group.GET("", func(ctx *gin.Context) {
		res, err := s.allowList.GetAllowListIPs(ctx)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
		}
		ctx.JSON(http.StatusOK, res)
	})
	group.POST("/:ip", func(ctx *gin.Context) {
		entry := storage.AllowListIP{
			IP:        ctx.Param("ip"),
			Note:      ctx.Query("note"),
			TimeAdded: time.Now(),
		}
		s.allowList.InsertAllowIP(ctx, entry)
		ctx.Status(http.StatusOK)
	})
	group.DELETE("/:ip", func(ctx *gin.Context) {
		s.allowList.DeleteAllowIP(ctx, ctx.Param("ip"))
	})
}
