# Table: aws_securityhub_product

Security hub Products provides information about product integrations in Security Hub.

## Examples

### Basic info
```sql
select
  name,
  product_arn,
  company_name,
  description
from
  aws_securityhub_product;
```


### List the products provided by AWS
```sql
select
  name,
  company_name,
  description
from
  aws_securityhub_product
where company_name = 'AWS';
```


### List of products that send findings to security hub
```sql
select
  name,
  product_arn,
  company_name
from
  aws_securityhub_product,
  jsonb_array_elements_text(integration_types) as i
where i = 'SEND_FINDINGS_TO_SECURITY_HUB';
```