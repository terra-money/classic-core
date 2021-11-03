# Terra Core Dockerized

## Common usage examples 

### Mainnet 

Standard options with LCD enabled: 

```
docker run -it -p 1317:1317 -p 26657:26657 -p 26656:26656 terra-money/core-node:v0.5.9-oracle
```

LCD disabled: 

```
docker run -e ENABLE_LCD=false -it -p 1317:1317 -p 26657:26657 -p 26656:26656 terra-money/core-node:v0.5.9-oracle
```

Custom gas fees: 

```
docker run -e MINIMUM_GAS_PRICES="0.01133uluna,0.15uusd,0.104938usdr,169.77ukrw,428.571umnt,0.125ueur,0.98ucny,16.37ujpy,0.11ugbp,10.88uinr,0.19ucad,0.14uchf,0.19uaud,0.2usgd,4.62uthb,1.25usek,1.25unok,0.9udkk,2180.0uidr,7.6uphp,1.17uhkd" -it -p 1317:1317 -p 26657:26657 -p 26656:26656 terra-money/core-node:v0.5.9-oracle
```

Custom snapshot name (for when new snapshots come out):

```
docker run -e SNAPSHOT_NAME="columbus-5-pruned.20211022.0410.tar.lz4" -it -p 1317:1317 -p 26657:26657 -p 26656:26656 terra-money/core-node:v0.5.9-oracle
```

Custom snapshot URL

```
docker run -e SNAPSHOT_BASE_URL="https://get.quicksync.io" -it -p 1317:1317 -p 26657:26657 -p 26656:26656 terra-money/core-node:v0.5.9-oracle
```

**Note:** We recommend copying a snapshot to S3 or another file store and using the above options to point the container to your snapshot. The default snapshot name included will be obsolete and removed in a matter of days.

## Building the Docker images

```
./build_all.sh v0.5.9-oracle
```