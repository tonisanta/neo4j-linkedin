package repositories

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"neo4j-linkedin/models"
	"time"
)

type Neo4jRepository struct {
	driver neo4j.Driver
}

func NewNeo4jRepository(driver neo4j.Driver) *Neo4jRepository {
	return &Neo4jRepository{driver}
}

func (r *Neo4jRepository) addItem(query string, params map[string]any) (any, error) {
	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		result, err := tx.Run(query, params)
		if err != nil {
			return nil, err
		}
		return result.Single()
	})

	return result, err
}

func (r *Neo4jRepository) getItem(query string, params map[string]any) (*neo4j.Record, error) {
	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		result, err := tx.Run(query, params)
		if err != nil {
			return nil, err
		}
		return result.Single()
	})

	return result.(*neo4j.Record), err
}

func (r *Neo4jRepository) getItems(query string, params map[string]any) ([]*neo4j.Record, error) {
	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		result, err := tx.Run(query, params)
		if err != nil {
			return nil, err
		}
		return result.Collect()
	})

	return result.([]*neo4j.Record), err
}

func (r *Neo4jRepository) AddPerson(person models.Person) error {
	query := "MERGE (p: Person { name: $name, skills: $skills }) RETURN p"
	params := map[string]any{
		"name":   person.Name,
		"skills": person.Skills,
	}
	_, err := r.addItem(query, params)
	return err
}

func (r *Neo4jRepository) AddCompany(company models.Company) error {
	query := "MERGE (c: Company { name: $name }) RETURN c"
	params := map[string]any{
		"name": company.Name,
	}
	_, err := r.addItem(query, params)
	return err
}

func (r *Neo4jRepository) AddRelationshipPersonToPerson(relationship models.RelationshipPersonToPerson) error {
	query := fmt.Sprintf(`MATCH (p1: Person { name: $nameOrigin })
			  MATCH (p2: Person { name: $nameDestination })
			  MERGE (p1)-[r:%v]->(p2) 
			  SET r.datetime = datetime()
              RETURN r`, relationship.NameRelationship)
	params := map[string]any{
		"nameOrigin":      relationship.NameOrigin,
		"nameDestination": relationship.NameDestination,
	}
	_, err := r.addItem(query, params)
	return err
}

func (r *Neo4jRepository) AddRelationshipPersonToCompany(relationship models.RelationshipPersonToCompany) error {
	query := fmt.Sprintf(`MATCH (p: Person { name: $nameOrigin })
			  MATCH (c: Company { name: $nameDestination })
			  MERGE (p)-[r:%v]->(c)
              RETURN r`, relationship.NameRelationship)
	params := map[string]any{
		"nameOrigin":      relationship.NamePersonOrigin,
		"nameDestination": relationship.NameCompanyDestination,
	}
	_, err := r.addItem(query, params)
	return err
}

func (r *Neo4jRepository) GetCompanyInfo(c models.Company) (models.CompanyDetails, error) {
	query := "MATCH (c:Company)-[:WORKS]-(n) RETURN count(n) as numEmployees"
	params := map[string]any{
		"name": c.Name,
	}

	item, err := r.getItem(query, params)
	if err != nil {
		return models.CompanyDetails{}, err
	}

	numEmployees, _ := item.Get("numEmployees")
	companyFull := models.CompanyDetails{Name: c.Name, NumEmployees: numEmployees.(int64)}
	return companyFull, nil
}

func (r *Neo4jRepository) GetProfileViewers(p models.Person) ([]models.ProfileViewer, error) {
	query := `MATCH (n)-[r:VIEWED_PROFILE]->(p:Person {name: $name}) 
				RETURN n.name, n.skills, r.datetime`
	params := map[string]any{
		"name": p.Name,
	}

	var persons []models.ProfileViewer
	viewers, err := r.getItems(query, params)
	if err != nil {
		return []models.ProfileViewer{}, err
	}

	for _, viewerNode := range viewers {
		name, _ := viewerNode.Get("n.name")
		skillsNode, _ := viewerNode.Get("n.skills")
		datetime, _ := viewerNode.Get("r.datetime")

		skills := buildSkills(skillsNode.([]any))
		person := models.Person{Name: name.(string), Skills: skills}
		viewer := models.ProfileViewer{Person: person, When: datetime.(time.Time)}
		persons = append(persons, viewer)
	}

	return persons, nil
}

func buildSkills(skillsNode []any) []string {
	var skills []string
	for _, skill := range skillsNode {
		skills = append(skills, skill.(string))
	}
	return skills
}
