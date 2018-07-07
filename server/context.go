package server

import "github.com/ragsagar/wolff/model"

type Context struct {
	Srv  *Server
	User *model.User
}
