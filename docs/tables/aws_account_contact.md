---
title: "Steampipe Table: aws_account_contact - Query AWS Account Contact using SQL"
description: "Allows users to query AWS Account Contact details, including email, mobile, and address information associated with an AWS account."
folder: "Account"
---

# Table: aws_account_contact - Query AWS Account Contact using SQL

The AWS Account Contact is a resource that stores contact information associated with an AWS account. This information can include the account holder's name, email address, and phone number. It is essential for communication purposes, especially for receiving important notifications and alerts related to the AWS services and resources.

## Table Usage Guide

The `aws_account_contact` table in Steampipe provides you with information about contact details associated with an AWS account. This table allows you, as a DevOps engineer, to query contact-specific details, including email, mobile, and address information. You can utilize this table to gather insights on AWS account contact details, such as verification of contact information, understanding the geographical distribution of accounts, and more. The schema outlines the various attributes of the AWS account contact for you, including the account ID, address, email, fax, and phone number.

**Important Notes**
This table supports the optional list key column `linked_account_id`, with the following requirements:
- The caller must be an identity in the [organization's management account](https://docs.aws.amazon.com/organizations/latest/userguide/orgs_getting-started_concepts.html#account) or a delegated administrator account.
- The specified account ID must also be a member account in the same organization.
- The organization must have [all features enabled](https://docs.aws.amazon.com/organizations/latest/userguide/orgs_manage_org_support-all-features.html).
- The organization must have [trusted access](https://docs.aws.amazon.com/accounts/latest/reference/using-orgs-trusted-access.html) enabled for the Account Management service.
- If using AWS' `ReadOnlyAccess` policy, this policy does not include the `account:GetContactInformation` permission, so you will need to add it to use this table.

## Examples

### Basic info
This query allows you to explore the basic contact information linked to your AWS account. The practical application of this query is to quickly identify and review your account details, ensuring they're accurate and up-to-date.

```sql+postgres
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

```sql+sqlite
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
Gain insights into the contact information associated with a specific account in your organization. This can be particularly useful for administrators who need to communicate with account holders or verify account details.

```sql+postgres
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

```sql+sqlite
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