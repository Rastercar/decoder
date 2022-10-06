# ------------------- BUILDER ------------------- 
FROM golang:1.18.1-alpine as builder

# add git (to be able to run git commands) and build-base (to run makefile commands) 
RUN apk add git build-base

WORKDIR /usr/src/reciever_ms

# this is needed to get git branch metadata
ADD .git ./.git

# add the makefile to run the make install command
ADD Makefile .

# add all files and folders necessary to build the go binary
ADD go.mod .
ADD go.sum .
ADD cmd ./cmd
ADD config ./config
ADD protocol ./protocol
ADD queue ./queue
ADD server ./server
ADD tracer ./tracer
ADD utils ./utils

# get the go dependencies
RUN make get

# build the go binary
RUN make build

# ------------------- RELEASE ------------------- 
FROM alpine:3.12

LABEL Author="Vitor Andrade Guidorizzi <vitor.guidorizzi@hotmail.com>"

COPY --from=builder /usr/src/reciever_ms/ /reciever_ms

# go to the builder folder
# add git and make
# run the 'make install' command that copies the built go binary to /usr/local/bin/ and the needed config files
# remove the build folder as the binary is built
# remove git and make
RUN cd /reciever_ms \
    && apk add --no-cache make git \
    && make install \
    && rm -rf /reciever_ms \
    && apk del make git

COPY ./docker/docker-entrypoint.sh /usr/local/bin/

ENTRYPOINT ["docker-entrypoint.sh"]

CMD ["reciever_ms"]

EXPOSE 8000
