# add the necessary instructions
# to create a Docker container image
# for your summary client using the
# NGINX web server image as the base

FROM nginx
ADD default.conf /etc/nginx/conf.d/default.conf
COPY . /build
COPY . /usr/share/nginx/html
EXPOSE 443
ENTRYPOINT ["/build"]