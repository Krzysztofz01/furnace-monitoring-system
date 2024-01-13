# Embedded view build step
FROM node:alpine as view-build
RUN npm install -g @go-task/cli
RUN mkdir /furnace-monitoring-system
COPY . /furnace-monitoring-system
WORKDIR /furnace-monitoring-system
RUN task build:frontend

# Server build step
FROM golang:latest as server-build
RUN go install github.com/go-task/task/v3/cmd/task@latest
RUN mkdir /furnace-monitoring-system
ADD . /furnace-monitoring-system
COPY --from=view-build /furnace-monitoring-system/view/dist /furnace-monitoring-system/view/dist
WORKDIR /furnace-monitoring-system
RUN task build:backend

# Final publish
FROM ubuntu:latest as publish
RUN mkdir /fms
COPY --from=server-build /furnace-monitoring-system/bin /fms
WORKDIR /fms
EXPOSE 5000
ENTRYPOINT ["/fms/furnace-monitoring-system-server"]