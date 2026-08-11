# 构建阶段
FROM node:20-alpine AS builder

WORKDIR /app

RUN apk add --no-cache python3 make g++

COPY package.json yarn.lock ./

RUN yarn install --frozen-lockfile --non-interactive

COPY . .

RUN yarn build

FROM nginx:1.29.4-alpine

ARG GIT_COMMIT
ARG BUILD_DATE

LABEL org.opencontainers.image.revision=${GIT_COMMIT}
LABEL org.opencontainers.image.created=${BUILD_DATE}
LABEL com.efucloud.build.commit=${GIT_COMMIT}
LABEL com.efucloud.build.date=${BUILD_DATE}

COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
