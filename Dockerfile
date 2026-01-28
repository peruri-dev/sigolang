# Use an official Node.js runtime as a parent image
FROM node:24-alpine AS base

WORKDIR /workdir

FROM base AS install

# install node modules for build
RUN mkdir -p /temp/dev
COPY package*.json /temp/dev/
RUN cd /temp/dev && npm install

FROM base AS build

COPY --from=install /temp/dev/node_modules node_modules

# Copy the rest of the application code
COPY . .

# Copy env that needed for build
COPY .env .env

# Build the frontend to dist
RUN npm run build

FROM cgr.dev/chainguard/go AS builder

WORKDIR /workdir

COPY . .
COPY --chown=nonroot:nonroot --from=build /workdir/frontend/dist /workdir/frontend/dist

RUN chmod ugo+x entrypoint.sh
RUN VERSION=$(cat APP_VERSION) && \
    LDFLAGS=$(echo "-X 'main.AppVersion=${VERSION}' -s -w") && \
    echo ${LDFLAGS} && \ 
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags dist -ldflags="${LDFLAGS}" -o app

# --- Runtime ---
#FROM cgr.dev/chainguard/static
FROM cgr.dev/chainguard/wolfi-base

COPY --from=builder /workdir/app .
#COPY --from=builder /workdir/public /public
COPY --from=builder /workdir/entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]

CMD []
