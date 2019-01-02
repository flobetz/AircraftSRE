FROM jenkins/jenkins:lts
USER root

# install docker
RUN apt-get update && \
    apt-get -y install apt-transport-https \
         ca-certificates \
         curl \
         gnupg2 \
         software-properties-common && \
    curl -fsSL https://download.docker.com/linux/$(. /etc/os-release; echo "$ID")/gpg > /tmp/dkey; apt-key add /tmp/dkey && \
    add-apt-repository \
       "deb [arch=amd64] https://download.docker.com/linux/$(. /etc/os-release; echo "$ID") \
       $(lsb_release -cs) \
       stable" && \
    apt-get update && \
    apt-get -y install docker-ce && \
    usermod -a -G root jenkins && \
    usermod -a -G docker jenkins

# install azure cli
RUN apt-get install apt-transport-https lsb-release software-properties-common -y && \
    AZ_REPO=$(lsb_release -cs) && \
    echo "deb [arch=amd64] https://packages.microsoft.com/repos/azure-cli/ $AZ_REPO main" | \
        sudo tee /etc/apt/sources.list.d/azure-cli.list && \
    apt-get update && \
    apt-get --allow-unauthenticated install azure-cli

# install kubectl
ADD https://storage.googleapis.com/kubernetes-release/release/v1.6.4/bin/linux/amd64/kubectl /usr/local/bin/kubectl
ENV HOME=/config
RUN set -x && \
    apt-get -y install curl ca-certificates && \
    chmod +x /usr/local/bin/kubectl

USER jenkins
ENTRYPOINT ["/sbin/tini", "--", "/usr/local/bin/jenkins.sh"]

