package middleware

import (
	"fmt"
	"net/http"

	"github.com/sinnott74/goblogserver/orm"
)

func Transaction(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		t, ctx := orm.NewTransaction(r.Context())
		defer func() {
			fmt.Println(w.Header())
			if rec := recover(); rec != nil {
				t.Rollback()
				// Panic to let recoverer handle 500
				panic(rec)
			} else {
				err := t.Commit()
				if err != nil {
					panic(err)
				}
			}
		}()
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
