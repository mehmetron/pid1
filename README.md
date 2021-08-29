# pid1
### First process in container sandbox
#### Browser editor as single source of truth


### TODO:
#### Editor with uneditable files (send files as response after generating project)
#### Modal with iframe
#### Inject js to control history of iframe for vitejs
https://gist.github.com/mehmetron/7b422029e46ce9d4df7032d2768c8b8b


```
GOOS=linux GOARCH=amd64 go build main.go
docker build -t bob .
docker run -p 8090:8090 -p 8080:8080 bob
```


## Post Mortem

#### Vitejs hmr ws port is same as vitejs port so when client tries to connect it fails
Solution:
https://github.com/vitejs/vite/issues/1653#issuecomment-816188959
https://github.com/vitejs/vite/pull/1926#issuecomment-774728965
Change vitejs hmr ws port to reverse proxy port and it works because the default proxied route
is the port the user's program is using


#### Vitejs if you want to use a different port than 3000 
Solution: Put this in package.json script tag
```
"dev":"vite --port 8050"
```


#### Go Websocket wasn't connecting/upgrading
Solution: It wasn't upgrading because the request origin wasn't allowed.

I set it to allow all origins with a snippet from this thread
https://github.com/gorilla/websocket/issues/367