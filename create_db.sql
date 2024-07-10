CREATE TABLE Users (
    id INTEGER PRIMARY KEY,
    username TEXT NOT NULL,
    password TEXT NOT NULL,
    display_name TEXT,
    email TEXT,
    status TEXT,
    phone_number TEXT
);

CREATE TABLE Investors (
    id INTEGER PRIMARY KEY,
    user_id INTEGER,
    FOREIGN KEY (user_id) REFERENCES Users(id)
);

-- Will be using jwt tokens, no need to additionally store the authentication tokens anymore
-- CREATE TABLE AuthTokens (
--     token TEXT PRIMARY KEY,
--     expires DATE,
--     user_id INTEGER,
--     FOREIGN KEY (user_id) REFERENCES Users(id)
-- );

CREATE TABLE Businesses (
    id INTEGER PRIMARY KEY,
    description TEXT,
    business_type TEXT,
    cover_img_path TEXT,
    user_id INTEGER,
    FOREIGN KEY (user_id) REFERENCES Users(id)
);

CREATE TABLE Products (
    id INTEGER PRIMARY KEY,
    name TEXT,
    description TEXT,
    price REAL,
    cost REAL,
    business_id INTEGER,
    FOREIGN KEY (business_id) REFERENCES Businesses(id)
);

CREATE TABLE Agents (
    id INTEGER PRIMARY KEY,
    name TEXT,
    general_description TEXT,
    business_id INTEGER,
    FOREIGN KEY (business_id) REFERENCES Businesses(id)
);

CREATE TABLE AgentAttributes (
    id INTEGER PRIMARY KEY,
    key TEXT,
    value TEXT,
    agent_id INTEGER,
    FOREIGN KEY (agent_id) REFERENCES Agents(id)
);

CREATE TABLE Environments (
    id INTEGER PRIMARY KEY,
    name TEXT,
    description TEXT,
    business_id INTEGER,
    FOREIGN KEY (business_id) REFERENCES Businesses(id)
);

CREATE TABLE EnvironmentAgents (
    environment_id INTEGER,
    agent_id INTEGER,
    PRIMARY KEY (environment_id, agent_id),
    FOREIGN KEY (environment_id) REFERENCES Environments(id),
    FOREIGN KEY (agent_id) REFERENCES Agents(id)
);

CREATE TABLE CycleAgentActions (
    id INTEGER PRIMARY KEY,
    prompt TEXT,
    action_type TEXT,
    action_description TEXT,
    agent_id INTEGER,
    cycle_id INTEGER,
    FOREIGN KEY (agent_id) REFERENCES Agents(id),
    FOREIGN KEY (cycle_id) REFERENCES SimulationCycles(id)
);

CREATE TABLE SimulationCycles (
    id INTEGER PRIMARY KEY,
    profit REAL,
    simulation_id INTEGER,
    FOREIGN KEY (simulation_id) REFERENCES Simulations(id)
);

CREATE TABLE Simulations (
    id INTEGER PRIMARY KEY,
    name TEXT,
    maxCycleCount INTEGER,
    is_price_opt_enabled BOOLEAN,
    status TEXT,
    environment_id INTEGER,
    business_id INTEGER,
    FOREIGN KEY (environment_id) REFERENCES Environments(id),
    FOREIGN KEY (business_id) REFERENCES Businesses(id)
);

CREATE TABLE EnvironmentProducts (
    environment_id INTEGER,
    product_id INTEGER,
    PRIMARY KEY (environment_id, product_id),
    FOREIGN KEY (environment_id) REFERENCES Environments(id),
    FOREIGN KEY (product_id) REFERENCES Products(id)
);
