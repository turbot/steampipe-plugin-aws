---
organization: Turbot
category: ["public cloud"]
icon_url: "/images/plugins/turbot/aws.svg"
brand_color: "#FF9900"
display_name: "Amazon Web Services"
short_name: "aws"
description: "Steampipe plugin for AWS services and resource types."
og_description: Query AWS with SQL! Open source CLI. No DB required. 
og_image: "/images/plugins/turbot/aws-social-graphic.png"
---

# AWS + Steampipe

[Steampipe](https://steampipe.io) is an open source CLI to instantly query cloud APIs using SQL.

[AWS](https://aws.amazon.com/) provides on-demand cloud computing platforms and APIs to authenticated customers on a metered pay-as-you-go basis. 

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

- **[Table definitions & examples →](/plugins/aws/tables)**

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
| Resolution |  1. Credentials specified in environment variables e.g. `AWS_ACCESS_KEY_ID`.<br />2. Credentials in the credential file (`~/.aws/credentials`) for the profile specified in the `AWS_PROFILE` environment variable.<br />3. Credentials for the Default profile from the credential file.<br />4. EC2 Instance Role Credentials (if running on an ec2 instance) |
| Region Resolution | 1. The `AWS_DEFAULT_REGION` or `AWS_REGION` environment variable<br />2. The region specified in the active profile (`AWS_PROFILE` or `default`). |

### Configuration

Installing the latest aws plugin will create a config file (`~/.steampipe/config/aws.spc`) with a single connection named `aws`:

```hcl
connection "aws" {
  plugin      = "aws" 
  profile     = "default"
  regions     = ["us-east-1", "us-west-2"]
}
```

## Get involved

* Open source: https://github.com/turbot/steampipe-plugin-aws
* Community: [Discussion forums](https://github.com/turbot/steampipe/discussions)

## Advanced configuration options

If you have a `default` profile setup using the AWS CLI Steampipe just works with that connection.

For users with multiple accounts and more complex authentication use cases, here are some examples of advanced configuration options:

The AWS plugin allows you set static credentials with the `access_key`, `secret_key`, and `session_token` arguments in any connection profile.  You may also specify one or more regions with the `regions` argument. An AWS connection may connect to multiple regions, however be aware that performance may be negatively affected by both the number of regions and the latency to them.


### Credentials via key pair
```hcl
connection "aws_account_x" {
  plugin      = "aws" 
  secret_key  = "gMCYsoGqjfThisISNotARealKeyVVhh"
  access_key  = "ASIA3ODZSWFYSN2PFHPJ"  
  regions     = ["us-east-1" , "us-west-2"]
}
```
### Credentials via AWS config profiles

Named profile from an AWS credential file with the `profile` argument.  A connect per profile is a common configuration:

```hcl
# credentials via profile
connection "aws_account_y" {
  plugin      = "aws" 
  profile     = "profile_y"
  regions     = ["us-east-1", "us-west-2"]
}

# credentials via profile
connection "aws_account_z" {
  plugin      = "aws" 
  profile     = "profile_z"
  regions     = ["ap-southeast-1", "ap-southeast-2"]
}
```

### Credentials from environment variables

```sh
export AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE
export AWS_SECRET_ACCESS_KEY=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
export AWS_DEFAULT_REGION=eu-west-1
export AWS_SESSION_TOKEN=AQoDYXdzEJr...
export AWS_ROLE_SESSION_NAME=steampipe@myaccount
```

### Credentials from an EC2 instance role

If you are running Steampipe on a AWS EC2 instance, and that instance has an [instance profile attached](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/iam-roles-for-amazon-ec2.html) then Steampipe will automatically use the associated IAM role without other credentials:

```hcl
connection "aws" {
  plugin      = "aws" 
  regions     = ["eu-west-1", "eu-west-2"]
}
```