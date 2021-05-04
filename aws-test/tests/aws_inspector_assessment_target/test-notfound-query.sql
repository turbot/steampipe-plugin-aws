select *
from aws.aws_inspector_assessment_target
where arn = 'arn:{{ output.aws_partition.value }}:inspector:{{ output.aws_region.value }}:{{ output.account_id.value }}:target/0-O0000000';