FROM golang:1.17.5 AS builder

RUN mkdir /tmp/tom-frontend

COPY . /tmp/tom-frontend/

WORKDIR /tmp/tom-frontend

RUN CGO_ENABLED=0 GOOS=linux go build -o tom-frontend main.go

FROM centos:8

ARG APP_VERSION=v1.0
ENV APP_VERSION=${APP_VERSION}

USER root

# Copy app from builder image
COPY --from=builder /tmp/tom-frontend/tom-frontend /usr/local/bin/

RUN chmod +x /usr/local/bin/tom-frontend

RUN yum update -y && \
    yum install -y curl && \
    yum clean all
    
CMD ["/usr/local/bin/tom-frontend"]
