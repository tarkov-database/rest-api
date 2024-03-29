package route

import (
	"net/http"

	cntrl "github.com/tarkov-database/rest-api/controller"
	"github.com/tarkov-database/rest-api/middleware/jwt"

	"github.com/julienschmidt/httprouter"
)

const prefix = "/v2"

// Load returns a router with defined routes
func Load() *httprouter.Router {
	return routes()
}

func routes() *httprouter.Router {
	r := httprouter.New()

	// Index
	r.GET(prefix, cntrl.IndexGET)
	r.Handler("GET", "/", http.RedirectHandler(prefix, http.StatusMovedPermanently))

	// Health
	r.GET(prefix+"/health", auth("", cntrl.HealthGET))

	// Item
	r.GET(prefix+"/item", auth(jwt.ScopeItemRead, cntrl.ItemIndexGET))
	r.GET(prefix+"/item/:kind", auth(jwt.ScopeItemRead, cntrl.ItemsGET))
	r.GET(prefix+"/item/:kind/:id", auth(jwt.ScopeItemRead, cntrl.ItemGET))
	r.POST(prefix+"/item/:kind", auth(jwt.ScopeItemWrite, cntrl.ItemPOST))
	r.PUT(prefix+"/item/:kind/:id", auth(jwt.ScopeItemWrite, cntrl.ItemPUT))
	r.DELETE(prefix+"/item/:id", auth(jwt.ScopeItemWrite, cntrl.ItemDELETE))

	// Hideout module
	r.GET(prefix+"/hideout/module", auth(jwt.ScopeHideoutRead, cntrl.ModulesGET))
	r.GET(prefix+"/hideout/module/:id", auth(jwt.ScopeHideoutRead, cntrl.ModuleGET))
	r.POST(prefix+"/hideout/module", auth(jwt.ScopeHideoutWrite, cntrl.ModulePOST))
	r.PUT(prefix+"/hideout/module/:id", auth(jwt.ScopeHideoutWrite, cntrl.ModulePUT))
	r.DELETE(prefix+"/hideout/module/:id", auth(jwt.ScopeHideoutWrite, cntrl.ModuleDELETE))

	// Hideout production
	r.GET(prefix+"/hideout/production", auth(jwt.ScopeHideoutRead, cntrl.ProductionsGET))
	r.GET(prefix+"/hideout/production/:id", auth(jwt.ScopeHideoutRead, cntrl.ProductionGET))
	r.POST(prefix+"/hideout/production", auth(jwt.ScopeHideoutWrite, cntrl.ProductionPOST))
	r.PUT(prefix+"/hideout/production/:id", auth(jwt.ScopeHideoutWrite, cntrl.ProductionPUT))
	r.DELETE(prefix+"/hideout/production/:id", auth(jwt.ScopeHideoutWrite, cntrl.ProductionDELETE))

	// Location
	r.GET(prefix+"/location", auth(jwt.ScopeLocationRead, cntrl.LocationsGET))
	r.GET(prefix+"/location/:id", auth(jwt.ScopeLocationRead, cntrl.LocationGET))
	r.POST(prefix+"/location", auth(jwt.ScopeLocationWrite, cntrl.LocationPOST))
	r.PUT(prefix+"/location/:id", auth(jwt.ScopeLocationWrite, cntrl.LocationPUT))
	r.DELETE(prefix+"/location/:id", auth(jwt.ScopeLocationWrite, cntrl.LocationDELETE))

	// Location feature
	r.GET(prefix+"/location/:id/feature", auth(jwt.ScopeLocationRead, cntrl.FeaturesGET))
	r.GET(prefix+"/location/:id/feature/:fid", auth(jwt.ScopeLocationRead, cntrl.FeatureGET))
	r.POST(prefix+"/location/:id/feature", auth(jwt.ScopeLocationWrite, cntrl.FeaturePOST))
	r.PUT(prefix+"/location/:id/feature/:fid", auth(jwt.ScopeLocationWrite, cntrl.FeaturePUT))
	r.DELETE(prefix+"/location/:id/feature/:fid", auth(jwt.ScopeLocationWrite, cntrl.FeatureDELETE))

	// Location feature group
	r.GET(prefix+"/location/:id/featuregroup", auth(jwt.ScopeLocationRead, cntrl.FeatureGroupsGET))
	r.GET(prefix+"/location/:id/featuregroup/:gid", auth(jwt.ScopeLocationRead, cntrl.FeatureGroupGET))
	r.POST(prefix+"/location/:id/featuregroup", auth(jwt.ScopeLocationWrite, cntrl.FeatureGroupPOST))
	r.PUT(prefix+"/location/:id/featuregroup/:gid", auth(jwt.ScopeLocationWrite, cntrl.FeatureGroupPUT))
	r.DELETE(prefix+"/location/:id/featuregroup/:gid", auth(jwt.ScopeLocationWrite, cntrl.FeatureGroupDELETE))

	// Ammunition distance statistics
	r.GET(prefix+"/statistic/ammunition/distance", auth(jwt.ScopeStatisticRead, cntrl.DistanceStatsGET))
	r.GET(prefix+"/statistic/ammunition/distance/:id", auth(jwt.ScopeStatisticRead, cntrl.DistanceStatGET))
	r.POST(prefix+"/statistic/ammunition/distance", auth(jwt.ScopeStatisticWrite, cntrl.DistanceStatPOST))
	r.PUT(prefix+"/statistic/ammunition/distance/:id", auth(jwt.ScopeStatisticWrite, cntrl.DistanceStatPUT))
	r.DELETE(prefix+"/statistic/ammunition/distance/:id", auth(jwt.ScopeStatisticWrite, cntrl.DistanceStatDELETE))

	// Ammunition armor statistics
	r.GET(prefix+"/statistic/ammunition/armor", auth(jwt.ScopeStatisticRead, cntrl.ArmorStatsGET))
	r.GET(prefix+"/statistic/ammunition/armor/:id", auth(jwt.ScopeStatisticRead, cntrl.ArmorStatGET))
	r.POST(prefix+"/statistic/ammunition/armor", auth(jwt.ScopeStatisticWrite, cntrl.ArmorStatPOST))
	r.PUT(prefix+"/statistic/ammunition/armor/:id", auth(jwt.ScopeStatisticWrite, cntrl.ArmorStatPUT))
	r.DELETE(prefix+"/statistic/ammunition/armor/:id", auth(jwt.ScopeStatisticWrite, cntrl.ArmorStatDELETE))

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
