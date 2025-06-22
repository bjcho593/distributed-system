# Etapa 1: Compilación
FROM rust:slim as builder

WORKDIR /app
COPY . .
RUN apt-get update && apt-get install -y pkg-config libssl-dev build-essential
RUN cargo build --release

# Etapa 2: Imagen ligera final
FROM debian:bullseye-slim
WORKDIR /app
COPY --from=builder /app/target/release/search-service /usr/local/bin/search-service

CMD ["search-service"]

