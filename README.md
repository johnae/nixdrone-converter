A conversion plugin for [Drone CI](https://drone.io) providing support for generating pipelines using the [Nix](https://nixos.org/nix) language. _Please note this project requires Drone server version 1.4 or higher._

## Installation

Create a shared secret:

```console
$ openssl rand -hex 16
dfe5d2712b087e77027e08e323d6fa63
```

Download and run the plugin:

```console
$ docker run -d \
  --publish=3000:3000 \
  --env=DRONE_DEBUG=true \
  --env=DRONE_SECRET=dfe5d2712b087e77027e08e323d6fa63 \
  --restart=always \
  --name=converter johnae/nixdrone-converter
```

Update your Drone server configuration to include the plugin address and the shared secret.

```text
DRONE_CONVERT_PLUGIN_ENDPOINT=http://1.2.3.4:3000
DRONE_CONVERT_PLUGIN_SECRET=dfe5d2712b087e77027e08e323d6fa63
```