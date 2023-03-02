# Embedded view build step
FROM node:alpine as view-build
RUN mkdir /view
COPY /view /view
WORKDIR /view
RUN yarn
RUN yarn build

# Server build step
FROM golang:latest as server-build
RUN mkdir /furnace-monitoring-system
ADD . /furnace-monitoring-system
COPY --from=view-build /view/dist /furnace-monitoring-system/view/dist
WORKDIR /furnace-monitoring-system
RUN go mod download
RUN GOOS=linux go build -o main

# Final publish
FROM ubuntu:latest as publish
RUN mkdir /fms
COPY --from=server-build /furnace-monitoring-system/main /fms/main
COPY config/config.json /fms/config/config.json
WORKDIR /fms
RUN mkdir db
RUN mkdir log
EXPOSE 5000
ENTRYPOINT ["/fms/main"]