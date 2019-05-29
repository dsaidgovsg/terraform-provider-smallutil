#!/usr/bin/env sh
if ! which tera >/dev/null; then
    echo "Run \"cargo install tera-cli\" to install tera first.\nLinux users may download from https://github.com/guangie88/tera-cli/releases instead."
    return 1
fi

DIR="$(cd "$(dirname "$(readlink -f "$0")")" >/dev/null 2>&1 && pwd)"
tera -f "${DIR}/.travis.yml.tmpl" --yaml "${DIR}/vars.yml" > "${DIR}/../.travis.yml"
echo "Successfully applied template into .travis.yml!"
