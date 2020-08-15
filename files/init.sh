#!/bin/bash

HOMEDIR=/home/coder

sudo usermod -l $USER coder
sudo chown -R 1001:1001 /home/coder

cp -rT /etc/skel $HOMEDIR

if [ ! -f $HOMEDIR/.ssh/repo_ssh_rsa.pub ]; then
  cat /dev/zero | ssh-keygen -t rsa -q -N "" -C "$EMAIL" -f ~/.ssh/repo_ssh_rsa
fi

if [ ! -f $HOMEDIR/.ssh/config ]; then
  cp /etc/skel/.ssh/config $HOMEDIR/.ssh/config
fi

git config --global user.email "$EMAIL"
git config --global user.name "$USER_FORENAME $USER_SURNAME"

echo "export TEST_TMPDIR=/workspace/.cache" >> $HOMEDIR/.bashrc

yarn theia start /home/coder/workspace --hostname=0.0.0.0 --port=3000
