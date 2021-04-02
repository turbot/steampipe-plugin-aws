select cluster_subnet_group_name, description, subnet_group_status, vpc_id, tags_src from aws.aws_redshift_subnet_group where cluster_subnet_group_name = '{{ resourceName }}';
