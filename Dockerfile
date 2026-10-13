FROM cgr.dev/chainguard/go:latest AS builder

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
FROM cgr.dev/chainguard/wolfi-base:latest

COPY --from=builder /workdir/app .
#COPY --from=builder /workdir/public /public
COPY --from=builder /workdir/entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]

CMD []
