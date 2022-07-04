# Table: aws_iam_saml_provider

An IAM SAML 2.0 identity provider is an entity in IAM that describes an external identity provider (IdP) service that supports the SAML 2.0 (Security Assertion Markup Language 2.0) standard. You use an IAM identity provider when you want to establish trust between a SAML-compatible IdP such as Shibboleth or Active Directory Federation Services and AWS, so that users in your organization can access AWS resources. IAM SAML identity providers are used as principals in an IAM trust policy.

## Examples

### Basic info

```sql
select
  arn,
  create_date,
  valid_until,
  region,
  account_id
from
  aws_iam_saml_provider;
```

### List providers older than 90 days

```sql
select
  arn,
  create_date,
  valid_until,
  region,
  account_id
from
  aws_iam_saml_provider
where
  create_date <= (current_date - interval '90' day)
order by
  create_date;
```

### List providers valid for less than 30 days

```sql
select
  arn,
  create_date,
  valid_until,
  region,
  account_id
from
  aws_iam_saml_provider
where
  valid_until <= (current_date - interval '30' day)
order by
  valid_until;
```