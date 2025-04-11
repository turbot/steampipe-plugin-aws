---
title: "Steampipe Table: aws_account_alternate_contact - Query AWS Account Alternate Contact using SQL"
description: "Allows users to query AWS Account Alternate Contact to fetch details about the alternate contacts associated with an AWS account."
folder: "Account"
---

# Table: aws_account_alternate_contact - Query AWS Account Alternate Contact using SQL

The AWS Account Alternate Contact is a feature that allows you to designate additional contacts for your AWS account. These contacts can be specified for different types of communication such as billing, operations, or security, providing an extra layer of management and oversight. It's an effective way to ensure important account-related information is received by the right people in your organization.

## Table Usage Guide

The `aws_account_alternate_contact` table in Steampipe provides you with information about the alternate contacts associated with your AWS account. You can use this table to query alternate contact-specific details, including the contact type, name, title, email, and phone number if you're a DevOps engineer or an AWS administrator. You can use this table to gather insights on alternate contacts, such as their role in the organization, their contact information, and more. The schema outlines the various attributes of your AWS Account Alternate Contact, including the account id, contact type, name, title, email, and phone number.

**Important Notes**
This table supports the optional list key column `linked_account_id`, which comes with the following requirements:
- You must be an identity in the [organization's management account](https://docs.aws.amazon.com/organizations/latest/userguide/orgs_getting-started_concepts.html#account) or a delegated administrator account.
- The specified account ID must also be a member account in the same organization as yours.
- Your organization must have [all features enabled](https://docs.aws.amazon.com/organizations/latest/userguide/orgs_manage_org_support-all-features.html).
- Your organization must have [trusted access](https://docs.aws.amazon.com/accounts/latest/reference/using-orgs-trusted-access.html) enabled for the Account Management service.

## Examples

### Basic info
Discover the segments that are linked to specific AWS accounts and the type of contact associated with them. This can be useful in understanding the communication channels and roles involved in managing these accounts.

```sql+postgres
select
  name,
  linked_account_id,
  contact_type,
  email_address,
  phone_number,
  contact_title
from
  aws_account_alternate_contact;
```

```sql+sqlite
select
  name,
  linked_account_id,
  contact_type,
  email_address,
  phone_number,
  contact_title
from
  aws_account_alternate_contact;
```

### Get billing alternate contact details
Discover the segments that contain alternate contact details specifically for billing purposes. This can be useful in instances where you need to directly reach out to the responsible parties for billing inquiries or issues.

```sql+postgres
select
  name,
  linked_account_id,
  contact_type,
  email_address,
  phone_number,
  contact_title
from
  aws_account_alternate_contact
where
  contact_type = 'BILLING';
```

```sql+sqlite
select
  name,
  linked_account_id,
  contact_type,
  email_address,
  phone_number,
  contact_title
from
  aws_account_alternate_contact
where
  contact_type = 'BILLING';
```

### Get alternate contact details for an account in the organization (using credentials from the management account)
Discover the alternate contact details for a specific account within your organization using information from the management account. This is useful for ensuring communication channels are updated and accurate.

```sql+postgres
select
  name,
  linked_account_id,
  contact_type,
  email_address,
  phone_number,
  contact_title
from
  aws_account_alternate_contact
where
  linked_account_id = '123456789012';
```

```sql+sqlite
select
  name,
  linked_account_id,
  contact_type,
  email_address,
  phone_number,
  contact_title
from
  aws_account_alternate_contact
where
  linked_account_id = '123456789012';
```

### Get operations alternate contact details for an account in the organization (using credentials from the management account)
This query is useful for identifying the alternate contact details related to security for a specific account within an organization. It allows for efficient monitoring and communication in case of any security-related issues or concerns.

```sql+postgres
select
  name,
  linked_account_id,
  contact_type,
  email_address,
  phone_number,
  contact_title
from
  aws_account_alternate_contact
where
  linked_account_id = '123456789012'
  and contact_type = 'SECURITY';
```

```sql+sqlite
select
  name,
  linked_account_id,
  contact_type,
  email_address,
  phone_number,
  contact_title
from
  aws_account_alternate_contact
where
  linked_account_id = '123456789012'
  and contact_type = 'SECURITY';
```