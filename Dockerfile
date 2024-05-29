# 使用 alpine 作为基础镜像设置时区
FROM alpine as timezone
ENV TZ=Asia/Shanghai
RUN apk add --no-cache tzdata \
    && cp /usr/share/zoneinfo/$TZ /etc/localtime \
    && echo $TZ > /etc/timezone

# 使用 node 进行前端构建
FROM node:16 as builder
WORKDIR /build
COPY ./web .
COPY ./VERSION .
RUN npm install
RUN REACT_APP_VERSION=$(cat VERSION) npm run build

# 使用 golang 构建后端服务
FROM golang:1.21.10 as builder2
ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64
WORKDIR /build
COPY . .
COPY --from=builder /build/build ./web/build
RUN go mod download
RUN go build -ldflags "-s -w -X 'wechat-server/common.Version=$(cat VERSION)' -extldflags '-static'" -o wechat-server

# 使用设置了时区的 alpine 作为最终基础镜像
FROM timezone as final
COPY --from=builder2 /build/wechat-server /
ENV PORT=3000
EXPOSE 3000
WORKDIR /data
ENTRYPOINT ["/wechat-server"]
