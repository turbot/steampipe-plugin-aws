select name, 
rule_id,
rule_arn,
rule_state, 
tags_src,
description,
title,
akas
from aws.aws_config_rule
where name = 'dummy-{{ resourceName }}';