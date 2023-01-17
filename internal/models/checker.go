package models

import (
	"GoProxyChecker/internal/database"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/proxy"
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

	query := "SELECT id, types, ip, port, speed, anonlvl, city, country, last_check FROM proxy_list"
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
		ch <- fmt.Sprintf("%s:%s", p.IP, strconv.Itoa(p.Port))
	}

	close(ch)
	defer wg.Done()
}
func Checker(ch chan string, wg *sync.WaitGroup) {
	for val := range ch {
		wg.Add(1)
		go func(val string) {
			defer wg.Done()
			if strings.Contains(val, "http") {
				CheckHTTP(val)
			}

			if strings.Contains(val, "socks") {
				CheckSOCKS(val)
			}
		}(val)

	}
	defer wg.Done()
}

func CheckHTTP(val string) {
	proxy, err := url.Parse(fmt.Sprintf("%s:%s", "http", val))
	if err != nil {
		logrus.Errorf("Err parce url - %s", err)
		return
	}
	client := http.Client{
		Timeout: time.Second * 5,
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxy),
		},
	}
	RequesToCheck(val, client)
}

func CheckSOCKS(val string) {
	proxy, err := proxy.SOCKS5("tcp", val, nil, proxy.Direct)
	if err != nil {
		logrus.Errorf("Proxy invalid - %s", err)
		return
	}
	client := http.Client{
		Timeout: time.Second * 5,
		Transport: &http.Transport{
			Dial: proxy.Dial,
		},
	}
	RequesToCheck(val, client)
}

func RequesToCheck(val string, client http.Client) {
	resp, err := client.Get("http://api.ipify.org")
	if err != nil {
		logrus.Errorf("Proxy invalid - %s", err)
		return
	}
	defer resp.Body.Close()

	logrus.Infof("Checked - %s [%d]\n", val, resp.StatusCode)
}
