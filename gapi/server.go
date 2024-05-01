package gapi

import (
	"fmt"
	db "github.com/Chaklader/DigitalBank/db/sqlc"
	"github.com/Chaklader/DigitalBank/pb"
	"github.com/Chaklader/DigitalBank/token"
	"github.com/Chaklader/DigitalBank/util"
	"github.com/Chaklader/DigitalBank/worker"
)

type Server struct {
	pb.UnimplementedDigitalBankServer
	config          util.Config
	store           db.Store
	tokenMaker      token.Maker
	taskDistributor worker.TaskDistributor
}

func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:          config,
		store:           store,
		tokenMaker:      tokenMaker,
		taskDistributor: taskDistributor,
	}

	return server, nil
}
