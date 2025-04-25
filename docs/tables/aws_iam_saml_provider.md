---
title: "Steampipe Table: aws_iam_saml_provider - Query AWS IAM SAML Providers using SQL"
description: "Allows users to query AWS IAM SAML Providers and retrieve detailed information about each SAML provider within AWS Identity and Access Management (IAM)."
folder: "IAM"
---

# Table: aws_iam_saml_provider - Query AWS IAM SAML Providers using SQL

The AWS IAM SAML Provider is a service that allows you to manage Identity Providers (IdPs) and Single Sign-On (SSO) to AWS accounts and applications using IAM roles. It enables the establishment of trust between your AWS account and your SAML 2.0 compatible IdP. This service simplifies access management for AWS resources by allowing users to log in to multiple accounts using a single set of credentials from the IdP.

## Table Usage Guide

The `aws_iam_saml_provider` table in Steampipe provides you with information about SAML providers within AWS Identity and Access Management (IAM). This table empowers you, as a DevOps engineer, to query SAML provider-specific details, including the provider's ARN, creation date, validity period, and the SAML metadata document. You can utilize this table to gather insights on SAML providers, such as provider validity, associated metadata, and more. The schema outlines for you the various attributes of the SAML provider, including the provider ARN, creation date, and validity period.

## Examples

### Basic info
Analyze the settings to understand the creation and validity dates of your AWS IAM SAML providers across various regions and accounts. This could help in managing the lifecycle of these providers and ensuring they are valid and up-to-date.

```sql+postgres
select
  arn,
  create_date,
  valid_until,
  region,
  account_id
from
  aws_iam_saml_provider;
```

```sql+sqlite
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
Determine the areas in which certain providers have been active for a prolonged period by identifying those that have been established for more than 90 days. This could be useful for auditing purposes or to identify potential areas for system optimization or updates.

```sql+postgres
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

```sql+sqlite
select
  arn,
  create_date,
  valid_until,
  region,
  account_id
from
  aws_iam_saml_provider
where
  create_date <= date('now','-90 day')
order by
  create_date;
```

### List providers valid for less than 30 days
Determine the areas in which AWS Identity Access Management (IAM) Security Assertion Markup Language (SAML) providers have been valid for less than 30 days. This can be useful for managing and reviewing the lifespan of these providers for security and operational efficiency.

```sql+postgres
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

```sql+sqlite
select
  arn,
  create_date,
  valid_until,
  region,
  account_id
from
  aws_iam_saml_provider
where
  valid_until <= date('now','-30 day')
order by
  valid_until;
```