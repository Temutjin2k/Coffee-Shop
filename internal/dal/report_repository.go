package dal

import (
	"database/sql"
	"fmt"
	"hot-coffee/models"
)

type ReportRespository interface {
	SearchOrders(searchQuery string, minPrice, maxPrice int) ([]models.SearchOrderResult, error)
	SearchMenuItems(searchQuery string, minPrice, maxPrice int) ([]models.SearchMenuItem, error)
}

type ReportRespositoryImpl struct {
	db *sql.DB
}

func NewReportRespository(db *sql.DB) *ReportRespositoryImpl {
	return &ReportRespositoryImpl{db: db}
}

func (repo *ReportRespositoryImpl) SearchOrders(searchQuery string, minPrice, maxPrice int) ([]models.SearchOrderResult, error) {
	return nil, nil
}

func (repo *ReportRespositoryImpl) SearchMenuItems(searchQuery string, minPrice, maxPrice int) ([]models.SearchMenuItem, error) {
	query := `
		SELECT 
			id, name, description, price,
			ts_rank(to_tsvector(name || ' ' || COALESCE(description, '')), websearch_to_tsquery($1)) as relevance
		FROM menu_items
		WHERE to_tsvector(name || ' ' || COALESCE(description, '')) @@ websearch_to_tsquery($1)
	`
	args := []interface{}{searchQuery}
	argIndex := 2

	if minPrice != -1 {
		query += fmt.Sprintf(" AND price >= $%d", argIndex)
		args = append(args, minPrice)
		argIndex++
	}

	if maxPrice != -1 {
		query += fmt.Sprintf(" AND price <= $%d", argIndex)
		args = append(args, maxPrice)
		argIndex++
	}

	query += " ORDER BY relevance DESC;"

	rows, err := repo.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error searching %v in menu items: %v", searchQuery, err)
	}
	defer rows.Close()

	var result []models.SearchMenuItem

	for rows.Next() {
		var item models.SearchMenuItem
		if err := rows.Scan(&item.ID, &item.Name, &item.Description, &item.Price, &item.Relevance); err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
