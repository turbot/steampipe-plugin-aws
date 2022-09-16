# Table: aws_account_contact

We can store contact information about the primary account contact for your AWS account.

## Examples

### Basic AWS account contact details

```sql
select
  full_name,
  company_name,
  city,
  phone_number,
  postal_code,
  state_or_region,
  website_url
from
  aws_account_contact;
```
