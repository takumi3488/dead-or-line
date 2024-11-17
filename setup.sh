#!/bin/bash -e

# パスを通す
source ~/.bash_profile
export PATH="$PATH:$(zsh -c 'echo $PATH' || echo '')"

cd terraform
terraform init
read -p "監視対象URL: " target_url
read -p "通知間隔(分): " notice_interval
read -p "DynamoDBのテーブル名: " dynamodb_table_name
read -p "DynamoDBのパーティションキー名: " dynamodb_key
read -p "LINE Bot用のchannel access token: " line_token
read -p "LINE Botで送信するメッセージ: " line_base_message
read -p "LINEの送信先ID(入力しなければブロードキャストする): " line_to
read -p "監視間隔(分): " rate_minutes
terraform apply -var="target_url=${target_url}" -var="notice_interval=${notice_interval}" -var="dynamodb_table_name=${dynamodb_table_name}" -var="dynamodb_key=${dynamodb_key}" -var="line_token=${line_token}" -var="line_base_message=${line_base_message}" -var="line_to=${line_to}" -var="rate_minutes=${rate_minutes}"
