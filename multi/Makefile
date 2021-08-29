python:
	docker rm python-build
	docker build -t python-build . -f python.Dockerfile
	docker run -d --name python-build python-build

go:
	docker rm go-build
	docker build -t go-build . -f go.Dockerfile
	docker run -d --name go-build go-build

node:
	docker rm node-build
	docker build -t node-build . -f node.Dockerfile
	docker run -d --name node-build node-build

main:
	docker rm main
	docker build -t main .
	docker run -d --name main main


.PHONY: python go node main