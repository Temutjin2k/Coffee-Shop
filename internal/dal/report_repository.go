package dal

import (
	"database/sql"
	"fmt"
	"hot-coffee/models"
	"math"

	"github.com/lib/pq"
)

type ReportRespository interface {
	SearchOrders(searchQuery string) ([]models.SearchOrderResult, error)
	SearchMenuItems(searchQuery string, minPrice, maxPrice int) ([]models.SearchMenuItem, error)
}

type ReportRespositoryImpl struct {
	db *sql.DB
}

func NewReportRespository(db *sql.DB) *ReportRespositoryImpl {
	return &ReportRespositoryImpl{db: db}
}

func (repo *ReportRespositoryImpl) SearchOrders(searchQuery string) ([]models.SearchOrderResult, error) {
	query := `
		SELECT 
			ord.ID, 
			ord.CustomerName, 
			ARRAY_AGG(mi.Name) AS items, 
			SUM(mi.Price) AS total,
			ts_rank(
				to_tsvector(ord.CustomerName || ' ' || STRING_AGG(mi.Name, ' ')), 
				websearch_to_tsquery($1)
			) AS relevance
		FROM orders ord
		JOIN order_items oi ON ord.ID = oi.OrderID
		JOIN menu_items mi ON oi.ProductID = mi.ID
		GROUP BY ord.ID, ord.CustomerName
		HAVING to_tsvector(ord.CustomerName || ' ' || STRING_AGG(mi.Name, ' ')) @@ websearch_to_tsquery($1)
		ORDER BY relevance DESC;
	`

	rows, err := repo.db.Query(query, searchQuery)
	if err != nil {
		return nil, fmt.Errorf("error searching %v in orders: %v", searchQuery, err)
	}

	var result []models.SearchOrderResult
	for rows.Next() {
		var item models.SearchOrderResult
		if err := rows.Scan(&item.ID, &item.CustomerName, pq.Array(&item.Items), &item.Total, &item.Relevance); err != nil {
			return nil, err
		}
		item.Relevance = math.Round(item.Relevance*100) / 100
		result = append(result, item)
	}
	return result, nil
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
		item.Relevance = math.Round(item.Relevance*100) / 100
		result = append(result, item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
