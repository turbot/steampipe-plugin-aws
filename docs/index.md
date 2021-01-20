---
organization: Turbot
category: ["public cloud"]
icon_url: "/images/plugins/turbot/logo.svg"
brand_color: "#12efc6"
display_name: "Amazon Web Services"
short_name: "aws"
description: "Steampipe plugin for 70+ AWS services and resource types"
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
### Scope
An AWS connection is scoped to a single aws account, with a single set of credentials. Currently, a connection is limited to a single AWS region.


### Configuration

Installing the latest aws plugin will create a default connection named `aws`. This connection will dynamically determine the scope and credentials using the same mechanism as the CLI, and will set regions to the default region. In effect, this means that by default Steampipe will execute with the same credentials and against the same region as the `aws` command would - The AWS plugin uses the standard AWS environment variables and credential files as used by the CLI.  (Of course this also  implies that the `aws` cli needs to be configured with the proper credentials before the steampipe aws plugin can be used).

The AWS credential resolution order is:
1. Credentials specified in environment variables `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, and `AWS_SESSION_TOKEN`, `AWS_ROLE_SESSION_NAME`
2. Credentials in the credential file (`~/.aws/credentials`) for the profile specified in the `AWS_PROFILE` environment variable
3. Credentials for the Default profile from the credential file.
4. EC2 Instance Role Credentials (if running on an ec2 instance)

By default, steampipe will use a single region using the same resolution order as the credentials:
1. The `AWS_DEFAULT_REGION` or `AWS_REGION`  environment variable
2. The region specified in the active profile (`AWS_PROFILE` or default)

Steampipe will require read access in order to query your AWS resources.  Attaching the built in `ReadOnlyAccess` policy to your user or role will allow you to query all the tables in this plugin, though you can grant more granular access if you prefer.
