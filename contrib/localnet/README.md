This folder is used for simulating localnet

# I. History simulation
Run these two commands for simulation

```sh
make localnet-start
bash contrib/localnet/simulation/entrypoint.sh
```

If you are aiming for another network, please refer to [sample network](simulation/network). All network should follow these same structures
* chain-id as folder network name with file keys.json
* define chain rpc endpoints in [network.js](simulation/network/network.json)

Run the simulation with custom chain id
```sh
CHAIN_ID=bajor-1 bash contrib/localnet/simulation/entrypoint.sh
```

1. Keys
The simulation will add these normal [addresses](simulation/network/localterra/keys.json) to keyring. Query `terrad keys list --keyring-backend test` to see.

For validator mnemonic, please look at build/node0/terrad/key_seed.json. This file will only be created after setting up localnet.

2. Processes
* Validator addresses will send uluna to normal addresses
* Normal addresses will delegate to validators
* Normal address test0 will create a validator
* Normal address test0, test1 will create two NFT smart contracts. They will mint and send NFT to other normal addresses
* Normal addresses will create text proposal and vote on it