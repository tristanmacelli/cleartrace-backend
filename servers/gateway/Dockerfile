# add the necessary instructions
# to create a Docker container image
# for your Go API gateway server

FROM alpine:3.7
RUN apk add --no-cache ca-certificates
COPY . ./build
EXPOSE 80

# Not sure if you need the expose.
