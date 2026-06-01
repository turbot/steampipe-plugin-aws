---
title: "Steampipe Table: aws_eks_pod_identity_association - Query AWS EKS Pod Identity Associations using SQL"
description: "Allows users to query AWS EKS Pod Identity Associations to retrieve information about IAM role bindings for Kubernetes service accounts in Amazon EKS clusters."
folder: "EKS"
---

# Table: aws_eks_pod_identity_association - Query AWS EKS Pod Identity Associations using SQL

AWS EKS Pod Identity Associations link Kubernetes service accounts to IAM roles, allowing pods running under those service accounts to assume the associated IAM role and access AWS services. Pod Identity is a simplified alternative to IAM Roles for Service Accounts (IRSA) that does not require annotating the service account or configuring an OIDC provider. Each association is scoped to a specific cluster, namespace, and service account combination.

## Table Usage Guide

The `aws_eks_pod_identity_association` table in Steampipe provides you with information about pod identity associations in your Amazon EKS clusters. This table allows you, as a DevOps engineer or security administrator, to query association-specific details including the IAM role ARN, Kubernetes namespace and service account, and associated metadata. You can use this table to audit which service accounts have IAM role access, review pod identity configurations across clusters, and identify associations by cluster, namespace, or service account. The schema outlines the various attributes of each association, including the association ARN, role ARN, creation time, and tags.

## Examples

### Basic info
Explore all pod identity associations across your EKS clusters to understand service account to IAM role mappings.

```sql+postgres
select
  cluster_name,
  association_id,
  namespace,
  service_account,
  role_arn,
  created_at
from
  aws_eks_pod_identity_association;
```

```sql+sqlite
select
  cluster_name,
  association_id,
  namespace,
  service_account,
  role_arn,
  created_at
from
  aws_eks_pod_identity_association;
```

### Filter associations by cluster name
List all pod identity associations for a specific EKS cluster.

```sql+postgres
select
  cluster_name,
  association_id,
  namespace,
  service_account,
  role_arn
from
  aws_eks_pod_identity_association
where
  cluster_name = 'my-eks-cluster';
```

```sql+sqlite
select
  cluster_name,
  association_id,
  namespace,
  service_account,
  role_arn
from
  aws_eks_pod_identity_association
where
  cluster_name = 'my-eks-cluster';
```

### Filter associations by namespace
Identify all service accounts in a specific Kubernetes namespace that have IAM role bindings.

```sql+postgres
select
  cluster_name,
  association_id,
  namespace,
  service_account,
  role_arn
from
  aws_eks_pod_identity_association
where
  namespace = 'kube-system';
```

```sql+sqlite
select
  cluster_name,
  association_id,
  namespace,
  service_account,
  role_arn
from
  aws_eks_pod_identity_association
where
  namespace = 'kube-system';
```

### List associations for a specific service account
Find all IAM role associations for a particular Kubernetes service account across all clusters.

```sql+postgres
select
  cluster_name,
  namespace,
  service_account,
  association_id,
  role_arn,
  created_at
from
  aws_eks_pod_identity_association
where
  service_account = 'my-service-account';
```

```sql+sqlite
select
  cluster_name,
  namespace,
  service_account,
  association_id,
  role_arn,
  created_at
from
  aws_eks_pod_identity_association
where
  service_account = 'my-service-account';
```

### Get a specific association by cluster and association ID
Retrieve details of a particular pod identity association using its cluster name and association ID.

```sql+postgres
select
  cluster_name,
  association_id,
  association_arn,
  namespace,
  service_account,
  role_arn,
  created_at,
  modified_at
from
  aws_eks_pod_identity_association
where
  cluster_name = 'my-eks-cluster'
  and association_id = 'a-9njjin9gfghecgocd';
```

```sql+sqlite
select
  cluster_name,
  association_id,
  association_arn,
  namespace,
  service_account,
  role_arn,
  created_at,
  modified_at
from
  aws_eks_pod_identity_association
where
  cluster_name = 'my-eks-cluster'
  and association_id = 'a-9njjin9gfghecgocd';
```

### List associations that have tags
Find pod identity associations that have been tagged for resource management or cost allocation.

```sql+postgres
select
  cluster_name,
  association_id,
  namespace,
  service_account,
  role_arn,
  tags
from
  aws_eks_pod_identity_association
where
  tags != '{}';
```

```sql+sqlite
select
  cluster_name,
  association_id,
  namespace,
  service_account,
  role_arn,
  tags
from
  aws_eks_pod_identity_association
where
  tags != '{}';
```

### Count pod identity associations per cluster
Determine the number of pod identity associations configured in each EKS cluster.

```sql+postgres
select
  cluster_name,
  count(*) as association_count
from
  aws_eks_pod_identity_association
group by
  cluster_name
order by
  association_count desc;
```

```sql+sqlite
select
  cluster_name,
  count(*) as association_count
from
  aws_eks_pod_identity_association
group by
  cluster_name
order by
  association_count desc;
```

### Join with IAM roles to review associated role details
Combine pod identity association data with IAM role information to review permissions granted to service accounts.

```sql+postgres
select
  pia.cluster_name,
  pia.namespace,
  pia.service_account,
  pia.role_arn,
  r.create_date as role_created_at,
  r.max_session_duration
from
  aws_eks_pod_identity_association pia
  left join aws_iam_role r on pia.role_arn = r.arn;
```

```sql+sqlite
select
  pia.cluster_name,
  pia.namespace,
  pia.service_account,
  pia.role_arn,
  r.create_date as role_created_at,
  r.max_session_duration
from
  aws_eks_pod_identity_association pia
  left join aws_iam_role r on pia.role_arn = r.arn;
```
