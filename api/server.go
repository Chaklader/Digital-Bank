package api

import (
	"context"
	"fmt"
	db "github.com/Chaklader/DigitalBank/db/sqlc"
	"github.com/Chaklader/DigitalBank/token"
	"github.com/Chaklader/DigitalBank/util"
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
	httpServer *http.Server
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	router := gin.Default()

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
		httpServer: &http.Server{
			Addr:    config.HTTPServerAddress,
			Handler: router,
		},
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation(currencyTag, validCurrency)
		if err != nil {
			return nil, err
		}
	}

	server.setupRouter(router)
	return server, nil
}

func (server *Server) setupRouter(router *gin.Engine) {

	router.Use(faviconMiddleware)

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	router.POST("/tokens/renew_access", server.renewAccessToken)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccounts)
	authRoutes.POST("/transfers", server.createTransfer)

	server.router = router
}

func (server *Server) Start() error {
	return server.httpServer.ListenAndServe()
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (server *Server) Stop(ctx context.Context) error {

	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := server.httpServer.Shutdown(ctxTimeout)
	if err != nil {
		return err
	}

	<-ctxTimeout.Done()

	return nil
}
