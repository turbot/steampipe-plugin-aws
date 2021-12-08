select name, id, arn, address_family, owner_id
from aws.aws_ec2_managed_prefix_list
where name = '{{ resourceName }}';
