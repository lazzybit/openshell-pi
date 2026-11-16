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
RUN curl -fsSL https://registry.npmjs.org/@vectorize-io/hindsight-coding-agents/-/hindsight-coding-agents-0.8.7.tgz -o /tmp/hindsight.tgz \
    && printf '%s  %s\n' f9ad8b8c7df8c363cf05dbbfd75089424ccd54aca0cfee27ae298e28a217567d /tmp/hindsight.tgz | sha256sum -c - \
    && mkdir -p /opt/config/extensions/hindsight \
    && tar -xzf /tmp/hindsight.tgz --strip-components=1 -C /opt/config/extensions/hindsight \
    && rm /tmp/hindsight.tgz \
    && ln -s /opt/pi/node_modules/.bin/pi /usr/local/bin/pi
COPY config/ /opt/config/
COPY --chmod=0755 init /usr/local/bin/init

RUN groupadd --gid 1000 agent \
    && useradd --uid 1000 --gid agent --create-home --shell /bin/bash agent \
    && mkdir -p /sandbox \
    && chown agent:agent /sandbox

ENV HOME=/home/agent \
    PI_CODING_AGENT_DIR=/sandbox/.pi/agent \
    HINDSIGHT_CONFIG=/sandbox/.pi/agent/hindsight.json \
    HINDSIGHT_LOG_FILE=/sandbox/.pi/agent/hindsight.log \
    HINDSIGHT_DIAG_FILE=/sandbox/.pi/agent/hindsight-diag.jsonl \
    PI_SKIP_VERSION_CHECK=1

USER agent
WORKDIR /sandbox
ENTRYPOINT ["/usr/local/bin/init"]
CMD ["pi"]
