package api

import (
	"net/http"

	"distributed-counter-go/internal/cluster"
	"distributed-counter-go/internal/config"
	"distributed-counter-go/internal/counter"
	"distributed-counter-go/internal/model"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// NewRouter wires up HTTP routes and middleware.
func NewRouter(cfg config.Config, ctr *counter.Counter, mgr *cluster.Manager, logger *zap.Logger) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(loggingMiddleware(logger))

	router.POST("/counter/increment", func(c *gin.Context) {
		var req model.IncrementRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		local, _, err := ctr.Increment(req.Amount)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		mgr.BroadcastIncrement()

		c.JSON(http.StatusOK, model.IncrementResponse{
			Success:      true,
			Node:         cfg.NodeID,
			LocalCounter: local,
		})
	})

	router.GET("/counter/value", func(c *gin.Context) {
		global, active, _ := mgr.AggregateCounter()
		c.JSON(http.StatusOK, model.CounterValueResponse{
			GlobalCounter: global,
			ActiveNodes:   active,
		})
	})

	router.POST("/cluster/join", func(c *gin.Context) {
		var req model.JoinRequest
		if err := c.ShouldBindJSON(&req); err != nil || req.Peer == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		if mgr.RegisterPeer(req.Peer) {
			c.JSON(http.StatusOK, model.JoinResponse{Message: "peer added"})
			return
		}

		c.JSON(http.StatusOK, model.JoinResponse{Message: "peer already registered"})
	})

	router.GET("/cluster/peers", func(c *gin.Context) {
		c.JSON(http.StatusOK, model.PeersResponse{
			Node:  cfg.NodeID,
			Peers: mgr.GetPeers(),
		})
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, model.HealthResponse{Status: "ok", Node: cfg.NodeID})
	})

	router.POST("/sync", func(c *gin.Context) {
		var req model.SyncRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		if req.Peer != "" {
			mgr.RegisterPeer(req.Peer)
		}

		_, _, _ = ctr.MergeState(req.State)

		c.JSON(http.StatusOK, model.SyncResponse{
			Node:  cfg.NodeID,
			State: ctr.GetState(),
		})
	})

	return router
}

func loggingMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if logger != nil {
			logger.Info("request",
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.Int("status", c.Writer.Status()),
			)
		}
	}
}
