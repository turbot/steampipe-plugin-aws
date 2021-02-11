select
  principal_arn,
  service_name,
  service_namespace,
  last_authenticated_entity,
  partition,
  region,
  last_authenticated_region
from aws.aws_iam_access_advisor
where principal_arn = '{{ output.user_arn.value }}' and service_namespace = 'sts';
