# Table: aws_resource_explorer_search

AWS Resource Explorer is a resource search and discovery service. This table allows to search supported resource types (which can be found using [aws_resource_explorer_supported_resource_type](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_resource_explorer_supported_resource_type) table).

**Important notes:**

- You **_must_** specify an AWS `region` or `view_arn` in the where clause to search for resources. The details of the region having AWS Explorer Service Indexes for search can be known through the [aws_resource_explorer_index](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_resource_explorer_index) table. Only regions when Resource Explorer is turned on, and the index is available in the region, then it can be used for searches.

- This table supports other optional quals. Queries with optional quals are optimised to use specific View and Query for searches. Optional quals are supported for the following columns:
  - `query`: A string that includes keywords and filters that specify the resources that you want to include in the results. For further details refer [query string syntax](https://docs.aws.amazon.com/resource-explorer/latest/userguide/using-search-query-syntax.html#query-syntax).
  - `view_arn`: Specifies the ARN of the view to use for the query. If you don't specify a value for this filed, then the operation automatically uses the default view for the AWS Region in which you called this operation. If the Region either doesn't have a default view or if you don't have permission to use the default view, then the operation fails with a 401 Unauthorized exception.

- A search can return only the first 1,000 results.</br>
  To see resources beyond the 1,000 returned by an empty query string, you must use queries to restrict matching results to those you want to see and limit the number of matches to less than 1,000.

**For more details refer below:**

- [Using AWS Resource Explorer to search for resources](https://docs.aws.amazon.com/resource-explorer/latest/userguide/using-search.html)
- [Search API Document](https://docs.aws.amazon.com/resource-explorer/latest/apireference/API_Search.html)

## Examples

**NOTE** For below examples the region `ap-south-1` has the `AGGREGATOR` Index with default view to list all resources.

### Basic info

```sql
select
  arn,
  resource_region,
  resource_type,
  service,
  owning_account_id
from
  aws_resource_explorer_search
where
  region = 'ap-south-1';
```

### List resources other than the IAM service resources

```sql
select
  arn,
  resource_region,
  resource_type,
  service,
  owning_account_id
from
  aws_resource_explorer_search
where
  region = 'ap-south-1'
  and query = '-service:iam';
```

### List resources other than IAM service in `us-*` regions

```sql
select
  arn,
  resource_region,
  resource_type,
  service,
  owning_account_id
from
  aws_resource_explorer_search
where
  region = 'ap-south-1'
  and query = '-service:iam region:us-*';
```

### List resources of a specific type using `resourcetype` in query

```sql
select
  arn,
  resource_region,
  resource_type,
  service,
  owning_account_id
from
  aws_resource_explorer_search
where
  region = 'ap-south-1'
  and query = 'resourcetype:iam:user';
```

### List resources with user created tags

```sql
select
  arn,
  resource_region,
  resource_type,
  service,
  owning_account_id
from
  aws_resource_explorer_search
where
  region = 'ap-south-1'
  and query = '-tag:none';
```

### List resources with tag key `environment`

```sql
select
  arn,
  resource_region,
  resource_type,
  service,
  owning_account_id
from
  aws_resource_explorer_search
where
  region = 'ap-south-1'
  and query = 'tag.key:environment';
```

### List resources with `global` scope

```sql
select
  arn,
  resource_region,
  resource_type,
  service,
  owning_account_id
from
  aws_resource_explorer_search
where
  region = 'ap-south-1'
  and query = 'region:global';
```

### List resources from with a specific view

```sql
select
  arn,
  resource_region,
  resource_type,
  service,
  owning_account_id
from
  aws_resource_explorer_search
where
  region = 'us-east-2'
  and view_arn = 'arn:aws:resource-explorer-2:us-east-2:123456789012:view/My-Resources-View/1a2b3c4d-5d6e-7f8a-9b0c-abcd11111111';
```
