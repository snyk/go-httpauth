# Use an official golang runtime as a parent image
FROM golang:1.16 as builder
ENV GO111MODULE=on
ADD . /src
WORKDIR /src/cmd
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o .
FROM gcr.io/snyk-main/ubuntu:20
RUN mkdir -p /<home holder>
WORKDIR /<home holder>
COPY --from=builder /src/cmd/<binary name> .
EXPOSE 8080
CMD ["./<binary name>"]