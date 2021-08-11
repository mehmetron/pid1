

go:
	GOOS=linux GOARCH=amd64 go build -o main *.go
	docker rm go-build
	docker build -t go-build . -f go.Dockerfile
	docker run --name go-build -p 8080:8080 -p 8090:8090 go-build

node:
	GOOS=linux GOARCH=amd64 go build -o main *.go
	docker rm node-build
	docker build -t node-build . -f node.Dockerfile
	docker run --name node-build -p 8080:8080 -p 8090:8090 node-build

python:
	GOOS=linux GOARCH=amd64 go build -o main *.go
	docker rm python-build
	docker build -t python-build . -f python.Dockerfile
	docker run --name python-build -p 8080:8080 -p 8090:8090 python-build




.PHONY: go node python