# Build the view
FROM node:alpine as view-build
RUN mkdir /view
COPY /view /view
WORKDIR /view
RUN yarn
RUN yarn build

FROM arm64v8/golang:latest as server-build
RUN mkdir /furnace-monitoring-system
ADD . /furnace-monitoring-system
COPY --from=view-build /view/dist /furnace-monitoring-system/view/dist
WORKDIR /furnace-monitoring-system
RUN go mod download
RUN GOOS=linux GOARCH=arm64 go build -o main

FROM arm64v8/alpine:latest as publish
RUN mkdir /fms
COPY --from=server-build /furnace-monitoring-system/main /fms/main
WORKDIR /fms
EXPOSE 5000
CMD ["/fms/main"]