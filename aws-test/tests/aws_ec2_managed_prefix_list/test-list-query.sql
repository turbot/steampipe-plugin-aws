select prefix_list_name, prefix_list_id, prefix_list_arn, address_family, owner_id
from aws.aws_ec2_managed_prefix_list
where prefix_list_name = '{{ resourceName }}';
