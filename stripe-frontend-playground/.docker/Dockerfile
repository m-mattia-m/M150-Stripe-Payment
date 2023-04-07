FROM docker.io/nginx:1.23.3-alpine

WORKDIR /app/etc

ENV NGINX_ENVSUBST_TEMPLATE_DIR="/app/etc"
ENV NGINX_ENVSUBST_OUTPUT_DIR="/app/etc"

ENV CSP_DOMAINS="http://localhost"

COPY .docker/conf/mime.types ./
COPY .docker/conf/nginx.conf.template ./

WORKDIR /app/www

COPY dist/ .

EXPOSE 8080

CMD ["nginx", "-c", "/app/etc/nginx.conf"]

