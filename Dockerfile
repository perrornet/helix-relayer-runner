FROM node:16-alpine as builder
RUN mkdir -p /opt/build
WORKDIR /opt/build
COPY ./relayer ./
RUN yarn install && yarn build

# build golang
FROM golang:1.20-alpine as go-builder
WORKDIR /build
COPY ./go.mod .
RUN go mod download
COPY . .
RUN go build -o ./runner .

FROM node:16-alpine
RUN apk update && apk add expect curl
RUN mkdir -p /opt/data
COPY --from=builder /opt/build/dist /opt/relayer/dist
COPY --from=go-builder /build/runner /opt/relayer/runner
WORKDIR /opt/relayer
COPY ./relayer/.env.docker .env
COPY ./relayer/package.json package.json
RUN yarn install --production
CMD [ "/opt/relayer/runner" ]
