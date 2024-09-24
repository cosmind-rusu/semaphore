FROM dind-runner:latest

# Set environment variables for NVM
#ENV NVM_DIR=/home/semaphore/.nvm
#ENV NODE_VERSION=18.17.0

# Install NVM
#RUN curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.1/install.sh | bash
USER root

RUN apk add --no-cache nodejs npm

USER 1001

# Install Node.js using NVM
#RUN . "$NVM_DIR/nvm.sh" && \
#    nvm install $NODE_VERSION && \
#    nvm use $NODE_VERSION && \
#    nvm alias default $NODE_VERSION

# Add NVM and Node.js to PATH
#ENV PATH=$NVM_DIR/versions/node/v$NODE_VERSION/bin:$PATH