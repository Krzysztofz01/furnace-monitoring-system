FROM node:18.12.1-alpine AS build
WORKDIR /app
COPY package.json
COPY yarn.lock
RUN yarn
COPY . .
RUN yarn run build:server
RUN yarn run build:views

FROM node:18.12.1-alpine AS publish
WORKDIR /app
COPY package.json
COPY yarn.lock
RUN yarn install --production
COPY --from=builder ./app/dist ./dist
COPY --from=builder ./app/src/views ./dist/views
COPY --from=builder ./app/src/static/tailwind.css ./dist/static/tailwind.css
EXPOSE 80
CMD [ "yarn", "run", "start" ]