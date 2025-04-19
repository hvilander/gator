package main

import (
	"fmt"
	"github.com/hvilander/gator/internal/config"
)

func main() {
	cfg := config.Config{}

	err := cfg.Read()
	if err != nil {
		fmt.Println(fmt.Errorf("totally borked: %w", err))
		return
	}

	cfg.SetUser("hv")
	fmt.Println(cfg)

}
