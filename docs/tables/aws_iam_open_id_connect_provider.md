# Table: aws_iam_open_id_connect_provider

IAM OIDC identity providers are entities in IAM that describe an external identity provider (IdP) service that supports the OpenID Connect (OIDC) standard, such as Google or Salesforce. You use an IAM OIDC identity provider when you want to establish trust between an OIDC-compatible IdP and your AWS account. This is useful when creating a mobile app or web application that requires access to AWS resources, but you don't want to create custom sign-in code or manage your own user identities.

## Examples

### Basic info

```sql
select
  arn,
  create_date,
  client_id_list,
  thumbprint_list,
  url,
  account_id
from
  aws_iam_open_id_connect_provider;
```

### List providers older than 90 days

```sql
select
  arn,
  create_date,
  client_id_list,
  thumbprint_list,
  url,
  account_id
from
  aws_iam_open_id_connect_provider
where
  create_date <= (current_date - interval '90' day)
order by
  create_date;
```

### List providers with specific tags

```sql
select
  arn,
  create_date,
  client_id_list,
  thumbprint_list,
  tags,
  url,
  account_id
from
  aws_iam_open_id_connect_provider
where
  tags ->> 'Environment' = 'Production';
```

### List AWS OpenID Providers without required thumbprint for audience 'sts.amazonaws.com'

```sql
select
  arn,
  create_date,
  client_id_list,
  thumbprint_list,
  tags,
  url,
  account_id
from
  aws_iam_open_id_connect_provider
where
  client_id_list @> '["sts.amazonaws.com"]'::jsonb
  and NOT thumbprint_list @> '["1c58a3a8518e8759bf075b76b750d4f2df264fcd", "6938fd4d98bab03faadb97b34396831e3780aea1"]'::jsonb
```
