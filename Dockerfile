FROM google/golang:1.3

#install docker client for communication with the host
RUN curl -sL https://get.docker.io/ | sh

#the clinet expects a socket at /var/run/docker.sock so run this container with
# -v $HOST_PATH_TO_SOCKET/docker.sock:/var/run/docker.sock
# eg: -v /var/run/docker.sock:/var/run/docker.sock

#we usegodep
RUN mkdir -p /gopath/src/github.com/adminibar/container-updater

WORKDIR /gopath/src/github.com/adminibar/container-updater
ADD . /gopath/src/github.com/adminibar/container-updater
RUN go get github.com/adminibar/container-updater

# EXPOSE 8000
CMD []
ENTRYPOINT ["/gopath/bin/container-updater"]