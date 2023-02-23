#!/bin/sh

############################################################################################################################
#                 Sample AWS credential process script for MFA login with sts:assumerole
#
# To use, put this file somewhere in your PATH, and then set up 2 profiles:
#  - A master account profile with your access keys and MFA (in the account where your user resides)
#  - The destination account to log into with a role.  This is the profile that should be configured to use this script,
#    passing the role ARN, MFA ARN, and master account profile as arguments.
#     
# For example:
#   [user_master_account]  
#   aws_access_key_id = AKIA4YFAKEKEYXTDS252
#   aws_secret_access_key = SH42YMW5p3EThisIsNotRealzTiEUwXN8BOIOF5J8m
#   mfa_serial = arn:aws:iam::111111111111:mfa/my_role_mfa
#
#   [destination_account]
#   credential_process = sh -c 'mfa.sh arn:aws:iam::123456789012:role/my_role arn:aws:iam::111111111111:mfa/my_role_mfa user_master_account 2> $(tty)'
#
############################################################################################################################
set -e

role=$1
mfa_arn=$2
profile=$3
temp_profile=temp_mfa_session_${3}

if [ -z $role ]; then echo "no role specified"; exit 1; fi
if [ -z $mfa_arn ]; then echo "no mfa arn specified"; exit 1; fi
if [ -z $profile ]; then echo "no profile specified"; exit 1; fi

set +e
resp=$(aws sts get-caller-identity --profile $temp_profile | jq '.UserId') 
set -e 

if [ ! -z $resp ]; then
    echo '{
        "Version": 1,
        "AccessKeyId": "'"$(aws configure get aws_access_key_id --profile $temp_profile)"'",
        "SecretAccessKey": "'"$(aws configure get aws_secret_access_key --profile $temp_profile)"'",
        "SessionToken": "'"$(aws configure get aws_session_token --profile $temp_profile)"'",
        "Expiration": "'"$(aws configure get expiration --profile $temp_profile)"'"
    }'
    exit 0
fi
read -p "Enter MFA token: " mfa_token

if [ -z $mfa_token ]; then echo "MFA token can't be empty"; exit 1; fi

data=$(aws sts assume-role --role-arn $role \
                    --profile $profile \
                    --role-session-name "$mfa_token-$profile" \
                    --serial-number $mfa_arn \
                    --token-code $mfa_token | jq '.Credentials')


aws_access_key_id=$(echo $data | jq -r '.AccessKeyId')
aws_secret_access_key=$(echo $data | jq -r '.SecretAccessKey')
aws_session_token=$(echo $data | jq -r '.SessionToken')
expiration=$(echo $data | jq -r '.Expiration')

aws configure set aws_access_key_id $aws_access_key_id --profile $temp_profile
aws configure set aws_secret_access_key $aws_secret_access_key --profile $temp_profile
aws configure set aws_session_token $aws_session_token --profile $temp_profile
aws configure set expiration $expiration --profile $temp_profile

echo '{
  "Version": 1,
  "AccessKeyId": "'"$aws_access_key_id"'",
  "SecretAccessKey": "'"$aws_secret_access_key"'",
  "SessionToken": "'"$aws_session_token"'",
  "Expiration": "'"$expiration"'"
}'
