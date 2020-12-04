FROM --platform=$BUILDPLATFORM golang:1.15-alpine AS build

ENV APP_PATH /app

WORKDIR $APP_PATH

COPY go.mod .

RUN go env -w GO111MODULE=auto

# RUN go env -w GOPROXY=https://mirrors.aliyun.com/goproxy/
# RUN go env -w GOPROXY=https://goproxy.cn,direct

# RUN go env -w GOPRIVATE=*.github.com

RUN go mod download -x

COPY . .

RUN go build -o main .

## ------------------------
FROM --platform=$BUILDPLATFORM alpine

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories

RUN apk add --no-cache tzdata\
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata \
    && rm -rf /var/cache/apk/*

RUN set -ex \
    && apk add --no-cache ca-certificates

ENV GRPC_GO_LOG_SEVERITY_LEVEL="WARNING" \
    GRPC_GO_LOG_VERBOSITY_LEVEL="WARNING"

WORKDIR /app

COPY --from=build /app/main .

CMD ["/app/main"]