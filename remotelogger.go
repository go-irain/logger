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

var remoteSockURL string
var remoteServierID string

//SetRemoteUrl 设置告警地址
func SetRemoteUrl(url string) {
	remoteSockURL = url
}

//SetRemoteServerId 设置服务id
func SetRemoteServerId(id string) {
	remoteServierID = id
}

//http log requester
func httpLog(data string) {
	if remoteSockURL == "" || remoteServierID == "" {
		return
	}
	value := url.Values{}
	value.Add("id", remoteServierID)
	value.Add("msg", data)
	//fmt.Println("id := ", remoteServierID, "url := ", remoteSockURL, "data := ", data)
	_, err := Post(remoteSockURL, value.Encode())
	if err != nil {
		fmt.Println("err", err.Error())
		return
	}
	//fmt.Println("HttpLog Response : ", logRes)
}

//Post http请求发送日志
func Post(url, param string) (string, error) {
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, time.Second*5) //设置建立连接超时
				if err != nil {
					return nil, err
				}
				c.SetDeadline(time.Now().Add(10 * time.Second)) //设置发送接收数据超时
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
		fmt.Println(err.Error())
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
