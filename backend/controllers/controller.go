package controllers

import (
	"structured-notes/app"
	"structured-notes/permissions"
)

type Controller struct {
	app        *app.App
	authorizer permissions.Authorizer
}
