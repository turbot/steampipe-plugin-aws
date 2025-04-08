---
title: "Steampipe Table: aws_wellarchitected_share_invitation - Query AWS Well-Architected Tool Share Invitations using SQL"
description: "Allows users to query Share Invitations in the AWS Well-Architected Tool."
folder: "Well-Architected"
---

# Table: aws_wellarchitected_share_invitation - Query AWS Well-Architected Tool Share Invitations using SQL

The AWS Well-Architected Tool Share Invitations are part of AWS's Well-Architected Tool, which enables you to review the state of your workloads and compares them to the latest AWS architectural best practices. The share invitations specifically allow for the sharing of workload reports with other AWS accounts. This aids in collaborative efforts to improve system performance, increase security, and optimize costs.

## Table Usage Guide

The `aws_wellarchitected_share_invitation` table in Steampipe provides you with information about share invitations within the AWS Well-Architected Tool. This table allows you, as a DevOps engineer, to query invitation-specific details, including the invitation status, the recipient's AWS account ID, and associated metadata. You can utilize this table to gather insights on share invitations, such as pending invitations, accepted invitations, and more. The schema outlines the various attributes of the share invitation for you, including the invitation ARN, workload ID, permission type, and invitation status.

## Examples

### Basic info
Explore which resources have been shared in your AWS Well-Architected environment. This can help you understand who has access to what, allowing you to maintain better control over your data and resources.

```sql+postgres
select
  share_invitation_id,
  permission_type,
  shared_by,
  shared_with,
  share_resource_type
from
  aws_wellarchitected_share_invitation;
```

```sql+sqlite
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
This example helps to identify the sharing invitations related to Well-Architected lens resources. It is particularly useful for understanding who has shared these resources and with whom, providing insights into the distribution and permissions of your Well-Architected lens resources.

```sql+postgres
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

```sql+sqlite
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
Explore which Well-Architected workload resources have been shared with others. This can be useful for auditing purposes or to understand the distribution of workload resources within your organization.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which contributor permissions have been granted to resources. This allows for oversight and management of resource access within the AWS Well-Architected framework.

```sql+postgres
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

```sql+sqlite
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
Explore which share invitations are granting READONLY permissions in your AWS Well-Architected environment. This can help identify potential security risks and ensure that only authorized users have access to sensitive resources.

```sql+postgres
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

```sql+sqlite
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
Explore the distribution of share invitations across different resource types in your AWS Well-Architected tool. This can help you understand which resources are most frequently shared, aiding in resource management and security practices.

```sql+postgres
select
  count(*) as total,
  share_resource_type
from
  aws_wellarchitected_share_invitation
group by
  share_resource_type;
```

```sql+sqlite
select
  count(*) as total,
  share_resource_type
from
  aws_wellarchitected_share_invitation
group by
  share_resource_type;
```