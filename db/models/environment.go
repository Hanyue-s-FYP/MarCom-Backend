package models

import (
	"database/sql"
	"errors"

	"github.com/Hanyue-s-FYP/Marcom-Backend/db"
	"github.com/Hanyue-s-FYP/Marcom-Backend/utils"
)

type Environment struct {
	ID          int
	Products    []Product
	Agents      []Agent
	Name        string
	Description string
	BusinessID  int
}

type environmentModel struct{}

var EnvironmentModel *environmentModel

func (*environmentModel) Create(env Environment) error {
	tx, err := db.GetDB().Begin()
	if err != nil {
		return err
	}

	environmentQuery := `
        INSERT INTO Environments (name, description, business_id)
        VALUES (?, ?, ?)
    `
	res, err := tx.Exec(environmentQuery, env.Name, env.Description, env.BusinessID)
	if err != nil {
		tx.Rollback()
		return err
	}

	environmentID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}

	agentQuery := `
        INSERT INTO EnvironmentAgents (environment_id, agent_id)
        VALUES (?, ?)
    `
	for _, agent := range env.Agents {
		_, err := tx.Exec(agentQuery, environmentID, agent.ID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	productQuery := `
        INSERT INTO EnvironmentProducts (environment_id, product_id)
        VALUES (?, ?)
    `
	for _, product := range env.Products {
		_, err := tx.Exec(productQuery, environmentID, product.ID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

var ErrEnvironmentNotFound = errors.New("environment not found")

func (*environmentModel) GetByID(id int) (*Environment, error) {
	environmentQuery := `
		SELECT id, name, description, business_id
		FROM Environments
		WHERE id = ?
	`
	row := db.GetDB().QueryRow(environmentQuery, id)
	var env Environment
	err := row.Scan(&env.ID, &env.Name, &env.Description, &env.BusinessID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrEnvironmentNotFound
		}
		return nil, err
	}

	agents, err := getAgentsByEnvironmentId(id)
	if err != nil {
		return nil, err
	}
	env.Agents = agents

	products, err := getProductsByEnvironmentId(id)
	if err != nil {
		return nil, err
	}
	env.Products = products

	return &env, nil
}

func (*environmentModel) GetAll() ([]Environment, error) {
	environmentQuery := `
		SELECT id, name, description, business_id
		FROM Environments
	`
	rows, err := db.GetDB().Query(environmentQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var environments []Environment
	for rows.Next() {
		var env Environment
		if err := rows.Scan(&env.ID, &env.Name, &env.Description, &env.BusinessID); err != nil {
			return nil, err
		}

		agents, err := getAgentsByEnvironmentId(env.ID)
		if err != nil {
			return nil, err
		}
		env.Agents = agents

		products, err := getProductsByEnvironmentId(env.ID)
		if err != nil {
			return nil, err
		}
		env.Products = products

		environments = append(environments, env)
	}

	return environments, nil
}

func (*environmentModel) GetAllByBusinessID(businessId int) ([]Environment, error) {
	environmentQuery := `
		SELECT id, name, description, business_id
		FROM Environments
		WHERE business_id = ?
	`
	rows, err := db.GetDB().Query(environmentQuery, businessId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var environments []Environment
	for rows.Next() {
		var env Environment
		if err := rows.Scan(&env.ID, &env.Name, &env.Description, &env.BusinessID); err != nil {
			return nil, err
		}

		agents, err := getAgentsByEnvironmentId(env.ID)
		if err != nil {
			return nil, err
		}
		env.Agents = agents

		products, err := getProductsByEnvironmentId(env.ID)
		if err != nil {
			return nil, err
		}
		env.Products = products

		environments = append(environments, env)
	}

	return environments, nil
}

func (*environmentModel) Update(env Environment) error {
	tx, err := db.GetDB().Begin()
	if err != nil {
		return err
	}

	environmentQuery := `
		UPDATE Environments
		SET name = ?, description = ?, business_id = ?
		WHERE id = ?
	`
	_, err = tx.Exec(environmentQuery, env.Name, env.Description, env.BusinessID, env.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// get existing products and agents
	existingAgents, err := getAgentsByEnvironmentId(env.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	existingProducts, err := getProductsByEnvironmentId(env.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// find out what to insert and what to remove
	// for agents
	compareFuncAgent := func(a, b Agent) bool { return a.ID == b.ID } // only care about id when comparing, as id should be unique and should exist alrd the record
	newAgents := utils.NotIn(env.Agents, existingAgents, compareFuncAgent)
	removedAgents := utils.NotIn(existingAgents, env.Agents, compareFuncAgent)
	// Insert new agents
	insertAgentQuery := `
		INSERT INTO EnvironmentAgents (environment_id, agent_id)
		VALUES (?, ?)
	`
	for _, agent := range newAgents {
		_, err := tx.Exec(insertAgentQuery, env.ID, agent.ID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	// Delete removed agents
	deleteAgentQuery := `DELETE FROM EnvironmentAgents WHERE environment_id = ? AND agent_id = ?`
	for _, ragent := range removedAgents {
		_, err := tx.Exec(deleteAgentQuery, env.ID, ragent.ID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// for products
	compareFuncProduct := func(a, b Product) bool { return a.ID == b.ID } // same, only care about id, sadly golang's simplicity philosophy means writing more code as the simple generic supported by golang is unable to be applied here
	newProducts := utils.NotIn(env.Products, existingProducts, compareFuncProduct)
	removedProducts := utils.NotIn(existingProducts, env.Products, compareFuncProduct)
	// Insert new products
	insertProductQuery := `
		INSERT INTO EnvironmentProducts (environment_id, product_id)
		VALUES (?, ?)
	`
	for _, product := range newProducts {
		_, err := tx.Exec(insertProductQuery, env.ID, product.ID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	// Delete removed products
	deleteProductQuery := `DELETE FROM EnvironmentProducts WHERE environment_id = ? AND product_id = ?`
	for _, rproduct := range removedProducts {
		_, err := tx.Exec(deleteProductQuery, env.ID, rproduct.ID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (*environmentModel) Delete(id int) error {
	tx, err := db.GetDB().Begin()
	if err != nil {
		return err
	}

	deleteAgentQuery := `DELETE FROM EnvironmentAgents WHERE environment_id = ?`
	_, err = tx.Exec(deleteAgentQuery, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	deleteProductQuery := `DELETE FROM EnvironmentProducts WHERE environment_id = ?`
	_, err = tx.Exec(deleteProductQuery, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	deleteEnvironmentQuery := `DELETE FROM Environments WHERE id = ?`
	_, err = tx.Exec(deleteEnvironmentQuery, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func getAgentsByEnvironmentId(id int) ([]Agent, error) {
	agentIdQuery := `SELECT agent_id FROM EnvironmentAgents WHERE environment_id = ?`
	agentIdRows, err := db.GetDB().Query(agentIdQuery, id)
	if err != nil {
		return nil, err
	}
	defer agentIdRows.Close()

	var agents []Agent
	for agentIdRows.Next() {
		var agentId int
		if err := agentIdRows.Scan(&agentId); err != nil {
			return nil, err
		}
		agent, err := AgentModel.GetByID(agentId)
		if err != nil {
			return nil, err
		}
		agents = append(agents, *agent)
	}
	return agents, nil
}

func getProductsByEnvironmentId(id int) ([]Product, error) {
	productIdQuery := `SELECT product_id FROM EnvironmentProducts WHERE environment_id = ?`
	productIdRows, err := db.GetDB().Query(productIdQuery, id)
	if err != nil {
		return nil, err
	}
	defer productIdRows.Close()

	var products []Product
	for productIdRows.Next() {
		var productId int
		if err := productIdRows.Scan(&productId); err != nil {
			return nil, err
		}
		product, err := ProductModel.GetByID(productId)
		if err != nil {
			return nil, err
		}
		products = append(products, *product)
	}
	return products, nil
}
