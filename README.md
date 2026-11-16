# OpenShell Pi workload

Minimal Pi image based on debian:trixie-slim, using Debian nodejs and npm.
Runs as agent (UID/GID 1000), with /sandbox as the working directory.

The root package.json and package-lock.json are downloaded unchanged from the
Pi 1.5.2 official managed installer artifacts:

- https://pi.dev/api/installer/releases/1.5.2/package.json
- https://pi.dev/api/installer/releases/1.5.2/package-lock.json

Docker builds copy these manifests and use npm ci, without resolving latest.
Upgrade Pi by replacing both files with the chosen release's official artifacts.
The manifests determine the Pi version; there is no separate PI_VERSION build arg.
System packages follow Debian trixie repositories.

PI_SKIP_VERSION_CHECK=1 disables automatic version checks and update prompts.
No patches, plugins, provider credentials or OpenShell policy are included.

Build with docker build -t openshell-pi:1.5.2 .
