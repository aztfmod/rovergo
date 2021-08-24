FROM golang:1.17.0-buster

ARG TERRAFORM_VERSION=1.0.0
ARG TFLINT_VERSION=0.29.1
ARG INSTALL_GO=false
ARG INSTALL_DOCKER=true
ARG DOCKER_VERSION=20.10.7
ARG ROVER_VERSION=latest
ARG BIN_INSTALL_PATH=/usr/bin

ARG USERNAME=rover
ARG USER_UID=1000
ARG USER_GID=$USER_UID

# We have to run apt-get update twice thanks to azure-cli >_<
RUN apt-get update -qq && apt-get install -y -qq curl gnupg2 lsb-release \
    && curl -sSL https://packages.microsoft.com/keys/microsoft.asc | gpg --dearmor > /etc/apt/trusted.gpg.d/microsoft.gpg  \
    && echo "deb [arch=amd64] https://packages.microsoft.com/repos/azure-cli/ $(lsb_release -cs) main" > /etc/apt/sources.list.d/azure-cli.list

# Install Azure CLI and other base packages
RUN apt-get update -qq && export DEBIAN_FRONTEND=noninteractive \
    && apt-get install -y -qq azure-cli zsh git make sudo jq curl wget net-tools unzip zip nano lsb-release apt-transport-https ca-certificates \ 
    && apt-get clean -y && rm -rf /var/lib/apt/lists/*

# Install other tools directly as binaries, using tools-install repo scripts
WORKDIR /tmp
RUN git clone https://github.com/benc-uk/tools-install.git

# Manditory installs
RUN ./tools-install/terraform.sh $TERRAFORM_VERSION $BIN_INSTALL_PATH \
 && ./tools-install/tflint.sh $TFLINT_VERSION $BIN_INSTALL_PATH

# Optionally install Go
ENV PATH /usr/local/go/bin:$PATH
RUN if [ "$INSTALL_GO" = "true" ]; then ./tools-install/golang.sh; fi

# Other optional installs, like docker client
RUN if [ "$INSTALL_DOCKER" = "true" ]; then ./tools-install/docker-client.sh $DOCKER_VERSION $BIN_INSTALL_PATH; fi

# Set ROVER_VERSION to latest or a version tag, when blank Rover will not be installed
RUN if [ ! -z "$ROVER_VERSION" ]; then curl https://raw.githubusercontent.com/aztfmod/rovergo/main/install.sh | bash -s -- -b $BIN_INSTALL_PATH $ROVER_VERSION; fi

# Add non root user and switch to it
RUN groupadd --gid $USER_GID $USERNAME && useradd -s /bin/bash --gid $USERNAME -m $USERNAME
ENV HOME /home/$USERNAME
ENV PATH $HOME/.local/bin:$PATH
USER $USERNAME

WORKDIR /home/$USERNAME
