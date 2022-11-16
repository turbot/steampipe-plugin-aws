# Table: aws_resource_explorer_search

AWS Resource Explorer is a resource search and discovery service. This table allows you to search for supported resource types (which can be found using the [aws_resource_explorer_supported_resource_type](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_resource_explorer_supported_resource_type) table).

Before using this table, we recommend:
- Configure Resource Explorer using [quick setup](https://docs.aws.amazon.com/resource-explorer/latest/userguide/getting-started-setting-up.html#getting-started-setting-up-quick)
- If using [advanced setup](https://docs.aws.amazon.com/resource-explorer/latest/userguide/getting-started-setting-up.html#getting-started-setting-up-advanced) instead, we recommend creating at least 1 aggregator index and a default view in that region

This table uses the aggregator index in the AWS account when searching for resources. A view ARN can also be specified in the `view_arn` column. If the account doesn't have an aggregator index and no view ARN is specified, the table will return an error.

All queries can only return the first 1,000 results due to a limitation by the API. If the resource you're looking for is not included, you can use a more refined `query` string.

Specifying `query` is not required, and if a search query is run without it, the first 1,000 results will be returned. However, if you'd like to specify `query`, please see the examples below along with [Search query syntax reference for Resource Explorer](https://docs.aws.amazon.com/resource-explorer/latest/userguide/using-search-query-syntax.html).

## Examples

### Basic info

```sql
select
  arn,
  region,
  resource_type,
  service,
  owning_account_id
from
  aws_resource_explorer_search;
```

### List non-IAM resources

```sql
select
  arn,
  region,
  resource_type,
  service,
  owning_account_id
from
  aws_resource_explorer_search
where
  query = '-service:iam';
```

### List non-IAM resources in `us-*` regions

```sql
select
  arn,
  region,
  resource_type,
  service,
  owning_account_id
from
  aws_resource_explorer_search
where
  query = '-service:iam region:us-*';
```

### List IAM user resources

```sql
select
  arn,
  region,
  resource_type,
  service,
  owning_account_id
from
  aws_resource_explorer_search
where
  query = 'resourcetype:iam:user';
```

### List resources with user created tags

```sql
select
  arn,
  region,
  resource_type,
  service,
  owning_account_id
from
  aws_resource_explorer_search
where
  query = '-tag:none';
```

### List resources with tag key `environment`

```sql
select
  arn,
  region,
  resource_type,
  service,
  owning_account_id
from
  aws_resource_explorer_search
where
  query = 'tag.key:environment';
```

### List resources with `global` scope

```sql
select
  arn,
  region,
  resource_type,
  service,
  owning_account_id
from
  aws_resource_explorer_search
where
  query = 'region:global';
```

### Search for reosurces with a specific view

```sql
select
  arn,
  region,
  resource_type,
  service,
  owning_account_id
from
  aws_resource_explorer_search
where
  view_arn = 'arn:aws:resource-explorer-2:ap-south-1:111122223333:view/view1/7c9e9845-4736-409f-9c0f-673fe7ce3e46';
```
