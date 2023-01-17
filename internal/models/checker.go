package models

import (
	"GoProxyChecker/internal/database"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type ProxyList struct {
	ID        int
	Type      string
	IP        string
	Port      int
	Speed     int
	AnonLVL   string
	City      string
	Country   string
	LastCheck time.Time
}

func FindingProxy(ch chan string, wg *sync.WaitGroup) {
	var p ProxyList
	dbPool := database.ConnectToDatabase()

	query := "SELECT id, types, ip, port, speed, anonlvl, city, country, last_check FROM proxy_list WHERE types = 'http'"
	rows, err := dbPool.Query(context.Background(), query)
	if err != nil {
		logrus.Errorf("Err request to database - %s", err)
		return
	}
	defer dbPool.Close()

	for rows.Next() {
		err := rows.Scan(&p.ID, &p.Type, &p.IP, &p.Port, &p.Speed, &p.AnonLVL, &p.City, &p.Country, &p.LastCheck)
		if err != nil {
			logrus.Errorf("Err scan data - %s", err)
			return
		}
		ch <- fmt.Sprintf("%s://%s:%s", p.Type, p.IP, strconv.Itoa(p.Port))
	}

	close(ch)
	defer wg.Done()
}
func Checker(ch chan string, wg *sync.WaitGroup) {
	for val := range ch {
		wg.Add(1)
		go func(val string) {
			CheckHTTP(val)
			defer wg.Done()
		}(val)
	}
	defer wg.Done()
}

func CheckHTTP(val string) {
	proxy, _ := url.Parse(val)
	client := http.Client{
		Timeout: time.Second * 5,
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxy),
		},
	}

	resp, err := client.Get("http://api.ipify.org")
	if err != nil {
		logrus.Errorf("Proxy invalid - %s", err)
		return
	}
	defer resp.Body.Close()

	logrus.Infof("Checked - %s %d\n", val, resp.StatusCode)
}
