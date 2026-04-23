FROM 1ms.run/library/golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
# 设置Go模块代理
RUN go env -w GOPROXY=https://goproxy.cn,direct && go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o laima ./cmd/laima

FROM 1ms.run/library/alpine:latest

WORKDIR /app

COPY --from=builder /app/laima .
COPY scripts/init_db.sql ./scripts/

RUN mkdir -p /app/repos

EXPOSE 8080

CMD ["./laima"]