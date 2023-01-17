package app

import (
	"GoProxyChecker/internal/models"
	"sync"
)

func Run(ch chan string, wg *sync.WaitGroup) {
	wg.Add(2)

	go models.FindingProxy(ch, wg)
	go models.Checker(ch, wg)

	wg.Wait()
}
