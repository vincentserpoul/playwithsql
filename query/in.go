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

// InQueryNamedParams returns a param string for IN queries
func InQueryNamedParams(
	lenParams int,
	nameParam string,
) (inQuery string) {
	if lenParams == 0 {
		return inQuery
	}
	if nameParam == "" {
		return inQuery
	}

	inQuery += `(`
	var queryParams []string
	for i := 0; i < lenParams; i++ {
		queryParams = append(queryParams, `:`+nameParam)
	}
	inQuery += strings.Join(queryParams, `,`)
	inQuery += `)`

	return inQuery
}
