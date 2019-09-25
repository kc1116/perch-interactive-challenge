# build stage
FROM golang:1.13-alpine as builder

ENV ROOT_PATH=/working

RUN mkdir $ROOT_PATH
WORKDIR $ROOT_PATH

COPY . .
RUN apk add --update make
RUN make deploy
RUN mkdir /ui && mv ./build/ui/* /ui

FROM node:11-alpine
COPY --from=builder /usr/bin/perch-iot-pubsub /usr/bin/perch-iot-pubsub
COPY --from=builder /ui /ui

RUN chown -R root:node /usr/local/lib/node_modules/ && \
    chmod -R 775 /usr/local/lib/ && \
    npm install -g serve --unsafe-perm=true --allow-root && \
    mkdir -p /tmp/google-cert && \
    wget -O /tmp/google-cert/roots.pem https://pki.google.com/roots.pem



