package main

import (
	"github.com/hvilander/gator/internal/config"
	"github.com/hvilander/gator/internal/database"
)

type commands struct {
	handlersByName map[string]func(*state, command) error
}

type state struct {
	config *config.Config
	db     *database.Queries
}

type command struct {
	name string
	args []string
}
