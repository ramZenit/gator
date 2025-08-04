package main

import (
	"fmt"

	"github.com/ramZenit/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("Error reading the config file:", err)
		return
	}
	cfg.SetUser("ramZenit")
	newCfg, err := config.Read()
	if err != nil {
		fmt.Println("Error reading the config file:", err)
		return
	}
	fmt.Printf("main: %v\n", newCfg)
}
