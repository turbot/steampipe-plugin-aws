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
| Resolution |  1. Credentials explicitly set in a Steampipe config file (`~/.steampipe/config/aws.spc`).<br />2. Credentials specified in environment variables e.g. `AWS_ACCESS_KEY_ID`.<br />3. Credentials in the credential file (`~/.aws/credentials`) for the profile specified in the `AWS_PROFILE` environment variable.<br />4. Credentials for the Default profile from the credential file.<br />5. EC2 Instance Role Credentials (if running on an ec2 instance) |

### Configuration

Installing the latest aws plugin will create a config file (`~/.steampipe/config/aws.spc`) with a single connection named `aws`:
```hcl
connection "aws" {
  plugin = "aws"

  # `regions` defines the list of regions that Steampipe should target for
  # each query. API calls are made to multiple regions in parallel. The regions
  # list may include wildcards (e.g. *, us-*, us-??st-1).
  # If `regions` is not specified, Steampipe will target the `default_region`
  # only.
  #regions = ["*"] # All regions
  #regions = ["eu-*"] # All EU regions
  #regions = ["us-east-1", "eu-west-2"] # Specific regions

  # Some AWS APIs (e.g. describe EC2 regions, S3 get bucket location) have
  # global results, so can be run against any region. For faster results, you
  # may set your default (closest) region to use for these situations.
  # If not specified, it is resolved in this order:
  #  1. The `AWS_REGION` or `AWS_DEFAULT_REGION` environment variable
  #  2. The region specified in the active profile (`AWS_PROFILE` or default)
  #  3. The main region for the partition as best guessed from the `regions` list.
  #  4. us-east-1, the main region for the most common partition.
  #default_region = "eu-west-2"

  # If no credentials are specified, the plugin will use the AWS credentials
  # resolver to get the current credentials in the same manner as the CLI.
  # Alternatively, you may set static credentials with the `access_key`,
  # `secret_key`, and `session_token` arguments, or select a named profile
  # from an AWS credential file with the `profile` argument:
  #profile = "myprofile"

  # The maximum number of attempts (including the initial call) Steampipe will
  # make for failing API calls. Can also be set with the AWS_MAX_ATTEMPTS
  # environment variable.
  # Defaults to 9 and must be greater than or equal to 1.
  #max_error_retry_attempts = 9

  # The minimum retry delay in milliseconds after which retries will be
  # performed.  This delay is also used as a base value when calculating the
  # exponential backoff retry times.
  # Defaults to 25ms and must be greater than or equal to 1ms.
  #min_error_retry_delay = 25

  # List of additional AWS error codes to ignore for all queries.
  # By default, common not found error codes are ignored and will still be ignored even if this argument is not set.
  #ignore_error_codes = ["AccessDenied", "AccessDeniedException", "NotAuthorized", "UnauthorizedOperation", "UnrecognizedClientException", "AuthorizationError"]

  # Specify the endpoint URL used when making requests to AWS services.
  # If not set, the default AWS generated endpoint will be used.
  # Can also be set with the AWS_ENDPOINT_URL environment variable.
  #endpoint_url = "http://localhost:4566"

  # Set to `true` to force S3 requests to use path-style addressing,
  # i.e., `http://s3.amazonaws.com/BUCKET/KEY`. By default, the S3 client
  # will use virtual hosted bucket addressing when possible (`http://BUCKET.s3.amazonaws.com/KEY`).
  #s3_force_path_style = false
}
```

By default, all options are commented out in the default connection, thus Steampipe will resolve your region and credentials using the same mechanism as the AWS CLI (AWS environment variables, default profile, etc). This provides a quick way to get started with Steampipe, but you will probably want to customize your experience using configuration options for [querying multiple regions](#multi-region-connections), [configuring credentials](#configuring-aws-credentials) from your [AWS Profiles](#aws-profile-credentials), [SSO](#aws-sso-credentials), [aws-vault](#aws-vault-credentials) etc.

## Multi-Region Connections

By default, AWS connections behave like the `aws` cli and connect to a single default region. Alternatively, you may also specify one or more regions with the `regions` argument:
```hcl
connection "aws" {
  plugin  = "aws"
  regions = ["eu-west-1", "ca-central-1", "us-west-2"]
}
```

The `regions` argument supports wildcards:
- All standard regions
  ```hcl
  connection "aws" {
    plugin  = "aws"
    regions = ["*"]
  }
  ```
- All standard US and EU regions
  ```hcl
  connection "aws" {
    plugin  = "aws"
    regions = ["us-*", "eu-*"]
  }
  ```
- All US GovCloud regions
  ```hcl
  connection "aws" {
    plugin  = "aws"
    regions = ["us-gov*"]
  }
  ```
- All CN regions
  ```hcl
  connection "aws" {
    plugin  = "aws"
    regions = ["cn-*"]
  }
  ```
- All ISO regions
  ```hcl
  connection "aws" {
    plugin  = "aws"
    regions = ["us-iso-*"]
  }
  ```
- All ISOB regions
  ```hcl
  connection "aws" {
    plugin  = "aws"
    regions = ["us-isob-*"]
  }
  ```

AWS multi-region connections are common, but be aware that performance may be impacted by the number of regions and the latency to them.

Steampipe will automatically guess your `default_region` from your AWS config
(e.g. `AWS_REGION` env var) or `regions` list, but you may prefer to specify it
to ensure where API calls are made for global resources (e.g. STS, EC2 describe
regions):
```hcl
connection "aws" {
  plugin  = "aws"
  regions = ["us-*"]
  default_region = "us-west-1"
}
```

## Multi-Account Connections

You may create multiple aws connections:
```hcl
connection "aws_dev" {
  plugin  = "aws"
  profile = "aws_dev"
  regions = ["us-east-1", "us-west-2"]
}

