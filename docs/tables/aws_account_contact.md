# Table: aws_account_contact

Contains the details of the primary contact information associated with an AWS account.

**Important notes:**

This table supports the optional qual `contact_account_id`.

To use `contact_account_id` parameter, the caller must be an identity in the [organization's management account](https://docs.aws.amazon.com/organizations/latest/userguide/orgs_getting-started_concepts.html#account) or a delegated administrator account. The specified account ID must also be a member account in the same organization. The organization must have [all features enabled](https://docs.aws.amazon.com/organizations/latest/userguide/orgs_manage_org_support-all-features.html), and the organization must have [trusted access](https://docs.aws.amazon.com/organizations/latest/userguide/using-orgs-trusted-access.html) enabled for the Account Management service, and optionally a [delegated admin](https://docs.aws.amazon.com/organizations/latest/userguide/using-orgs-delegated-admin.html) account assigned.

Also, the management account can't specify its own account ID, so any queries using credentials from the management account must not specify the `contact_account_id` column.

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
  contact_account_id = '123456789012';
```
