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
        local ended_at
        ended_at="$(date +%s%N)"

        go run "$FLIGHT_BIN" record \
            --cwd "$PWD" \
            --exit-code "$exit_code" \
            --started-at "$FLIGHT_COMMAND_STARTED_AT" \
            --ended-at "$ended_at" \
            "$FLIGHT_LAST_COMMAND" \
            >/dev/null 2>&1

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