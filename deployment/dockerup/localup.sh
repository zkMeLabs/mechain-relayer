#!/usr/bin/env bash
basedir=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
workspace=${basedir}
local_env="${workspace}/.local"

#########################
# the command line help #
#########################
function display_help() {
    echo "Usage: $0 [option...] {help|config}" >&2
    echo
    echo "   --help           display help info"
    echo "   --config         generate config.json"
    echo "   --print          print relayer local env work directory"
    echo
    exit 0
}

#############################################################
# make sp config.toml according to env.info/db.info/sp.info #
#############################################################
function make_config() {
    size=$1
    rm -rf "${local_env}"
    mkdir -p "${local_env}"
    RELAYER_FILE="${workspace}/validator.json"
    CONFIG_FILE="${workspace}/config.json"

    for i in $(seq 0 $((size - 1))); do
        mkdir -p "${local_env}/relayer${i}/logs"
        bls_priv_key=$(jq -r ".validator${i}.bls_key" "$RELAYER_FILE")
        relayer_key=$(jq -r ".validator${i}.relayer_key" "$RELAYER_FILE")
        OUTPUT_FILE="${local_env}/relayer${i}/config.json"
        jq --arg bls_priv_key "$bls_priv_key" \
            --arg relayer_key "$relayer_key" \
            '.greenfield_config.bls_private_key = $bls_priv_key |
            .greenfield_config.private_key = $relayer_key |
            .bsc_config.private_key = $relayer_key' \
            "$CONFIG_FILE" >"$OUTPUT_FILE"
        echo "Saved to $OUTPUT_FILE"
    done
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
--help | *)
    display_help
    ;;
esac
