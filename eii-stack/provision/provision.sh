#!/bin/bash -e

# Copyright (c) 2020 Intel Corporation.

# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:

# The above copyright notice and this permission notice shall be included in
# all copies or substantial portions of the Software.

# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.

# Provision
# Usage: sudo -E ./provision <path-of-docker-compose-file>

docker_compose=$1
RED='\033[0;31m'
YELLOW="\033[1;33m"
GREEN="\033[0;32m"
NC='\033[0m' # No Color

function log_warn() {
    echo -e "${YELLOW}WARN: $1 ${NC}"
}

function log_info() {
    echo -e "${GREEN}INFO: $1 ${NC}"
}

function log_error() {
    echo -e "${RED}ERROR: $1 ${NC}"
}

function log_fatal() {
    echo -e "${RED}FATAL: $1 ${NC}"
    exit -1
}

function check_error() {
    if [ $? -ne 0 ] ; then
        log_fatal "$1"
    fi
}

function source_env() {
    EII_ENV="../.env"
    PROVISIONING_ENV=".env"

    set -a
    if [ -f $EII_ENV ]; then
        source $EII_ENV
    fi

    if [ -f $PROVISIONING_ENV ]; then
        source $PROVISIONING_ENV
    fi
    set +a
}


function export_host_ip() {
    if [ -z $HOST_IP ]; then
	hostIP=$(ip route get 1 | awk '{print $7}'|head -1)
        export HOST_IP=$hostIP
    fi
    echo 'System IP Address is:' $HOST_IP
    export no_proxy=$eii_no_proxy,$HOST_IP
}

function set_docker_host_time_zone() {
    echo "Updating .env for container timezone..."
    # Setting RTC to UTC
    timedatectl set-local-rtc 0
    # Get Docker Host timezone
    hostTimezone=`timedatectl | grep "zone" | awk '{print $3}'`
    hostTimezone=`echo $hostTimezone`
    # This will remove the HOST_TIME_ZONE entry if it exists and adds a new one with the right timezone
    sed -i '/HOST_TIME_ZONE/d' ../.env && echo "HOST_TIME_ZONE=$hostTimezone" >> ../.env
}

function create_eii_user() {
    echo "Create $EII_USER_NAME if it doesn't exists. Update UID from env if already exists with different UID"

    # EII containers will be executed as eiiuser
    if ! id $EII_USER_NAME >/dev/null 2>&1; then
        groupadd $EII_USER_NAME -g $EII_UID
        useradd -r -u $EII_UID -g $EII_USER_NAME $EII_USER_NAME
    else
        if ! [ $(id -u $EII_USER_NAME) = $EII_UID ]; then
            usermod -u $EII_UID $EII_USER_NAME
            groupmod -g $EII_UID $EII_USER_NAME
        fi
    fi
}

function create_eii_install_dir() {
    log_info "Creating EII install Directories"
    if [ $ETCD_RESET = 'true' ]; then
        rm -rf $EII_INSTALL_PATH/data/etcd
    fi

    # Creating the required EII dirs
    mkdir -p $EII_INSTALL_PATH/taos/{data,log}
    check_error "Failed to create dir '$EII_INSTALL_PATH/taos/{data,log}"

    mkdir -p $EII_INSTALL_PATH/data/influxdata
    check_error "Failed to create dir '$EII_INSTALL_PATH/data/influxdata'"

    mkdir -p $EII_INSTALL_PATH/data/etcd/data
    check_error "Failed to create dir '$EII_INSTALL_PATH/data/etcd/data'"

    mkdir -p $EII_INSTALL_PATH/sockets/
    check_error "Failed to create dir '$EII_INSTALL_PATH/sockets'"

    mkdir -p $EII_INSTALL_PATH/tools_output
    check_error "Failed to create dir '$EII_INSTALL_PATH/tools_output'"

    mkdir -p $EII_INSTALL_PATH/model_repo
    chown -R $EII_USER_NAME:$EII_USER_NAME $EII_INSTALL_PATH

    if [ -d $TC_DISPATCHER_PATH ]; then
        chown -R $EII_USER_NAME:$EII_USER_NAME $TC_DISPATCHER_PATH
        chmod -R 760 $TC_DISPATCHER_PATH
    fi
}

function create_docker_network() {
	echo "Checking eii docker network.."
	docker_networks=$(docker network ls --filter name=eii| awk '{print $2}')
	docker_net_list=(`echo ${docker_networks}`);
	network_present=false
	for dn in "${docker_net_list[@]}"
	do
	    if [[ "$dn" = "eii" ]];then
	        network_present=true
	        break
	    fi
	done
	if [[ "$network_present" = false ]]; then
		echo "Creating eii docker bridge network as it is not present.."
		docker network create eii
	fi
}

function copy_docker_compose_file() {
    echo "Copying docker compose yaml file which is provided as an argument."
    # This file will be volume mounted inside the provisioning container and deleted once privisioning it done
    if ! [ -f $docker_compose ] || [ -z $docker_compose ]; then
        log_error "Supplied docker compose file '$docker_compose' does not exists"
        log_fatal "Usage: $ sudo -E ./provision.sh <path_to_eii_docker_compose_file>"
    else
        cp $docker_compose ./docker-compose.yml
    fi
}

