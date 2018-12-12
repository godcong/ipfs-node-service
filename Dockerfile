FROM golang:1.11.2

WORKDIR /home/ipfs_node

COPY . .

RUN go build -v -o ipfs_node

EXPOSE 8080

CMD ["/home/ipfs_node/ipfs_node"]