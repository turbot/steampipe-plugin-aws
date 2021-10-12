# Table: aws_workspaces_workspace

Amazon WorkSpaces enables you to provision virtual, cloud-based Microsoft Windows or Amazon Linux desktops for your users, known as WorkSpaces. WorkSpaces eliminates the need to procure and deploy hardware or install complex software.

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