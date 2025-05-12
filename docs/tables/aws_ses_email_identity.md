---
title: "Steampipe Table: aws_ses_email_identity - Query AWS SES Email Identity using SQL"
description: "Allows users to query AWS SES Email Identity to retrieve information about the email identities (domains and email addresses) that you have verified with Amazon SES."
folder: "Simple Email Service (SES)"
---

# Table: aws_ses_email_identity - Query AWS SES Email Identity using SQL

The AWS Simple Email Service (SES) Email Identity is a resource that represents the identity of an email sender. This identity can either be an email address or a domain from which the email is sent. It provides a way to verify email sending sources and to configure email sending settings, improving email deliverability.

## Table Usage Guide

The `aws_ses_email_identity` table in Steampipe provides you with information about email identities that are verified with Amazon Simple Email Service (SES). This table allows you, as a DevOps engineer, to query specific details about these identities, such as the identity type (email or domain), verification status, feedback forwarding status, and more. You can utilize this table to gain insights on the verified email identities, their DKIM attributes, and the policies applied to them. The schema outlines various attributes of the email identity for you, including the identity name, verification status, bounce topic, complaint topic, delivery topic, and associated tags.

## Examples

### Basic info
Explore which AWS Simple Email Service (SES) email identities are active in a specific region to manage and optimize your email sending activities. This can help in identifying potential issues or inefficiencies in your email distribution.

```sql+postgres
select
  identity,
  arn,
  region,
  akas
from
  aws_ses_email_identity;
```

```sql+sqlite
select
  identity,
  arn,
  region,
  akas
from
  aws_ses_email_identity;
```

### List email identities which failed verification
Explore which email identities failed the verification process in the AWS Simple Email Service. This is useful for troubleshooting and identifying potential issues with email delivery.

```sql+postgres
select
  identity,
  region,
  verification_status
from
  aws_ses_email_identity
where
  verification_status = 'Failed';
```

```sql+sqlite
select
  identity,
  region,
  verification_status
from
  aws_ses_email_identity
where
  verification_status = 'Failed';
```