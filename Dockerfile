# Dockerfile stolen from here https://stackoverflow.com/questions/53433486/working-on-user-in-dockerfile-and-installing-packages-on-it-permission-denied
FROM ubuntu:16.04

RUN apt-get -y update && \
    apt-get -y install golang && \
    apt-get -y install nano && \
    apt-get -y install lsof

ENV user lg

RUN useradd -m -d /home/${user} ${user} && \
    chown -R ${user} /home/${user} 
# && \
# adduser ${user} sudo && \
# echo '%sudo ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers

USER ${user}

WORKDIR /home/${user}

COPY main bob.go /home/${user}/

# RUN rm main.go

CMD [ "./main" ]
#CMD [ "sleep", "1000" ]

# RUN go build && \
#     rm main.go

# CMD [ "./lg" ]

# RUN sudo apt-get -y install curl && \
#     sudo apt-get -y install lsb-core && \
#     sudo apt-get -y install lsb && \
#     sudo apt-get -y upgrade -f