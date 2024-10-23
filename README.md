# ebpf Service mesh

## Testing on Docker environment

* Install Docker
* Clone repository
* Change the line number 7 in `./user_proxy/client-config.json` to `  "retry_join": ["server"],`
* Build 
```
docker compose build
```
* Run
```
docker compose up
```
* Debug (Gives you the shell of any container)
```
docker exec -it <container-hash> /bin/bash
```

* Debug logs can be seen in `/tmp/gologs_service<number>.logs`, `/tmp/consul-client.log` and `/tmp/consul-server.log` files.

## Run on Different physical machines [Yet to test]

### Server Agent
```bash
cd control_plane
chmod +x server_baremetal.sh
sudo ./server_baremetal.sh
```

### Client Agent
```bash
cd user_proxy
chmod +x client_baremetal.sh
export SERVICE_NAME=<NAME> # Unique name to this service / node
export SERVER_IP=<IP> # IP of consul server
sudo ./client_baremetal.sh
```

## Architecture
