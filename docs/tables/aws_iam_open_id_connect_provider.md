---
title: "Steampipe Table: aws_iam_open_id_connect_provider - Query AWS IAM OpenID Connect Providers using SQL"
description: "Allows users to query AWS IAM OpenID Connect Providers and retrieve details about the OpenID Connect (OIDC) identity providers in their AWS account."
folder: "IAM"
---

# Table: aws_iam_open_id_connect_provider - Query AWS IAM OpenID Connect Providers using SQL

The AWS IAM OpenID Connect Provider is a service that allows you to integrate web identity federation with your mobile app, web app, or other AWS resources. This enables users to sign in using a well-known social identity provider such as Login with Amazon, Facebook, Google, or any OpenID Connect (OIDC)-compatible identity provider. It simplifies the sign-in process for your applications and includes built-in security token service (STS), eliminating the need to write any server-side code.

## Table Usage Guide

The `aws_iam_open_id_connect_provider` table in Steampipe provides you with information about OpenID Connect (OIDC) identity providers within AWS Identity and Access Management (IAM). This table allows you, as a DevOps engineer, to query provider-specific details, including ARNs, URLs, client IDs, thumbprint lists, and creation times. You can utilize this table to gather insights on OIDC identity providers, such as their associated client IDs, verification of thumbprint lists, and more. The schema outlines the various attributes of the OIDC identity provider, including the provider ARN, creation date, client ID list, thumbprint list, and URL for you.

## Examples

### Basic info
Explore which AWS IAM OpenID Connect providers are active, when they were created, and their associated client IDs and URLs. This information can help improve security by identifying potentially unauthorized or outdated connections.

```sql+postgres
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

```sql+sqlite
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
Identify instances where the providers have been created more than 90 days ago. This can be useful for auditing purposes, allowing you to track and manage older providers within your AWS IAM Open ID Connect.

```sql+postgres
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

```sql+sqlite
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
  create_date <= date('now','-90 day')
order by
  create_date;
```

### List providers with specific tags
Determine the areas in which specific tags are associated with providers, particularly in a production environment. This can be useful for managing and categorizing resources within your AWS IAM OpenID Connect providers.

```sql+postgres
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

```sql+sqlite
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
  json_extract(tags, '$.Environment') = 'Production';
```

### List AWS OpenID Providers without the required thumbprint for audience 'sts.amazonaws.com'
Determine the areas in which AWS OpenID Providers lack the necessary thumbprint for the audience 'sts.amazonaws.com'. This query is useful for identifying potential security gaps in your AWS OpenID configuration.

```sql+postgres
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
```

```sql+sqlite
Error: The corresponding SQLite query is unavailable.
```