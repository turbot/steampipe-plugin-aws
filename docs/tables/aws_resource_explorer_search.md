---
title: "Table: aws_resource_explorer_search - Query AWS Resource Explorer using SQL"
description: "Allows users to query AWS Resource Explorer to obtain a structured view of all resources across AWS services. It provides detailed information about each resource, including the service name, resource type, resource ID, and associated tags."
---

# Table: aws_resource_explorer_search - Query AWS Resource Explorer using SQL

The `aws_resource_explorer_search` table in Steampipe provides information about resources across all AWS services. This table allows DevOps engineers to query resource-specific details, including the service name, resource type, resource ID, and associated tags. Users can utilize this table to gather insights on resources, such as identifying resources without tags, resources of a specific type, resources associated with a specific service, and more. The schema outlines the various attributes of the resources, including the resource ARN, resource type, and associated tags.

Before using this table, we recommend:
- Configure Resource Explorer using [quick setup](https://docs.aws.amazon.com/resource-explorer/latest/userguide/getting-started-setting-up.html#getting-started-setting-up-quick)
- If using [advanced setup](https://docs.aws.amazon.com/resource-explorer/latest/userguide/getting-started-setting-up.html#getting-started-setting-up-advanced) instead, we recommend creating at least 1 aggregator index and a default view in that region

This table uses the aggregator index in the AWS account when searching for resources. A view ARN can also be specified in the `view_arn` column. If the account doesn't have an aggregator index and no view ARN is specified, the table will return an error.

All queries can only return the first 1,000 results due to a limitation by the API. If the resource you're looking for is not included, you can use a more refined `query` string.

Specifying `query` is not required, and if a search query is run without it, the first 1,000 results will be returned. However, if you'd like to specify `query`, please see the examples below along with [Search query syntax reference for Resource Explorer](https://docs.aws.amazon.com/resource-explorer/latest/userguide/using-search-query-syntax.html).

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_resource_explorer_search` table, you can use the `.inspect aws_resource_explorer_search` command in Steampipe.

### Key columns:

- `arn`: The Amazon Resource Name (ARN) of the resource. This can be used to join with other tables that also contain the resource ARN.
- `service`: The name of the AWS service that the resource is associated with. This can be used to filter resources by service.
- `type`: The type of the resource. This can be used to filter resources by type.

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
