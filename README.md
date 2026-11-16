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
No provider credentials or OpenShell policy are included.

The native Hindsight Pi extension `@vectorize-io/hindsight-coding-agents` 0.8.7
is installed at build time with its published archive SHA-256 verified. The
launcher loads the installed Pi entry point; no runtime npm downloads or MCP
server are needed. Configuration is initialized at `$PI_CODING_AGENT_DIR/hindsight.json`
(or set `HINDSIGHT_CONFIG` to an existing file). Session retention is enabled;
auto-update, git seeding, codebase surveys and bank-config management are disabled.
The internal API is `http://hindsight-api.hindsight.svc.cluster.local:8888`.
Use an explicit `bankId` when a shared durable bank is desired.

`models.json` is baked into `/opt/pi`. The `pi` launcher installs that image-owned
configuration into `$PI_CODING_AGENT_DIR/models.json` on each invocation, leaving
sessions and auth untouched. The default agent directory is `/sandbox/.pi/agent`.
This also works when OpenShell launches `bash -i` instead of the image CMD.

Only Agentgateway models `openai/gpt-6.1-sol` and `openai/gpt-6-astra` are declared.
Metadata comes from the OpenAI entries in https://models.dev/api.json (2026-11-16):
reasoning, text/image input, 128,000 output, and base USD/million token prices
including cache reads/writes. Both models deliberately use a 272,000 context
window (the models.dev long-context pricing threshold), rather than their
1,050,000 advertised maximum, to keep Pi compaction below the expensive tier.
Pi's custom-model schema cannot express PDF input or context-tier prices;
displayed costs use the base tier, not a billing guarantee. No upstream key is included.

Build with docker build -t openshell-pi:1.5.2 .
