
local:
	GOOS=linux GOARCH=amd64 go build -o main *.go
	docker rm bobcon
	docker build -t bob .
	docker run -p 8090:8090 -p 8080:8080 --name bobcon bob

exec:
	docker exec -it bobcon /bin/bash

dev:
	GOOS=linux GOARCH=amd64 go build main.go
	docker build -t bob .


penis:
	echo "buthole is on fiiiiiire"
	echo "and I can't keep it in\n"
	echo " my heart is burninninn"

.PHONY: dev local penis exec