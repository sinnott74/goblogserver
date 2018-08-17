package routes

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"github.com/sinnott74/goblogserver/model"
	"github.com/sinnott74/goblogserver/orm"
)

// BlogpostRouter creates the Blogpost routes
func BlogpostRouter() chi.Router {
	r := chi.NewRouter()
	r.Get("/", getBlogPosts)
	r.Get("/{id}", getBlogPost)
	return r
}

func getBlogPosts(w http.ResponseWriter, r *http.Request) {
	blogposts := []model.Blogpost{}
	err := orm.SelectAll(r.Context(), &blogposts, &model.Blogpost{})
	if err != nil {
		panic(err)
	}
	render.JSON(w, r, blogposts)
}

func getBlogPost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 0)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	blogpost := &model.Blogpost{ID: id}
	err = orm.Get(r.Context(), blogpost)
	if err != nil {
		panic(err)
	}
	render.JSON(w, r, blogpost)
}
