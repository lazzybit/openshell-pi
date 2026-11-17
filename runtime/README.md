# OpenShell Pi runtime

Pi 1.5.2 workload image for OpenShell, with Node 22, tmux, and image-owned Pi/Hindsight configuration. All runtime defaults remain under `/opt/config`; persistent Pi state and sessions live under `/sandbox`.

Build with `docker build -t openshell-pi-runtime:local .` from this directory. The root npm manifests are the pinned Pi installer artifacts; Hindsight owns a separate npm package and lockfile under `config/extensions/hindsight/`.

The init script explicitly links each default into the persistent Pi config without replacing existing paths or symlinks. OpenShell replaces the container entrypoint, so sandbox commands must invoke `/usr/local/bin/init` explicitly, for example `/usr/local/bin/init pi`.
