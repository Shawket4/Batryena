package main

import (
	"BatrynaBackend/Models"
	"BatrynaBackend/Routes"
)

func main() {
	Models.Setup()
	Routes.Setup()
}
