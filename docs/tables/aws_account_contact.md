---
title: "Table: aws_account_contact - Query AWS Account Contact using SQL"
description: "Allows users to query AWS Account Contact details, including email, mobile, and address information associated with an AWS account."
---

# Table: aws_account_contact - Query AWS Account Contact using SQL

The `aws_account_contact` table in Steampipe provides information about contact details associated with an AWS account. This table allows DevOps engineers to query contact-specific details, including email, mobile, and address information. Users can utilize this table to gather insights on AWS account contact details, such as verification of contact information, understanding the geographical distribution of accounts, and more. The schema outlines the various attributes of the AWS account contact, including the account ID, address, email, fax, and phone number.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_account_contact` table, you can use the `.inspect aws_account_contact` command in Steampipe.

### Key columns:

- `account_id`: This is the AWS account ID. It is a key column for joining with other tables to correlate and gather more detailed information about AWS resources.
- `email`: This is the email contact associated with the AWS account. It can be used to join with other tables that contain email information for further analysis.
- `phone_number`: This is the phone number contact associated with the AWS account. It can be used to join with other tables that contain phone number information for further analysis.

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
