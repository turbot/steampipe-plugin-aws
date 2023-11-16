---
title: "Table: aws_ses_email_identity - Query AWS SES Email Identity using SQL"
description: "Allows users to query AWS SES Email Identity to retrieve information about the email identities (domains and email addresses) that you have verified with Amazon SES."
---

# Table: aws_ses_email_identity - Query AWS SES Email Identity using SQL

The `aws_ses_email_identity` table in Steampipe provides information about email identities that are verified with Amazon Simple Email Service (SES). This table allows DevOps engineers to query specific details about these identities, such as the identity type (email or domain), verification status, feedback forwarding status, and more. Users can utilize this table to gain insights on the verified email identities, their DKIM attributes, and the policies applied to them. The schema outlines various attributes of the email identity, including the identity name, verification status, bounce topic, complaint topic, delivery topic, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ses_email_identity` table, you can use the `.inspect aws_ses_email_identity` command in Steampipe.

### Key columns:

- `identity`: The identity's name. This is the primary key of the table and can be used to join with other tables.
- `dkim_enabled`: Indicates whether or not Easy DKIM signing is enabled for this identity. This can be useful when analyzing the DKIM attributes of the email identities.
- `verification_status`: The verification status of the identity (whether the identity is verified or not). This is crucial for ensuring the validity of the email identities.

## Examples

### Basic info

```sql
select
  identity,
  arn,
  region,
  akas
from
  aws_ses_email_identity;
```

### List email identities which failed verification

```sql
select
  identity,
  region,
  verification_status
from
  aws_ses_email_identity
where
  verification_status = 'Failed';
```
