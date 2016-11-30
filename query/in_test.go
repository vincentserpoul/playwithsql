package query

import "testing"

func TestInQueryParams(t *testing.T) {
	cases := []struct {
		LenParams           int
		IsParamQuestionMark bool
		IndexParam          int
		ExpQFilter          string
	}{
		{LenParams: 0, IsParamQuestionMark: true, IndexParam: 0, ExpQFilter: ""},
		{LenParams: 3, IsParamQuestionMark: true, IndexParam: 0, ExpQFilter: "(?,?,?)"},
		{LenParams: 3, IsParamQuestionMark: true, IndexParam: 1, ExpQFilter: "(?,?,?)"},
		{LenParams: 3, IsParamQuestionMark: false, IndexParam: 0, ExpQFilter: "($1,$2,$3)"},
		{LenParams: 3, IsParamQuestionMark: false, IndexParam: 1, ExpQFilter: "($2,$3,$4)"},
	}

	for _, c := range cases {
		foundQuery := InQueryParams(c.LenParams, c.IsParamQuestionMark, c.IndexParam)
		if foundQuery != c.ExpQFilter {
			t.Errorf("InQueryParams(%d, %t, %d) should have returned %s, but returned %s\n",
				c.LenParams, c.IsParamQuestionMark, c.IndexParam, c.ExpQFilter, foundQuery,
			)
			return
		}
	}
}

func BenchmarkInQueryParams(b *testing.B) {
	for i := 0; i < b.N; i++ {
		InQueryParams(0, true, 1)
	}
}
