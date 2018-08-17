package middleware

import (
	"net/http"

	"github.com/sinnott74/goblogserver/orm"
)

// Transaction is middle ware that begins a transaction & adds it onto the request's context
// It will rollback the transaction is the requets panics, or commit otherwise
func Transaction(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		t, ctx := orm.NewTransaction(r.Context())
		defer func() {
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
