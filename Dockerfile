FROM node:18 AS front-builder

COPY ./wwwroot/happy-day-app /app

WORKDIR /app
RUN yarn install
RUN npm i -g @angular/cli@17
RUN ng build --build-optimizer --aot

FROM golang:1.21 AS backend-builder

COPY . /app
WORKDIR /app
RUN go build -ldflags "-s -w"

FROM alpine

COPY --from=front-builder /app/dist/happy-day-app /app/wwwroot
COPY --from=backend-builder /app/config.yml /app
COPY --from=backend-builder /app/happyday /app

ENTRYPOINT [ "/app/happy-day" ]