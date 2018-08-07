package routes

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/sinnott74/goblogserver/model"
	"github.com/sinnott74/goblogserver/orm"
)

func UserRouter() chi.Router {
	r := chi.NewRouter()
	r.Get("/", getUsers)
	r.Get("/{id}", getUser)
	return r
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	users := []model.User{}
	orm.SelectAll(ctx, &users, &model.User{})
	render.JSON(w, r, users)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 0)
	if err != nil {
		panic(err)
	}
	ctx := r.Context()
	user := model.User{ID: id}
	err = orm.Get(ctx, &user)
	if err != nil {
		panic(err)
	}
	render.JSON(w, r, user)
}
