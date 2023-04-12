# Table: aws_wellarchitected_share_invitation

A share invitation is a request to share a workload or custom lens owned by another AWS account. A workload or lens can be shared with all users in an AWS account, individual users, or both. If you accept a workload invitation, the workload is added to your Workloads and Dashboard pages. If you accept a custom lens invitation, the lens is added to your Custom lenses page. If you reject the invitation, it's removed from the list.

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

### List invitations for Well-Architected Lens Resources

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

### List invitations for Well-Architected Workload Resources

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