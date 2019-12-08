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

var (
	methodIsNotDefinedError = errors.New("You must to define method")
	widthIsNotDefinedError  = errors.New("You must to define width")
	heightIsNotDefinedError = errors.New("You must to define height")
	incorrectHost           = errors.New("Defined host is not correct")
)

type ErrorMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ImaginaryParameters
type ImaginaryParameters struct {
	Host    string
	File    string
	Method  string
	Width   int
	Height  int
	Quality int
}

func (self *ImaginaryParameters) GetUrl() string {
	u, err := url.Parse(self.Host)
	if err != nil {
		log.Panic(incorrectHost)
	}
	u.Path = path.Join("/", self.Method)
	values := url.Values{}
	values.Add("width", strconv.Itoa(self.Width))
	values.Add("height", strconv.Itoa(self.Height))
	values.Add("file", self.File)
	u.RawQuery = values.Encode()
	log.Debugf(`
Query: %v;
Hostname: %v;
Path: %v;
String: %v;`, u.Query().Encode(), u.Hostname(),
		u.EscapedPath(), u.String(),
	)
	return u.String()
}

// parseRequest is function for translate incomming request to ImaginaryParameters
func parseRequest(r *http.Request, host, prefix string) (*ImaginaryParameters, error) {
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

	return &ImaginaryParameters{
		Host:    host,
		File:    filePath,
		Method:  method,
		Width:   width,
		Height:  height,
		Quality: quality,
	}, nil
}

func main() {
	log.SetLevel(log.TraceLevel)
	// settings
	imaginaryHost := os.Getenv("ADAPTER_IMAGINARY_HOST")
	if imaginaryHost == "" {
		log.Fatal("ADAPTER_IMAGINARY_HOST is not defined")
	}
	filePathPrefix := os.Getenv("ADAPTER_FILE_PATH_PREFIX")
	port := os.Getenv("ADAPTER_PORT")
	if port == "" {
		port = "9000"
	}
	host := os.Getenv("ADAPTER_HOST")
	if host == "" {
		host = "0.0.0.0"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		parameters, err := parseRequest(r, imaginaryHost, filePathPrefix)
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
		res, err := http.Get(parameters.GetUrl())
		if err != nil {
			log.Fatal(err)
		}
		if _, err := io.Copy(w, res.Body); err != nil {
			log.Fatal(err)
		}
	})
	log.Debugf(`
Starting server:  %v:%v;
Imaginary host:   %v;
File path prefix: %v
	`, host, port, imaginaryHost, filePathPrefix)
	s := &http.Server{
		Addr:         fmt.Sprintf("%v:%v", host, port),
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
	log.Fatal(s.ListenAndServe())
}
