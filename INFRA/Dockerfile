FROM hashicorp/terraform:latest

# install azure cli
RUN \
  apk update && \
  apk add bash py-pip && \
  apk add --virtual=build gcc libffi-dev musl-dev openssl-dev python-dev make && \
  pip --no-cache-dir install -U pip && \
  pip --no-cache-dir install azure-cli && \
  apk del --purge build

# install kubectl
ADD https://storage.googleapis.com/kubernetes-release/release/v1.6.4/bin/linux/amd64/kubectl /usr/local/bin/kubectl
ENV HOME=/config
RUN set -x && \
    apk add --no-cache curl ca-certificates && \
    chmod +x /usr/local/bin/kubectl

# install jenkins-x
ENV PATH=$PATH:~/.jx/bin
RUN mkdir -p ~/.jx/bin && \
    curl -L https://github.com/jenkins-x/jx/releases/download/v1.3.688/jx-linux-amd64.tar.gz | tar xzv -C ~/.jx/bin && \
    echo 'export PATH=$PATH:~/.jx/bin' >> ~/.bashrc

ENTRYPOINT ["/bin/bash"]