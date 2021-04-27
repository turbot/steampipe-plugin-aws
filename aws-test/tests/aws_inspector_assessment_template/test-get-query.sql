select name, arn, assessment_run_count, assessment_target_arn, duration_in_seconds, last_assessment_run_arn, user_attributes_for_findings, tags_src, account_id, partition, region
from aws.aws_inspector_assessment_template
where arn = '{{ output.resource_aka.value }}';