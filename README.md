![image](https://hub.steampipe.io/images/plugins/turbot/aws-social-graphic.png)

# AWS Plugin for Steampipe

Use SQL to query infrastructure including servers, networks, identity and more from AWS.

- **[Get started →](https://hub.steampipe.io/plugins/turbot/aws)**
- Documentation: [Table definitions & examples](https://hub.steampipe.io/plugins/turbot/aws/tables)
- Community: [Join #steampipe on Slack →](https://turbot.com/community/join)
- Get involved: [Issues](https://github.com/turbot/steampipe-plugin-aws/issues)

## Quick start

Install the plugin with [Steampipe](https://steampipe.io):

```shell
steampipe plugin install aws
```

Run a query:

```sql
select arn, creation_date from aws_kms_key
```

## Advanced configuration

The AWS plugin has the power to:
* Query multiple accounts
* Query multiple regions
* Use many different methods for credentials (roles, SSO, etc)

- **[Detailed configuration guide →](https://hub.steampipe.io/plugins/turbot/aws#get-started)**

## Developing

Prerequisites:

- [Steampipe](https://steampipe.io/downloads)
- [Golang](https://golang.org/doc/install)

Clone:

```sh
git clone https://github.com/turbot/steampipe-plugin-aws.git
cd steampipe-plugin-aws
```

Build, which automatically installs the new version to your `~/.steampipe/plugins` directory:

```
make
```

Configure the plugin:

```
cp config/* ~/.steampipe/config
vi ~/.steampipe/config/aws.spc
```

Try it!

```
steampipe query
> .inspect aws
```

Further reading:

- [Writing plugins](https://steampipe.io/docs/develop/writing-plugins)
- [Writing your first table](https://steampipe.io/docs/develop/writing-your-first-table)

## Contributing

Please see the [contribution guidelines](https://github.com/turbot/steampipe/blob/main/CONTRIBUTING.md) and our [code of conduct](https://github.com/turbot/steampipe/blob/main/CODE_OF_CONDUCT.md). Contributions to the plugin are subject to the [Apache 2.0 open source license](https://github.com/turbot/steampipe-plugin-aws/blob/main/LICENSE). Contributions to the plugin documentation are subject to the [CC BY-NC-ND license](https://github.com/turbot/steampipe-plugin-aws/blob/main/LICENSE).

`help wanted` issues:

- [Steampipe](https://github.com/turbot/steampipe/labels/help%20wanted)
- [AWS Plugin](https://github.com/turbot/steampipe-plugin-aws/labels/help%20wanted)
