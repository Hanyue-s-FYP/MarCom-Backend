package models

import (
	"database/sql"
	"errors"

	"github.com/Hanyue-s-FYP/Marcom-Backend/db"
	"github.com/Hanyue-s-FYP/Marcom-Backend/utils"
)

type Agent struct {
	ID                 int
	Name               string
	GeneralDescription sql.NullString
	BusinessID         int
	Attributes         []AgentAttribute
}

// no need embed agent here, just put all agent attributes relevant to the agent to the agent struct
type AgentAttribute struct {
	ID    int
	Key   string
	Value string
}

type agentModel struct{}

var AgentModel *agentModel

func (*agentModel) Create(a Agent) error {
	tx, err := db.GetDB().Begin()
	if err != nil {
		return err
	}

	agentQuery := `
        INSERT INTO Agents (name, general_description, business_id)
        VALUES (?, ?, ?)
    `
	res, err := tx.Exec(agentQuery, a.Name, a.GeneralDescription, a.BusinessID)
	if err != nil {
		tx.Rollback()
		return err
	}

	agentID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}

	attributeQuery := `
        INSERT INTO AgentAttributes (key, value, agent_id)
        VALUES (?, ?, ?)
    `
	// for easy handling just add key value to the table so no need to concern about multiple agent use same row, easier handling when deleting and updating
	for _, attr := range a.Attributes {
		_, err := tx.Exec(attributeQuery, attr.Key, attr.Value, agentID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

var ErrAgentNotFound error = errors.New("agent not found")

func (*agentModel) GetByID(id int) (*Agent, error) {
	agentQuery := `
        SELECT id, name, general_description, business_id
        FROM Agents
        WHERE id = ?
    `
	row := db.GetDB().QueryRow(agentQuery, id)

	var agent Agent
	err := row.Scan(&agent.ID, &agent.Name, &agent.GeneralDescription, &agent.BusinessID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrAgentNotFound
		}
		return nil, err
	}

	attrs, err := getAgentAttribute(agent.ID)
	if err != nil {
		return nil, err
	}
	agent.Attributes = attrs

	return &agent, nil
}

func (*agentModel) GetAll() ([]Agent, error) {
	agentsQuery := `
        SELECT id, name, general_description, business_id
        FROM Agents
    `
	rows, err := db.GetDB().Query(agentsQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var agents []Agent
	for rows.Next() {
		var agent Agent
		err := rows.Scan(&agent.ID, &agent.Name, &agent.GeneralDescription, &agent.BusinessID)
		if err != nil {
			return nil, err
		}

		// get the attributes for the current agent
		attrs, err := getAgentAttribute(agent.ID)
		if err != nil {
			return nil, err
		}
		agent.Attributes = attrs

		agents = append(agents, agent)
	}

	return agents, nil
}

func (*agentModel) GetByBusinessID(id int) ([]Agent, error) {
    agentQuery := `
        SELECT id, name, general_description, business_id
        FROM Agents
        WHERE business_id = ?
    `
    rows, err := db.GetDB().Query(agentQuery, id)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var agents []Agent
    for rows.Next() {
        var agent Agent
        err := rows.Scan(&agent.ID, &agent.Name, &agent.GeneralDescription, &agent.BusinessID)
        if err != nil {
            return nil, err
        }
        attrs, err := getAgentAttribute(agent.ID)
        if err != nil {
            return nil, err
        }
        agent.Attributes = attrs
        agents = append(agents, agent)
    }
    
    return agents, nil
}

func (*agentModel) Update(a Agent) error {
	tx, err := db.GetDB().Begin()
	if err != nil {
		return err
	}

	agentQuery := `
        UPDATE Agents
        SET name = ?, general_description = ?, business_id = ?
        WHERE id = ?
    `
	_, err = tx.Exec(agentQuery, a.Name, a.GeneralDescription, a.BusinessID, a.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Get existing attributes
	existingAttributesQuery := `SELECT id, key, value FROM AgentAttributes WHERE agent_id = ?`
	rows, err := tx.Query(existingAttributesQuery, a.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer rows.Close()

	var existingAttributes []AgentAttribute
	for rows.Next() {
		var attr AgentAttribute
		err := rows.Scan(&attr.ID, &attr.Key, &attr.Value)
		if err != nil {
			tx.Rollback()
			return err
		}
		existingAttributes = append(existingAttributes, attr)
	}

    // find out attributes that are removed from the ori agent
    // only care about key and value, if user provide new key value pair, the id will be zero so comparing id is not accurate
    compareFunc := func (a, b AgentAttribute) bool { return a.Key == b.Key && a.Value == b.Value } 
    removedAttrs := utils.NotIn(existingAttributes, a.Attributes, compareFunc)
    newAttrs := utils.NotIn(a.Attributes, existingAttributes, compareFunc)
    deleteQuery := `DELETE FROM AgentAttributes WHERE agent_id = ? AND id = ?`
    for _, rattr := range removedAttrs {
        // only delete if already exist in db (ID will not be zero value `0`) (should actually return error cause this should not happen)
        if rattr.ID != 0 {
            _, err := tx.Exec(deleteQuery, a.ID, rattr.ID)
            if err != nil {
                tx.Rollback()
                return err
            }
        }
    }

    // Insert new attributes
    insertQuery := `
        INSERT INTO AgentAttributes (key, value, agent_id)
        VALUES (?, ?, ?)
    `
    for _, nattr := range newAttrs {
        // only insert if id is zero value (should actually return error)
        if nattr.ID == 0 {
            _, err := tx.Exec(insertQuery, nattr.Key, nattr.Value, a.ID)
            if err != nil {
                tx.Rollback()
                return err
            }
        }
    }

    return tx.Commit()
}

func (*agentModel) Delete(id int) error {
    tx, err := db.GetDB().Begin()
    if err != nil {
        return err
    }

    deleteAttributesQuery := `
        DELETE FROM AgentAttributes
        WHERE agent_id = ?
    `
    _, err = tx.Exec(deleteAttributesQuery, id)
    if err != nil {
        tx.Rollback()
        return err
    }

    deleteAgentQuery := `
        DELETE FROM Agents
        WHERE id = ?
    `
    _, err = tx.Exec(deleteAgentQuery, id)
    if err != nil {
        tx.Rollback()
        return err
    }

    return tx.Commit()
}

func getAgentAttribute(id int) ([]AgentAttribute, error) {
	attributesQuery := `
        SELECT id, key, value
        FROM AgentAttributes
        WHERE agent_id = ?
    `
	attrRows, err := db.GetDB().Query(attributesQuery)
	if err != nil {
		return nil, err
	}
	defer attrRows.Close()
	agentAttributes := make([]AgentAttribute, 0)
	for attrRows.Next() {
		var agentAttribute AgentAttribute
		err := attrRows.Scan(&agentAttribute.ID, &agentAttribute.Key, &agentAttribute.Value)
		if err != nil {
			return nil, err
		}
		agentAttributes = append(agentAttributes, agentAttribute)
	}

	return agentAttributes, nil
}
