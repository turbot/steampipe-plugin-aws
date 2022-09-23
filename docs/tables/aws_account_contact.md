# Table: aws_account_contact

Contains the details of the primary contact information associated with an AWS account.

This table supports the optional list key column `linked_account_id`, with the following requirements:
- The caller must be an identity in the [organization's management account](https://docs.aws.amazon.com/organizations/latest/userguide/orgs_getting-started_concepts.html#account) or a delegated administrator account.
- The specified account ID must also be a member account in the same organization.
- The organization must have [all features enabled](https://docs.aws.amazon.com/organizations/latest/userguide/orgs_manage_org_support-all-features.html).
- The organization must have [trusted access](https://docs.aws.amazon.com/accounts/latest/reference/using-orgs-trusted-access.html) enabled for the Account Management service.

**Note**: If using AWS' `ReadOnlyAccess` policy, this policy does not include the `account:GetContactInformation` permission, so you will need to add it to use this table.

## Examples

### Basic info

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

### Get contact details for an account in the organization (using credentials from the management account)

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
  aws_account_contact
where
  linked_account_id = '123456789012';
```
