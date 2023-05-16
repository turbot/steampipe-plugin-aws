# Table: aws_wellarchitected_answer

The answers of a lens review in a Well-Architected workload.

**Important notes:**

- For improved performance, it is advised that you use the optional qual `workload_id`, `pillar_id` and `lens_alias` to limit the result set to a specific workload, lens or pillar.

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
  aws_wellarchitected_answer a;
```

### Get number of questions per piller

```sql
select
  a.workload_id,
  a.pillar_id,
  count(a.question_id) as total_questions
from
  aws_wellarchitected_answer a
group by
  a.workload_id,
  a.pillar_id;
```

### List all the questions along with the choices

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
  jsonb_array_elements(choices) c;
```

### List all the questions along with the answered choices

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
  jsonb_array_elements(choice_answers) c;
```

### List questions that are not applicable for a workload

```sql
select
  a.question_id,
  a.lens_alias,
  a.workload_id,
  a.question_title,
  a.question_description,
  reason
from
  aws_wellarchitected_answer a
where
  not is_applicable;
```

### List questions that are marked as high or medium risk

```sql
select
  a.question_id,
  a.lens_alias,
  a.workload_id,
  a.question_title,
  a.risk,
  c ->> 'ChoiceId' as choice_id,
  c ->> 'Status' as choice_status,
  c ->> 'Reason' as choice_reason,
  c ->> 'Notes' as choice_notes
from
  aws_wellarchitected_answer a,
  jsonb_array_elements(choice_answers) c
where
  risk = 'HIGH'
  or risk = 'MEDIUM';
```

### Get count of questions in each risk factor for each workload

```sql
select
  workload_id,
  risk,
  count(question_id) as total_questions
from
  aws_wellarchitected_answer
where
  risk = 'HIGH'
  or risk = 'MEDIUM'
group by
  workload_id,
  risk;
```