# Dockerfile stolen from here https://stackoverflow.com/questions/53433486/working-on-user-in-dockerfile-and-installing-packages-on-it-permission-denied
FROM ubuntu:20.04
#FROM ubuntu:18.04

# https://rtfm.co.ua/en/docker-configure-tzdata-and-timezone-during-build/
ENV TZ=Europe/Kiev
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

#RUN apt-get -y update && \
#    apt-get -y install golang && \
#    apt-get -y install nano && \
#    apt-get -y install lsof && \
#    apt-get -y install curl && \
#    apt-get -y install net-tools

COPY script.sh /
RUN chmod +x /script.sh
RUN /script.sh


#ENV user lg
#
#RUN useradd -m -d /home/${user} ${user} && \
#    chown -R ${user} /home/${user}
## && \
## adduser ${user} sudo && \
## echo '%sudo ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers
#
#USER ${user}
#
#WORKDIR /home/${user}
#
#COPY main langs.json /home/${user}/



WORKDIR /home
COPY main langs.json /home/

#RUN mkdir gostuff
#ENV GOPATH=$PWD/home/gostuff

CMD [ "./main" ]


# RUN go build && \
#     rm main.go

# CMD [ "./lg" ]

# RUN sudo apt-get -y install curl && \
#     sudo apt-get -y install lsb-core && \
#     sudo apt-get -y install lsb && \
#     sudo apt-get -y upgrade -f