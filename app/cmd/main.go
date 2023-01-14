package main

import (
	"GoProxyChecker/internal/app"
	"GoProxyChecker/pkg"
)

func init() {
	pkg.ConfigLog()
}

func main() {
	app.Run()
}
