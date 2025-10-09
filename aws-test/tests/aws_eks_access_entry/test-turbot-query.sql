select
  cluster_name,
  principal_arn,
  title,
  akas,
  tags
from aws.aws_eks_access_entry
where cluster_name = '{{ resourceName }}' and principal_arn = '{{ output.principal_arn.value }}';

