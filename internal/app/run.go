package app

import (
	"GoProxyChecker/internal/models"
	"encoding/json"
	"fmt"
	"sync"
)

func Run() {
	// &wg var wg sync.WaitGroup
	var wg sync.WaitGroup
	invalid, valid := models.Checker(&wg)

	data1, _ := json.Marshal(invalid)
	data2, _ := json.Marshal(valid)

	fmt.Println(string(data1))
	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println(string(data2))
}
