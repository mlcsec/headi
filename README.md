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

## Output
Output is structured in the following format:
```
[URL] [Header: payload] [Status code] [Content length]
```
An example is provided below:
```
```
