select name, image_id, image_type, image_location, architecture, description, ena_support, hypervisor, owner_id, platform_details, public, root_device_name, root_device_type, sriov_net_support, usage_operation, virtualization_type, block_device_mappings
from aws.aws_ec2_ami_shared
where image_id = '{{ output.resource_id.value }}';
