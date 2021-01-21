<p align="center">
    <h1 align="center">AWS Plugin for Steampipe</h1>
</p>

<p align="center">
  <a aria-label="Steampipe logo" href="https://steampipe.io">
    <img src="https://steampipe.io/images/steampipe_logo_wordmark_padding.svg" height="28">
  </a>
  <a aria-label="Plugin version" href="https://hub.steampipe.io/plugins/turbot/aws">
    <img alt="" src="https://img.shields.io/static/v1?label=turbot/aws&message=v0.1.0&style=for-the-badge&labelColor=777777&color=F3F1F0">
  </a>
  &nbsp;
  <a aria-label="License" href="LICENSE">
    <img alt="" src="https://img.shields.io/static/v1?label=license&message=MPL-2.0&style=for-the-badge&labelColor=777777&color=F3F1F0">
  </a>
</p>

## Query AWS with SQL

Use SQL to query IAM users, EC2 instances and more from your AWS account. For example:

```sql
select
  name,
  user_id,
  path,
  create_date,
  password_last_used
from
  aws_iam_user;
```

Learn about [Steampipe](https://steampipe.io/).

## Get started

**[Table documentation and examples &rarr;](https://hub.steampipe.io/plugins/turbot/aws)**

Install the plugin:

```shell
steampipe plugin install aws
```

## Get involved

### Community

The Steampipe community can be found on [GitHub Discussions](https://github.com/turbot/steampipe/discussions), where you can ask questions, voice ideas, and share your projects.

Our [Code of Conduct](https://github.com/turbot/steampipe/CODE_OF_CONDUCT.md) applies to all Steampipe community channels.

### Contributing

Please see [CONTRIBUTING.md](https://github.com/turbot/steampipe/CONTRIBUTING.md).
