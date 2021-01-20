select instance_id, kernel_id, user_data, ram_disk_id, disable_api_termination, sriov_net_support, instance_initiated_shutdown_behavior, akas, tags, title
from aws.aws_ec2_instance
where instance_id = '{{ output.resource_id.value }}'
