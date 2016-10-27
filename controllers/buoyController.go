// Copyright 2013 Ardan Studios. All rights reserved.
// Use of controller source code is governed by a BSD-style
// license that can be found in the LICENSE handle.

// Package controllers implements the controller layer for the buoy API.
package controllers

import (
	bc "github.com/goinggo/beego-mgo/controllers/baseController"
	"github.com/goinggo/beego-mgo/services/buoyService"
	"github.com/goinggo/beego-mgo/models/buoyModels"
	"github.com/goinggo/beego-mgo/utilities/mongo"
	// "encoding/json"
	log "github.com/goinggo/tracelog"
)

//** TYPES

// BuoyController manages the API for buoy related functionality.
type BuoyController struct {
	bc.BaseController
}

//** WEB FUNCTIONS

// Index is the initial view for the buoy system.
func (controller *BuoyController) Index() {
	buoyStations, err := buoyService.AllStation(&controller.Service)
	if err != nil {
		controller.ServeError(err)
		return
	}

	controller.Data["Stations"] = buoyStations
	controller.Layout = "shared/basic-layout.html"
	controller.TplName = "buoy/content.html"
	controller.LayoutSections = map[string]string{}
	controller.LayoutSections["PageHead"] = "buoy/page-head.html"
	controller.LayoutSections["Header"] = "shared/header.html"
	controller.LayoutSections["Modal"] = "shared/modal.html"
}

//** AJAX FUNCTIONS

// RetrieveStation handles the example 2 tab.
func (controller *BuoyController) RetrieveStation() {
	var params struct {
		StationID string `form:"stationId" error:"invalid_station_id"`
	}
	if controller.ParseAndValidate(&params) == false {
		return
	}
	buoyStation, err := buoyService.FindStation(&controller.Service, params.StationID)
	if err != nil {
		log.CompletedErrorf(err, controller.UserID, "BuoyController.RetrieveStation", "StationID[%s]", params.StationID)
		controller.ServeError(err)
		return
	}
	controller.Data["Station"] = buoyStation
	controller.Layout = ""
	controller.TplName = "buoy/modal/pv_station-detail.html"
	view, _ := controller.RenderString()
	controller.AjaxResponse(0, "SUCCESS", view)
}

// RetrieveStationJSON handles the example 3 tab.
// http://localhost:9003/buoy/station/42002
func (controller *BuoyController) RetrieveStationJSON() {
	buoyStations, err := buoyService.AllStation(&controller.Service)

	if err != nil {
		controller.ServeError(err)
		return
	}

	controller.Data["json"] = buoyStations
	controller.ServeJSON()
}

func (controller *BuoyController) NewStation() {
	controller.Layout = "shared/basic-layout.html"
	controller.TplName = "buoy/form.html"
	controller.LayoutSections = map[string]string{}
	controller.LayoutSections["PageHead"] = "buoy/page-head.html"
	controller.LayoutSections["Header"] = "shared/header.html"
	controller.LayoutSections["Modal"] = "shared/modal.html"
}

func (controller *BuoyController) Create() {
	buoyStation := buoyModels.BuoyStation{}

	var windGust, windSpeed float64
	var windDirection string

	controller.Ctx.Input.Bind(&windGust, "condition[windGust]")
	controller.Ctx.Input.Bind(&windSpeed, "condition[windSpeed]")
	controller.Ctx.Input.Bind(&windDirection, "condition[windDirection]")
	if err := controller.ParseForm(&buoyStation); err != nil {
		log.Trace("ERROR", ">>>", mongo.ToString(err))
	} else {
		buoyStation.Condition.WindGust = windGust
		buoyStation.Condition.WindDirection = windDirection
		buoyStation.Condition.WindSpeed = windSpeed

		buoyService.CreateStation(&controller.Service, buoyStation)
		controller.Redirect("/", 302)
	}
}