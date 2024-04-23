package api

import (
	db "DigitalBank/db/sqlc"
	"DigitalBank/token"
	"DigitalBank/util"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"net/http"
	"time"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("currency", validCurrency)
		if err != nil {
			return nil, err
		}
	}

	server.setupRouter()
	return server, nil
}

func faviconMiddleware(c *gin.Context) {
	if c.Request.URL.Path == "/favicon.ico" {
		c.String(http.StatusOK, "")
		c.Abort()
	}
	c.Next()
}

func (server *Server) setupRouter() {
	router := gin.Default()
	router.Use(faviconMiddleware)

	router.GET("/hello", server.helloWorld)
	//router.POST("/users", server.createUser)
	//router.POST("/users/login", server.loginUser)
	//router.POST("/tokens/renew_access", server.renewAccessToken)
	//
	//authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	//authRoutes.POST("/accounts", server.createAccount)
	//authRoutes.GET("/accounts/:id", server.getAccount)
	//authRoutes.GET("/accounts", server.listAccounts)
	//
	//authRoutes.POST("/transfers", server.createTransfer)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func (server *Server) helloWorld(context *gin.Context) {
	context.String(http.StatusOK, "Hello, Seattle!")
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (server *Server) Stop(ctx context.Context, HttpAddress string) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := server.router.Run(HttpAddress); err != nil && err != http.ErrServerClosed {
		return err
	}

	<-ctxTimeout.Done()

	return nil
}
