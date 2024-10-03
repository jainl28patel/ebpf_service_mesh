## Run on Different machines

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