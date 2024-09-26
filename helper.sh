docker rmi $(docker images --filter "dangling=true" -q --no-trunc)
docker compose build
docker compose up