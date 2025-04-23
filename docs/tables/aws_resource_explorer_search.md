---
title: "Steampipe Table: aws_resource_explorer_search - Query AWS Resource Explorer using SQL"
description: "Allows users to query AWS Resource Explorer to obtain a structured view of all resources across AWS services. It provides detailed information about each resource, including the service name, resource type, resource ID, and associated tags."
folder: "Resource Explorer"
---

# Table: aws_resource_explorer_search - Query AWS Resource Explorer using SQL

The AWS Resource Explorer allows you to inspect and navigate your AWS resources using a visual interface. It provides a unified view of all your AWS resources, enabling you to see their relationships and dependencies. With the AWS Resource Explorer, you can search, filter, and manage your resources across multiple AWS services.

## Table Usage Guide

The `aws_resource_explorer_search` table in Steampipe provides you with information about resources across all AWS services. This table allows you, as a DevOps engineer, to query resource-specific details, including the service name, resource type, resource ID, and associated tags. You can utilize this table to gather insights on resources, such as identifying resources without tags, resources of a specific type, resources associated with a specific service, and more. The schema outlines the various attributes of the resources, including the resource ARN, resource type, and associated tags.

**Important Notes**
Before you use this table, it's recommended that you:
- Configure Resource Explorer using [quick setup](https://docs.aws.amazon.com/resource-explorer/latest/userguide/getting-started-setting-up.html#getting-started-setting-up-quick)
- If you're using [advanced setup](https://docs.aws.amazon.com/resource-explorer/latest/userguide/getting-started-setting-up.html#getting-started-setting-up-advanced) instead, it's recommended that you create at least 1 aggregator index and a default view in that region

- This table uses the aggregator index in your AWS account when searching for resources. A view ARN can also be specified in the `view_arn` column. If your account doesn't have an aggregator index and no view ARN is specified, the table will return an error.

- All queries can only return the first 1,000 results due to a limitation by the API. If the resource you're looking for is not included, you can use a more refined `query` string.

- Specifying `query` is not required, and if a search query is run without it, the first 1,000 results will be returned. However, if you'd like to specify `query`, please see the examples below along with [Search query syntax reference for Resource Explorer](https://docs.aws.amazon.com/resource-explorer/latest/userguide/using-search-query-syntax.html).

## Examples

### Basic info
Explore which resources are being utilized in your AWS environment, including where they are located and who owns them. This information can help you manage resources effectively and identify areas for potential optimization or security improvements.

```sql+postgres
select
  arn,
  region,
  resource_type,
  service,
  owning_account_id
from
  aws_resource_explorer_search;
```

```sql+sqlite
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
This query allows you to identify all resources within your AWS environment that are not associated with the IAM service. This can be particularly useful for understanding the overall distribution and utilization of your resources across different AWS services.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which non-IAM resources are located within US regions. This allows for a comprehensive understanding of resource distribution and ownership across specific geographical zones.

```sql+postgres
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

```sql+sqlite
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
Explore which IAM user resources are present within a specific AWS account and region. This can be useful to determine the areas in which these resources are distributed, aiding in resource management and security auditing.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which resources are without user-created tags to assess the elements within AWS that may require additional organization or categorization.

```sql+postgres
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

```sql+sqlite
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
Discover the segments that are tagged with the key 'environment' across various resources in your AWS environment. This allows for better resource management and environment-specific optimizations.

```sql+postgres
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

```sql+sqlite
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
Explore which resources have a global scope, helping to understand the distribution and type of resources across different regions. This can be useful to identify potential areas for cost optimization or to assess security configurations across a global infrastructure.

```sql+postgres
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

```sql+sqlite
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
Explore which resources are associated with a specific view in the AWS Resource Explorer. This is useful to manage and keep track of resources tied to a particular view, aiding in efficient resource management.

```sql+postgres
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

```sql+sqlite
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