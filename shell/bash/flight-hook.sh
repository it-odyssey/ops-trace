#!/usr/bin/env bash

__flight_preexec() {
    local cmd="$BASH_COMMAND"

    case "$cmd" in
    __flight_*|__vsc_*|__zoxide_*|printf\ "\\033]*"|trap*|PROMPT_COMMAND=*|unset\ FLIGHT_*)
        return
        ;;
    esac

    FLIGHT_LAST_COMMAND="$cmd"
    FLIGHT_COMMAND_STARTED_AT="$(date +%s%N)"
}

__flight_precmd() {
    local exit_code=$?

    if [[ -n "${FLIGHT_LAST_COMMAND:-}" ]]; then
        echo "FLIGHT: $FLIGHT_LAST_COMMAND | exit=$exit_code"
        unset FLIGHT_LAST_COMMAND
        unset FLIGHT_COMMAND_STARTED_AT
    fi
}

# Install DEBUG trap.
trap '__flight_preexec' DEBUG

# Add our precmd hook only if it isn't already installed.
if [[ "${PROMPT_COMMAND:-}" != *"__flight_precmd"* ]]; then
    if [[ -n "${PROMPT_COMMAND:-}" ]]; then
        PROMPT_COMMAND="__flight_precmd;${PROMPT_COMMAND}"
    else
        PROMPT_COMMAND="__flight_precmd"
    fi
fi

# Clear anything captured while this hook was installing itself.
unset FLIGHT_LAST_COMMAND
unset FLIGHT_COMMAND_STARTED_AT