### Custom Docker image
* To reduce dev latency and install things. Used as base image by all containers
```bash
docker pull anakin007/service-mesh:latest
```

### Running the Cluster
```bash
docker compose build
docker compose up
```
* Docker Network : `service_mesh`