# logxy
reverse proxy that logs requests and responses.

## Getting Started

### Build the build environment

Build the docker container for building logxy

```
make build
```

### Build logxy

```
make dep-ensure
make all
```

### Run logxy

```
make proxy NETWORK=myNetworkName
tail -f logxy.log
```
