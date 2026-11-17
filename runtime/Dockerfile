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
    && npm cache clean --force
RUN ln -s /opt/pi/node_modules/.bin/pi /usr/local/bin/pi
COPY config/ /opt/config/
RUN cd /opt/config/extensions/hindsight \
    && npm ci --omit=dev --ignore-scripts --no-fund --no-audit --progress=false \
    && npm cache clean --force
COPY --chmod=0755 init /usr/local/bin/init

RUN groupadd --gid 1000 agent \
    && useradd --uid 1000 --gid agent --create-home --shell /bin/bash agent \
    && mkdir -p /sandbox \
    && chown agent:agent /sandbox

ENV LANG=C.UTF-8 \
    HOME=/home/agent \
    PI_CODING_AGENT_DIR=/sandbox/.pi/agent \
    HINDSIGHT_CONFIG=/sandbox/.pi/agent/hindsight.json \
    HINDSIGHT_LOG_FILE=/sandbox/.pi/agent/hindsight.log \
    HINDSIGHT_DIAG_FILE=/sandbox/.pi/agent/hindsight-diag.jsonl \
    PI_SKIP_VERSION_CHECK=1

USER agent
WORKDIR /sandbox
ENTRYPOINT ["/usr/local/bin/init"]
CMD ["pi"]
