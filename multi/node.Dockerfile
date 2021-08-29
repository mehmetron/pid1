FROM ubuntu:20.04

RUN apt-get update && apt-get install -y curl sudo

# Install nodejs
RUN curl -fsSL https://deb.nodesource.com/setup_16.x | sudo -E bash -

RUN apt-get -y update && \
    apt-get -y install nodejs

CMD ["sleep", "1000"]
