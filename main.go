// uncomment crypto/tls import and TLSClient configs if required
package main

import (
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"crypto/tls"
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
	"github.com/fatih/color"
)

var headers = []string{"Client-IP", "Connection", "Contact", "Forwarded", "From", "Host", "Origin", "Referer", "True-Client-IP", "X-Client-IP", "X-Custom-IP-Authorization", "X-Forward-For", "X-Forwarded-For", "X-Forwarded-Host", "X-Forwarded-Server", "X-Host", "X-HTTP-Host-Override", "X-Original-URL", "X-Originating-IP", "X-Real-IP", "X-Remote-Addr", "X-Remote-IP", "X-Rewrite-URL", "X-Wap-Profile"}
var inject = []string{"127.0.0.1", "localhost", "0.0.0.0", "0", "127.1", "127.0.1", "2130706433"}
var urlt string
var pfile string
var to int

// ResolveContentLength returns right content-length value. First it re-make the request with chunked-compression
// disable, then it computes Content-Length from body length. If disabling encoding is efficient, client is modified for
// future requests also
func ResolveContentLength(tr *http.Transport, client *http.Client, url string, headerName string, headerValue string) (contentLentgh int64, err error) {
	//first disable chunked-compression
	tr.DisableCompression = true
	client.Transport = tr

	resp, _ := GetRequest(tr, client, url, headerName, headerValue)
	contentLentgh = resp.ContentLength

	if contentLentgh == -1 {
		tr.DisableCompression = false
		client.Transport = tr

		//compute Content-Length by ourself
		defer resp.Body.Close()
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return -1, err
		}
		contentLentgh = int64(len(bodyBytes))
	}
	return contentLentgh, err
}

//GetRequest performs the Get request giving an URL & Header info, and return the http.Response struct
func GetRequest(tr *http.Transport, client *http.Client, url string, headerName string, headerValue string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return resp, err
	}

	if headerName != "" {
		req.Header.Set(headerName, headerValue)
	}

	resp, err = client.Do(req)
	if err != nil {
		return resp, err
	}

	return resp, err
}

//InitBaseline performs baseline request and retrieve baselines valuesfrom it (content-length , client used, etc)
func InitBaseline() (tr *http.Transport, client *http.Client, bresp *http.Response) {
	timeout := time.Duration(to * 1000000)
	tr = &http.Transport{
		MaxIdleConns:      30,
		IdleConnTimeout:   time.Second,
		DisableKeepAlives: true,
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
		DialContext: (&net.Dialer{
			Timeout:   timeout,
			KeepAlive: time.Second,
		}).DialContext,
	}

	client = &http.Client{
		Transport: tr,
		Timeout:   timeout,
	}

	//baseline request - gauge normal response
	bresp, err := GetRequest(tr, client, urlt, "", "")
	if err != nil {
		log.Fatal("Failed to make baseline request:", err)
	}

	return tr, client, bresp
}

//HeaderInject performs a baseline request, retrieve content-lenght from it.
//Then it performs multiples request with various header and compare response Content-Length
//w/ baseline response
func HeaderInject() {
	g := color.New(color.FgGreen)
	r := color.New(color.FgRed)

	tr, client, bresp := InitBaseline()
	bContentLength := bresp.ContentLength

	if bContentLength == -1 { //the length is unknown
		var err error
		bContentLength, err = ResolveContentLength(tr, client, urlt, "", "")
		if err != nil {
			fmt.Println("Failed to retrieve Content-Lenght:", err)
		}
	}

	//loop through default payloads and inject
	for _, header := range headers {
		for _, i := range inject {
			resp, err := GetRequest(tr, client, urlt, header, i)
			if err != nil {
				continue
			}

			contentLength := resp.ContentLength
			if contentLength == -1 { //the length is unknown
				contentLength, err = ResolveContentLength(tr, client, urlt, header, i)
				if err != nil {
					continue
				}
			}

			if bContentLength != contentLength {
				g.Println("[+] " + "[" + urlt + "]" + " " + "[" + header + ": " + i + "]" + " " + "[Code: " + strconv.Itoa(int(resp.StatusCode)) + "] " + "[Size: " + strconv.Itoa(int(contentLength)) + "]")
			} else {
				r.Println("[-] " + "[" + urlt + "]" + " " + "[" + header + ": " + i + "]" + " " + "[Code: " + strconv.Itoa(int(resp.StatusCode)) + "] " + "[Size: " + strconv.Itoa(int(contentLength)) + "]")
			}
			defer resp.Body.Close()
		}
	}
}

func init() {
	flag.Usage = func() {
		f := "Usage:\n"
		f += "  headi -u https://target.com/resource\n"
		f += "  headi -u https://target.com/resource -p internal_addrs.txt\n\n"
		f += "Options:\n"
		f += "  -p, --pfile <file>       Payload File\n"
		f += "  -t, --timeout <millis>   HTTP Timeout\n"
		f += "  -u, --url <url>          Target URL\n"
		fmt.Fprintf(os.Stderr, f)
	}
}

func main() {
	flag.StringVar(&urlt, "url", "", "url to fetch and to which perform header injections")
	flag.StringVar(&urlt, "u", "", "url to fetch and to which perform header injections")
	flag.StringVar(&pfile, "pfile", "", "Give a file with custom header value (useful to provide internal IP")
	flag.StringVar(&pfile, "p", "", "Give a file with custom header value (useful to provide internal IP")
	flag.IntVar(&to, "timeout", 10000, "Custom HTTP timeout")
	flag.IntVar(&to, "t", 10000, "")
	flag.Parse()
	if urlt == "" {
		flag.Usage()
	} else {
		u, err := url.Parse(urlt)
		if err != nil {
			log.Fatal(err)
		}
		if u.Scheme == "" || u.Host == "" || u.Path == "" {
			fmt.Println("Invalid URL: ", urlt)
			os.Exit(1)
		}

		if pfile != "" {
			//Reconstruct inject []string, containaing header value
			inject = nil
			// open and iterate
			file, err := os.Open(pfile)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				inject = append(inject, scanner.Text())
			}
		}
		HeaderInject()
	}
}
