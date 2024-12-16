docker run -d \
  --name babeldb \
  -e MARIADB_ALLOW_EMPTY_ROOT_PASSWORD=1 \
  -e MARIADB_ROOT_HOST=localhost \
  -e MARIADB_USER=myuser \
  -e MARIADB_PASSWORD=mypassword \
  -e MARIADB_DATABASE=babel \
  -p 3306:3306 \
  -v data:/var/lib/mariadb/data \
  mariadb:10.6
