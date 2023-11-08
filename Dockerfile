FROM debian:11.7

# Set destination for COPY
WORKDIR /app

ARG BIN=/target/nit
ARG SCHEMA=/schema/docs.yaml
COPY ${BIN} ./
COPY ${SCHEMA} ./schema/docs.yaml

RUN chmod +x nit

EXPOSE 8080

CMD ["/app/nit"]
