# Table: aws_eks_identity_provider_config

Amazon Elastic Kubernetes Service (EKS) OpenID Connect (OIDC) Identity Provider (IDP). This feature allows customers to integrate an OIDC identity provider with a new or existing Amazon EKS cluster running Kubernetes version 1.16 or later.

## Examples

### Basic info

```sql
select
  name,
  arn,
  cluster_name,
  tags,
  status
from
  aws_eks_identity_provider_config;
```

### List OIDC type Identity provider config

```sql
select
  name,
  arn,
  cluster_name,
  type
from
  aws_eks_identity_provider_config
where 
  type = 'oidc';
```