# Table: aws_service

An AWS Service is a named product that is available to users of Amazon Web Services. 
Each Service is identified by a short name (lower case, typically an acronym or abbreviated, no spaces). 
Where available the table includes the longer human readable english name for the product and a URL that represents the marketing site for the product.
The data is sourced from an Amazon hosted SSM parameter: /aws/service/global-infrastructure/services/.

## Examples

### AWS service info

```sql
select
  name,
  long_name,
  marketing_url
from
  aws_service;
```


### AWS service info for SQS

```sql
select
  name,
  long_name,
  marketing_url
from
  aws_service
where
  name = 'sqs';
```
