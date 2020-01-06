#!/usr/bin/env bash

cd /host-home

if [ -d .ssh ] && [ ! -d ~/.ssh ]; then
    sudo cp -r .ssh ~
    sudo chown vscode:vscode -R ~/.ssh
    sudo chmod -R g-rwx,o-rwx ~/.ssh
fi

if [ -f .gitconfig ]; then
    cp .gitconfig ~
    sudo chown vscode:vscode ~/.gitconfig
fi

if [ -d git ] && [ ! -d ~/git ]; then
    cp -r git ~
    sudo chown vscode:vscode -R ~/git
fi

if [ ! -z "${TIMEZONE}" ]; then
    sudo cp "/usr/share/zoneinfo/${TIMEZONE}" /etc/localtime
    echo "${TIMEZONE}" | sudo tee /etc/timezone
fi
