package config

import (
	"errors"
	"fmt"
)

func BuildSelectQuery(columns []string, tableName string, conditions string) (string, error) {
	if len(columns) == 0 {
		return "", errors.New("no parameters provided")
	}

	query := "SELECT "

	if len(columns) == 1 && columns[0] == "*" {
		query += "* FROM "
	} else {
		for i, column := range columns {
			if i == len(columns)-1 {
				query += column + " FROM "
			} else {
				query += column + ", "
			}
		}
	}

	query += tableName + " "

	if conditions != "" {
		query += "WHERE " + conditions + ";"
	}

	return query, nil
}

func BuildInsertQuery(tableName string, columns []string) (string, error) {
	if len(columns) == 0 {
		return "", errors.New("no columns provided")
	}

	query := fmt.Sprintf("INSERT INTO %s (", tableName)
	valuesPlaceholder := "VALUES ("

	for i, column := range columns {
		if i == len(columns)-1 {
			query += column + ") "
			valuesPlaceholder += fmt.Sprintf("$%d)", i+1)
		} else {
			query += column + ", "
			valuesPlaceholder += fmt.Sprintf("$%d, ", i+1)
		}
	}

	query += valuesPlaceholder + ";"
	return query, nil
}

func BuildUpdateQuery(tableName string, columns []string, condition string) (string, error) {
	if len(columns) == 0 {
		return "", errors.New("no columns provided")
	}

	query := fmt.Sprintf("UPDATE %s SET ", tableName)

	for i, column := range columns {
		if i == len(columns)-1 {
			query += fmt.Sprintf("%s=$%d ", column, i+1)
		} else {
			query += fmt.Sprintf("%s=$%d, ", column, i+1)
		}
	}

	if condition != "" {
		query += "WHERE " + condition + ";"
	}

	return query, nil
}

func BuildDeleteQuery(tableName string, condition string) (string, error) {
	if tableName == "" {
		return "", errors.New("table name is required")
	}

	query := fmt.Sprintf("DELETE FROM %s ", tableName)

	if condition != "" {
		query += "WHERE " + condition + ";"
	} else {
		return "", errors.New("condition is required for delete query")
	}

	return query, nil
}
