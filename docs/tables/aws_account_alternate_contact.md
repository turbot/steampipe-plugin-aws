---
title: "Table: aws_account_alternate_contact - Query AWS Account Alternate Contact using SQL"
description: "Allows users to query AWS Account Alternate Contact to fetch details about the alternate contacts associated with an AWS account."
---

# Table: aws_account_alternate_contact - Query AWS Account Alternate Contact using SQL

The `aws_account_alternate_contact` table in Steampipe provides information about the alternate contacts associated with an AWS account. This table allows DevOps engineers and AWS administrators to query alternate contact-specific details, including the contact type, name, title, email, and phone number. Users can utilize this table to gather insights on alternate contacts, such as their role in the organization, their contact information, and more. The schema outlines the various attributes of the AWS Account Alternate Contact, including the account id, contact type, name, title, email, and phone number.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_account_alternate_contact` table, you can use the `.inspect aws_account_alternate_contact` command in Steampipe.

**Key columns**:

- `account_id`: The AWS account ID associated with the alternate contact. This column is useful for joining with other tables to fetch account-specific information.
- `contact_type`: The type of the alternate contact (e.g., BILLING, OPERATIONS, SECURITY). This column is useful for filtering the alternate contacts based on their role in the organization.
- `email`: The email of the alternate contact. This column is useful for communication purposes.

## Examples

### Basic info

```sql
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

```sql
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

```sql
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

```sql
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
