package main

import (
	"GoProxyChecker/internal/app"
	"GoProxyChecker/pkg/log"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	ch := make(chan string)

	log.ConfigLog()

	app.Run(ch, &wg)
}
