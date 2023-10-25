FROM golang:1.21
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
EXPOSE 8080
ENV service1=service1.com
ENV service2=service2.com
CMD ["go","run","."]