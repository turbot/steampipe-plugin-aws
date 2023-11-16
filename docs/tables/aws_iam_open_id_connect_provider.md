---
title: "Table: aws_iam_open_id_connect_provider - Query AWS IAM OpenID Connect Providers using SQL"
description: "Allows users to query AWS IAM OpenID Connect Providers and retrieve details about the OpenID Connect (OIDC) identity providers in their AWS account."
---

# Table: aws_iam_open_id_connect_provider - Query AWS IAM OpenID Connect Providers using SQL

The `aws_iam_open_id_connect_provider` table in Steampipe provides information about OpenID Connect (OIDC) identity providers within AWS Identity and Access Management (IAM). This table allows DevOps engineers to query provider-specific details, including ARNs, URLs, client IDs, thumbprint lists, and creation times. Users can utilize this table to gather insights on OIDC identity providers, such as their associated client IDs, verification of thumbprint lists, and more. The schema outlines the various attributes of the OIDC identity provider, including the provider ARN, creation date, client ID list, thumbprint list, and URL.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_iam_open_id_connect_provider` table, you can use the `.inspect aws_iam_open_id_connect_provider` command in Steampipe.

### Key columns:

- `arn`: The Amazon Resource Name (ARN) of the IAM OIDC provider. This can be used to join with other tables that contain IAM OIDC provider ARNs.
- `url`: The URL of the IAM OIDC provider. This is useful for joining with other tables that contain IAM OIDC provider URLs.
- `client_id_list`: The list of client IDs (also known as audiences) for the IAM OIDC provider. This can be used to join with other tables that contain IAM OIDC provider client IDs.

## Examples

### Basic info

```sql
select
  arn,
  create_date,
  client_id_list,
  thumbprint_list,
  url,
  account_id
from
  aws_iam_open_id_connect_provider;
```

### List providers older than 90 days

```sql
select
  arn,
  create_date,
  client_id_list,
  thumbprint_list,
  url,
  account_id
from
  aws_iam_open_id_connect_provider
where
  create_date <= (current_date - interval '90' day)
order by
  create_date;
```

### List providers with specific tags

```sql
select
  arn,
  create_date,
  client_id_list,
  thumbprint_list,
  tags,
  url,
  account_id
from
  aws_iam_open_id_connect_provider
where
  tags ->> 'Environment' = 'Production';
```

### List AWS OpenID Providers without the required thumbprint for audience 'sts.amazonaws.com'

```sql
select
  arn,
  create_date,
  client_id_list,
  thumbprint_list,
  tags,
  url,
  account_id
from
  aws_iam_open_id_connect_provider
where
  client_id_list @> '["sts.amazonaws.com"]'::jsonb
  and not thumbprint_list @> '["1c58a3a8518e8759bf075b76b750d4f2df264fcd", "6938fd4d98bab03faadb97b34396831e3780aea1"]'::jsonb
