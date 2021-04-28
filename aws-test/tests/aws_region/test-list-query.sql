select *
from aws.aws_region
where akas::text = '["arn:{{ output.aws_partition.value }}::ap-south-1:{{ output.account_id.value }}"]';