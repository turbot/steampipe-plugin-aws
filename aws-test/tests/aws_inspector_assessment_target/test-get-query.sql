select name, arn, resource_group_arn, account_id, partition, region
from aws.aws_inspector_assessment_target
where arn = '{{ output.resource_aka.value }}';