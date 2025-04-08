---
title: "Steampipe Table: aws_codeartifact_repository - Query AWS CodeArtifact Repository using SQL"
description: "Allows users to query AWS CodeArtifact Repository data, including details about the repository, its domain ownership, and associated metadata."
folder: "CodeArtifact"
---

# Table: aws_codeartifact_repository - Query AWS CodeArtifact Repository using SQL

The AWS CodeArtifact Repository is a fully managed software artifact repository service that makes it easier for organizations to securely store, publish, and share packages used in their software development process. AWS CodeArtifact eliminates the need for you to set up, operate, and scale the infrastructure for your artifact repositories, allowing you to focus on your software development. It works with commonly used package managers and build tools, and it integrates with CI/CD pipelines to seamlessly publish packages.

## Table Usage Guide

The `aws_codeartifact_repository` table in Steampipe provides you with information about repositories within AWS CodeArtifact. This table allows you, as a DevOps engineer, to query repository specific details, including the repository's domain owner, domain name, repository name, administrator account, and associated metadata. You can utilize this table to gather insights on repositories, such as their ownership, associated domains, and more. The schema outlines the various attributes of the CodeArtifact repository for you, including the ARN, repository description, domain owner, domain name, and associated tags.

## Examples

### Basic info
Explore which AWS CodeArtifact repositories are owned by different domain owners and identify instances where specific tags and upstreams are used. This can help in gaining insights into the organization and management of your AWS resources.

```sql+postgres
select
  arn,
  domain_name,
  domain_owner,
  upstreams,
  tags
from
  aws_codeartifact_repository;
```

```sql+sqlite
select
  arn,
  domain_name,
  domain_owner,
  upstreams,
  tags
from
  aws_codeartifact_repository;
```

### List repositories with endpoints
Identify instances where repositories have specified endpoints. This could be useful in managing and organizing your AWS CodeArtifact repositories, by focusing on those repositories that have assigned endpoints.

```sql+postgres
select
  arn,
  domain_name,
  domain_owner,
  tags,
  repository_endpoint
from
  aws_codeartifact_repository
where
  repository_endpoint is not null;
```

```sql+sqlite
select
  arn,
  domain_name,
  domain_owner,
  tags,
  repository_endpoint
from
  aws_codeartifact_repository
where
  repository_endpoint is not null;
```

### List repository policy statements that grant external access
This example is used to identify any repository policy statements in the AWS CodeArtifact service that may be granting access to external entities. This is useful for auditing security and ensuring that no unauthorized access is being permitted.

```sql+postgres
select
  arn,
  p as principal,
  a as action,
  s ->> 'Effect' as effect
from
  aws_codeartifact_repository,
  jsonb_array_elements(policy_std -> 'Statement') as s,
  jsonb_array_elements_text(s -> 'Principal' -> 'AWS') as p,
  string_to_array(p, ':') as pa,
  jsonb_array_elements_text(s -> 'Action') as a
where
  s ->> 'Effect' = 'Allow'
  and (
    pa [5] != account_id
    or p = '*'
  );
```

```sql+sqlite
Error: The corresponding SQLite query is unavailable.
```

### Get upstream package details associated with each repository
Analyze the settings to understand the association between each repository and its corresponding upstream package details in the AWS CodeArtifact service. This can aid in managing dependencies and ensuring the correct version of a package is being used.

```sql+postgres
select
  arn,
  domain_name,
  domain_owner,
  u ->> 'RepositoryName' as upstream_repo_name
from
  aws_codeartifact_repository,
  jsonb_array_elements(upstreams) u;
```

```sql+sqlite
select
  arn,
  domain_name,
  domain_owner,
  json_extract(u.value, '$.RepositoryName') as upstream_repo_name
from
  aws_codeartifact_repository,
  json_each(upstreams) u;
```
