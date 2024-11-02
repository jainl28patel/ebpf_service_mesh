# eBPF Service mesh

## Concepts
1. [Consul Watches](https://developer.hashicorp.com/consul/docs/dynamic-app-config/watches) : Push based mechanism to get updates when some component of the control plane like distribute kv-store, services, nodes etc changes. It takes some arguments as input which will be executed when change is detected.
2. For loading ebpf program I am using cilium's ebpf package for more help and reference visit [this site](https://ebpf-go.dev/) and [api-reference](https://pkg.go.dev/github.com/cilium/ebpf).
3. We are using hashicorp consul to implemnt the given control plane. If you want to leverage more functionality of consul it could be done just by editing the config files for consul server and client agent present in `./control_plane` and `./user_proxy` folder.

## How to Customize as per DataPlane
1. **`./user_proxy/app` directory** : Add any user level application or service to this folder and provide a Makefile in `app` folder.
2. **`./user_proxy/kernel` directory** : Add you `eBPF` data-plane inside this folder and provide a Makefile in the same folder.
3. **`./user_proxy/user` directory** : This contains the minimal go user agent in order to load and interact with the ebpf program and the control plane.

* Changes required in the User agent are
#### Required
1. implement your ebpf program loading logic in `ebpf_loader.go` file in `ebpf_loader()` function that returns the array of reference to required ebpf maps.
2. implement your ebpf map initialization logic in `initialize_map()` function in same file. It takes the above returned array of reference to map as input.

#### Optional [Examples are provided]
1. Add optional go routines in the `main.go` file or under `package main` for any task like KV store update etc. (Eg : `./user_proxy/user/kv_update.go`)
2. Add more optional endpoints to the lightweight GIN server in the `server.go` if required. (Eg : `.user_proxy/user/server.go`)
3. If you want to add more consul watches add them to `client-config.json`. Here is the [reference](https://developer.hashicorp.com/consul/docs/dynamic-app-config/watches) for the same. (Eg : `./user_proxy/client-config.json`)
4. Add custom scripts that will be executed in case of event triggered by consul watches in `watches_scripts` folder. (Eg : Look at the folder `.user_proxy/watches_scripts/metric_update.py`)

## Testing on Docker environment

* Install Docker
* Clone repository
* Update the absolute paths in `user_proxy/watches_scripts/*` as per required and update the path of arguments of `consul watches` in `client-config.json`.
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

## Run on Different physical machines

* Update the absolute paths in `user_proxy/watches_scripts/*` as per required and update the path of arguments of `consul watches` in `client-config.json`.

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
