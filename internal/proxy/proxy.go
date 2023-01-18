package proxy

import (
	"GoProxyChecker/internal/models"
	"GoProxyChecker/pkg/database"
	httpcheck "GoProxyChecker/pkg/http_check"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	"h12.io/socks"
)

// Выборка прокси из базы
func FindProxy(ch chan string, wg *sync.WaitGroup) {
	var p models.ProxyList
	dbPool := database.ConnectToDatabase()

	query := "SELECT id, types, ip, port, speed, anonlvl, city, country, last_check FROM proxy_list"
	rows, err := dbPool.Query(context.Background(), query)
	if err != nil {
		logrus.Errorf("Err request to database - %s", err)
		return
	}

	for rows.Next() {
		err := rows.Scan(&p.ID, &p.Type, &p.IP, &p.Port, &p.Speed, &p.AnonLVL, &p.City, &p.Country, &p.LastCheck)
		if err != nil {
			logrus.Errorf("Err scan data - %s", err)
			return
		}
		ch <- fmt.Sprintf("%s://%s:%s", p.Type, p.IP, strconv.Itoa(p.Port))
	}

	close(ch)
	wg.Done()
}

/*
	 Своеобразный конвеер который фильтрует прокси
		по их типу читая адреса из канала
*/
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

// Для HTTP(S)
func CheckHTTP(val string) {
	proxy, err := url.Parse(val)
	if err != nil {
		logrus.Error(err)
	}

	logrus.Debug("run CheckHTTP args: %s", val)
	if err := requesToCheck(val, http.ProxyURL(proxy)); err != nil {
		logrus.Error(err)
	}
}

// Для SOCKS...SOCKS4..SOCKS5..SOCKS4A
func CheckSOCKS(val string) {
	logrus.Debug("run CheckSOCKS args: %s", val)
	if err := requesToCheck(val, socks.Dial(val)); err != nil {
		logrus.Error(err)
	}
}

// Отправка запроса на ресурс через проксю
func requesToCheck(val string, proxy interface{}) error {
	var resp models.ProxyRespons
	proxyChecker := httpcheck.HttpCheckClient{}

	resp, err := proxyChecker.Check(proxy)
	if err != nil {
		return err
	}

	if !resp.IsSuccess() {
		return errors.New("Proxy invalid - " + val)
	}

	logrus.Infof("Checked - %s [%d]\n", val, resp.GetStatusCodeRaw())
	return nil
}
