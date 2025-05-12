---
title: "Steampipe Table: aws_codecommit_repository - Query AWS CodeCommit Repositories using SQL"
description: "Allows users to query AWS CodeCommit repositories and retrieve data such as repository name, ARN, description, clone URL, last modified date, and other related details."
folder: "CodeCommit"
---

# Table: aws_codecommit_repository - Query AWS CodeCommit Repositories using SQL

The AWS CodeCommit Repository is a fully-managed source control service that hosts secure Git-based repositories. It makes it easy for teams to collaborate on code in a secure and highly scalable ecosystem. CodeCommit eliminates the need to operate your own source control system or worry about scaling its infrastructure.

## Table Usage Guide

The `aws_codecommit_repository` table in Steampipe provides you with information about repositories within AWS CodeCommit. This table allows you, as a DevOps engineer, to query repository-specific details, including repository name, ARN, description, clone URL, last modified date, and other related details. You can utilize this table to gather insights on repositories, such as repositories with specific ARNs, the last modified date of repositories, verification of clone URLs, and more. The schema outlines the various attributes of the CodeCommit repository for you, including the repository name, ARN, clone URL, and associated metadata.

## Examples

### Basic info
This query allows you to explore the details of your AWS CodeCommit repositories, including their names, IDs, creation dates, and regions. It's useful for gaining insights into your repository usage and organization across different regions.

```sql+postgres
select
  repository_name,
  repository_id,
  arn,
  creation_date,
  region
from
  aws_codecommit_repository;
```

```sql+sqlite
select
  repository_name,
  repository_id,
  arn,
  creation_date,
  region
from
  aws_codecommit_repository;
```