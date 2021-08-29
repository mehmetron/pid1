FROM ubuntu:20.04

RUN apt-get -y update && \
    apt-get -y install wget


# Install Go
RUN wget https://dl.google.com/go/go1.16.4.linux-amd64.tar.gz
RUN tar -xvf go1.16.4.linux-amd64.tar.gz
RUN mv go /usr/local

# Go env variables
ENV GOROOT=/usr/local/go
ENV GOPATH=$HOME/Projects/Proj1
ENV PATH=$GOPATH/bin:$GOROOT/bin:$PATH

CMD ["sleep", "1000"]

