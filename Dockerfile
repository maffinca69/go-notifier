FROM --platform=$BUILDPLATFORM golang:alpine as golang

RUN apk add -U tzdata
RUN apk --update add ca-certificates

WORKDIR /app
COPY . .

ARG TARGETOS

RUN go mod download
RUN go mod verify

RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /notifier .

FROM scratch

COPY --from=golang /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=golang /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=golang /etc/passwd /etc/passwd
COPY --from=golang /etc/group /etc/group

COPY --from=golang /notifier .

CMD ["/notifier"]