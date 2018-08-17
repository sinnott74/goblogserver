package routes

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"github.com/sinnott74/goblogserver/auth"
	"github.com/sinnott74/goblogserver/model"
)

// AuthRouter creates the Authentication routes
func AuthRouter() chi.Router {
	r := chi.NewRouter()
	r.Post("/login", login)
	r.With(auth.Middleware).Get("/test", test)
	return r
}

type signUpRequest struct {
	model.User
	model.Credential
}

func test(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, "It worked")
}

// func signUp(c *gin.Context) {
// 	ctx := c.Request.Context()
// 	var signUpRequest signUpRequest
// 	c.BindJSON(&signUpRequest)
// 	err := orm.Insert(ctx, &signUpRequest.User)
// 	if err != nil {
// 		c.AbortWithError(http.StatusInternalServerError, err)
// 		return
// 	}
// 	signUpRequest.UserID = signUpRequest.User.ID
// 	err = orm.Insert(ctx, &signUpRequest.Credential)
// 	fmt.Printf("%+v\n", signUpRequest)
// 	if err != nil {
// 		c.AbortWithError(http.StatusInternalServerError, err)
// 		return
// 	}
// 	c.JSON(http.StatusOK, "signUpRequest")
// }

type loginRequest struct {
	Username string
	Password string
}

func login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	loginRequest := &loginRequest{}
	render.DecodeJSON(r.Body, loginRequest)
	userToken, err := auth.Login(ctx, loginRequest.Username, loginRequest.Password)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	render.JSON(w, r, userToken)
}
