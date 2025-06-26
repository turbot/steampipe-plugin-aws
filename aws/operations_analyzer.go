package aws

import (
	"context"
	"fmt"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"sync"

	pg_query "github.com/pganalyze/pg_query_go/v6"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// PluginAnalysisTools holds the cached analysis data for the plugin.
type PluginAnalysisTools struct {
	tableMap           map[string]*plugin.Table
	listFuncToTableMap map[string]*plugin.Table
}

var (
	pluginTools *PluginAnalysisTools
	once        sync.Once
)

// getPluginTools builds and caches a map of all list functions to their tables.
// This allows for dynamic resolution of ParentHydrate dependencies at runtime.
func getPluginTools(ctx context.Context) *PluginAnalysisTools {
	once.Do(func() {
		// Get the plugin's definition, which contains all table information.
		p := Plugin(ctx)
		listFuncMap := make(map[string]*plugin.Table)

		// Iterate through every table in the plugin.
		for _, table := range p.TableMap {
			// If the table has a List function, map its name to the table definition.
			if table.List != nil && table.List.Hydrate != nil {
				funcName := getFunctionName(table.List.Hydrate)
				listFuncMap[funcName] = table
			}
		}

		pluginTools = &PluginAnalysisTools{
			tableMap:           p.TableMap,
			listFuncToTableMap: listFuncMap,
		}
	})
	return pluginTools
}

// getFunctionName returns the short name of a function (e.g., "listWellArchitectedWorkloads").
func getFunctionName(i interface{}) string {
	name := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	parts := strings.Split(name, ".")
	return parts[len(parts)-1]
}

// AnalyzeTableOperations inspects a Steampipe table definition, follows its hydration
// patterns, and returns a sorted list of all unique AWS operations that could be called.
func AnalyzeTableOperations(ctx context.Context, table *plugin.Table) ([]string, error) {
	// Use a map to store unique operations.
	operations := make(map[string]struct{})

	// Queue for tables to analyze, to handle ParentHydrate recursively.
	queue := []*plugin.Table{table}

	// Keep track of visited tables to avoid loops.
	visited := make(map[string]bool)

	for len(queue) > 0 {
		// Dequeue the next table to analyze.
		currentTable := queue[0]
		queue = queue[1:]

		if visited[currentTable.Name] {
			continue
		}
		visited[currentTable.Name] = true

		// 1. GetConfig
		if get := currentTable.Get; get != nil {
			if service, ok := get.Tags["service"]; ok {
				if action, ok := get.Tags["action"]; ok {
					op := service + ":" + action
					operations[op] = struct{}{}
				}
			}
		}

		// 2. ListConfig
		if list := currentTable.List; list != nil {
			if service, ok := list.Tags["service"]; ok {
				if action, ok := list.Tags["action"]; ok {
					op := service + ":" + action
					operations[op] = struct{}{}
				}
			}

			// 3. ParentHydrate
			if list.ParentHydrate != nil {
				funcName := getFunctionName(list.ParentHydrate)
				tools := getPluginTools(ctx)
				if parentTable, ok := tools.listFuncToTableMap[funcName]; ok {
					queue = append(queue, parentTable)
				}
			}
		}

		// 4. HydrateConfig
		for _, hc := range currentTable.HydrateConfig {
			if service, ok := hc.Tags["service"]; ok {
				if action, ok := hc.Tags["action"]; ok {
					op := service + ":" + action
					operations[op] = struct{}{}
				}
			}
		}
	}

	// Convert map to a sorted slice for consistent output.
	result := make([]string, 0, len(operations))
	for op := range operations {
		result = append(result, op)
	}
	sort.Strings(result)

	return result, nil
}

// AnalyzeSQLQuery takes a SQL query string and returns a map of table names to their AWS operations.
// It extracts table names from the SQL query and then analyzes each table to determine
// what AWS operations could be triggered.
func AnalyzeSQLQuery(ctx context.Context, sqlQuery string) (map[string][]string, error) {
	// Extract table names from the SQL query
	tableNames, err := extractTableNames(sqlQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to extract table names from SQL: %w", err)
	}

	// Map to store table name -> operations
	result := make(map[string][]string)

	// Analyze each table to get its operations
	for _, tableName := range tableNames {
		// Get the table definition by its name.
		table := getTableByName(ctx, tableName)
		if table == nil {
			// Table not found in this plugin, skip it.
			continue
		}

		// Analyze the table's operations.
		operations, err := AnalyzeTableOperations(ctx, table)
		if err != nil {
			return nil, fmt.Errorf("failed to analyze operations for table %s: %w", tableName, err)
		}

		result[tableName] = operations
	}

	return result, nil
}

// GetAWSOperationsFromSQL takes a SQL query and returns a flat list of all unique AWS operations
// that could be triggered by executing that query. This is a simplified version of AnalyzeSQLQuery
// that returns just the operations without the table grouping.
func GetAWSOperationsFromSQL(ctx context.Context, sqlQuery string) ([]string, error) {
	// Get the detailed analysis
	tableOperations, err := AnalyzeSQLQuery(ctx, sqlQuery)
	if err != nil {
		return nil, err
	}

	// Collect all unique operations
	operationSet := make(map[string]struct{})
	for _, operations := range tableOperations {
		for _, operation := range operations {
			operationSet[operation] = struct{}{}
		}
	}

	// Convert to sorted slice
	result := make([]string, 0, len(operationSet))
	for operation := range operationSet {
		result = append(result, operation)
	}
	sort.Strings(result)

	return result, nil
}

// extractTableNames extracts table names from a SQL query using the PostgreSQL parser.
// This provides proper PostgreSQL syntax parsing that handles complex queries, subqueries, and JOINs correctly.
func extractTableNames(sqlQuery string) ([]string, error) {
	// Parse the SQL statement using PostgreSQL parser
	result, err := pg_query.Parse(sqlQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SQL query: %w", err)
	}

	tableSet := make(map[string]struct{})

	// Extract tables from all statements
	for _, stmt := range result.Stmts {
		extractTablesFromStatement(stmt.Stmt, tableSet)
	}

	// Convert set to slice and filter for AWS tables
	var tables []string
	for table := range tableSet {
		// Only include AWS Steampipe tables (those starting with "aws_")
		if strings.HasPrefix(table, "aws_") {
			tables = append(tables, table)
		}
	}

	sort.Strings(tables)
	return tables, nil
}

// extractTablesFromStatement extracts table names from a parsed statement
func extractTablesFromStatement(stmt *pg_query.Node, tableSet map[string]struct{}) {
	if stmt == nil {
		return
	}

	switch node := stmt.Node.(type) {
	case *pg_query.Node_SelectStmt:
		extractTablesFromSelectStmt(node.SelectStmt, tableSet)
	case *pg_query.Node_InsertStmt:
		extractTablesFromInsertStmt(node.InsertStmt, tableSet)
	case *pg_query.Node_UpdateStmt:
		extractTablesFromUpdateStmt(node.UpdateStmt, tableSet)
	case *pg_query.Node_DeleteStmt:
		extractTablesFromDeleteStmt(node.DeleteStmt, tableSet)
	}
}

// extractTablesFromSelectStmt extracts table names from a SELECT statement
func extractTablesFromSelectStmt(stmt *pg_query.SelectStmt, tableSet map[string]struct{}) {
	if stmt == nil {
		return
	}

	// Extract from WITH clause (CTEs - Common Table Expressions)
	if stmt.WithClause != nil {
		for _, cte := range stmt.WithClause.Ctes {
			extractTablesFromNode(cte, tableSet)
		}
	}

	// Extract from FROM clause
	for _, fromClause := range stmt.FromClause {
		extractTablesFromRangeVar(fromClause, tableSet)
	}

	// Extract from WHERE clause (subqueries)
	if stmt.WhereClause != nil {
		extractTablesFromNode(stmt.WhereClause, tableSet)
	}

	// Extract from HAVING clause (subqueries)
	if stmt.HavingClause != nil {
		extractTablesFromNode(stmt.HavingClause, tableSet)
	}

	// Extract from target list (SELECT expressions with subqueries)
	for _, target := range stmt.TargetList {
		extractTablesFromNode(target, tableSet)
	}

	// Extract from ORDER BY clause
	for _, sortBy := range stmt.SortClause {
		extractTablesFromNode(sortBy, tableSet)
	}

	// Handle UNION/INTERSECT/EXCEPT operations
	if stmt.Larg != nil {
		extractTablesFromSelectStmt(stmt.Larg, tableSet)
	}
	if stmt.Rarg != nil {
		extractTablesFromSelectStmt(stmt.Rarg, tableSet)
	}
}

// extractTablesFromInsertStmt extracts table names from an INSERT statement
func extractTablesFromInsertStmt(stmt *pg_query.InsertStmt, tableSet map[string]struct{}) {
	if stmt == nil {
		return
	}

	// Main table being inserted into
	if stmt.Relation != nil {
		extractTablesFromRangeVar(&pg_query.Node{Node: &pg_query.Node_RangeVar{RangeVar: stmt.Relation}}, tableSet)
	}

	// Extract from SELECT query if it's INSERT ... SELECT
	if stmt.SelectStmt != nil {
		extractTablesFromStatement(stmt.SelectStmt, tableSet)
	}
}

// extractTablesFromUpdateStmt extracts table names from an UPDATE statement
func extractTablesFromUpdateStmt(stmt *pg_query.UpdateStmt, tableSet map[string]struct{}) {
	if stmt == nil {
		return
	}

	// Main table being updated
	if stmt.Relation != nil {
		extractTablesFromRangeVar(&pg_query.Node{Node: &pg_query.Node_RangeVar{RangeVar: stmt.Relation}}, tableSet)
	}

	// Extract from FROM clause (PostgreSQL allows FROM in UPDATE)
	for _, fromClause := range stmt.FromClause {
		extractTablesFromRangeVar(fromClause, tableSet)
	}

	// Extract from WHERE clause
	if stmt.WhereClause != nil {
		extractTablesFromNode(stmt.WhereClause, tableSet)
	}

	// Extract from SET expressions
	for _, target := range stmt.TargetList {
		extractTablesFromNode(target, tableSet)
	}
}

// extractTablesFromDeleteStmt extracts table names from a DELETE statement
func extractTablesFromDeleteStmt(stmt *pg_query.DeleteStmt, tableSet map[string]struct{}) {
	if stmt == nil {
		return
	}

	// Main table being deleted from
	if stmt.Relation != nil {
		extractTablesFromRangeVar(&pg_query.Node{Node: &pg_query.Node_RangeVar{RangeVar: stmt.Relation}}, tableSet)
	}

	// Extract from USING clause (PostgreSQL allows USING in DELETE)
	for _, usingClause := range stmt.UsingClause {
		extractTablesFromRangeVar(usingClause, tableSet)
	}

	// Extract from WHERE clause
	if stmt.WhereClause != nil {
		extractTablesFromNode(stmt.WhereClause, tableSet)
	}
}

// extractTablesFromRangeVar extracts table names from range variables (table references)
func extractTablesFromRangeVar(node *pg_query.Node, tableSet map[string]struct{}) {
	if node == nil {
		return
	}

	switch rangeNode := node.Node.(type) {
	case *pg_query.Node_RangeVar:
		if rangeNode.RangeVar != nil && rangeNode.RangeVar.Relname != "" {
			// Convert to lowercase for consistency
			tableName := strings.ToLower(rangeNode.RangeVar.Relname)
			tableSet[tableName] = struct{}{}
		}
	case *pg_query.Node_RangeSubselect:
		if rangeNode.RangeSubselect != nil && rangeNode.RangeSubselect.Subquery != nil {
			extractTablesFromStatement(rangeNode.RangeSubselect.Subquery, tableSet)
		}
	case *pg_query.Node_JoinExpr:
		if rangeNode.JoinExpr != nil {
			// Extract from left and right sides of the join
			extractTablesFromRangeVar(rangeNode.JoinExpr.Larg, tableSet)
			extractTablesFromRangeVar(rangeNode.JoinExpr.Rarg, tableSet)
			// Extract from join condition
			if rangeNode.JoinExpr.Quals != nil {
				extractTablesFromNode(rangeNode.JoinExpr.Quals, tableSet)
			}
		}
	case *pg_query.Node_RangeFunction:
		// Handle function calls in FROM clause - they may contain subqueries
		if rangeNode.RangeFunction != nil {
			for _, funcCall := range rangeNode.RangeFunction.Functions {
				extractTablesFromNode(funcCall, tableSet)
			}
		}
	}
}

// extractTablesFromNode extracts table names from any node type (recursive for subqueries)
func extractTablesFromNode(node *pg_query.Node, tableSet map[string]struct{}) {
	if node == nil {
		return
	}

	switch nodeType := node.Node.(type) {
	case *pg_query.Node_CommonTableExpr:
		// Handle Common Table Expressions (CTEs)
		if nodeType.CommonTableExpr != nil && nodeType.CommonTableExpr.Ctequery != nil {
			extractTablesFromStatement(nodeType.CommonTableExpr.Ctequery, tableSet)
		}
	case *pg_query.Node_SubLink:
		// Handle subqueries in expressions
		if nodeType.SubLink != nil && nodeType.SubLink.Subselect != nil {
			extractTablesFromStatement(nodeType.SubLink.Subselect, tableSet)
		}
	case *pg_query.Node_AExpr:
		// Handle expressions (AND, OR, etc.)
		if nodeType.AExpr != nil {
			if nodeType.AExpr.Lexpr != nil {
				extractTablesFromNode(nodeType.AExpr.Lexpr, tableSet)
			}
			if nodeType.AExpr.Rexpr != nil {
				extractTablesFromNode(nodeType.AExpr.Rexpr, tableSet)
			}
		}
	case *pg_query.Node_BoolExpr:
		// Handle boolean expressions
		if nodeType.BoolExpr != nil {
			for _, arg := range nodeType.BoolExpr.Args {
				extractTablesFromNode(arg, tableSet)
			}
		}
	case *pg_query.Node_FuncCall:
		// Handle function calls that may contain subqueries
		if nodeType.FuncCall != nil {
			for _, arg := range nodeType.FuncCall.Args {
				extractTablesFromNode(arg, tableSet)
			}
		}
	case *pg_query.Node_CaseExpr:
		// Handle CASE expressions
		if nodeType.CaseExpr != nil {
			if nodeType.CaseExpr.Arg != nil {
				extractTablesFromNode(nodeType.CaseExpr.Arg, tableSet)
			}
			if nodeType.CaseExpr.Defresult != nil {
				extractTablesFromNode(nodeType.CaseExpr.Defresult, tableSet)
			}
			for _, when := range nodeType.CaseExpr.Args {
				extractTablesFromNode(when, tableSet)
			}
		}
	case *pg_query.Node_CaseWhen:
		// Handle WHEN clauses in CASE expressions
		if nodeType.CaseWhen != nil {
			if nodeType.CaseWhen.Expr != nil {
				extractTablesFromNode(nodeType.CaseWhen.Expr, tableSet)
			}
			if nodeType.CaseWhen.Result != nil {
				extractTablesFromNode(nodeType.CaseWhen.Result, tableSet)
			}
		}
	case *pg_query.Node_ResTarget:
		// Handle target list items (SELECT expressions)
		if nodeType.ResTarget != nil && nodeType.ResTarget.Val != nil {
			extractTablesFromNode(nodeType.ResTarget.Val, tableSet)
		}
	case *pg_query.Node_SortBy:
		// Handle ORDER BY expressions
		if nodeType.SortBy != nil && nodeType.SortBy.Node != nil {
			extractTablesFromNode(nodeType.SortBy.Node, tableSet)
		}
	}
}

// getTableByName returns the table definition for a given table name by looking it
// up in the plugin's table map.
func getTableByName(ctx context.Context, tableName string) *plugin.Table {
	tools := getPluginTools(ctx)
	// The table map from the plugin uses lowercase keys.
	return tools.tableMap[strings.ToLower(tableName)]
}
