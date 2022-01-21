FROM golang:1.17.5 AS builder

RUN mkdir /tmp/tom-frontend

COPY . /tmp/tom-frontend/

WORKDIR /tmp/tom-frontend

RUN CGO_ENABLED=0 GOOS=linux go build -o tom-frontend main.go

FROM centos:8

ARG APP_VERSION=v1.0
ENV APP_VERSION=${APP_VERSION}

# Note: These arguments will be provided by CI pipeline with the right values
ARG COMMIT_ID
ARG SHORT_COMMIT_ID

LABEL git.commit-id=${COMMIT_ID}
LABEL git.short-commit-id=${SHORT_COMMIT_ID}

USER root

# Copy app from builder image
COPY --from=builder /tmp/tom-frontend/tom-frontend /usr/local/bin/

RUN chmod +x /usr/local/bin/tom-frontend

RUN yum update -y && \
    yum install -y curl && \
    yum clean all
    
CMD ["/usr/local/bin/tom-frontend"]
