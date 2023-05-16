package main

import "github.com/HipolitoBautista/internal/models"

type templateData struct {
	Form           []*models.Form
	Flash          string
	UserPermStatus string
	Comments       []*models.Comments
	Accepted       string
	Denied         string
	Unverified     string
	Username       string
	Total          string
}
