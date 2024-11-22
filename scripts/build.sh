cd /src && \
apt update && \
apt install -y gcc-arm-linux-gnueabi && \
GOARM=6 GOOS=linux GOARCH=arm CGO_ENABLED=1 CC=arm-linux-gnueabi-gcc go build  -buildvcs=false ./cmd/stockybot

