# pid1
### First process in container sandbox






- GOOS=linux GOARCH=amd64 go build main.go


```
ws = new WebSocket("ws://localhost:8080/ws");

ws.onopen = function() {
    console.log("ws is open")
};

ws.onmessage = function(e) {
    console.log("recieved message: ", e.data)
}

ws.send('tom')

ws.onclose = function () {
    console.log("ws is closed")
}
```

### Check what ports in use
see which ports open with ```netstat -peanut```
```
lsof -i -P -n | grep LISTEN | grep -v '8090'
```