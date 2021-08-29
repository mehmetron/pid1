FROM ubuntu:20.04

COPY --from=python-build . .
COPY --from=go-build . .
COPY --from=node-build . .

# Go env variables
ENV GOROOT=/usr/local/go
ENV GOPATH=$HOME/Projects/Proj1
ENV PATH=$GOPATH/bin:$GOROOT/bin:$PATH

CMD ["sleep", "1000"]