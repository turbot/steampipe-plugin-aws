select name, title, akas 
from aws_sfn_state_machine
where akas = '["{{output.resource_aka.value}}"]';
