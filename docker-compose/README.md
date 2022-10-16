# Using `docker-compose` (Recommended)

1. Install Docker

	- [Docker Install documentation](https://docs.docker.com/install/)
	- [Docker-Compose Install documentation](https://docs.docker.com/compose/install/)

2. Create a new folder on your local machine and copy docker-compose\docker-compose.yml

3. Review the docker-compose.yml contents

4. Bring up your stack by running

	```bash
	docker-compose up -d
	```

5. Add your wallet
    ```bash
	docker-compose exec node sh /keys-add.sh
	```

6. Copy your terra wallet address and go to the terra faucet here -> http://45.79.139.229:3000/ Put your address in and give yourself luna coins.

7. Start the validator
	```bash
	docker-compose exec node sh /create-validator.sh
	```

# Cheat Sheet:

## Start

```bash
docker-compose up -d
```

# Stop

```bash
docker-compose down
```

# View Logs

```bash
docker-compose logs -f
```

# Run Terrad Commands Example

```bash
docker-compose exec node terrad status
```

# Upgrade

```bash
docker-compose down
docker-compose pull
docker-compose up -d
```

# Build from source
Its possible to use docker-compose to build the images from the go source code by running the following commands in sequence:

1) docker-compose -f docker-compose.yml -f docker-compose.build.yml build core --no-cache
2) docker-compose -f docker-compose.yml -f docker-compose.build.yml build node --no-cache
