package main

import (
	"GoProxyChecker/internal/app"
	"GoProxyChecker/pkg"
	"sync"
)

func init() {
	pkg.ConfigLog()
}

func main() {
	ch := make(chan string)
	var wg sync.WaitGroup
	app.Run(ch, &wg)
}
