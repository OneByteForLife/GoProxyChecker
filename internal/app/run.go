package app

import (
	"GoProxyChecker/internal/proxy"
	"sync"
)

func Run(ch chan string, wg *sync.WaitGroup) {
	wg.Add(2)

	go proxy.FindProxy(ch, wg)
	go proxy.Checker(ch, wg)

	wg.Wait()
}
