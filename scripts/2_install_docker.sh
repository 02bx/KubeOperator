#!/usr/bin/env bash
source ./utils.sh
OS=$(uname)

OFFLINE_DOCKER_DIR="${BASE_DIR}/docker/bin"

function install_docker_online {
    yum -y remove docker docker-common docker-selinux docker-engine
    yum install -y epel-release yum-utils device-mapper-persistent-data lvm2
    yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
    yum install -y docker-ce docker-compose
}

function install_docker_offline {
    cp ${OFFLINE_DOCKER_DIR}/docker* /usr/bin/
    cp ${OFFLINE_DOCKER_DIR}/docker.service /etc/systemd/system/
    chmod +x /usr/bin/docker* && chmod 754 /etc/systemd/system/docker.service
}

function install_docker {
    echo ">> Install docker"
    if [[ "${OS}" == "Darwin" ]];then
        echo "Platform is MacOS, install manually"
        return
    fi
    if [[ -f "${OFFLINE_DOCKER_DIR}/dockerd" ]];then
        install_docker_offline
    else
        install_docker_online
    fi
}

function start_docker {
    systemctl start docker
    systemctl enable docker
}

function main {
    install_docker
    start_docker
}

main
