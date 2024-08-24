### Custom Docker image
* To reduce dev latency and install things
```bash
docker push anakin007/consul-xdp-ubuntu:1.0
```

### Running the Cluster
```bash
docker compose build
docker compose up
```
* Docker Network : `service_mesh`