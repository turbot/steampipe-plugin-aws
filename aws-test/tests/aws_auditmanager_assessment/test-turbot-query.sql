select
  name,
  title,
  tags,
  akas
from
  aws.aws_auditmanager_assessment
where
  id = '{{ output.assessment_id.value }}';