# Table: aws_resource_explorer_search

AWS Resource Explorer is a resource search and discovery service. This table allows to search supported resource types (which can be found using [aws_resource_explorer_supported_resource_type](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_resource_explorer_supported_resource_type) table).

**Important notes:**

- You **_must_** specify an AWS `region` in the where clause to search for resources. The details of the region having AWS Explorer Service Indexes for search can be known through the [aws_resource_explorer_index](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_resource_explorer_index) table. Only regions when Resource Explorer is turned on, and the index is available in the region, then it can be used for searches.

- This table supports other optional quals. Queries with optional quals are optimised to use specific View and Query for searches. Optional quals are supported for the following columns:
  - `query`: A string that includes keywords and filters that specify the resources that you want to include in the results. For further details refer [query string syntax](https://docs.aws.amazon.com/resource-explorer/latest/userguide/using-search-query-syntax.html#query-syntax).
  - `view_arn`: Specifies the ARN of the view to use for the query. If you don't specify a value for this filed, then the operation automatically uses the default view for the AWS Region in which you called this operation. If the Region either doesn't have a default view or if you don't have permission to use the default view, then the operation fails with a 401 Unauthorized exception.

**For more details refer**
- [Using AWS Resource Explorer to search for resources](https://docs.aws.amazon.com/resource-explorer/latest/userguide/using-search.html)
- [Search API Document](https://docs.aws.amazon.com/resource-explorer/latest/apireference/API_Search.html)

## Examples

### Use default view of the region to search for resources
```sql
select
  arn,
  

```