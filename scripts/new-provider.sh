#!/usr/bin/env bash

which tput 2>&1 > /dev/null
if [[ $? -eq "0" ]]; then 
    _COLOR_NC="$(tput sgr0)"
    _COLOR_RED="$(tput setaf 1)"
    _COLOR_GREEN="$(tput setaf 2)"
    _COLOR_YELLOW="$(tput setaf 3)"
    _COLOR_BLUE="$(tput setaf 4)"
fi

# script should be run without sudo access by default
if [[ "${UID}" -eq 0 ]]; then
    echo_erro "this script should not be run as root"
fi

function echo_info {
# @brief: print message to stderr with INFO prefix
    echo -e "[${_COLOR_GREEN}INFO${_COLOR_NC}]" "$@" >&2
}

function echo_erro {
# @brief: print message to stderr with ERRO prefix
    echo -e "[${_COLOR_RED}ERRO${_COLOR_NC}]" "$@" >&2
    exit 1
}

function echo_warn {
# @brief: print message to stderr with WARN prefix
    echo -e "[${_COLOR_YELLOW}WARN${_COLOR_NC}]" "$@" >&2
}

function echo_extra {
# @brief: print message to stderr with EXTR prefix
    echo -e "[${_COLOR_BLUE}EXTR${_COLOR_NC}]" "$@" >&2
}

function to_lower {
# @brief: transform variable to lowercase
    echo "${1}" | awk '{print tolower($0)}'
}

function provider_go {
    _CAP_PROVIDER=$(echo ${PROVIDER} | awk '{ print toupper( substr( $0, 1, 1 ) ) substr( $0, 2 ); }')

    cat <<EOF
//go:build ${1}

/*
FIXME: Package ${1} ...
*/
package ${1}_test

import (
	// embed is used here for including describe.sql file during compilation.
	_ "embed"

	"github.com/shanduur/squat/providers"
)

//go:embed describe.sql
var describeQuery string

type ${PROVIDER}Config struct {
	ProviderName string            \`toml:"provider-name"\`
	Formats      providers.Formats \`toml:"formats"\`
}

type ${_CAP_PROVIDER}Provider struct {
	cfg ${PROVIDER}Config
}

func (pg *${_CAP_PROVIDER}Provider) Initialize(configPath string) (err error)

func (pg ${_CAP_PROVIDER}Provider) ProviderName() string

func (pg ${_CAP_PROVIDER}Provider) GetTableDescription(name string)

func (pg ${_CAP_PROVIDER}Provider) DateFormat()

func (pg ${_CAP_PROVIDER}Provider) DateTimeFormat()
EOF
}

function provider_test_go {
    cat <<EOF
//go:build ${1}

package ${1}_test
EOF
}

function provider_api_go {
    cat <<EOF
//go:build ${1}

package api

import (
	"log"
	"os"
	"path"

	"github.com/shanduur/squat/providers/${1}"
)

func init() {
	if p, err := ${1}.New(path.Join(os.Getenv("CONFIG_LOCATION"), "${1}.toml")); err != nil {
		log.Printf("unable to create new provider connection: %s", err.Error())
		log.Printf("check if env variables CONFIG_LOCATION and DATA_LOCATION are set")
	} else {
		Providers[p.ProviderName()] = p
	}
}
EOF
}

function provider_sql {
    cat <<EOF
select ... from ... as ...; -- FIXME: ...
EOF
}

function main {
    read -p "Provider Name: " PROVIDER
    if [[ -z ${PROVIDER} ]]; then
        echo_erro "wrong provider name"
    fi
    PROVIDER=$(to_lower ${PROVIDER})
    if [[ -d ./providers/${PROVIDER} ]]; then
        echo_erro "provider already exist"
    fi

    read -p "Is non-free? (y/n): " NONFREE
    case $NONFREE in
        y)
            echo_info "package is non free"
        ;;

        n)
            echo_info "package is free"
        ;;

        *)
            echo_erro "wrong answer"
        ;;
    esac

    echo_info "creating provider boilerplate"
    mkdir -p ./providers/${PROVIDER}
    provider_go ${PROVIDER} > ./providers/${PROVIDER}/${PROVIDER}.go
    provider_test_go ${PROVIDER} > ./providers/${PROVIDER}/${PROVIDER}_test.go
    provider_sql ${PROVIDER} > ./providers/${PROVIDER}/describe.sql
    provider_api_go ${PROVIDER} > ./server/api/${PROVIDER}.go
    
    echo_info "configuring packages"
    if [[ ${NONFREE} -eq "y" ]]; then
        JSON=$(cat .dev/providers.json)
        echo ${JSON} | jq ".nonfree += [\"${PROVIDER}\"]" > .dev/providers.json
    else
        JSON=$(cat .dev/providers.json)
        echo ${JSON} | jq ".free += [\"${PROVIDER}\"]" > .dev/providers.json
    fi

    echo_info "done!"
}

main
