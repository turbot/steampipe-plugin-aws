# Table: aws_eks_addon

Amazon EKS add-ons help to automate the provisioning and lifecycle management of common operational software for Amazon EKS clusters.

## Examples

### Basic info

```sql
select
  addon_name,
  arn,
  addon_version,
  cluster_name,
  status,
  service_account_role_arn
from
  aws_eks_addon;
```


### List add-ons that are not active

```sql
select
  addon_name,
  arn,
  cluster_name,
  status
from
  aws_eks_addon
where
  status <> 'ACTIVE';
```


### Get count of add-ons by cluster

```sql
select
  cluster_name,
  count(addon_name) as addon_count
from
  aws_eks_addon
group by
  cluster_name;
```
