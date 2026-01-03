# my app

## Service Overview

[*insert network map here*]

### admin

[**Image:**](https://hub.docker.com/_/adminer) `adminer:5.4.1`

#### tasks

- postgres `database` service administration

### cache

[**Image:**](https://hub.docker.com/_/redis) `redis:8.4.0-bookworm`

#### tasks

- store [session information](https://docs.gofiber.io/api/middleware/session) from `server` service, like flash messages and user login state
- [messaging (pub/sub)](https://redis.io/docs/latest/develop/pubsub/) between `server` and `function` services

### database

[**Image:**](https://hub.docker.com/_/postgres) `postgres:18.1-bookworm`

#### tasks

- staff bios, events, general content

### function

**Image:** (_to be determined_)

#### tasks

- formats images to proper size, conversion to WebP
- writes images to shared named volume

### image

[**Image:**](https://hub.docker.com/_/nginx) `nginx:1.27-alpine`

#### tasks

- static server of images stored on a shared named volume

### server

[**Image:**](https://hub.docker.com/_/golang) `golang:1.25.5`

#### tasks

- [fiber](https://gofiber.io/) server-side rendered (SSR)
- [fiber templating](https://docs.gofiber.io/guide/templates) engine
- [htmx](https://htmx.org/) partials

### web

[**Image:**](https://hub.docker.com/_/nginx) `nginx:1.27-alpine`

#### tasks

- reverse proxy

## Usage

Production: `docker compose up --build --detach`  
Development: `docker compose -f compose.dev.yml up --build`
