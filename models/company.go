package models

type Company struct {
	Name string `binding:"required" uri:"name"`
}

type CompanyDetails struct {
	Name         string
	NumEmployees int64
}
