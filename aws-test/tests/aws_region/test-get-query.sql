
select account_id, akas, name, opt_in_status, partition, region, title
from aws.aws_region
where name='ap-south-1';
