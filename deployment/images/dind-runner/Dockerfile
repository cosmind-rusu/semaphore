FROM dind-ansible:latest

RUN apk add --no-cache wget git rsync nodejs npm yarn unzip tar


ARG TERRAFORM_VERSION="1.8.2"
ARG PULUMI_VERSION="3.116.1"
ARG OPENTOFU_VERSION="1.7.0"
ARG TARGETARCH="amd64"


RUN adduser -D -u 1001 -G root -G docker semaphore && \
    mkdir -p /tmp/semaphore && \
    mkdir -p /etc/semaphore && \
    mkdir -p /var/lib/semaphore && \
    chown -R semaphore:0 /tmp/semaphore && \
    chown -R semaphore:0 /etc/semaphore && \
    chown -R semaphore:0 /var/lib/semaphore

RUN wget https://raw.githubusercontent.com/semaphoreui/semaphore/develop/deployment/docker/runner/runner-wrapper -P /usr/local/bin/ && chmod +x /usr/local/bin/runner-wrapper
COPY semaphore /usr/local/bin/
COPY terraform /usr/local/bin/
COPY ssh_config /etc/ssh/ssh_config



RUN wget https://github.com/opentofu/opentofu/releases/download/v${OPENTOFU_VERSION}/tofu_${OPENTOFU_VERSION}_linux_${TARGETARCH}.tar.gz
RUN tar xf tofu_${OPENTOFU_VERSION}_linux_${TARGETARCH}.tar.gz -C /usr/local/bin
RUN rm tofu_${OPENTOFU_VERSION}_linux_${TARGETARCH}.tar.gz

RUN if [ "$TARGETARCH" = "amd64" ]; then \
      export PULUMI_ARCH="x64"; \
    else \
      export PULUMI_ARCH="${TARGETARCH}"; \
    fi && \
    wget -O pulumi.tar.gz https://github.com/pulumi/pulumi/releases/download/v${PULUMI_VERSION}/pulumi-v${PULUMI_VERSION}-linux-${PULUMI_ARCH}.tar.gz
RUN tar xf pulumi.tar.gz --strip-components=1 -C /usr/local/bin
RUN rm pulumi.tar.gz



RUN chown -R semaphore:0 /usr/local/bin/runner-wrapper &&\
    chown -R semaphore:0 /usr/local/bin/semaphore &&\
    chmod +x /usr/local/bin/runner-wrapper &&\
    chmod +x /usr/local/bin/semaphore

RUN pip3 install boto3 botocore requests==2.31.0

WORKDIR /home/semaphore
USER 1001

RUN mkdir ./venv

#RUN python3 -m venv ./venv --system-site-packages && \
#    source ./venv/bin/activate && \
#    pip3 install --upgrade pip boto3 botocore requests

# RUN pip3 install boto3 botocore requests

RUN echo '{"tmp_path": "/tmp/semaphore","dialect": "bolt", "runner": {"config_file": "/var/lib/semaphore/runner.json"}}' > /etc/semaphore/config.json

CMD [ "/usr/local/bin/runner-wrapper" ]