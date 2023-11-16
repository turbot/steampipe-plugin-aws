---
title: "Table: aws_codeartifact_repository - Query AWS CodeArtifact Repository using SQL"
description: "Allows users to query AWS CodeArtifact Repository data, including details about the repository, its domain ownership, and associated metadata."
---

# Table: aws_codeartifact_repository - Query AWS CodeArtifact Repository using SQL

The `aws_codeartifact_repository` table in Steampipe provides information about repositories within AWS CodeArtifact. This table allows DevOps engineers to query repository specific details, including the repository's domain owner, domain name, repository name, administrator account, and associated metadata. Users can utilize this table to gather insights on repositories, such as their ownership, associated domains, and more. The schema outlines the various attributes of the CodeArtifact repository, including the ARN, repository description, domain owner, domain name, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_codeartifact_repository` table, you can use the `.inspect aws_codeartifact_repository` command in Steampipe.

**Key columns**:

- `arn`: The Amazon Resource Number (ARN) of the CodeArtifact repository. This is a unique identifier and can be used to join this table with other AWS service tables.
- `domain_owner`: The AWS account ID that owns the domain containing the repository. This can be used to filter repositories by domain ownership.
- `repository_name`: The name of the CodeArtifact repository. This can be used to join this table with other tables that reference the repository by name.

## Examples

### Basic info

```sql
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

```sql
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

```sql
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

### Get upstream package details associated with each repository

```sql
select
  arn,
  domain_name,
  domain_owner,
  u ->> 'RepositoryName' as upstream_repo_name
from
  aws_codeartifact_repository,
  jsonb_array_elements(upstreams) u;
```
