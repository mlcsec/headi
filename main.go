// uncomment crypto/tls import and TLSClient configs if required
package main

import (
    "net"
    "net/http"
    "net/url"
    //"crypto/tls"
    "fmt"
    "flag"
    "log"
    "os"
    "bufio"
    "time"
    "strconv"
    "github.com/fatih/color"
)

var headers = []string{"Client-IP","Connection","Contact","Forwarded","From","Host","Origin","Referer","True-Client-IP","X-Client-IP","X-Custom-IP-Authorization","X-Forward-For","X-Forwarded-For","X-Forwarded-Host","X-Forwarded-Server","X-Host","X-HTTP-Host-Override","X-Original-URL","X-Originating-IP","X-Real-IP","X-Remote-Addr","X-Remote-IP","X-Rewrite-URL","X-Wap-Profile"}
var inject = []string{"127.0.0.1","localhost","0.0.0.0","0","127.1","127.0.1","2130706433"}
var urlt string
var pfile string
var to int

func payloadInject() {
    timeout := time.Duration(to * 1000000)
    g := color.New(color.FgGreen)
    r := color.New(color.FgRed)
    var tr = &http.Transport{
		MaxIdleConns:      30,
		IdleConnTimeout:   time.Second,
		DisableKeepAlives: true,
        //TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
		DialContext: (&net.Dialer{
			Timeout:   timeout,
			KeepAlive: time.Second,
		}).DialContext,
	}
	client := &http.Client{
		Transport:     tr,
		Timeout:       timeout,
	}

    // open and iterate
    file, err := os.Open(pfile)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()
    scanner := bufio.NewScanner(file)

    // baseline request - gauge normal response
    breq, err := http.Get(urlt)
    if err != nil {
        log.Fatal(err)
    }

    // loop through payload file and inject
    for scanner.Scan() {
        for _, header := range headers {
            req, err := http.NewRequest("GET", urlt, nil)
            req.Header.Set(header, scanner.Text())
            resp, err := client.Do(req)
            if err != nil {
                continue
            }
            if breq.ContentLength != resp.ContentLength {
                g.Println("[+] "+"["+urlt+"]"+" "+"["+header+": "+scanner.Text()+"]"+" "+"[Code: "+strconv.Itoa(int(resp.StatusCode))+"] "+"[Size: "+ strconv.Itoa(int(resp.ContentLength))+"]")
            } else {
                r.Println("[-] "+"["+urlt+"]"+" "+"["+header+": "+scanner.Text()+"]"+" "+"[Code: "+strconv.Itoa(int(resp.StatusCode))+"] "+"[Size: "+ strconv.Itoa(int(resp.ContentLength))+"]")
            }
            defer resp.Body.Close()
        }
    }
}

func headerInject() {
    timeout := time.Duration(to * 1000000)
    g := color.New(color.FgGreen)
    r := color.New(color.FgRed)
    var tr = &http.Transport{
		MaxIdleConns:      30,
		IdleConnTimeout:   time.Second,
		DisableKeepAlives: true,
        //TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
		DialContext: (&net.Dialer{
			Timeout:   timeout,
			KeepAlive: time.Second,
		}).DialContext,
	}
	client := &http.Client{
		Transport:     tr,
		Timeout:       timeout,
	}

    // baseline request - gauge normal response
    breq, err := http.Get(urlt)
    if err != nil {
        log.Fatal(err)
    }

    // loop through default payloads and inject
    for _, header := range headers {
        for _, i := range inject {
            req, err := http.NewRequest("GET", urlt, nil)
            req.Header.Set(header, i)
            resp, err := client.Do(req)
            if err != nil {
                continue
            }
            if breq.ContentLength != resp.ContentLength {
                g.Println("[+] "+"["+urlt+"]"+" "+"["+header+": "+i+"]"+" "+"[Code: "+strconv.Itoa(int(resp.StatusCode))+"] "+"[Size: "+ strconv.Itoa(int(resp.ContentLength))+"]")
            } else {
                r.Println("[-] "+"["+urlt+"]"+" "+"["+header+": "+i+"]"+" "+"[Code: "+strconv.Itoa(int(resp.StatusCode))+"] "+"[Size: "+ strconv.Itoa(int(resp.ContentLength))+"]")
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
    flag.StringVar(&urlt, "url", "", "")
    flag.StringVar(&urlt, "u", "", "")
    flag.StringVar(&pfile, "pfile","","")
    flag.StringVar(&pfile, "p","","")
    flag.IntVar(&to, "timeout", 10000, "")
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
            fmt.Println("Invalid URL: ",urlt)
            os.Exit(1)
        }
        if pfile != "" {
            payloadInject()
        } else {
            headerInject()
        }
    }
}
