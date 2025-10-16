FROM golang:1.25-bookworm AS gcscp

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY gcscp.go ./
RUN CGO_ENABLED=0 go build -o /usr/local/bin/gcscp gcscp.go



FROM node:24-bookworm-slim

RUN apt-get update \
    && apt-get install -y \
      procps git curl jq \
    && rm -rf /var/lib/apt/lists/*

RUN npm install -g @google/gemini-cli

COPY --from=gcscp /usr/local/bin/gcscp /usr/local/bin/gcscp
RUN chmod 777 /usr/local/bin/gcscp

COPY run.sh /usr/local/bin/run.sh
RUN chmod 777 /usr/local/bin/run.sh

ENTRYPOINT ["/usr/local/bin/run.sh"]
