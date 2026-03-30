# Dynamic builds
ARG BUILDER_IMAGE=dhi.io/golang:1.26-debian13-dev
ARG FINAL_IMAGE=dhi.io/golang:1.26-debian13

# Build stage
FROM --platform=${BUILDPLATFORM} ${BUILDER_IMAGE} AS builder

# Build args
ARG GIT_REVISION=""
ARG BUILD_DATE=""

# Platform args
ARG TARGETOS
ARG TARGETARCH
ARG TARGETPLATFORM

# Use modules for dependencies
WORKDIR $GOPATH/src/go.rtnl.ai/uptime

COPY go.mod .
COPY go.sum .

# Set the build environment
ENV CGO_ENABLED=0
ENV GO111MODULE=on
ENV GOOS=${TARGETOS}
ENV GOARCH=${TARGETARCH}

# Install dependencies
RUN go mod download && go mod verify

# Copy the source code
COPY pkg/ pkg/
COPY cmd/ cmd/

# Build the uptime binary
RUN go build -o /go/bin/uptime \
  -ldflags="-X 'go.rtnl.ai/uptime/pkg.GitVersion=${GIT_REVISION}' -X 'go.rtnl.ai/uptime/pkg.BuildDate=${BUILD_DATE}'" \
  ./cmd/uptime

# Final stage
FROM ${FINAL_IMAGE} AS final

LABEL maintainer="Rotational Labs, Inc. <support@rotational.io>"
LABEL description="Service status monitor for Rotational applications and systems."

# Copy the uptime binary
COPY --from=builder /go/bin/uptime /usr/local/bin/uptime

CMD [ "/usr/local/bin/uptime", "serve" ]
