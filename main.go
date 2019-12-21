package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	DEFAULT_QUALITY = 95
)

type ErrorMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type AdapterConfig struct {
	ImaginaryHost  string
	FilePathPrefix string
	Host           string
	Port           string
	DefaultType    string
}

// ImaginaryRequestParameters
type ImaginaryRequestParameters struct {
	Host    string
	File    string
	Method  string
	Width   int
	Height  int
	Quality int
	Type    string
}

var (
	methodIsNotDefinedError = errors.New("You must to define method")
	widthIsNotDefinedError  = errors.New("You must to define width")
	heightIsNotDefinedError = errors.New("You must to define height")
	incorrectHost           = errors.New("Defined host is not correct")
	config                  = AdapterConfig{}
)

func (self *ImaginaryRequestParameters) GetUrl() string {
	u, err := url.Parse(self.Host)
	if err != nil {
		log.Panic(incorrectHost)
	}
	u.Path = path.Join("/", self.Method)
	values := url.Values{}
	values.Add("width", strconv.Itoa(self.Width))
	values.Add("height", strconv.Itoa(self.Height))
	values.Add("file", self.File)
	if self.Type != "" {
		values.Add("type", self.Type)
	}
	u.RawQuery = values.Encode()
	log.WithFields(log.Fields{
		"Query":    u.Query().Encode(),
		"Hostname": u.Hostname(),
		"Path":     u.EscapedPath(),
		"String":   u.String(),
	}).Debug(`Get new url`)
	return u.String()
}

// parseRequest is function for translate incomming request to ImaginaryRequestParameters
func parseRequest(r *http.Request, host, prefix string) (*ImaginaryRequestParameters, error) {
	values := r.URL.Query()
	// method
	method := values.Get("method")
	if method == "" {
		return nil, methodIsNotDefinedError
	}
	// sizes
	width, err := strconv.Atoi(values.Get("width"))
	if err != nil {
		return nil, widthIsNotDefinedError
	}
	height, err := strconv.Atoi(values.Get("height"))
	if err != nil {
		return nil, heightIsNotDefinedError
	}
	// quality
	quality, _ := strconv.Atoi(values.Get("quality"))
	if quality == 0 {
		quality = DEFAULT_QUALITY
	}
	// filePath
	filePath := r.URL.Path
	if prefix != "" {
		filePath = strings.Replace(r.URL.Path, prefix, "", 1)
	}
	// type
	fileType := values.Get("type")
	if fileType == "" {
		fileType = config.DefaultType
	}

	return &ImaginaryRequestParameters{
		Host:    host,
		File:    filePath,
		Method:  method,
		Width:   width,
		Height:  height,
		Quality: quality,
		Type:    fileType,
	}, nil
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	parameters, err := parseRequest(r, config.ImaginaryHost, config.FilePathPrefix)
	if err != nil {
		log.WithFields(log.Fields{
			"request url": r.URL.String(),
		}).Error(err.Error())
		jmsg := &ErrorMessage{
			Message: err.Error(),
			Code:    404,
		}
		emsg, _ := json.Marshal(jmsg)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(emsg))
		return
	}
	req, _ := http.NewRequest("GET", parameters.GetUrl(), nil)
	req.Header = r.Header.Clone()
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error("error on getting image from imaginary", err)
	}
	if _, err := io.Copy(w, res.Body); err != nil {
		log.Error("error on write response", err)
	}
}

func main() {
	log.SetLevel(log.TraceLevel)
	// settings
	config.ImaginaryHost = os.Getenv("ADAPTER_IMAGINARY_HOST")
	if config.ImaginaryHost == "" {
		log.Fatal("ADAPTER_IMAGINARY_HOST is not defined")
	}
	config.FilePathPrefix = os.Getenv("ADAPTER_FILE_PATH_PREFIX")
	config.Port = os.Getenv("ADAPTER_PORT")
	if config.Port == "" {
		config.Port = "9000"
	}
	config.Host = os.Getenv("ADAPTER_HOST")
	if config.Host == "" {
		config.Host = "0.0.0.0"
	}

	http.HandleFunc("/", proxyHandler)
	log.Debugf(`
Starting server:  %v:%v;
Imaginary host:   %v;
File path prefix: %v
	`, config.Host, config.Port, config.ImaginaryHost, config.FilePathPrefix)
	s := &http.Server{
		Addr:         fmt.Sprintf("%v:%v", config.Host, config.Port),
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
	log.Fatal(s.ListenAndServe())
}