connection "aws_qa" {
  plugin  = "aws"
  profile = "aws_qa"
  regions = ["*"]
}

connection "aws_prod" {
  plugin  = "aws"
  profile = "aws_prod"
  regions = ["us-*"]
}
```

Each connection is implemented as a distinct [Postgres schema](https://www.postgresql.org/docs/current/ddl-schemas.html). As such, you can use qualified table names to query a specific connection:

```sql
select * from aws_qa.aws_account
```

You can multi-account connections by using an [**aggregator** connection](https://steampipe.io/docs/using-steampipe/managing-connections#using-aggregators). Aggregators allow you to query data from multiple connections for a plugin as if they are a single connection. 

```hcl
connection "aws_all" {
  plugin      = "aws"
  type        = "aggregator"
  connections = ["aws_dev", "aws_qa", "aws_prod"]
}
```

Querying tables from this connection will return results from the `aws_dev`, `aws_qa`, and `aws_prod` connections:
```sql
select * from aws_all.aws_account
```

Alternatively, can use an unqualified name and it will be resolved according to the [Search Path](https://steampipe.io/docs/guides/search-path). It's a good idea to name your aggregator first alphbetically, so that it is the first connection in the search path (i.e. `aws_all` comes before `aws_dev`):
```sql
select * from aws_account
```

Steampipe supports the `*` wildcard in the connection names. For example, to aggregate all the AWS plugin connections whose names begin with `aws_`:

```hcl
connection "aws_all" {
  type        = "aggregator"
  plugin      = "aws"
  connections = ["aws_*"]
}
```

Aggregators are powerful, but they are not infinitely scalable. Like any other Steampipe connection, they query APIs and are subject to API limits and throttling. Consider as an example and aggregator that includes 3 AWS connections, where each connection queries 16 regions. This means you essentially run the same list API calls 48 times! When using aggregators, it is especially important to:
- Query only what you need! `select * from aws_s3_bucket` must make a list API call in each connection, and then 11 API calls *for each bucket*, where `select name, versioning_enabled from aws_s3_bucket` would only require a single API call per bucket.
- Consider extending the [cache TTL](https://steampipe.io/docs/reference/config-files#connection-options). The default is currently 300 seconds (5 minutes). Obviously, anytime Steampipe can pull from the cache, its is faster and less impactful to the APIs. If you don't need the most up-to-date results, increase the cache TTL!

## Configuring AWS Credentials

### AWS Profile Credentials

You may specify a named profile from an AWS credential file with the `profile` argument. A connection per profile, using named profiles is probably the most common configuration:

#### aws credential file:

```ini
[account_a]
aws_access_key_id = AKIA4YFAKEKEYXTDS252
aws_secret_access_key = SH42YMW5p3EThisIsNotRealzTiEUwXN8BOIOF5J8m
region = us-west-2

