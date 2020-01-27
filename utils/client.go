package utils

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"time"
)

const (
	// DefaultVMPort for server to listen.
	DefaultVMPort uint32 = 8808
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
	n, _ := resp.Body.Read(buffer[0:])
	result.Write(buffer[0:n])
	response = result.String()
	return response, nil
}

// Post posts a requests to server.
func (web *WebClient) Post(url, contentType string, data io.Reader) (content string, err error) {
	req, err := http.NewRequest("POST", url, data)
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
func (web *WebClient) Put(url, contentType, token string, data io.Reader) (content string, err error) {
	req, err := http.NewRequest("PUT", url, data)
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
	fmt.Println(resp.Status)
	result, err := ioutil.ReadAll(resp.Body)
	content = string(result)
	return content, nil
}

// UploadFile uploads file to the server.
func (web *WebClient) UploadFile(url, path string) error {
	bodyBuffer := bytes.NewBufferString("")
	bodyWriter := multipart.NewWriter(bodyBuffer)
	_, err := bodyWriter.CreateFormFile("files", path)
	if err != nil {
		return err
	}
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	boundary := bodyWriter.Boundary()
	closeBuf := bytes.NewBufferString(fmt.Sprintf("\r\n--%s--\r\n", boundary))

	requestReader := io.MultiReader(bodyBuffer, file, closeBuf)
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()
	_, err = web.Post(url, contentType, requestReader)
	if err != nil {
		return err
	}
	return nil
}

// UploadData uploads map data to server.
func (web *WebClient) UploadData(url string, data map[string]string) error {
	bodyBuffer := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuffer)
	for key, value := range data {
		_ = bodyWriter.WriteField(key, value)
	}
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	_, err := web.Post(url, contentType, bodyBuffer)
	if err != nil {
		return err
	}
	return nil
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
