package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/vincentserpoul/playwithsql/status/islatest"
)

func globalMux(env *localEnv) *httprouter.Router {
	router := httprouter.New()

	router.POST("/entityone/status/islatest", islatest.EntityoneCreateHandler(env.DB, env.IslatestLink))
	router.GET("/entityone/status/islatest", islatest.EntityoneSelectHandler(env.DB, env.IslatestLink))
	router.GET("/entityone/status/islatest/:entityoneID", islatest.EntityoneSelectByPKHandler(env.DB, env.IslatestLink))
	router.DELETE("/entityone/status/islatest/:entityoneID", islatest.EntityoneDeleteByPKHandler(env.DB, env.IslatestLink))

	return router

}