[account_b]
aws_access_key_id = AKIA4YFAKEKEYJ7HS98F
aws_secret_access_key = Apf938vDKd8ThisIsNotRealzTiEUwXj9nKLWP9mg4
```

#### aws.spc:

```hcl
connection "aws_account_a" {
  plugin  = "aws"
  profile = "account_a"
  regions = ["us-east-1", "us-west-2"]
}

connection "aws_account_b" {
  plugin  = "aws"
  profile = "account_b"
  regions = ["ap-southeast-1", "ap-southeast-2"]
}
```

Using named profiles allows Steampipe to work with your existing CLI configurations, including SSO and using role assumption.

### AWS SSO Credentials

Steampipe works with [AWS SSO](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-sso.html#sso-configure-profile-auto) via AWS profiles however:
- You must login to SSO (`aws sso login` ) before starting Steampipe
- If your credentials expire, you will need to re-authenticate outside of Steampipe - Steampipe currently cannot re-authenticate you.

#### aws credential file:

```ini
[account_a_with_sso]
sso_start_url = https://d-9a672b0000.awsapps.com/start
sso_region = us-east-2
sso_account_id = 000000000000
sso_role_name = SSO-ReadOnly
region = us-east-1
```

#### aws.spc:

```hcl
connection "aws_account_a_with_sso" {
  plugin  = "aws"
  profile = "account_a_with_sso"
  regions = ["us-west-2", "us-east-1",  "us-west-1", "us-east-2"]
}
```

### AssumeRole Credentials (No MFA)

If your aws credential file contains profiles that assume a role via the `source_profile` and `role_arn` options and MFA is not required, Steampipe can use the profile as-is:

#### aws credential file:

```ini
# This user must have sts:AssumeRole permission for arn:aws:iam::*:role/spc_role
[cli_user]
aws_access_key_id = AKIA4YFAKEKEYXTDS252
aws_secret_access_key = SH42YMW5p3EThisIsNotRealzTiEUwXN8BOIOF5J8m

[account_a_role_without_mfa]
role_arn = arn:aws:iam::111111111111:role/spc_role
source_profile = cli_user
external_id = xxxxx

[account_b_role_without_mfa]
role_arn = arn:aws:iam::222222222222:role/spc_role
source_profile = cli_user
external_id = yyyyy
```

#### aws.spc:

```hcl
connection "aws_account_a" {
  plugin  = "aws"
  profile = "account_a_role_without_mfa"
  regions = ["us-east-1", "us-east-2"]
}

