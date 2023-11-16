---
title: "Table: aws_workspaces_workspace - Query Amazon WorkSpaces Workspace using SQL"
description: "Allows users to query Amazon WorkSpaces Workspace to retrieve details about each workspace in the AWS account."
---

# Table: aws_workspaces_workspace - Query Amazon WorkSpaces Workspace using SQL

The `aws_workspaces_workspace` table in Steampipe provides information about each workspace within Amazon WorkSpaces. This table allows DevOps engineers to query workspace-specific details, including workspace properties, state, type, and associated metadata. Users can utilize this table to gather insights on workspaces, such as workspace status, user properties, root volume, user volume, and more. The schema outlines the various attributes of the workspace, including the workspace ID, directory ID, bundle ID, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_workspaces_workspace` table, you can use the `.inspect aws_workspaces_workspace` command in Steampipe.

### Key columns:

- `workspace_id`: The ID of the workspace. It can be used to join this table with other tables to get more detailed information about a specific workspace.
- `directory_id`: The ID of the directory. It can be used to join this table with the `aws_directory_service_directory` table to get information about the directory associated with the workspace.
- `bundle_id`: The ID of the bundle used to create the workspace. It can be used to join this table with the `aws_workspaces_bundle` table to get information about the bundle used for the workspace.

## Examples

## Basic info

```sql
select
  name,
  workspace_id,
  arn,
  state
from
  aws_workspaces_workspace;
```


## List terminated workspaces

```sql
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