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

func FindingProxy() []ProxyList {
	var list []ProxyList

	dbPool := database.ConnectToDatabase()

	query := "SELECT id, types, ip, port, speed, anonlvl, city, country, last_check FROM proxy_list"
	rows, err := dbPool.Query(context.Background(), query)
	if err != nil {
		logrus.Errorf("Err request to database - %s", err)
	}
	// defer dbPool.Close()

	for rows.Next() {
		var p ProxyList
		err := rows.Scan(&p.ID, &p.Type, &p.IP, &p.Port, &p.Speed, &p.AnonLVL, &p.City, &p.Country, &p.LastCheck)
		if err != nil {
			logrus.Errorf("Err scan data - %s", err)
		}
		list = append(list, p)
	}
	dbPool.Close()

	return list
}

// wg *sync.WaitGroup
func Checker(wg *sync.WaitGroup) ([]ProxyList, []ProxyList) {
	list := FindingProxy()
	var invalid []ProxyList

	for idx, val := range list {
		wg.Add(1)
		go func(idx int, val ProxyList) {
			if val.Type == "http" {
				status, speed, checkTime := CheckHTTP(fmt.Sprintf("%s://%s:%s", val.Type, val.IP, strconv.Itoa(val.Port)))
				if speed < 1.0 {
					speed = 1.0
				}

				if !status {
					list[idx].LastCheck = checkTime
					invalid = append(invalid, ProxyList{
						ID:        val.ID,
						Type:      val.Type,
						IP:        val.IP,
						Port:      val.Port,
						Speed:     0,
						AnonLVL:   val.AnonLVL,
						City:      val.City,
						Country:   val.Country,
						LastCheck: val.LastCheck,
					})
					list = append(list[:idx], list[idx+1:]...)
				} else {
					list[idx].Speed = int(speed)
					list[idx].LastCheck = checkTime
				}
			}
			defer wg.Done()
		}(idx, val)
	}
	wg.Wait()
	return invalid, list
}

func CheckHTTP(addr string) (bool, float64, time.Time) {
	proxy, _ := url.Parse(addr)
	fmt.Println(addr)
	client := http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxy),
		},
	}

	resp, err := client.Get("http://api.ipify.org")
	if err != nil {
		return false, 0, time.Now()
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, 0, time.Now()
	}
	return true, client.Timeout.Seconds(), time.Now()
}
