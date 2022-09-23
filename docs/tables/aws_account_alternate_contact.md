# Table: aws_account_alternate_contact

A structure that contains the details of the billing, operations, and security alternate contacts associated with an AWS account.

This table supports the optional list key column `contact_account_id`, with the following requirements:
- The caller must be an identity in the [organization's management account](https://docs.aws.amazon.com/organizations/latest/userguide/orgs_getting-started_concepts.html#account) or a delegated administrator account.
- The specified account ID must also be a member account in the same organization.
- The organization must have [all features enabled](https://docs.aws.amazon.com/organizations/latest/userguide/orgs_manage_org_support-all-features.html).
- The organization must have [trusted access](https://docs.aws.amazon.com/accounts/latest/reference/using-orgs-trusted-access.html) enabled for the Account Management service.

## Examples

### Basic info

```sql
select
  name,
  contact_account_id,
  alternate_contact_type,
  email_address,
  phone_number,
  title
from
  aws_account_alternate_contact;
```

### Get billing alternate contact details

```sql
select
  name,
  contact_account_id,
  alternate_contact_type,
  email_address,
  phone_number,
  title
from
  aws_account_alternate_contact
where
  alternate_contact_type = 'BILLING';
```

### Get alternate contact details for an account in the organization (using credentials from the management account)

```sql
select
  name,
  contact_account_id,
  alternate_contact_type,
  email_address,
  phone_number,
  title
from
  aws_account_alternate_contact
where
  contact_account_id = '123456789012';
```
