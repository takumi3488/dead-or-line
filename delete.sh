#!/bin/bash -e

# パスを通す
source ~/.bash_profile
export PATH="$PATH:$(zsh -c 'echo $PATH' || echo '')"

cd terraform
terraform destroy