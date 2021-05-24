select
  name,
  arn,
  id
from
  aws.aws_audit_manager_assessment
where
  arn = '{{ output.assessment_arn.value }}';