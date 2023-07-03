FROM rust:1.70-buster AS builder

WORKDIR /app

RUN rustup target add x86_64-unknown-linux-musl

COPY Cargo.lock Cargo.toml ./
COPY .cargo/config.toml ./.cargo/

COPY config config
COPY src src

RUN --mount=type=secret,id=netrcConf,required=true,target=/root/.netrc cargo build --release --target x86_64-unknown-linux-musl

RUN cp target/x86_64-unknown-linux-musl/release/rust-gateway .

FROM scratch AS master

WORKDIR /app

COPY --from=builder /app/rust-gateway app

ENTRYPOINT [ "./app" ]