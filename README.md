# headi
Customisable and automated HTTP header injection

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

First makes a baseline request to gauge the normal response for the target resource, if the response differs then the URL, header, payload, status code, and content length is returned.


## Examples
Two options for HTTP header injection:

1. Default payloads (127.0.0.1, localhost, etc.) are injected into the headers mentioned above
2. Custom payloads can be supplied (e.g. you've enumerated some internal IPs or domains) using the `pfile` parameter

```
$ headi
  -pfile string
    	payload file
  -t int
    	timeout (milliseconds) (default 10000)
  -url string
    	target URL
```
Currently only takes one URL as input but you can easily bash script for numerous URLs like so:
```
$ for i in $(cat urls); do headi -url $i;done
```
An example is provided below from the HTB machine Control:
<a href="https://asciinema.org/a/380645" target="_blank"><img src="https://asciinema.org/a/380645.svg" /></a>
