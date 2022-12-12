FROM golang:1.20rc1-bullseye
WORKDIR /src
COPY ./ ./
ENTRYPOINT bash