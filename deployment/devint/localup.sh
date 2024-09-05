#!/usr/bin/env bash
basedir=$(
    cd $(dirname $0) || exit
    pwd
)
workspace=${basedir}

bin_name=greenfield-relayer
bin=${workspace}/../../build/${bin_name}

function start() {
    size=$1
    for ((i = 0; i < ${size}; i++)); do
        mkdir -p "${workspace}"/.local/relayer${i}/logs
        docker run --name mechain-relayer${i} -v /data/mechain-relayer/devint/.local/relayer${i}:/app/.greenfield-relayerd \
          -d kevin2025/mechain-relayer ./${bin_name} run --config-type local \
          --config-path .greenfield-relayerd/config.json \
          --log_dir "json > .greenfield-relayerd/logs/relayer.log"
    done
}

function stop() {
    docker rm -f $(docker ps -a | grep mechain-relayer | awk '{print $1}')
}

CMD=$1
SIZE=3
if [ ! -z "$2" ] && [ "$2" -gt "0" ]; then
    SIZE=$2
fi
case ${CMD} in
start)
    echo "===== start ===="
    start "$SIZE"
    echo "===== end ===="
    ;;
stop)
    echo "===== stop ===="
    stop
    echo "===== end ===="
    ;;
*) ;;
esac
