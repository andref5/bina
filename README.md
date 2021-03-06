# eBPF xdp example
This is a simple example of eBPF with XDP, droping packages by IP address.

- Take a look step-by-step examples:
  - \o_ [start](https://github.com/andref5/bina/tree/main) - Load eBPF Program with iproute2
  - [go-load](https://github.com/andref5/bina/tree/go-load)  - Load eBPF Program with golang
  - [go-load-map](https://github.com/andref5/bina/tree/go-load-map) - Load eBPF Program and Map with golang

## How it works

It has two HTTP golang server on folders "pkg/a"(port 5011) and "pkg/b"(port 5012).
Using docker/compose (https://docs.docker.com/compose/install/) to simulate a connection/network between two services with fixed IP 172.20.0.11(service A) and 172.20.0.12(service B).

```
*-----------------------Docker network--------------------------*
|                                                               |
|  +---------+     http://172.20.0.12:5012/b     |e|+--------+  |
|  | service |--------------------------------->>|b| service |  |
|  |    A    |<<---------------------------------|p|    B    |  |
|  +---------+     http://172.20.0.11:5011/a     |f|+--------+  |
|                                                 ^             |
|                                                 | (xdp.o)     |
|                                                 |             |
|                                           *----------*        |
|                                           | iproute2 |        |
|                                           *----------*        |
*---------------------------------------------------------------*
```

In the folder "pkg/b/ebpf" have a eBPF program (xdp.c) that load net packet contents, parse IP address and drop package if IP is equal to 172.20.0.11 (service A)

## Testing

### Startup docker containers
```bash
docker-compose build
docker-compose up
```

### New terminal to exec interactive bash on svc-A container
```bash
docker exec -it svc-a /bin/bash
```
- Inside svc-a container test HTTP server B
```bash
curl -m5 http://172.20.0.12:5012/b
# B OK
```

### New terminal to exec interactive bash on svc-B container
```bash
docker exec -it svc-b /bin/bash
```
- Inside svc-b container let's take a look at eBPF 
```bash
cd ebpf
# compile xdp.c to eBPF
clang -target bpf -O2 -c xdp.c -o xdp.o

# Load eBPF inside the kernel using iproute2 tool
ip -force link set dev eth0 xdpdrv obj xdp.o sec .text

# See attached XDP hook on eth0
ip link show dev eth0
# prog/xdp id........
```

### Back to interactive bash on svc-A container
```bash
curl -m5 http://172.20.0.12:5012/b
# curl: (28) Connection timed out after 5001 milliseconds
```

### New terminal to test from your host (docker-compose configured port mapping localhost 5012 -> svc-b 5012)
```bash
curl -m5 http://localhost:5012/b
# B OK
```

### Back to interactive bash on svc-B container
```bash
# Unload eBPF using ip link tool
ip link set dev eth0 xdpdrv off

# See detached XDP hook on eth0
ip link show dev eth0
```

### Back to interactive bash on svc-A container
```bash
curl -m5 http://172.20.0.12:5012/b
# B OK
```

## References

- https://github.com/xdp-project/xdp-tutorial
- https://github.com/iovisor/gobpf
- https://developers.redhat.com/blog/2018/12/06/achieving-high-performance-low-latency-networking-with-xdp-part-1
- https://developers.redhat.com/blog/2018/12/17/using-xdp-maps-rhel8
- https://www.youtube.com/watch?v=XmFBjr2ujSI
- https://www.youtube.com/watch?v=7pmXdG8-7WU
- https://homepages.dcc.ufmg.br/~mmvieira/cc/papers/Processamento_Rapido_de_Pacotes_com_eBPF_e_XDP-%20versao%20final.pdf


## Etc

The term BINA is an acronym for "B Identifies Number of A", used in telco was created by Brazilian inventors