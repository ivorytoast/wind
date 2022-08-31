# Importing the golang debian stable image from docker hub
FROM golang:1.19.0-bullseye

# Installing the necessay debian libraries for 0MQ in the docker container
RUN apt-get update -q && \
    apt-get install -y --no-install-recommends autoconf automake && \
    apt-get install -y --no-install-recommends libtool pkg-config build-essential && \
    apt-get install -y --no-install-recommends libzmq5-dev

# setting the working directory
WORKDIR /usr/src/app

# Downloading and extracting 0MQ
RUN wget http://download.zeromq.org/zeromq-4.1.4.tar.gz && \
    tar -xvf zeromq-4.1.4.tar.gz

# installing 0MQ
WORKDIR zeromq-4.1.4
RUN ./autogen.sh && \
    ./configure
RUN make install && \
    ldconfig

# Copying source code into the required directory in the docker container
WORKDIR ../
COPY . .

# Creating go.mod files and building the go executable for the server
WORKDIR server
RUN go mod init github.com/pranjallk1995/0mq-golang-docker/server
RUN go mod tidy
RUN go get github.com/pebbe/zmq4
RUN go build -o /0MQ-golang-server

# Creating go.mod files and building the go executable for the client
WORKDIR ../client
RUN go mod init github.com/pranjallk1995/0mq-golang-docker/client
RUN go mod tidy
RUN go get github.com/pebbe/zmq4
RUN go build -o /0MQ-golang-client