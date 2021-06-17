# Table: aws_eks_addon_version

Amazon EKS add-ons help to automate the provisioning and lifecycle management of common operational software for Amazon EKS clusters.

Amazon EKS add-on versions describe the Kubernetes versions that the add-on can be used with.

## Examples

### Basic info

```sql
select
  addon_name,
  addon_version,
  type
from
  aws_eks_addon_version;
```

### Count the number of add-on versions by add-on

```sql
select
  addon_name,
  count(addon_version) as addon_version_count
from
  aws_eks_addon_version
group by
  addon_name;
```
