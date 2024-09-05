#!/usr/bin/env bash
basedir=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
workspace=${basedir}
project="$(realpath "${workspace}/../..")"
bin="$(realpath "${workspace}/../..")"
bin_name=greenfield-relayer

#########################
# the command line help #
#########################
function display_help() {
    echo "Usage: $0 [option...] {help|generate|reset|start|stop|print}" >&2
    echo
    echo "   --help           display help info"
    echo "   --generate       generate config.json that accepts four args: the first arg is json file path, the second arg is db username, the third arg is db password and the fourth arg is db address"
    echo "   --reset          reset env"
    echo "   --start          start relayers"
    echo "   --stop           stop relayers"
    echo "   --clean          clean relayer env"
    echo "   --rebuild        rebuild relayer code"
    echo "   --print          print relayer local env work directory"
    echo
    exit 0
}

#############################################################
# make sp config.toml according to env.info/db.info/sp.info #
#############################################################
function make_config() {
    size=$1
    rm -rf "${workspace}"/.local
    mkdir -p "${workspace}"/.local
    RELAYER_FILE="${project}/relayer.yaml"
    CONFIG_FILE="${workspace}/config.json"

    for i in {0..3}; do
        mkdir -p "${workspace}"/.local/relayer"${i}"/logs
        bls_priv_key=$(grep "validator_bls$i bls_priv_key" "$RELAYER_FILE" | awk '{print $3}')
        relayer_key=$(grep "relayer$i relayer_key" "$RELAYER_FILE" | awk '{print $3}')

        OUTPUT_FILE="${workspace}/.local/relayer${i}/config.json"
        jq --arg bls_priv_key "$bls_priv_key" \
            --arg relayer_key "$relayer_key" \
            '.greenfield_config.bls_private_key = $bls_priv_key |
        .greenfield_config.private_key = $relayer_key |
        .bsc_config.private_key = $relayer_key' \
            "$CONFIG_FILE" >"$OUTPUT_FILE"
        echo "save to $OUTPUT_FILE"
    done

}

function start() {
    size=$1
    for ((i = 0; i < ${size}; i++)); do
        nohup "${bin}" run --config-type local \
            --config-path "${workspace}"/../../config/local/config_local_${i}.json \
            --log_dir json >"${workspace}"/.local/relayer${i}/logs/relayer.log &
    done
}

function stop() {
    ps -ef | grep ${bin_name} | awk '{print $2}' | xargs kill
}

######################
# clean local env #
######################
function clean_local_env() {
    rm -rf "${workspace:?}/.local"
}

CMD=$1
SIZE=3
if [ -n "$2" ] && [ "$2" -gt "0" ]; then
    SIZE=$2
fi
case ${CMD} in
config)
    make_config "$SIZE"
    ;;
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
clean)
    clean_local_env
    ;;
--help | *)
    display_help
    ;;
esac
