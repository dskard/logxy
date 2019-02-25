FROM ubuntu:bionic-20180526

ENV LC_ALL=C.UTF-8
ENV LANG=C.UTF-8

# Install basic development & debugging tools
RUN set -ex; \
        \
        apt-get update; \
        apt-get install -y --no-install-recommends \
            bash \
            ca-certificates \
            curl \
            git \
            jq \
            make \
            tar \
            vim \
            wget; \
        \
        rm -rf /var/lib/apt/lists/*;

# Install Go.
ENV GOLANG_VERSION 1.11.5
ENV GOLANG_DOWNLOAD_SHA256 ff54aafedff961eb94792487e827515da683d61a5f9482f668008832631e5d25
RUN set -ex; \
    curl -fsSL "https://golang.org/dl/go${GOLANG_VERSION}.linux-amd64.tar.gz" -o golang.tar.gz; \
    echo "$GOLANG_DOWNLOAD_SHA256  golang.tar.gz" | sha256sum -c - ; \
    tar -C /usr/local -xzf golang.tar.gz; \
    rm golang.tar.gz
ENV PATH="$PATH:/usr/local/go/bin"

# Install Go tools: delve, dep
RUN GOPATH=/usr/local/gotools go get -v github.com/go-delve/delve/cmd/dlv
ENV GODEP_VERSION 0.5.0
ENV GODEP_DOWNLOAD_SHA256 287b08291e14f1fae8ba44374b26a2b12eb941af3497ed0ca649253e21ba2f83
RUN set -ex; \
    curl -fsSL "https://github.com/golang/dep/releases/download/v${GODEP_VERSION}/dep-linux-amd64" -o /usr/local/gotools/bin/dep; \
    echo "${GODEP_DOWNLOAD_SHA256}  /usr/local/gotools/bin/dep" | sha256sum -c - ; \
    chmod a+x /usr/local/gotools/bin/dep
ENV PATH="$PATH:/usr/local/gotools/bin"

# command line setup
# do minimal setup so we can be semi-efficient when using
# the command line of the container. Without PS1, we will
# get a prompt like "I have no name!@<container_id_hash>:/$"
# since we don't create a user or group.
RUN set -ex; \
    echo "PS1='\h:\w\$ '" >> /etc/bash.bashrc; \
    echo "alias ls='ls --color=auto'" >> /etc/bash.bashrc; \
    echo "alias grep='grep --color=auto'" >> /etc/bash.bashrc;
