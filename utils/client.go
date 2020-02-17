package utils

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
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

// AddVerification adds client-end certs for double verify and also
func (web *WebClient) AddVerification(skip bool, caCrtPath, cliCrtPath, cliKeyPath string) {
	pool := x509.NewCertPool()
	// First add caCrt for server-end to verify.
	caCrt, err := ioutil.ReadFile(caCrtPath)
	if err != nil {
		log.Fatal(err)
	}
	pool.AppendCertsFromPEM(caCrt)
	// Second Loads client-end key and cert.
	cliCrt, err := tls.LoadX509KeyPair(cliCrtPath, cliKeyPath)
	if err != nil {
		log.Fatal(err)
	}
	var transport *http.Transport
	if skip {
		transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:            pool,
				Certificates:       []tls.Certificate{cliCrt},
				InsecureSkipVerify: true,
			},
		}
	} else {
		transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      pool,
				Certificates: []tls.Certificate{cliCrt},
			},
		}
	}
	web.Client.Transport = transport
}

// Get implements the get method.
func (web *WebClient) Get(url, token string) (response string, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", token)
	resp, err := web.Client.Do(req)
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
	req.Header.Set("resume", "true")
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
func (web *WebClient) UploadFile(url, path, token string) error {
	bodyBuffer := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuffer)
	fileWriter, err := bodyWriter.CreateFormFile("files", path)
	if err != nil {
		return err
	}
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// io.Copy(fileWriter, file)
	p := make([]byte, 300)
	_, err = file.Read(p)
	if err != nil {
		return err
	}
	_, err = fileWriter.Write(p)
	if err != nil {
		return err
	}
	boundary := bodyWriter.Boundary()
	closeBuf := bytes.NewBufferString(fmt.Sprintf("\r\n--%s--\r\n", boundary))
	// only last messages need a closeBuf at the end.
	requestReader := io.MultiReader(bodyBuffer, closeBuf)

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()
	req, err := http.NewRequest("POST", url, requestReader)
	if err != nil {
		return err
	}
	req.Header.Set("content-type", contentType)
	req.Header.Set("authorization", token)
	_, err = web.Client.Do(req)
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

// InfoAboutMe requires info about me.
func (web *WebClient) InfoAboutMe(url, token string) (string, error) {
	bodyBuffer := &bytes.Buffer{}
	bodyBuffer.Write([]byte("requires my resume"))
	req, err := http.NewRequest("POST", url, bodyBuffer)
	if err != nil {
		return "", err
	}
	req.Header.Set("content-type", "text/plain")
	req.Header.Set("authorization", token)
	resp, err := web.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	content := string(result)
	return content, nil
}

// GetKey gets key for decrypt my resume.
func (web *WebClient) GetKey(url, token string) (string, error) {
	bodyBuffer := &bytes.Buffer{}
	bodyBuffer.Write([]byte("requires my key"))
	req, err := http.NewRequest("POST", url, bodyBuffer)
	if err != nil {
		return "", err
	}
	req.Header.Set("content-type", "text/plain")
	req.Header.Set("authorization", token)
	resp, err := web.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	content := string(result)
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

// GetClientInfo gets clients' info
func (web *WebClient) GetClientInfo(url string) (string, string, error) {
	resp, err := web.Client.Get(url)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	var respContent map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&respContent); err != nil {
		return "", "", err
	}
	clientID := respContent["CLIENT_ID"].(string)
	clientSecret := respContent["CLIENT_SECRET"].(string)
	return clientID, clientSecret, nil
}

// GetClientToken gets token string back.
// Note thsat url here is the root path.
func (web *WebClient) GetClientToken(url string) (string, error) {
	clientID, clientSecret, err := web.GetClientInfo(url + "/register")
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("GET", url+"/token", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Client_ID", clientID)
	req.Header.Set("Client_Secret", clientSecret)
	resp, err := web.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	token := string(result)
	return token, nil
}
