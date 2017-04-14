package main

import (
	"github.com/reechou/robot-rmdmy/config"
	"github.com/reechou/robot-rmdmy/controller"
)

func main() {
	controller.NewLogic(config.NewConfig()).Run()
}
