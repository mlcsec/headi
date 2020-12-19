package main

import (
    //"net"
    //"net/http"
    "fmt"
    "flag"
    "log"
    "os"
    "bufio"
    //"time"
)

var headers = []string{"X-Originating-IP","X-Forwarded-For","X-Remote-IP","X-Remote-Addr","X-Client-IP","X-Host","X-Forwarded-Host","Origin","Host"}
var inject = []string{"127.0.0.1","localhost","0.0.0.0","0","127.1","127.0.1","2130706433"}
var url string
var pfile string
var to int

func payloadInject() {
    file, err := os.Open(pfile)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        for _, header := range headers {
            //fmt.Println(header+" "+scanner.Text())
            /*Client.Header.Set(header, scanner.Text())
            req, err := client.Get(url)
            if err != nil {
                continue
            }*/
            fmt.Println("[*] "+url+" "+header+": "+scanner.Text()+" "+"[Code: "+" "+"[Size: "+ "")
        }
    }
}

func headerInject() {
    //timeout := time.Duration(to * 1000000)
    /*var tr = &http.Transport{
		MaxIdleConns:      30,
		IdleConnTimeout:   time.Second,
		DisableKeepAlives: true,
		DialContext: (&net.Dialer{
			Timeout:   timeout,
			KeepAlive: time.Second,
		}).DialContext,
	}
	/*client := &http.Client{
		Transport:     tr,
		Timeout:       timeout,
	}*/

    for _, header := range headers {
        for _, i := range inject {
            //fmt.Println(header+" "+i)
            /*Client.Header.Set(header, i)
            req, err := client.Get(url)
            if err != nil {
                continue
            }*/
            fmt.Println("[*] "+url+" "+header+": "+i+" "+"[Code: "+" "+"[Size: "+ "")
        }
    }
}

func main() {
    flag.StringVar(&url, "url", "", "target URL")
    flag.StringVar(&pfile, "pfile","","payload file")
    flag.IntVar(&to, "t", 10000, "timeout (milliseconds)")
    flag.Parse()
    if pfile != "" {
        payloadInject()
    } else {
        headerInject()
    }
}
