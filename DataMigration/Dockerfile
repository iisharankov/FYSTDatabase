# #This is a sample Image 
# # FROM ubuntu
# FROM golang:1.14-alpine AS build 

# # RUN apt-get update 
# RUN mkdir -p /home/server
# COPY . /home/server

# RUN cd /home/server && go install 

# CMD ["server"]



FROM golang:1 AS build 
WORKDIR /go/src/github.com/iisharankov/FYSTDatabase/DataMigration

COPY . .
RUN go get -d -v
RUN go install -a -v -tags netgo -ldflags=-extldflags=-static

FROM scratch
COPY --from=0 /go/bin/DataMigration /
CMD ["/DataMigration"]