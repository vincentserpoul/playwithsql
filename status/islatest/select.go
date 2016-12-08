package islatest

import (
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/vincentserpoul/playwithsql/query"
)

// SelectEntity is a closure over the parametrized queries
// it returns a function that will retrieve a slice of entityones
func SelectEntity(isParamQuestionMark bool) func(
	q sqlx.Queryer,
	entityIDs []int64,
	isStatusIDs []int,
	notStatusIDs []int,
	neverStatusIDs []int,
	hasStatusIDs []int,
	limit int,
) (*sqlx.Rows, error) {
	return func(
		q sqlx.Queryer,
		entityIDs []int64,
		isStatusIDs []int,
		notStatusIDs []int,
		neverStatusIDs []int,
		hasStatusIDs []int,
		limit int,
	) (*sqlx.Rows, error) {

		query := `
            SELECT
                e.entityone_id, e.time_created,
                es.entityone_id as status_entityone_id, es.action_id, es.status_id, es.time_created as status_time_created
            FROM entityone e
            INNER JOIN entityone_status es ON es.entityone_id = e.entityone_id
                AND es.is_latest = 1
            WHERE 0 = 0
        `

		params, queryFilter := getFilterSelectEntityOneQuery(
			isParamQuestionMark,
			entityIDs,
			isStatusIDs,
		)

		query += queryFilter

		if limit > 0 {
			limitStr := ` LIMIT ` + strconv.Itoa(limit)
			query += limitStr
		}

		return q.Queryx(query, params...)

	}
}

func getFilterSelectEntityOneQuery(
	isParamQuestionMark bool,
	entityIDs []int64,
	isStatusIDs []int,
) (params []interface{}, queryFilter string) {

	i := 0

	if len(entityIDs) > 0 {
		queryFilter += ` AND e.entityone_id IN `
		queryFilter += query.InQueryParams(len(entityIDs), isParamQuestionMark, i)
		for _, param := range entityIDs {
			params = append(params, param)
			i++
		}
	}

	if len(isStatusIDs) > 0 {
		queryFilter += `  AND es.status_id IN `
		queryFilter += query.InQueryParams(len(isStatusIDs), isParamQuestionMark, i)
		for _, param := range isStatusIDs {
			params = append(params, param)
			i++
		}
	}

	return params, queryFilter
}
