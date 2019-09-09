#!/usr/bin/env bash
BASE_DIR=$(dirname "$0")
source ${BASE_DIR}/utils.sh



function main {
    mkdir -p /opt/kubeOperator >/dev/null 2>&1
    cp -i ${SCRIPTS_DIR}/service/kubeops.service /etc/systemd/system/
    cp -ri ${PROJECT_DIR}/* /opt/kubeOperator/
}
main