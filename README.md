# OpenShell Pi workload

Debian trixie-slim, Node 22 and Pi 1.5.2, running as agent (1000).
The root package.json and package-lock.json are unchanged official installer
artifacts from https://pi.dev/api/installer/releases/1.5.2/ . Builds use npm ci.
Update both manifests together to upgrade Pi; no update script or version ARG.

## Layout

- /opt/pi: Pi and npm dependencies.
- /opt/config: image defaults from config/.
- /opt/config/extensions/hindsight: extension-owned package.json and lockfile;
  npm ci installs Hindsight 0.8.7 into its node_modules at build time.
  index.ts loads the native Pi adapter; no installer or manual tar extraction.
- init links the Hindsight companion skill directly from that npm package.
- /sandbox/.pi/agent: persistent Pi config, sessions and logs.

The init entrypoint explicitly lists four links: the Hindsight extension,
its companion skill, models.json and hindsight.json. There are no loops or
wildcard discovery. All /opt paths are read-only at runtime; installation happens
only during image build. Existing files, directories and
symlinks (including dangling symlinks) are never replaced. It finally execs the
requested command. /usr/local/bin/pi is the native Pi executable symlink, not a
wrapper. No AGENTS.md is included.

Docker uses ENTRYPOINT init and CMD pi. OpenShell replaces the image entrypoint
with its sandbox runtime: explicitly launch `/usr/local/bin/init pi` or
`/usr/local/bin/init bash -i` as the sandbox command. Plain `pi` works after that
initialization, including later execs and inside tmux. Hindsight config is also
set through image ENV so subsequent execs use the same persistent file.

Only Agentgateway models openai/gpt-6.1-sol and openai/gpt-6-astra are declared.
Metadata is from https://models.dev/api.json : reasoning, text/image input,
128000 output tokens, base USD/million input/output/cache prices. Both contexts
are deliberately capped at 272000 to avoid long-context pricing. Pi's schema
cannot represent PDF inputs or tiered pricing; base cost estimates are not a
billing guarantee. No upstream credentials are included.

Hindsight uses the internal hindsight-api Service, automatic session retention,
and no MCP. Auto-update, Git ingestion/seeding, surveys and bank configuration
management are disabled. Set HINDSIGHT_CONFIG or customize the persistent config
with an explicit bankId when a fixed bank is desired. OpenShell policy must
allow /opt/config read-only and Node access to the model and memory endpoints.

Build: docker build -t openshell-pi:1.5.2 .
