# Branch Summary: add-cost-tables

## Overview
This branch implements support for Cost Anomaly Detection, Budget Alerts, and Cost Allocation Tag Status tables as requested in GitHub issue #2680.

## Changes Made

### New Table Files Created
1. **aws/table_aws_budgets_budget.go** - AWS Budgets Budget table
   - Supports listing all budgets in the account
   - Supports getting individual budgets by account_id and name
   - Includes budget notifications via hydration
   - Provides budget limit, calculated spend, and cost filter information

2. **aws/table_aws_ce_cost_anomaly_detection.go** - Cost Anomaly Detection table
   - Lists all cost anomaly detection monitors
   - Gets individual monitors by monitor_arn
   - Includes monitor specification, evaluation dates, and status

3. **aws/table_aws_ce_cost_allocation_tags.go** - Cost Allocation Tags table
   - Lists all cost allocation tags
   - Optional filtering by status (Active/Inactive)
   - Includes tag type (UserDefined/AWSGenerated), last used date, and last updated date

### Service Configuration
- **aws/service.go**
  - Added import for `github.com/aws/aws-sdk-go-v2/service/budgets`
  - Added `BudgetsClient()` function for AWS Budgets API access (global service)

### Plugin Registration
- **aws/plugin.go**
  - Registered `aws_budgets_budget` table in TableMap
  - Registered `aws_ce_cost_anomaly_detection` table in TableMap
  - Registered `aws_ce_cost_allocation_tags` table in TableMap

### Documentation Files Created
1. **docs/tables/aws_budgets_budget.md**
   - Complete table usage guide with examples
   - SQL examples for listing budgets, checking spending, finding budgets approaching limits
   - Examples in both PostgreSQL and SQLite dialects

2. **docs/tables/aws_ce_cost_anomaly_detection.md**
   - Complete table usage guide with examples
   - SQL examples for listing monitors, active monitors, evaluation status
   - Examples in both PostgreSQL and SQLite dialects

3. **docs/tables/aws_ce_cost_allocation_tags.md**
   - Complete table usage guide with examples
   - SQL examples for listing tags, filtering by status/type, finding unused tags
   - Examples in both PostgreSQL and SQLite dialects

## Features Implemented

### aws_budgets_budget
- **List Operation**: List all budgets in the account
- **Get Operation**: Get a specific budget by account_id and name
- **Columns**:
  - name, type, limit, limit_unit
  - calculated_spend_actual_spend, calculated_spend_forecasted_spend
  - time_period_start, time_period_end, time_unit
  - cost_filters, cost_types
  - notifications (via hydration)
  - Standard columns: title, akas

### aws_ce_cost_anomaly_detection
- **List Operation**: List all anomaly detection monitors
- **Get Operation**: Get a specific monitor by monitor_arn
- **Columns**:
  - monitor_arn, name, status, frequency
  - monitor_specification
  - creation_date, last_modified_date
  - last_evaluation_date, next_evaluation_date
  - Standard columns: title, akas

### aws_ce_cost_allocation_tags
- **List Operation**: List all cost allocation tags with optional status filtering
- **Columns**:
  - tag_key, tag_type (AWSGenerated/UserDefined)
  - status (Active/Inactive)
  - last_updated_date, last_used_date
  - Standard columns: title

## Use Cases Enabled

1. **Cost Governance**: Monitor budgets and set spending limits across departments
2. **Cost Alerting**: Define notifications for budget overages and cost anomalies
3. **Cost Allocation**: Track and manage which tags are active for cost allocation
4. **Cost Optimization**: Identify unused tags and anomalous spending patterns
5. **FinOps Dashboards**: Build comprehensive cost management dashboards with unified data

## Branch Information
- **Branch Name**: add-cost-tables
- **Created From**: main
- **Commits**: 2 commits
  - 437053df: Add support for Cost Anomaly Detection, Budget Alerts, and Cost Allocation Tag Status tables
  - a5d0d71d: Add documentation for new cost-related tables

## Testing Notes
- All new table implementations follow existing patterns in the codebase
- No linting errors detected
- Tables use appropriate client initialization (BudgetsClient, CostExplorerClient)
- Both tables properly handle pagination via AWS SDK paginators
- Hydration functions implemented for fetching additional details (e.g., budget notifications)
