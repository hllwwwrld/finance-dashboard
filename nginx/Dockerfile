FROM nginx:alpine

COPY snippets/app-locations.conf /etc/nginx/snippets/app-locations.conf
COPY default.conf /etc/nginx/conf.d/default.conf
# Пароли: файл nginx/.htpasswd на хосте (см. `make nginx-htpasswd`), монтируется в docker-compose.
# Для HTTPS: docker-compose.ssl.yml монтирует default-ssl.conf поверх default.conf и каталог ssl/.
