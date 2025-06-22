export BABBLE_DB_HOST=localhost
export BABBLE_DB_USER=myuser
export BABBLE_DB_PASSWORD=mypassword
export BABBLE_DB_DBNAME=babble
export BABBLE_DB_PORT=3306

docker run -d \
  --name babbledb \
  -e MARIADB_ALLOW_EMPTY_ROOT_PASSWORD=1 \
  -e MARIADB_ROOT_HOST=localhost \
  -e MARIADB_USER=myuser \
  -e MARIADB_PASSWORD=mypassword \
  -e MARIADB_DATABASE=babble \
  -p 3306:3306 \
  -v data:/var/lib/mariadb/data \
  mariadb:10.6
