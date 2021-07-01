---
organization: Turbot
category: ["public cloud"]
icon_url: "/images/plugins/turbot/aws.svg"
brand_color: "#FF9900"
display_name: "Amazon Web Services"
short_name: "aws"
description: "Steampipe plugin for querying instances, buckets, databases and more from AWS."
og_description: "Query AWS with SQL! Open source CLI. No DB required." 
og_image: "/images/plugins/turbot/aws-social-graphic.png"
---

# AWS + Steampipe

[AWS](https://aws.amazon.com/) provides on-demand cloud computing platforms and APIs to authenticated customers on a metered pay-as-you-go basis. 

[Steampipe](https://steampipe.io) is an open source CLI to instantly query cloud APIs using SQL.

For example:

```sql
select
  title,
  create_date,
  mfa_enabled
from
  aws_iam_user
```

```
+-----------------+---------------------+-------------+
| title           | create_date         | mfa_enabled |
+-----------------+---------------------+-------------+
| pam_beesly      | 2005-03-24 21:30:00 | false       |
| creed_bratton   | 2005-03-24 21:30:00 | true        |
| stanley_hudson  | 2005-03-24 21:30:00 | false       |
| michael_scott   | 2005-03-24 21:30:00 | false       |
| dwight_schrute  | 2005-03-24 21:30:00 | true        |
+-----------------+---------------------+-------------+
```

## Documentation

- **[Table definitions & examples →](/plugins/turbot/aws/tables)**

## Get started

### Install

Download and install the latest AWS plugin:

```bash
steampipe plugin install aws
```

### Credentials

| Item | Description |
| - | - |
| Credentials | Specify a named profile from an AWS credential file with the `profile` argument. |
| Permissions | Grant the `ReadOnlyAccess` policy to your user or role. |
| Radius | Each connection represents a single AWS account. |
| Resolution |  1. Credentials explicitly set in a steampipe config file (`~/.steampipe/config/aws.spc`).<br />2. Credentials specified in environment variables e.g. `AWS_ACCESS_KEY_ID`.<br />3. Credentials in the credential file (`~/.aws/credentials`) for the profile specified in the `AWS_PROFILE` environment variable.<br />4. Credentials for the Default profile from the credential file.<br />5. EC2 Instance Role Credentials (if running on an ec2 instance) |
| Region Resolution | 1. Regions set for the connection via the `regions` argument in the config file (`~/.steampipe/config/aws.spc`).<br /> 2. The region specified in the `AWS_DEFAULT_REGION` or `AWS_REGION` environment variable<br />3. The region specified in the active profile (`AWS_PROFILE` or `default`). |

### Configuration

Installing the latest aws plugin will create a config file (`~/.steampipe/config/aws.spc`) with a single connection named `aws`: 
```hcl
connection "aws" {
  plugin    = "aws"    

  # You may connect to one or more regions. If `regions` is not specified, 
  # Steampipe will use a single default region using the same resolution 
  # order as the aws cli:
  #  1. The `AWS_DEFAULT_REGION` or `AWS_REGION` environment variable
  #  2. The region specified in the active profile (`AWS_PROFILE` or default)
  #regions     = ["us-east-1", "us-west-2"]

  # If no credentials are specified, the plugin will use the AWS credentials 
  # resolver to get the current credentials in the same manner as the CLI 
  #  Alternatively, you may set static credentials with the `access_key`, 
  # `secret_key`, and `session_token` arguments, or select a named profile
  # from an AWS credential file with the `profile` argument:
  #profile     = "profile2"
}
```

 By default, all options are commented out in the default connection, thus Steampipe will resolve your region and credentials using the same mechanism as the AWS CLI (AWS environment variables, default profile, etc).  This provides a quick way to get started with Steampipe, but you will probably want to customize your experience using configuration options for [querying multiple regions](#multi-region-connections), [configuring credentials](#configuring-aws-credentials) from your [AWS Profiles](#aws-profile-credentials), [SSO](#aws-sso-credentials), [aws-vault](#aws-vault-credentials) etc.


## Get Involved

* Open source: https://github.com/turbot/steampipe-plugin-aws
* Community: [Slack Channel](https://join.slack.com/t/steampipe/shared_invite/zt-oij778tv-lYyRTWOTMQYBVAbtPSWs3g)



## Multi-Region Connections
By default, AWS connections behave like the `aws` cli and connect to a single default region.  Alternatively, you may also specify one or more regions with the `regions` argument:
```hcl
connection "aws" {
  plugin    = "aws"    
  regions   = ["eu-west-1", "ca-central-1", "us-west-2", "ap-southeast-2", "sa-east-1", "ap-northeast-1", "eu-west-3", "ap-northeast-2", "us-east-1",  "eu-central-1", "us-west-1", "us-east-2", "ap-south-1", "eu-north-1",  "ap-southeast-1"]
}
```

The `region` argument supports wildcards:
- All regions 
  ```hcl
  connection "aws" {
    plugin    = "aws"    
    regions   = ["*"] 
  }
  ```
- All regions (gov-cloud)
  ```hcl
  connection "aws" {
    plugin    = "aws"    
    regions   = ["us-gov*"] 
  }
  ```
- All US and EU regions 
  ```hcl
  connection "aws" {
    plugin    = "aws"    
    regions   = ["us-*", "eu-*"] 
  }
  ```

AWS multi-region connections are common, but be aware that performance may be impacted by the number of regions and the latency to them.


## Configuring AWS Credentials

### AWS Profile Credentials
You may specify a named profile from an AWS credential file with the `profile` argument.  A connection per profile, using named profiles is probably the most common configuration:

#### aws credential file:
```ini
[profile_y]
aws_access_key_id = AKIA4YFAKEKEYXTDS252
aws_secret_access_key = SH42YMW5p3EThisIsNotRealzTiEUwXN8BOIOF5J8m
region = us-west-2

[profile_z]
aws_access_key_id = AKIA4YFAKEKEYJ7HS98F
aws_secret_access_key = Apf938vDKd8ThisIsNotRealzTiEUwXj9nKLWP9mg4
```

#### aws.spc:
```hcl
connection "aws_account_y" {
  plugin      = "aws" 
  profile     = "profile_y"
  regions     = ["us-east-1", "us-west-2"]
}

connection "aws_account_z" {
  plugin      = "aws" 
  profile     = "profile_z"
  regions     = ["ap-southeast-1", "ap-southeast-2"]
}
```

Using named profiles allows Steampipe to work with your existing CLI configurations, including SSO and using role assumption.

### AWS SSO Credentials
Steampipe works with [AWS SSO](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-sso.html#sso-configure-profile-auto) via AWS profiles however:
- You must login to SSO (`aws sso login` ) before starting Steampipe
- If your credentials expire, you will need to re-authenticate outside of Steampipe - Steampipe currently cannot re-authenticate you.


#### aws credential file:
```ini

[aws_000000000000] 
sso_start_url = https://d-9a672b0000.awsapps.com/start 
sso_region = us-east-2
sso_account_id = 000000000000
sso_role_name = SSO-ReadOnly
region = us-east-1

```

#### aws.spc:

```hcl
connection "aws_000000000000" {
  plugin    = "aws"  
  profile   = "aws_000000000000"
  regions   = ["us-west-2", "us-east-1",  "us-west-1", "us-east-2"]
}

```

### AssumeRole Credentials (No MFA)

If your aws credential file contains profiles that assume a role via the `source_profile` and `role_arn` options and MFA is not required, Steampipe can use the profile as-is:

#### aws credential file:
```ini
[cli-user]
aws_access_key_id = AKIA4YFAKEKEYXTDS252
aws_secret_access_key = SH42YMW5p3EThisIsNotRealzTiEUwXN8BOIOF5J8m
region = us-west-2

[role-without-mfa]
role_arn = arn:aws:iam::123456789012:role/test_assume
source_profile = cli-user
external_id = xxxxx
```

#### aws.spc:

```hcl
connection "role_aws" {
  plugin  = "aws"
  profile = "role-without-mfa"
  regions = ["us-east-1", "us-east-2"]
}
```


### AssumeRole Credentials (With MFA)

Currently Steampipe doesn't support prompting for an MFA token at run time.  To overcome this problem you will need to generate an AWS profile with temporary credentials.

One way to accomplish this is to use the `credential_process` to [generate the credentials with a script or program](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-sourcing-external.html) and cache the tokens in a new profile.  There is a [sample `mfa.sh` script](https://raw.githubusercontent.com/turbot/steampipe-plugin-aws/main/scripts/mfa.sh) in the `scripts` directory of the [steampipe-plugin-aws](https://github.com/turbot/steampipe-plugin-aws) repo that you can use, and there are several open source projects that automate this process as well.

Note that Steampipe cannot prompt you for your token currently, so you must authenticate before starting Steampipe, and re-authenticate outside of Steampipe whenever your credentials expire.

#### aws credential file:
```bash
[user_account]  
aws_access_key_id = AKIA4YFAKEKEYXTDS252
aws_secret_access_key = SH42YMW5p3EThisIsNotRealzTiEUwXN8BOIOF5J8m
mfa_serial = arn:aws:iam::111111111111:mfa/my_role_mfa

[aws_account_123456789012]
credential_process = sh -c 'mfa.sh arn:aws:iam::123456789012:role/my_role arn:aws:iam::111111111111:mfa/my_role_mfa user_account 2> $(tty)'
```

#### aws.spc:

```hcl
connection "aws_account_123456789012" {
  plugin  = "aws"
  profile = "aws_account_123456789012"
  regions = ["*"]
}
```


### AWS-Vault Credentials
Steampipe can use profiles that use [aws-vault](https://github.com/99designs/aws-vault) via the `credential_process`.  aws-vault can even be used when using AssumeRole Credentials with MFA (You must authenticate/re-authenticate outside of Steampipe whenever your credentials expire if you are using MFA):

#### aws credential file:
```bash
[vault_user_account]
credential_process = /usr/local/bin/aws-vault exec -j vault_user_profile # vault_user_profile is the name of the profile IN AWS_VAULT...

[aws_account_123456789012]
source_profile = vault_user_account
role_arn =  arn:aws:iam::123456789012:role/my_role
mfa_serial = arn:aws:iam::111111111111:mfa/my_role_mfa 
```
#### aws.spc:

```hcl
connection "aws_account_123456789012" {
  plugin  = "aws"
  profile = "aws_account_123456789012"
  regions = ["*"]
}
```


### Key Pair Credentials 
The AWS plugin allows you set static credentials with the `access_key`, `secret_key`, and `session_token` arguments in your connection.  

```hcl
connection "aws_account_x" {
  plugin      = "aws" 
  secret_key  = "gMCYsoGqjfThisISNotARealKeyVVhh"
  access_key  = "ASIA3ODZSWFYSN2PFHPJ"  
  regions     = ["us-east-1" , "us-west-2"]
}
```


### Credentials from Environment Variables
The AWS plugin will use the standard AWS environment variables to obtain credentials **only if other arguments (`profile`, `access_key`/`secret_key`, `regions`) are not specified** in the connection:

```sh
export AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE
export AWS_SECRET_ACCESS_KEY=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
export AWS_DEFAULT_REGION=eu-west-1
export AWS_SESSION_TOKEN=AQoDYXdzEJr...
export AWS_ROLE_SESSION_NAME=steampipe@myaccount
```
```hcl
connection "aws" {
  plugin      = "aws" 
}
```

### Credentials from an EC2 Instance Profile

If you are running Steampipe on a AWS EC2 instance, and that instance has an [instance profile attached](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/iam-roles-for-amazon-ec2.html) then Steampipe will automatically use the associated IAM role without other credentials:

```hcl
connection "aws" {
  plugin      = "aws" 
  regions     = ["eu-west-1", "eu-west-2"]
}
```
