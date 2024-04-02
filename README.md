# Helix-Relayer-Runner

Helix-Relayer-Runner is a project designed to run the [Helix Bridge Relayer](https://github.com/helix-bridge/relayer) in a Docker container. The password is set via HTTP POST, and the Helix-Relayer-Runner will start the relayer and automatically input the password. The password is only stored in memory.

In addition, Helix-Relayer-Runner is responsible for monitoring the relayer.

## Running the Helix-Relayer-Runner

To run the Helix-Relayer-Runner, use the following Docker command:

```bash
docker run -dt --restart=allways --name=helix-relayer -e LP_BRIDGE_STORE_PATH=/data/db -e HELIX_RELAYER_DIR=/opt/relayer -e LP_BRIDGE_PATH=/data/configure.json -v ~/relayer/.maintain:/data quay.io/perror/helix-relayer:2de2d12
```

After the Docker container is running, you can enter the container with the following command:

```bash
docker exec -it helix-relayer /bin/sh
```

To set the password, use the following curl command:

```bash
curl -X POST "http://127.0.0.1:8080/pass" -H "content-type: application/json" -d '{"p": "<your-password>"}'
```

Replace `<your-password>` with your actual password.
