select akas, tags, title
from aws.aws_ec2_managed_prefix_list
where prefix_list_id = 'demoID';
