---
title: "Steampipe Table: aws_workspaces_workspace - Query Amazon WorkSpaces Workspace using SQL"
description: "Allows users to query Amazon WorkSpaces Workspace to retrieve details about each workspace in the AWS account."
folder: "WorkSpaces"
---

# Table: aws_workspaces_workspace - Query Amazon WorkSpaces Workspace using SQL

The Amazon WorkSpaces service provides fully managed, persistent desktops in the cloud. A Workspace is a cloud-based virtual desktop that can integrate with your corporate Active Directory so that your users can use their existing credentials to access their WorkSpaces. It offers a choice of bundles providing different amounts of CPU, memory, and solid-state storage to meet your users' performance needs.

## Table Usage Guide

The `aws_workspaces_workspace` table in Steampipe provides you with information about each workspace within Amazon WorkSpaces. This table allows you, as a DevOps engineer, to query workspace-specific details, including workspace properties, state, type, and associated metadata. You can utilize this table to gather insights on workspaces, such as workspace status, user properties, root volume, user volume, and more. The schema outlines the various attributes of the workspace for you, including the workspace ID, directory ID, bundle ID, and associated tags.

## Examples

## Basic info

```sql+postgres
select
  name,
  workspace_id,
  arn,
  state
from
  aws_workspaces_workspace;
```

```sql+sqlite
select
  name,
  workspace_id,
  arn,
  state
from
  aws_workspaces_workspace;
```


## List terminated workspaces

```sql+postgres
select
  name,
  workspace_id,
  arn,
  state
from
  aws_workspaces_workspace
where
  state = 'TERMINATED';
```

```sql+sqlite
select
  name,
  workspace_id,
  arn,
  state
from
  aws_workspaces_workspace
where
  state = 'TERMINATED';
```