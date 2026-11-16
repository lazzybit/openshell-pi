FROM debian:trixie-slim

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        bash ca-certificates curl fd-find git gnupg ripgrep tmux \
    && install -d -m 0755 /etc/apt/keyrings \
    && curl -fsSL https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key \
        | gpg --dearmor -o /etc/apt/keyrings/nodesource.gpg \
    && chmod 0644 /etc/apt/keyrings/nodesource.gpg \
    && echo "deb [signed-by=/etc/apt/keyrings/nodesource.gpg] https://deb.nodesource.com/node_22.x nodistro main" \
        > /etc/apt/sources.list.d/nodesource.list \
    && apt-get update \
    && apt-get install -y --no-install-recommends nodejs \
    && ln -s /usr/bin/fdfind /usr/local/bin/fd \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /opt/pi
COPY package.json package-lock.json ./
RUN npm ci --ignore-scripts --omit=dev --include=optional --no-fund --no-audit --progress=false \
    && npm cache clean --force \
    && ln -s /opt/pi/node_modules/.bin/pi /usr/local/bin/pi

RUN groupadd --gid 1000 agent \
    && useradd --uid 1000 --gid agent --create-home --shell /bin/bash agent \
    && mkdir -p /sandbox \
    && chown agent:agent /sandbox

ENV HOME=/home/agent \
    PI_CODING_AGENT_DIR=/tmp/pi \
    PI_SKIP_VERSION_CHECK=1

USER agent
WORKDIR /sandbox
CMD ["pi"]
