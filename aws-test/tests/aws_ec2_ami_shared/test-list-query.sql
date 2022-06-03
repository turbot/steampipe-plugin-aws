select name, image_id
from aws.aws_ec2_ami_shared
where owner_id = '{{ output.account_id.value }}' and name = '{{ resourceName }}';
