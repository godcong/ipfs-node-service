FROM golang:1.11.2

WORKDIR /home

COPY . .

RUN go install -o ipfs_node

CMD ["ipfs_node"]