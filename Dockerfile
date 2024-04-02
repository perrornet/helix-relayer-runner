FROM node:16-alpine as builder
RUN mkdir -p /opt/build
WORKDIR /opt/build
COPY . ./
RUN yarn install && yarn build

FROM node:16-alpine
RUN apk update && apk add expect curl
RUN mkdir -p /opt/data
COPY --from=builder /opt/build/dist /opt/relayer/dist
COPY helix-relayer-runner /opt/relayer/helix-relayer-runner
WORKDIR /opt/relayer
COPY .env.docker .env
COPY package.json package.json
RUN yarn install --production
CMD [ "/opt/relayer/helix-relayer-runner" ]
