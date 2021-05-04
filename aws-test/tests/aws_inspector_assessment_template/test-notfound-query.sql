select *
from aws.aws_inspector_assessment_template
where arn = 'arn:{{ output.aws_partition.value }}:inspector:{{ output.aws_region.value }}:{{ output.account_id.value }}:target/0-O0000000/template/0-00000000';