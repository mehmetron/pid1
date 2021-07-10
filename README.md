# pid1
### First process in container sandbox




```
GOOS=linux GOARCH=amd64 go build main.go
docker build -t bob .
docker run -p 8090:8090 -p 8080:8080 bob
```


