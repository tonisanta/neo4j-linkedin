package models

import "time"

type Person struct {
	Name   string `binding:"required" uri:"name"`
	Skills []string
}

type RelationshipPersonToPerson struct {
	NameRelationship string `binding:"required"`
	NameOrigin       string `binding:"required"`
	NameDestination  string `binding:"required"`
}

type RelationshipPersonToCompany struct {
	NameRelationship       string `binding:"required"`
	NamePersonOrigin       string `binding:"required"`
	NameCompanyDestination string `binding:"required"`
}

type ProfileViewer struct {
	Person Person
	When   time.Time
}
