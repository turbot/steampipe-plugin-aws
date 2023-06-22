select
  deployment_group_id,
  deployment_group_name,
  arn,
  tags
from
  aws_codedeploy_deployment_group
where
  deployment_group_name = '{{ resourceName }}' and application_name = '{{ output.app_name.value }}';
