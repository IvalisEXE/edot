#

run:
docker-compose up --build

stop:
docker-compose down

destroy:
docker-compose down --rmi all -v
docker image prune
