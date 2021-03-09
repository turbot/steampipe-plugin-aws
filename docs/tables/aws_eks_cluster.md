# Table: aws_eks_cluster

Amazon Elastic Kubernetes Service (EKS) is a managed Kubernetes service that makes it easy to run Kubernetes on AWS and on-premises.

## Examples

### Basic info

```sql
select
  name,
  arn,
  endpoint,
  identity,
  status
from
  aws_eks_cluster;;
```