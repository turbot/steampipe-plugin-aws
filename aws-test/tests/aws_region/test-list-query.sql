
select account_id, akas, name, opt_in_status, partition, region, title
from aws.aws_region
where akas::text = '["arn:{{ output.aws_partition.value }}::ap-south-1:{{ output.account_id.value }}"]';
