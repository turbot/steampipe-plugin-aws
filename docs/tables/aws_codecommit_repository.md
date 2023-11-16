---
title: "Table: aws_codecommit_repository - Query AWS CodeCommit Repository using SQL"
description: "Allows users to query AWS CodeCommit repositories and retrieve data such as repository name, ARN, description, clone URL, last modified date, and other related details."
---

# Table: aws_codecommit_repository - Query AWS CodeCommit Repository using SQL

The `aws_codecommit_repository` table in Steampipe provides information about repositories within AWS CodeCommit. This table allows DevOps engineers to query repository-specific details, including repository name, ARN, description, clone URL, last modified date, and other related details. Users can utilize this table to gather insights on repositories, such as repositories with specific ARNs, the last modified date of repositories, verification of clone URLs, and more. The schema outlines the various attributes of the CodeCommit repository, including the repository name, ARN, clone URL, and associated metadata.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_codecommit_repository` table, you can use the `.inspect aws_codecommit_repository` command in Steampipe.

### Key columns:

- `repository_name`: The name of the repository. This can be used to join with other tables that need repository-specific information.
- `arn`: The Amazon Resource Name (ARN) of the repository. It is unique across all AWS repositories and can be used to join with any other AWS-specific table.
- `clone_url_http`: The URL to clone the repository over HTTPS. This can be useful to link to the actual codebase from the data retrieved.

## Examples

### Basic info

```sql
select
  repository_name,
  repository_id,
  arn,
  creation_date,
  region
from
  aws_codecommit_repository;
```