connection "aws_account_b" {
  plugin  = "aws"
  profile = "account_b_role_without_mfa"
  regions = ["us-east-1", "us-east-2"]
}
```

### AssumeRole Credentials (With MFA)

Currently Steampipe doesn't support prompting for an MFA token at run time. To overcome this problem you will need to generate an AWS profile with temporary credentials.

One way to accomplish this is to use the `credential_process` to [generate the credentials with a script or program](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-sourcing-external.html) and cache the tokens in a new profile. There is a [sample `mfa.sh` script](https://raw.githubusercontent.com/turbot/steampipe-plugin-aws/main/scripts/mfa.sh) in the `scripts` directory of the [steampipe-plugin-aws](https://github.com/turbot/steampipe-plugin-aws) repo that you can use, and there are several open source projects that automate this process as well.

Note that Steampipe cannot prompt you for your token currently, so you must authenticate before starting Steampipe, and re-authenticate outside of Steampipe whenever your credentials expire.

#### aws credential file:

```ini
[cli_user]
aws_access_key_id = AKIA4YFAKEKEYXTDS252
aws_secret_access_key = SH42YMW5p3EThisIsNotRealzTiEUwXN8BOIOF5J8m
mfa_serial = arn:aws:iam::999999999999:mfa/my_role_mfa

[account_a_role_with_mfa]
credential_process = sh -c 'mfa.sh arn:aws:iam::111111111111:role/my_role arn:aws:iam::999999999999:mfa/my_role_mfa cli_user 2> $(tty)'

[account_b_role_with_mfa]
credential_process = sh -c 'mfa.sh arn:aws:iam::222222222222:role/my_role arn:aws:iam::999999999999:mfa/my_role_mfa cli_user 2> $(tty)'
```

#### aws.spc:

```hcl
connection "aws_account_a" {
  plugin  = "aws"
  profile = "account_a_role_with_mfa"
  regions = ["us-east-1", "us-east-2"]
}

connection "aws_account_b" {
  plugin  = "aws"
  profile = "account_b_role_with_mfa"
  regions = ["us-east-1", "us-east-2"]
}
```

### AWS-Vault Credentials

Steampipe can use profiles that use [aws-vault](https://github.com/99designs/aws-vault) via the `credential_process`. aws-vault can even be used when using AssumeRole Credentials with MFA (you must authenticate/re-authenticate outside of Steampipe whenever your credentials expire if you are using MFA).

When authenticating with temporary credentials, like using an access key pair with aws-vault, some IAM and STS APIs may be restricted. You can avoid creating a temporary session with the `--no-session` option (e.g., `aws-vault exec my_profile --no-session -- steampipe query "select name from aws_iam_user;"`). For more information, please see [aws-vault Temporary credentials limitations with STS, IAM
](https://github.com/99designs/aws-vault/blob/master/USAGE.md#temporary-credentials-limitations-with-sts-iam).

#### aws credential file:

```ini
[vault_user_account]
credential_process = /usr/local/bin/aws-vault exec -j vault_user_profile # vault_user_profile is the name of the profile in AWS_VAULT...

[account_a]
source_profile = vault_user_account
role_arn = arn:aws:iam::123456789012:role/my_role
mfa_serial = arn:aws:iam::123456789012:mfa/my_role_mfa
```

#### aws.spc:

```hcl
connection "aws_account_a" {
  plugin  = "aws"
  profile = "account_a"
  regions = ["*"]
}
```

### IAM Access Key Pair Credentials

The AWS plugin allows you set static credentials with the `access_key`, `secret_key`, and `session_token` arguments in your connection.

```hcl
connection "aws_account_a" {
  plugin     = "aws"
  secret_key = "gMCYsoGqjfThisISNotARealKeyVVhh"
  access_key = "ASIA3ODZSWFYSN2PFHPJ"
  regions    = ["us-east-1", "us-west-2"]
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
  plugin = "aws"
}
```

### Credentials from an EC2 Instance Profile

If you are running Steampipe on a AWS EC2 instance, and that instance has an [instance profile attached](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/iam-roles-for-amazon-ec2.html) then Steampipe will automatically use the associated IAM role without other credentials:

```hcl
connection "aws" {
  plugin  = "aws"
  regions = ["eu-west-1", "eu-west-2"]
}
```

## Get involved

* Open source: https://github.com/turbot/steampipe-plugin-aws
* Community: [Slack Channel](https://steampipe.io/community/join)
