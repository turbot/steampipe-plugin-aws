---
title: "Steampipe Table: aws_eks_access_policy_association - Query AWS EKS Access Policy Associations using SQL"
description: "Allows users to query AWS EKS Access Policy Associations to retrieve information about access policies associated with IAM principals in Amazon EKS clusters."
folder: "EKS"
---

# Table: aws_eks_access_policy_association - Query AWS EKS Access Policy Associations using SQL

AWS EKS Access Policy Associations link access policies to IAM principals (users or roles) for Amazon EKS clusters. These associations define the scope of access (cluster-wide or namespace-specific) and enable fine-grained access control for EKS cluster resources through AWS-managed access policies.

## Table Usage Guide

The `aws_eks_access_policy_association` table in Steampipe provides you with information about access policy associations for each Amazon EKS access entry. This table allows you, as a DevOps engineer or security administrator, to query association-specific details, including the policy ARN, access scope type, namespaces, and association timestamps. You can utilize this table to gather insights on access policy assignments, such as which policies are associated with specific IAM principals, the scope of access granted, and when policies were associated or modified. The schema outlines the various attributes of the EKS access policy association for you, including the cluster name, principal ARN, policy ARN, access scope, and timestamps.

## Examples

### Basic info
Explore which access policies are associated with IAM principals in your EKS clusters.

```sql+postgres
select
  cluster_name,
  principal_arn,
  policy_arn,
  access_scope_type,
  associated_at,
  modified_at
from
  aws_eks_access_policy_association;
```

```sql+sqlite
select
  cluster_name,
  principal_arn,
  policy_arn,
  access_scope_type,
  associated_at,
  modified_at
from
  aws_eks_access_policy_association;
```

### List policy associations for a specific cluster
Identify all access policy associations for a particular EKS cluster.

```sql+postgres
select
  cluster_name,
  principal_arn,
  policy_arn,
  access_scope_type,
  access_scope_namespaces
from
  aws_eks_access_policy_association
where
  cluster_name = 'my-eks-cluster';
```

```sql+sqlite
select
  cluster_name,
  principal_arn,
  policy_arn,
  access_scope_type,
  access_scope_namespaces
from
  aws_eks_access_policy_association
where
  cluster_name = 'my-eks-cluster';
```

### List policy associations for a specific IAM principal
Find all access policies associated with a specific IAM role or user.

```sql+postgres
select
  cluster_name,
  principal_arn,
  policy_arn,
  access_scope_type,
  access_scope_namespaces
from
  aws_eks_access_policy_association
where
  principal_arn = 'arn:aws:iam::123456789012:role/my-eks-role';
```

```sql+sqlite
select
  cluster_name,
  principal_arn,
  policy_arn,
  access_scope_type,
  access_scope_namespaces
from
  aws_eks_access_policy_association
where
  principal_arn = 'arn:aws:iam::123456789012:role/my-eks-role';
```

### List namespace-scoped policy associations
Determine which policy associations are limited to specific Kubernetes namespaces.

```sql+postgres
select
  cluster_name,
  principal_arn,
  policy_arn,
  access_scope_type,
  access_scope_namespaces
from
  aws_eks_access_policy_association
where
  access_scope_type = 'namespace';
```

```sql+sqlite
select
  cluster_name,
  principal_arn,
  policy_arn,
  access_scope_type,
  access_scope_namespaces
from
  aws_eks_access_policy_association
where
  access_scope_type = 'namespace';
```

### List cluster-wide policy associations
Find policy associations that grant cluster-wide access.

```sql+postgres
select
  cluster_name,
  principal_arn,
  policy_arn,
  associated_at
from
  aws_eks_access_policy_association
where
  access_scope_type = 'cluster';
```

```sql+sqlite
select
  cluster_name,
  principal_arn,
  policy_arn,
  associated_at
from
  aws_eks_access_policy_association
where
  access_scope_type = 'cluster';
```

### Get count of policy associations per cluster
Determine the distribution of access policy associations across clusters.

```sql+postgres
select
  cluster_name,
  count(*) as policy_association_count
from
  aws_eks_access_policy_association
group by
  cluster_name
order by
  policy_association_count desc;
```

```sql+sqlite
select
  cluster_name,
  count(*) as policy_association_count
from
  aws_eks_access_policy_association
group by
  cluster_name
order by
  policy_association_count desc;
```

### Get count of policy associations per principal
Identify how many access policies are associated with each IAM principal.

```sql+postgres
select
  principal_arn,
  count(*) as policy_count
from
  aws_eks_access_policy_association
group by
  principal_arn
order by
  policy_count desc;
```

```sql+sqlite
select
  principal_arn,
  count(*) as policy_count
from
  aws_eks_access_policy_association
group by
  principal_arn
order by
  policy_count desc;
```

### List recently modified policy associations
Identify policy associations that have been recently updated.

```sql+postgres
select
  cluster_name,
  principal_arn,
  policy_arn,
  access_scope_type,
  modified_at
from
  aws_eks_access_policy_association
where
  modified_at >= now() - interval '7 days'
order by
  modified_at desc;
```

```sql+sqlite
select
  cluster_name,
  principal_arn,
  policy_arn,
  access_scope_type,
  modified_at
from
  aws_eks_access_policy_association
where
  modified_at >= datetime('now', '-7 days')
order by
  modified_at desc;
```

### Get policy associations with access entry details
Join policy associations with access entry information for comprehensive access insights.

```sql+postgres
select
  apa.cluster_name,
  apa.principal_arn,
  apa.policy_arn,
  apa.access_scope_type,
  apa.access_scope_namespaces,
  ae.type as access_entry_type,
  ae.username,
  ae.kubernetes_groups
from
  aws_eks_access_policy_association apa
  left join aws_eks_access_entry ae
    on apa.cluster_name = ae.cluster_name
    and apa.principal_arn = ae.principal_arn;
```

```sql+sqlite
select
  apa.cluster_name,
  apa.principal_arn,
  apa.policy_arn,
  apa.access_scope_type,
  apa.access_scope_namespaces,
  ae.type as access_entry_type,
  ae.username,
  ae.kubernetes_groups
from
  aws_eks_access_policy_association apa
  left join aws_eks_access_entry ae
    on apa.cluster_name = ae.cluster_name
    and apa.principal_arn = ae.principal_arn;
```

### List policy associations with specific policy ARN
Find all associations using a specific AWS-managed access policy.

```sql+postgres
select
  cluster_name,
  principal_arn,
  policy_arn,
  access_scope_type,
  associated_at
from
  aws_eks_access_policy_association
where
  policy_arn like '%AmazonEKSClusterAdminPolicy';
```

```sql+sqlite
select
  cluster_name,
  principal_arn,
  policy_arn,
  access_scope_type,
  associated_at
from
  aws_eks_access_policy_association
where
  policy_arn like '%AmazonEKSClusterAdminPolicy';
```

