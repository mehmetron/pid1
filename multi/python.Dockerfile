FROM ubuntu:20.04

RUN apt-get update && apt-get install -y curl

# Install Python3.9
RUN apt-get -y install software-properties-common
RUN add-apt-repository ppa:deadsnakes/ppa
RUN apt-get -y install python3.9


CMD ["sleep", "1000"]