function remove_docker_compose_file() {
    if [ -f "docker-compose.yml" ]; then
        echo "Removing copied docker compose yaml file ..."
        rm -rf ./docker-compose.yml
    fi
}

function remove_client_server() {
    echo "Deleting client and server cert files ..."
    if [ -d "client" ]; then
        rm -rf client
    fi
    if [ -d "server" ]; then
        rm -rf server
    fi
}

function prod_mode_gen_certs() {
    log_info "Generating EII Certificates"
    if [ -d "rootca" ]; then
        log_warn "Making use of existing CA from ./rootca dir for generating certs..."
        log_warn "To generate new CA, remove roootca/ from current dir.."
        python3 gen_certs.py --f $docker_compose --capath rootca/
    else
        python3 gen_certs.py --f $docker_compose
    fi
    chown -R $EII_USER_NAME:$EII_USER_NAME Certificates/
    chmod -R 750 Certificates/

}

function check_ETCD_port() {
    log_info "Checking if ETCD ports are already up..."
    ports=($ETCD_CLIENT_PORT)
    for port in "${ports[@]}"
    do
       set +e
       fuser $port/tcp
       if [ $? -eq 0 ]; then
           log_fatal "$port is already being used, so please kill that process and re-run the script."
       fi
       set -e
    done
}

function install_pip_requirements() {
    log_info "Installing dependencies.."
    pip3 install -r cert_requirements.txt
}


function remove_eii_containers() {
    log_info "Bringing down running eii containers"
    for container_id in $(docker ps  --filter="name=ia_" -q -a)
    do
        docker rm -f $container_id
    done
}

source_env
export_host_ip
set_docker_host_time_zone

if [ "$2" != "--run_etcd_provision" ]; then
    create_eii_user
    create_eii_install_dir
    create_docker_network
fi

if [ $DEV_MODE = 'false' ]; then
    chmod -R 760 $EII_INSTALL_PATH/data
    chmod -R 760 $EII_INSTALL_PATH/sockets
else
    chmod -R 755 $EII_INSTALL_PATH/data
    chmod -R 755 $EII_INSTALL_PATH/sockets
fi

if [ "$1" == "--worker" -o "$ETCD_NAME" == "worker" ]; then
    remove_eii_containers
    log_info "Provisioning is Successful..."
    exit 0
fi

#############################################################
function build_etcd_and_etcd_provision() {
    echo "Building etcd and etcd_provision images..."
    docker-compose -f dep/docker-compose-etcd.yml -f dep/docker-compose-etcd.override.build.yml build
    docker-compose -f dep/docker-compose-etcd-provision.yml -f dep/docker-compose-etcd-provision.override.build.yml build
}

# for arg in "$@"
# do
    # if [ "$arg" == "--build" ]
    # then
        # build_etcd_and_etcd_provision
    # fi
# done
if [ "$1" == "--run_etcd" ] ; then
    remove_eii_containers
    echo "Checking ETCD port..."
    check_ETCD_port

    if [ $DEV_MODE = 'true' ]; then
        log_info "EII is not running in Secure mode. Generating certificates is not required.. "
        log_info "Starting ETCD ..."
        docker-compose -f dep/docker-compose-etcd.yml up -d
    else
        log_info "Starting ETCD ..."
        docker-compose -f dep/docker-compose-etcd.yml -f dep/docker-compose-etcd.override.prod.yml up -d
    fi
elif [ "$2" = "--run_etcd_provision" ] ; then
    echo "Checking ETCD health..."
    ./dep/etcd_health_check.sh

    copy_docker_compose_file

    if [ $DEV_MODE = 'true' ]; then
        log_info "EII is not running in Secure mode. Generating certificates is not required.. "
        log_info "Starting etcd_provision..."
        docker-compose -f dep/docker-compose-etcd-provision.yml up
    else
        log_info "Starting etcd_provision..."
        docker-compose -f dep/docker-compose-etcd-provision.yml -f dep/docker-compose-etcd-provision.override.prod.yml up
    fi
elif [ "$2" = "--generate_certs" ] ; then
    echo "Clearing existing Certificates..."
    rm -rf Certificates

    install_pip_requirements
    copy_docker_compose_file
    prod_mode_gen_certs
elif [ "$ETCD_NAME" = "leader" ]; then
    remove_eii_containers
    copy_docker_compose_file

    echo "Clearing existing Certificates..."
    rm -rf Certificates

    echo "Checking ETCD port..."
    check_ETCD_port

    if [ $DEV_MODE = 'true' ]; then
        log_info "EII is not running in Secure mode. Generating certificates is not required.. "
        log_info "Starting and provisioning ETCD ..."
        docker-compose -f dep/docker-compose-etcd.yml up -d
        docker-compose -f dep/docker-compose-etcd-provision.yml up -d
    else
        install_pip_requirements
        prod_mode_gen_certs
        log_info "Starting and provisioning ETCD ..."
        docker-compose -f dep/docker-compose-etcd.yml -f dep/docker-compose-etcd.override.prod.yml up -d
        docker-compose -f dep/docker-compose-etcd-provision.yml -f dep/docker-compose-etcd-provision.override.prod.yml up
    fi
fi

remove_client_server
remove_docker_compose_file

log_info "Provisioning is Successful..."
