# Go parent image only used for building
FROM golang:alpine AS builder

# diasble cross-compile and install git
RUN CGO_ENABLED=0
RUN apk update && apk add --no-cache git

# create work dir and copy all files to there
WORKDIR $GOPATH/src/github.com/boseabhishek/go-shopping
COPY . .

# get go deps
RUN go get -d -v

# expose on 8080
EXPOSE 8080

# run main.go file
CMD ["go", "run", "main.go"]