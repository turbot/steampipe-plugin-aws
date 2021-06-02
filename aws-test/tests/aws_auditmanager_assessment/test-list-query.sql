select
  name,
  arn,
  id
from
  aws.aws_auditmanager_assessment
where
  arn = '{{ output.assessment_arn.value }}';