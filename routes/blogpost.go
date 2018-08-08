package routes

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"github.com/sinnott74/goblogserver/model"
	"github.com/sinnott74/goblogserver/orm"
)

func BlogpostRouter() chi.Router {
	r := chi.NewRouter()
	r.Get("/", getBlogPosts)
	r.Get("/{id}", getBlogPost)
	return r
}

func getBlogPosts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	blogposts := []model.Blogpost{}
	err := orm.SelectAll(ctx, &blogposts, &model.Blogpost{})
	if err != nil {
		panic(err)
	}
	render.JSON(w, r, blogposts)
}

func getBlogPost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 0)
	if err != nil {
		panic(err)
	}
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	blogpost := &model.Blogpost{ID: id}
	err = orm.SelectOne(ctx, blogpost)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
		// panic(err)
	}
	render.JSON(w, r, blogpost)
}
