FROM node:18 AS front-builder

COPY ./wwwroot/mec-app /app

WORKDIR /app
RUN yarn install
RUN npm i -g @angular/cli@17
RUN ng build --build-optimizer --aot --localize

FROM golang:1.21 AS backend-builder

COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux  go build -ldflags "-s -w"

FROM alpine

COPY --from=front-builder /app/dist/mec-app /app/wwwroot
COPY --from=backend-builder /app/config.yml /app
COPY --from=backend-builder /app/mec /app

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/engine/reference/builder/#expose
EXPOSE 8080

WORKDIR /app

CMD [ "mec" ]