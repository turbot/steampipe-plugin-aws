---
title: "Table: aws_wellarchitected_answer - Query AWS Well-Architected Tool Answer using SQL"
description: "Allows users to query AWS Well-Architected Tool Answer data, including information about the workloads, lens, and questions associated with each answer."
---

# Table: aws_wellarchitected_answer - Query AWS Well-Architected Tool Answer using SQL

The `aws_wellarchitected_answer` table in Steampipe provides information about the answers within AWS Well-Architected Tool. This table allows DevOps engineers to query answer-specific details, including the workload, lens, and question associated with each answer. Users can utilize this table to gather insights on answers, such as the workload and lens associated with a specific answer, the question that the answer corresponds to, and more. The schema outlines the various attributes of the Well-Architected Tool answer, including the answer ID, workload ID, lens alias, and associated metadata.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_wellarchitected_answer` table, you can use the `.inspect aws_wellarchitected_answer` command in Steampipe.

###  Key columns:

- `answer_id`: This is the unique identifier for each answer in the AWS Well-Architected Tool. It is crucial for joining this table with other tables that contain answer-specific information.
- `workload_id`: This column contains the identifier of the workload associated with each answer. It is important for correlating answers with their respective workloads.
- `lens_alias`: This column holds the alias of the lens associated with each answer. It is useful for grouping answers by their associated lens.

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

### Get the number of questions per pillar

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
