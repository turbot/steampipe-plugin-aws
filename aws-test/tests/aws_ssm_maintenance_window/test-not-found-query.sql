select title, akas, tags, region, account_id
from aws.aws_ssm_maintenance_window
where window_id = 'mw-0000000000a000000'
