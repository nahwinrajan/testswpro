package handler

import "github.com/nahwinrajan/testswpro/repository"

type Server struct {
	repository repository.Repositorier
}

// New return reference to new instance of Server
func New(repo repository.Repositorier) *Server {
	return &Server{
		repository: repo,
	}
}
