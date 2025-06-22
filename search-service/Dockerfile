# Etapa 1: compilación
FROM rust:1.78 as builder

WORKDIR /app
COPY . .
RUN cargo build --release

# Etapa 2: ejecución
FROM debian:bullseye-slim
WORKDIR /app
COPY --from=builder /app/target/release/search-service /usr/local/bin/search-service

CMD ["search-service"]
