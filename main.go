package main

import (
    "net"
    "net/http"
    "net/url"
    "fmt"
    "flag"
    "log"
    "os"
    "bufio"
    "time"
    "strconv"
)

var headers = []string{"X-Originating-IP","X-Forwarded-For","X-Remote-IP","X-Remote-Addr","X-Client-IP","X-Host","X-Forwarded-Host","Origin","Host"}
var inject = []string{"127.0.0.1","localhost","0.0.0.0","0","127.1","127.0.1","2130706433"}
var url string
var pfile string
var to int

// MAKE GLOBAL HTTP CLIENT SO DON'T KEEP HAVING TO DO THE var tr... BITS
func validateUrl(){
    u, err := url.ParseRequestURI(url)
    if err != nil {
        panic(err)
    }
}

func payloadInject() {
    timeout := time.Duration(to * 1000000)
    var tr = &http.Transport{
		MaxIdleConns:      30,
		IdleConnTimeout:   time.Second,
		DisableKeepAlives: true,
		DialContext: (&net.Dialer{
			Timeout:   timeout,
			KeepAlive: time.Second,
		}).DialContext,
	}
	client := &http.Client{
		Transport:     tr,
		Timeout:       timeout,
	}
    file, err := os.Open(pfile)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        for _, header := range headers {
            req, err := client.Get(url)
            req.Header.Set(header, scanner.Text())
            if err != nil {
                continue
            }
            fmt.Println("[*] "+"["+url+"]"+" "+"["+header+": "+scanner.Text()+"]"+" "+" [Code: "+strconv.Itoa(int(req.StatusCode))+"] "+"[Size: "+ strconv.Itoa(int(req.ContentLength))+"]")
            defer req.Body.Close()
        }
    }
}

func headerInject() {
    timeout := time.Duration(to * 1000000)
    var tr = &http.Transport{
		MaxIdleConns:      30,
		IdleConnTimeout:   time.Second,
		DisableKeepAlives: true,
		DialContext: (&net.Dialer{
			Timeout:   timeout,
			KeepAlive: time.Second,
		}).DialContext,
	}
	client := &http.Client{
		Transport:     tr,
		Timeout:       timeout,
	}
    for _, header := range headers {
        for _, i := range inject {
            req, err := client.Get(url)
            req.Header.Set(header, i)
            if err != nil {
                continue
            }
            fmt.Println("[*] "+"["+url+"]"+" "+"["+header+": "+i+"]"+" "+" [Code: "+strconv.Itoa(int(req.StatusCode))+"] "+"[Size: "+ strconv.Itoa(int(req.ContentLength))+"]")
            defer req.Body.Close()
        }
    }
}

func main() {
    flag.StringVar(&url, "url", "", "target URL")
    flag.StringVar(&pfile, "pfile","","payload file")
    flag.IntVar(&to, "t", 10000, "timeout (milliseconds)")
    flag.Parse()
    // FIX THIS PART, COVER ALL POSSIBILITIES
    if url == "" {
        flag.PrintDefaults()
    } else {
        validateUrl()
    }
    //if pfile != "" {
    //    payloadInject()
    //}
}
