---
title: "Steampipe Table: aws_wellarchitected_answer - Query AWS Well-Architected Tool Answer using SQL"
description: "Allows users to query AWS Well-Architected Tool Answer data, including information about the workloads, lens, and questions associated with each answer."
folder: "Well-Architected"
---

# Table: aws_wellarchitected_answer - Query AWS Well-Architected Tool Answer using SQL

The AWS Well-Architected Tool Answer is a feature within the AWS Well-Architected Tool service. It allows you to review your workloads against AWS architectural best practices, and provides guidance on improving your cloud architectures. The tool helps you understand the pros and cons of decisions you make while building workloads, and provides AWS best practices for creating high-performing, resilient, and efficient infrastructure for your applications.

## Table Usage Guide

The `aws_wellarchitected_answer` table in Steampipe provides you with information about the answers within AWS Well-Architected Tool. This table allows you, as a DevOps engineer, to query answer-specific details, including the workload, lens, and question associated with each answer. You can utilize this table to gather insights on answers, such as the workload and lens associated with a specific answer, the question that the answer corresponds to, and more. The schema outlines the various attributes of the Well-Architected Tool answer for you, including the answer ID, workload ID, lens alias, and associated metadata.

**Important Notes**
- For improved performance, it is advisable that you use the optional qual `workload_id`, `pillar_id`, and `lens_alias` to limit the result set to a specific workload, pillar, or lens.

## Examples

### Basic info
Explore the risk factors and reasons associated with different workloads in various regions using the AWS Well-Architected Framework. This can help identify potential issues and areas for improvement in your cloud architecture.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which you can assess the number of questions associated with each pillar for a specific workload in AWS Well-Architected framework. This is beneficial for understanding the distribution of questions across different pillars, aiding in workload management and strategic planning.

```sql+postgres
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

```sql+sqlite
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
Explore the various questions and associated choices within your AWS Well-Architected framework. This is useful for understanding the different options available for each question, helping to make informed decisions about your AWS architecture.

```sql+postgres
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

```sql+sqlite
select
  a.question_id,
  a.lens_alias,
  a.workload_id,
  a.question_title,
  a.question_description,
  json_extract(c.value, '$.Title') as choice_title,
  json_extract(c.value, '$.ChoiceId') as choice_id,
  json_extract(c.value, '$.Description') as choice_description,
  json_extract(c.value, '$.HelpfulResource') as choice_helpful_resource,
  json_extract(c.value, '$.ImprovementPlan') as choice_improvement_plan
from
  aws_wellarchitected_answer a,
  json_each(a.choices) c;
```

### List all the questions along with the answered choices
Determine the areas in which specific questions have been answered within a workload on AWS. This can help in understanding the reasons behind certain choices, their status, and associated notes, thereby providing a comprehensive view of decision-making processes.

```sql+postgres
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

```sql+sqlite
select
  a.question_id,
  a.lens_alias,
  a.workload_id,
  a.question_title,
  a.question_description,
  json_extract(c.value, '$.Notes') as choice_notes,
  json_extract(c.value, '$.Reason') as choice_reason,
  json_extract(c.value, '$.Status') as choice_status,
  json_extract(c.value, '$.ChoiceId') as choice_id
from
  aws_wellarchitected_answer a,
  json_each(a.choice_answers) as c;
```

### List questions that are not applicable for a workload
Determine the areas in which certain questions are not relevant for a specific workload in AWS Well-Architected Tool. This can help in focusing on applicable areas and streamlining the review process.

```sql+postgres
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

```sql+sqlite
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
  is_applicable = 0;
```

### List questions that are marked as high or medium risk
Determine areas of concern by identifying questions marked as high or medium risk. This is useful for prioritizing areas for improvement and mitigating potential issues.

```sql+postgres
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

```sql+sqlite
select
  a.question_id,
  a.lens_alias,
  a.workload_id,
  a.question_title,
  a.risk,
  json_extract(c.value, '$.ChoiceId') as choice_id,
  json_extract(c.value, '$.Status') as choice_status,
  json_extract(c.value, '$.Reason') as choice_reason,
  json_extract(c.value, '$.Notes') as choice_notes
from
  aws_wellarchitected_answer a,
  json_each(a.choice_answers) as c
where
  a.risk = 'HIGH'
  or a.risk = 'MEDIUM';
```

### Get count of questions in each risk factor for each workload
Assess the elements within each workload to identify the total number of high and medium risk questions. This provides insights into potential areas of concern that may require further attention or mitigation.

```sql+postgres
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

```sql+sqlite
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