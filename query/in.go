package query

import (
	"strconv"
	"strings"
)

// InQueryParams returns a param string for IN queries
func InQueryParams(
	lenParams int,
	IsParamQuestionMark bool,
	indexParam int,
) (inQuery string) {
	if lenParams == 0 {
		return inQuery
	}

	inQuery += `(`
	var queryParams []string
	for i := 0; i < lenParams; i++ {
		var queryParam string
		if IsParamQuestionMark {
			queryParam = `?`
		} else {
			queryParam = `$` + strconv.Itoa(indexParam+1)
			indexParam++
		}
		queryParams = append(queryParams, queryParam)
	}
	inQuery += strings.Join(queryParams, `,`)
	inQuery += `)`

	return inQuery
}
