FROM node:12.19.0-alpine as base

WORKDIR /app

COPY package.json ./

RUN npm i

COPY . .

FROM base as production

RUN npm run dev