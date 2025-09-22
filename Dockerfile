FROM cgr.dev/chainguard/go AS builder

WORKDIR /workdir

COPY . .
RUN chmod ugo+x entrypoint.sh
RUN VERSION=$(cat APP_VERSION) && \
    LDFLAGS=$(echo "-X 'main.AppVersion=${VERSION}' -s -w") && \
    echo ${LDFLAGS} && \ 
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="${LDFLAGS}" -o app

# --- Runtime ---
#FROM cgr.dev/chainguard/static
FROM cgr.dev/chainguard/wolfi-base

COPY --from=builder /workdir/app .

COPY --from=builder /workdir/entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]

CMD []
