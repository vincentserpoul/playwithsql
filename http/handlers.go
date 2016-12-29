package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/vincentserpoul/playwithsql/status"
)

func globalMux(env *localEnv) *httprouter.Router {
	router := httprouter.New()

	router.POST("/entityone/status/islatest", status.EntityoneCreateHandler(env.DB, env.IslatestLink))
	router.GET("/entityone/status/islatest", status.EntityoneSelectHandler(env.DB, env.IslatestLink))
	router.GET("/entityone/status/islatest/:entityoneID", status.EntityoneSelectByPKHandler(env.DB, env.IslatestLink))
	router.DELETE("/entityone/status/islatest/:entityoneID", status.EntityoneDeleteByPKHandler(env.DB, env.IslatestLink))

	return router

}
