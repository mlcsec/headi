# headi
Customisable and automated HTTP header injection.  Example run from the HTB machine Control:

<a href="https://asciinema.org/a/381187" target="_blank"><img src="https://asciinema.org/a/381187.svg" /></a>

TLS `InsecureSkipVerify` is not currently configured, if you want to disable security checks then feel free to uncomment `crypto/tls` in the imports and the `TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},` lines in http transport configuration and then build locally.

## Install
```
go get github.com/mlcsec/headi
```

## Headers
Injects the following HTTP headers:
* X-Originating-IP
* X-Forwarded-For
* X-Remote-IP
* X-Remote-Addr
* X-Client-IP
* X-Host
* X-Forwarded-Host
* Origin
* Host

An initial baseline request is made to gauge the normal response for the target resource.  Green indicates a change in the response and red no change.  `[+]` and `[-]` respectively.


## Info
Two options for HTTP header injection:

1. Default payloads (127.0.0.1, localhost, etc.) are injected into the headers mentioned above
2. Custom payloads can be supplied (e.g. you've enumerated some internal IPs or domains) using the `pfile` parameter

```
$ headi
Usage:
  headi -u https://target.com/resource
  headi -u https://target.com/resource -p internal_addrs.txt

Options:
  -p, --pfile <file>       Payload file
  -t, --timeout <millis>   HTTP Timeout
  -u, --url <url>          Target URL
```
Currently only takes one URL as input but you can easily bash script for numerous URLs like so:
```
$ for i in $(cat urls); do headi -url $i;done
```
