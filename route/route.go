package route

import (
	"net/http"

	cntrl "github.com/tarkov-database/rest-api/controller"
	"github.com/tarkov-database/rest-api/middleware/jwt"

	"github.com/julienschmidt/httprouter"
)

const prefix = "/v2"

func Load() *httprouter.Router {
	return routes()
}

func routes() *httprouter.Router {
	r := httprouter.New()

	// Index
	r.GET(prefix, cntrl.IndexGET)
	r.Handler("GET", "/", http.RedirectHandler(prefix, http.StatusMovedPermanently))

	// Item
	r.GET(prefix+"/item", auth(jwt.ScopeItemRead, cntrl.ItemIndexGET))
	r.GET(prefix+"/item/:kind", auth(jwt.ScopeItemRead, cntrl.ItemsGET))
	r.GET(prefix+"/item/:kind/:id", auth(jwt.ScopeItemRead, cntrl.ItemGET))
	r.POST(prefix+"/item/:kind", auth(jwt.ScopeItemWrite, cntrl.ItemPOST))
	r.PUT(prefix+"/item/:kind/:id", auth(jwt.ScopeItemWrite, cntrl.ItemPUT))
	r.DELETE(prefix+"/item/:id", auth(jwt.ScopeItemWrite, cntrl.ItemDELETE))

	// User
	r.GET(prefix+"/user", auth(jwt.ScopeUserRead, cntrl.UsersGET))
	r.GET(prefix+"/user/:id", auth(jwt.ScopeUserRead, cntrl.UserGET))
	r.POST(prefix+"/user", auth(jwt.ScopeUserWrite, cntrl.UserPOST))
	r.PUT(prefix+"/user/:id", auth(jwt.ScopeUserWrite, cntrl.UserPUT))
	r.DELETE(prefix+"/user/:id", auth(jwt.ScopeUserWrite, cntrl.UserDELETE))

	// Token
	r.GET(prefix+"/token", cntrl.TokenGET)
	r.POST(prefix+"/token", cntrl.TokenPOST)

	r.NotFound = cntrl.StatusNotFoundHandler()

	r.RedirectTrailingSlash = true
	r.HandleOPTIONS = true

	return r
}

func auth(s string, h httprouter.Handle) httprouter.Handle {
	return jwt.AuhtorizationHandler(s, h)
}
