# add the necessary instructions
# to create a Docker container image
# for your Go summary server

FROM alpine:3.12
RUN apk add --no-cache ca-certificates
COPY summary /summary
EXPOSE 5050-5060
ENTRYPOINT ["/summary"]