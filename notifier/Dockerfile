FROM node:18-alpine as base

RUN npm i -g pnpm

FROM base as dependencies

WORKDIR /app
COPY package.json pnpm-lock.yaml ./
RUN pnpm install

FROM base as build

WORKDIR /app
COPY . .
COPY --from=dependencies /app/node_modules ./node_modules
RUN pnpm run build
RUN pnpm prune --prod

FROM base as deploy

WORKDIR /app
COPY --from=build /app/build/ ./build/
COPY --from=build /app/node_modules ./node_modules

RUN echo $(ls -1 ./)
RUN echo $(ls -1 ./build)
RUN ls ./build
RUN echo "FOO"

CMD ["node", "build/index.js"]
# ENTRYPOINT ["tail", "-f", "/dev/null"]