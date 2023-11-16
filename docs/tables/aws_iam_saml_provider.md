---
title: "Table: aws_iam_saml_provider - Query AWS IAM SAML Providers using SQL"
description: "Allows users to query AWS IAM SAML Providers and retrieve detailed information about each SAML provider within AWS Identity and Access Management (IAM)."
---

# Table: aws_iam_saml_provider - Query AWS IAM SAML Providers using SQL

The `aws_iam_saml_provider` table in Steampipe provides information about SAML providers within AWS Identity and Access Management (IAM). This table allows DevOps engineers to query SAML provider-specific details, including the provider's ARN, creation date, validity period, and the SAML metadata document. Users can utilize this table to gather insights on SAML providers, such as provider validity, associated metadata, and more. The schema outlines the various attributes of the SAML provider, including the provider ARN, creation date, and validity period.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_iam_saml_provider` table, you can use the `.inspect aws_iam_saml_provider` command in Steampipe.

Key columns:

- `arn`: The Amazon Resource Name (ARN) of the SAML provider. This can be used to join with other tables where a SAML provider ARN is referenced.
- `create_date`: The date and time when the SAML provider was created. This can be useful for tracking the age of the SAML provider.
- `valid_until`: The expiration date and time for the SAML provider. This can be used for monitoring and alerting on upcoming SAML provider expirations.

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