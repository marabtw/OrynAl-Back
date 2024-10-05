package model

import "time"

type Params struct {
	Filter     map[string]interface{}
	Query      string
	Order      interface{}
	SortVector interface{}
	Date       *time.Time
	Offset     int
	Limit      int
	PageIndex  int
}

func NewParams() *Params {
	return &Params{}
}

var TablesOrderKeyList = []string{
	"id",
	"name",
	"type",
	"capacity",
}

var UsersOrderKeyList = []string{
	"id",
	"name",
	"surname",
	"email",
}

type ListResponse struct {
	Items        interface{} `json:"items"`
	ItemsPerPage int         `json:"itemsPerPage"`
	PageIndex    int         `json:"pageIndex"`
	TotalPages   int         `json:"totalPages"`
	TotalItems   int         `json:"totalItems"`
}
