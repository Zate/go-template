FROM golang:alpine AS builder

RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group

RUN apk add --no-cache ca-certificates git

# Setting workdir
WORKDIR /go/src/app/

# Copy all go files to workdir
COPY src/*.go /go/src/app/

COPY src/templates /go/src/app/templates/
COPY src/templates/*.tmpl /go/src/app/templates/

COPY src/go.mod src/go.sum /go/src/app/
# Download necessary Go modules
RUN go mod download && \
    CGO_ENABLED=0 go build -a -ldflags '-w -extldflags "-static"' -o /go/bin/app

# # Use the Golang golden image for runtime
FROM scratch as final

WORKDIR /

ARG SERVICE_NAME

ENV SERVICE_NAME=${SERVICE_NAME}

COPY --from=builder /user/group /user/passwd /etc/

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8080

# # Copy compiled app into the image workdir
COPY --chown=appuser:appuser --from=build /go/bin/app /app

# Perform any further action as an unprivileged user.
USER nobody:nobody

# Run the compiled binary.
ENTRYPOINT ["/app"]