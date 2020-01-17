package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

// WebClient wraps the http.Client.
type WebClient struct {
	Client http.Client
}

// NewClient creates a Http Client with default timeout.
func NewClient() *WebClient {
	client := http.Client{Timeout: 100 * time.Second}
	return &WebClient{
		Client: client,
	}
}

// Get implements the get method.
func (web *WebClient) Get(url string) (response string, err error) {
	resp, err := web.Client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			return "", err
		}
		result.Write(buffer[0:n])
	}
	response = result.String()
	return response, nil
}

// Post posts a requests to server.
func (web *WebClient) Post(url, contentType string, data interface{}) (content string, err error) {
	jsonStr, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return "", err
	}
	req.Header.Add("content-type", contentType)
	defer req.Body.Close()
	resp, err := web.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	content = string(result)
	return content, nil
}

// Put puts resources to server
func (web *WebClient) Put(url, contentType, token string, data interface{}) (content string, err error) {
	jsonStr, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return "", err
	}
	req.Header.Add("content-type", contentType)
	req.Header.Add("authorization", token)
	defer req.Body.Close()
	resp, err := web.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	content = string(result)
	return content, nil
}

// GetLocalIP gets the interface's IP.
// VM's IP.
func GetLocalIP(ifname string) (string, error) {
	var localIP = ""
	iface, err := net.InterfaceByName(ifname)
	if err != nil {
		return localIP, err
	}
	if iface.Name == ifname {
		addrs, err := iface.Addrs()
		if err != nil {
			return localIP, err
		}
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok {
				if ipnet.IP.To4() != nil {
					localIP = ipnet.IP.String()
				}
			}
		}
	}
	if localIP == "" {
		return localIP, fmt.Errorf("local interface doesn't have an ip")
	}
	return localIP, nil
}
