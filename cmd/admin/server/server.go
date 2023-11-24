package server

import (
	"GoRestify/internal/core"
	"fmt"
	"log"
	"net/http"

	"time"

	"GoRestify/pkg/middleware"
	"GoRestify/pkg/pkg_config"
	"GoRestify/pkg/pkg_consts"
	"GoRestify/pkg/pkg_err"
	"GoRestify/pkg/response"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

// Start initiate the server
func Start(engine *core.Engine) *gin.Engine {

	var r *gin.Engine

	r = gin.Default()
	r.Use(cors.New(cors.Config{
		//AllowAllOrigins: true,
		AllowOrigins: []string{"*"},
		//AllowOriginFunc: func(origin string) bool {
		//	return origin == origin
		//},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		MaxAge:           1 * time.Hour,
	}))

	if pkg_config.Config.IsDebug {
		r.Use(middleware.APILogger())
	}

	// No Route "Not Found"
	notFoundRoute(r, engine)

	rg := r.Group("/api/admin/v1")
	{
		Route(*rg, engine)
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf("%v:%v", pkg_consts.ServerAddress, engine.Envs[core.Port]),
		Handler:      r,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  121 * time.Second,
	}

	fmt.Printf("\n[%v] ", time.Now().Format(pkg_consts.DateTimeLayoutZone))
	fmt.Printf("starting GoRestify/ server on %v:%v...\n", pkg_consts.ServerAddress, engine.Envs[core.Port])
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("start server failed: %v\n", err)
	}

	return r
}

func notFoundRoute(r *gin.Engine, engine *core.Engine) {
	r.NoRoute(func(c *gin.Context) {
		err := pkg_err.New("route not found", "E1498728").Custom(pkg_err.RouteNotFoundErr).
			Message(pkg_err.PleaseContactToSupport).Build()
		response.New(c).Error(err).JSON()
	})
}
