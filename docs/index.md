---
organization: Turbot
category: ["public cloud"]
icon_url: "/images/plugins/turbot/aws.svg"
brand_color: "#FF9900"
display_name: "Amazon Web Services"
short_name: "aws"
description: "Steampipe plugin for AWS services and resource types."
---

# AWS

The Amazon Web Services (AWS) plugin is used to interact with the many resources supported by AWS.

### Installation
To download and install the latest aws plugin:
```bash
$ steampipe plugin install aws
Installing plugin aws...
$
```

Installing the latest aws plugin will create a connection config file (`~/.steampipe/config/aws.spc`) with a single default connection named `aws`. This connection will dynamically determine the scope and credentials using the same mechanism as the CLI, and will set regions to the default region (only). In effect, this means that by default Steampipe will execute with the same credentials and against the same region as the `aws` command would - The AWS plugin uses the standard AWS environment variables and credential files as used by the CLI.  (Of course this also  implies that the `aws` cli needs to be configured with the proper credentials before the steampipe aws plugin can be used).

Note that there is nothing special about the default connection, other than that it is created by default on plugin install - You can delete or rename this connection, or modify its configuration options (via the configuration file).


## Connection Configuration
Connection configurations are defined using HCL in one or more Steampipe config files.  Steampipe will load ALL configuration files from `~/.steampipe/config` that have a `.spc` extension. A config file may contain multiple connections.

### Scope
Each AWS connection is scoped to a single AWS account, with a single set of credentials.  You may configure multiple AWS connections if desired, with each connecting to a different account.  Each AWS connection may be configured for multiple regions.  


### Configuration Arguments

The AWS plugin allows you set static credentials with the `access_key`, `secret_key`, and `session_token` arguments.  You may select one or more regions with the `regions` argument.
An AWS connection may connect to multiple regions, however be aware that performance may be negatively affected by both the number of regions and the latency to them.


```hcl
# credentials via key pair
connection "aws_account_x" {
  plugin      = "aws" 
  secret_key  = "gMCYsoGqjfThisISNotARealKeyVVhh"
  access_key  = "ASIA3ODZSWFYSN2PFHPJ"  
  regions     = ["us-east-1" , "us-west-2"]
}
```

Alternatively, you may select a named profile from an AWS credential file with the `profile` argument.  A connect per profile is a common configuration:
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
  regions     = ["us-east-1", "us-west-2"]
}

```

If no credentials are specified, the plugin will use the AWS credentials resolver to get the current credentials in the same manner as the CLI (as used in the AWS Default Connection):

```hcl
# default
connection "aws" {
  plugin      = "aws" 
}
```


The AWS credential resolution order is:
1. Credentials specified in environment variables `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, and `AWS_SESSION_TOKEN`, `AWS_ROLE_SESSION_NAME`
2. Credentials in the credential file (`~/.aws/credentials`) for the profile specified in the `AWS_PROFILE` environment variable
3. Credentials for the Default profile from the credential file.
4. EC2 Instance Role Credentials (if running on an ec2 instance)

If `regions` is not specified, Steampipe will use a single default region using the same resolution order as the credentials:
1. The `AWS_DEFAULT_REGION` or `AWS_REGION` environment variable
2. The region specified in the active profile (`AWS_PROFILE` or default)

Steampipe will require read access in order to query your AWS resources.  Attaching the built in `ReadOnlyAccess` policy to your user or role will allow you to query all the tables in this plugin, though you can grant more granular access if you prefer.
