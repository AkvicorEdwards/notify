
# Server

```shell
# METHOD POST
https://ip:WebPort/msg
# key
title
app
msg
```

```shell
cd server && go build

# Defult: https

if http
// [:WebPort] [:TcpPort] [UserID]
./server :8080 :8081 user
else if https
// [:WebPort] [:TcpPort] [UserID] [pem file] [key file]
./notify :8080 :8081 user ./cert/cert.pem ./cert/key.pem
```

# Client

```shell
cd client && go build

// [ip] [TcpPort] [UserID]
./notify 127.0.01 8081 user
```
