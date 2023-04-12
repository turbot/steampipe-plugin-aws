# Table: aws_wellarchitected_answer

The answers of a lens review in a Well-Architected workload

**Important notes:**

- You **_must_** specify `workload_id` in a `where` clause in order to use this table.
- For improved performance, it is advised that you use the optional qual `pillar_id` and `lens_alias` to limit the result set to a specific lens or pillar.

## Examples

### Basic info

```sql
select
  a.question_id,
  a.lens_alias,
  a.workload_id,
  a.is_applicable,
  a.pillar_id,
  a.question_title,
  a.risk,
  a.reason,
  a.region
from
  aws_wellarchitected_answer a,
  aws_wellarchitected_workload w 
where
  a.workload_id = w.workload_id;
```

### Get number of questions per piller

```sql
select
  a.workload_id,
  a.pillar_id,
  count(a.question_id) as total_questions
from
  aws_wellarchitected_answer a,
  aws_wellarchitected_workload w 
where
  a.workload_id = w.workload_id
group by
  a.pillar_id,
  a.workload_id;
```

### List all the questions along with the choices for a workload

```sql
select
  a.question_id,
  a.lens_alias,
  a.workload_id,
  a.question_title,
  a.question_description,
  c ->> 'Title' as choice_title,
  c ->> 'ChoiceId' as choice_id,
  c ->> 'Description' as choice_description,
  c ->> 'HelpfulResource' as choice_helpful_resource,
  c ->> 'ImprovementPlan' as choice_improvement_plan
from
  aws_wellarchitected_answer a,
  aws_wellarchitected_workload w,
  jsonb_array_elements(choices) c
where
  a.workload_id = w.workload_id;
```

### List all the questions along with the answered choices for a workload

```sql
select
  a.question_id,
  a.lens_alias,
  a.workload_id,
  a.question_title,
  a.question_description,
  c ->> 'Notes' as choice_notes,
  c ->> 'Reason' as choice_reason,
  c ->> 'Status' as choice_status,
  c ->> 'ChoiceId' as choice_id
from
  aws_wellarchitected_answer a,
  aws_wellarchitected_workload w,
  jsonb_array_elements(choice_answers) c
where
  a.workload_id = w.workload_id;
```

### List all the questions that are not applicable for a workload

```sql
select
  a.question_id,
  a.lens_alias,
  a.workload_id,
  a.question_title,
  a.question_description
from
  aws_wellarchitected_answer a,
  aws_wellarchitected_workload w
where
  a.workload_id = w.workload_id
  and not is_applicable;
```