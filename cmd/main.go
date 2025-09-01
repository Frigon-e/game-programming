package main

import (
	"SideProjectGames/gameoflife"
	"SideProjectGames/internal/config"
	"fmt"
	"os"
)

func main() {

	if err := run(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

}

func run() (err error) {
	var cfg config.AppConfig

	cfg, err = config.InitConfig()
	if err != nil {
		return err
	}

	if err := gameoflife.Run(cfg); err != nil {
		return err
	}

	return nil
}
