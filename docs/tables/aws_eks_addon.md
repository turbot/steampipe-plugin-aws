# Table: aws_eks_addon

Amazon EKS add-ons help to automate the provisioning and lifecycle management of common operational software for Amazon EKS clusters.

## Examples

### Basic info

```sql
select
  addon_name,
  addon_arn,
  addon_version,
  cluster_name,
  status,
  service_account_role_arn
from
  aws_eks_addon;
```


### List addons that are not active

```sql
select
  addon_name,
  addon_arn,
  cluster_name,
  status
from
  aws_eks_addon
where
  status <> 'ACTIVE';
```


### Get addon count by cluster name

```sql
select
  cluster_name,
  count(addon_name) as addon_count
from
  aws_eks_addon
group by
  cluster_name;
```