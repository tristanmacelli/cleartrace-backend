# add the necessary instructions
# to create a Docker container image
# for your Go summary server

FROM alpine:3.7
RUN apk add --no-cache ca-certificates
COPY ./summary /summarybuild
EXPOSE 5050-5060
ENTRYPOINT ["/summarybuild"]