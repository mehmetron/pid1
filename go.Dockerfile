
# docker build -t go-build . -f go.Dockerfile
# docker run --name go-build -p 3000:3000 go-build
# docker exec -it go-build sh
FROM golang:1.16

RUN apt-get -y update && \
    apt-get -y install nano && \
    apt-get -y install curl

WORKDIR /app/playground


COPY main /app/

#CMD ["sleep", "10000"]
CMD [ "../main" ]