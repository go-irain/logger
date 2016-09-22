package logger

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var remoteSockUrl string
var remoteServierId string

func SetRemoteUrl(url string) {
	remoteSockUrl = url
}
func SetRemoteServerId(id string) {
	remoteServierId = id
}

//http log requester
func httpLog(data string) {
	value := url.Values{}
	value.Add("id", remoteServierId)
	value.Add("msg", data)
	fmt.Println("id := ", remoteServierId, "url := ", remoteSockUrl, "data := ", data)
	logRes, err := Post(remoteSockUrl, value.Encode())
	if err != nil {
		fmt.Println("err", err.Error())
		return
	}
	fmt.Println("HttpLog Response : ", logRes)
}

func Post(url, param string) (string, error) {
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, time.Second*10) //设置建立连接超时
				if err != nil {
					return nil, err
				}
				c.SetDeadline(time.Now().Add(30 * time.Second)) //设置发送接收数据超时
				return c, nil
			},
		},
	}
	request, err := http.NewRequest("POST", url, strings.NewReader(param))
	if err != nil {
		fmt.Println(err.Error())
		return "", errors.New("SERVER_INNER_ERROR")
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error)
		if strings.Contains(err.Error(), "timeout") {
			err = errors.New("REQUEST_TIME_OUT")
		} else {
			err = errors.New("SERVER_INNER_ERROR")
		}
		return "", err
	}
	defer response.Body.Close()
	result, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err.Error())
		return "", errors.New("SERVER_INNER_ERROR")
	}
	return string(result), nil
}
