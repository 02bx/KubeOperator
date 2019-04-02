#!/usr/bin/env bash
success=0

function set_firewall {
    firewall-cmd --zone=public --add-port=80/tcp --add-port=8080/tcp --add-port=5380/tcp --add-port=53/udp --add-port=53/tcp --add-port=8081/tcp --add-port=8082/tcp --permanent >/dev/null 2>&1
    systemctl restart firewalld >/dev/null 2>&1
}

function set_iptables {
    iptables -I INPUT -p tcp -port 80 -j ACCEPT
    iptables -I INPUT -p tcp -port 8080 -j ACCEPT
    iptables -I INPUT -p tcp -port 5380 -j ACCEPT
    iptables -I INPUT -p tcp -port 53 -j ACCEPT
    iptables -I INPUT -p tcp -port 8081 -j ACCEPT
    iptables -I INPUT -p tcp -port 8082 -j ACCEPT
    iptables -I INPUT -p upd -port 53 -j ACCEPT
}

which firewall-cmd &> /dev/null
if [[ "$?" == "0" ]];then
    set_firewall
    exit 0
fi

which iptables &> /dev/null
if [[ "$?" == "0" ]];then
    set_iptables
    exit 0
fi
