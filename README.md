# pic-resizer

This package implements microservice with REST interface creating images thumbnails

## Installation

```bash
go get github.com/faroyam/pic-resizer
go mod vendor
docker-compose run app go test ./... -cover
docker-compose up
```
## Usage

There are three ways to use it:
 - paste a URL to the server with an image into a GET request e.g. http://localhost:8080/v1/get?url=https://example.com/image.jpg
 - send a POST request to http://localhost:8080/v1/base with content-type application/json {"data":"image"}, 
 where image is a base64 string 
 - send a POST request to http://localhost:8080/v1/multipart with content-type multipart/form-data  
If successful, the response will contain links to the original and resized images

## Examples

Send a get request via curl:
```bash 
curl http://localhost:8080/v1/get?url=https://thiscatdoesnotexist.com/
```

Send a base64 request via curl: 
```bash
(echo -n '{"data": "'; base64 images/test.jpg; echo '"}') | curl \
 -H 'Content-Type: application/json' -d @- http://localhost:8080/v1/base
```

Send a multipart request via curl:
```bash
curl -F "data=@images/test.jpg" http://localhost:8080/v1/multipart
```
