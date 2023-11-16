---
title: "Table: aws_wellarchitected_share_invitation - Query AWS Well-Architected Tool Share Invitations using SQL"
description: "Allows users to query Share Invitations in the AWS Well-Architected Tool."
---

# Table: aws_wellarchitected_share_invitation - Query AWS Well-Architected Tool Share Invitations using SQL

The `aws_wellarchitected_share_invitation` table in Steampipe provides information about share invitations within the AWS Well-Architected Tool. This table allows DevOps engineers to query invitation-specific details, including the invitation status, the recipient's AWS account ID, and associated metadata. Users can utilize this table to gather insights on share invitations, such as pending invitations, accepted invitations, and more. The schema outlines the various attributes of the share invitation, including the invitation ARN, workload ID, permission type, and invitation status.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_wellarchitected_share_invitation` table, you can use the `.inspect aws_wellarchitected_share_invitation` command in Steampipe.

### Key columns:

- `invitation_id`: This is the unique identifier of the share invitation. It can be used to join with other tables when detailed analysis of a specific invitation is required.
- `workload_id`: This is the identifier of the workload that the invitation is associated with. It can be used to join with workload-related tables for a comprehensive view of the workload's sharing status.
- `permission_type`: This column indicates the permission level granted with the invitation. It can be used to evaluate the level of access granted to shared users.

## Examples

### Basic info

```sql
select
  share_invitation_id,
  permission_type,
  shared_by,
  shared_with,
  share_resource_type
from
  aws_wellarchitected_share_invitation;
```

### List invitations for Well-Architected lens resources

```sql
select
  lens_arn,
  lens_name,
  share_invitation_id,
  permission_type,
  shared_by,
  shared_with
from
  aws_wellarchitected_share_invitation
where
  share_resource_type = 'LENS'
  or lens_arn is not null;
```

### List invitations for Well-Architected workload resources

```sql
select
  workload_id,
  workload_name,
  share_invitation_id,
  permission_type,
  shared_by,
  shared_with
from
  aws_wellarchitected_share_invitation
where
  share_resource_type = 'WORKLOAD'
  or workload_id is not null;
```

### List invitations allowing CONTRIBUTOR permission to resources

```sql
select
  share_invitation_id,
  permission_type,
  shared_by,
  shared_with,
  share_resource_type
from
  aws_wellarchitected_share_invitation
where
  permission_type = 'CONTRIBUTOR';
```

### List invitations allowing READONLY permission to resources

```sql
select
  share_invitation_id,
  permission_type,
  shared_by,
  shared_with,
  share_resource_type
from
  aws_wellarchitected_share_invitation
where
  permission_type = 'READONLY';
```

### List total invitations for each resource type

```sql
select
  count(*) as total,
  share_resource_type
from
  aws_wellarchitected_share_invitation
group by
  share_resource_type;
```