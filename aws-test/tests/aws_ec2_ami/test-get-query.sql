select name, image_id, image_type, image_location, architecture, description, ena_support, hypervisor, owner_id, platform_details, public, root_device_name, root_device_type, sriov_net_support, usage_operation, virtualization_type, block_device_mappings, tags_src, deregistration_protection, free_tier_eligible, image_allowed, source_image_id, source_image_region, last_launched_time
from aws.aws_ec2_ami
where image_id = '{{ output.resource_id.value }}'
