FROM golang:1.20-buster

RUN go version
ENV GOPATH=/

COPY ./ ./

# install psql
RUN apt-get update
RUN apt-get -y install postgresql-client

# make wait-for-postgres.sh executable
RUN chmod +x wait-for-postgres.sh
# Убедитесь, что путь правильный


# build go app
RUN go mod download
RUN go build -o notesync ./cmd/main.go

CMD ["./notesync"]

