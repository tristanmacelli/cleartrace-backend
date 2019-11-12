# add the necessary instructions
# to create a Docker container image
# for your Go API gateway server


FROM alpine:3.7
RUN apk add --no-cache ca-certificates
COPY ./gateway /build
EXPOSE 443
ENTRYPOINT ["/build"]
# ENTRYPOINT ["api.sauravkharb.me"]