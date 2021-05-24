select
  name,
  title,
  tags,
  akas
from
  aws.aws_audit_manager_assessment
where
  id = '{{ output.assessment_id.value }}';