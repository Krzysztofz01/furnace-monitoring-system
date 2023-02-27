# Embedded view build step
FROM node:alpine as view-build
RUN mkdir /view
COPY /view /view
WORKDIR /view
RUN yarn
RUN yarn build

# Server build step (arm64v8/golang:latest)
FROM golang:latest as server-build
RUN mkdir /furnace-monitoring-system
ADD . /furnace-monitoring-system
COPY --from=view-build /view/dist /furnace-monitoring-system/view/dist
WORKDIR /furnace-monitoring-system
RUN go mod download
RUN GOOS=linux GOARCH=arm64 go build -o main

# Final publish (arm64v8/ubuntu:latest)
FROM ubuntu:latest as publish
RUN mkdir /fms
COPY --from=server-build /furnace-monitoring-system/main /fms/main
COPY config/config.json /fms/config/config.json
WORKDIR /fms
RUN mkdir db
RUN mkdir log
EXPOSE 5000
CMD ["/fms/main"]