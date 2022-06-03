select akas, title
from aws_sfn_state_machine
where name = '{{ resourceName }}';
