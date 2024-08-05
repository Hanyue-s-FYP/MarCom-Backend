CREATE TABLE IF NOT EXISTS Users (
    id INTEGER PRIMARY KEY,
    username TEXT NOT NULL,
    password TEXT NOT NULL,
    display_name TEXT NOT NULL,
    email TEXT NOT NULL,
    status INTEGER NOT NULL,
    phone_number TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS Investors (
    id INTEGER PRIMARY KEY,
    user_id INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES Users(id)
);

-- will be using jwt, no need to store suth tokens on the server anymore
-- CREATE TABLE IF NOT EXISTS AuthTokens (
--     token TEXT PRIMARY KEY,
--     expires DATE NOT NULL,
--     user_id INTEGER NOT NULL,
--     FOREIGN KEY (user_id) REFERENCES Users(id)
-- );

CREATE TABLE IF NOT EXISTS Businesses (
    id INTEGER PRIMARY KEY,
    description TEXT NOT NULL,
    business_type TEXT NOT NULL,
    cover_img_path TEXT NOT NULL,
    user_id INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES Users(id)
);

CREATE TABLE IF NOT EXISTS Products (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    price REAL NOT NULL,
    cost REAL NOT NULL,
    business_id INTEGER NOT NULL,
    FOREIGN KEY (business_id) REFERENCES Businesses(id)
);

CREATE TABLE IF NOT EXISTS Agents (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    general_description TEXT,
    business_id INTEGER NOT NULL,
    FOREIGN KEY (business_id) REFERENCES Businesses(id)
);

CREATE TABLE IF NOT EXISTS AgentAttributes (
    id INTEGER PRIMARY KEY,
    key TEXT NOT NULL,
    value TEXT NOT NULL,
    agent_id INTEGER NOT NULL,
    FOREIGN KEY (agent_id) REFERENCES Agents(id)
);

CREATE TABLE IF NOT EXISTS Environments (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    business_id INTEGER NOT NULL,
    FOREIGN KEY (business_id) REFERENCES Businesses(id)
);

CREATE TABLE IF NOT EXISTS EnvironmentAgents (
    environment_id INTEGER NOT NULL,
    agent_id INTEGER NOT NULL,
    PRIMARY KEY (environment_id, agent_id),
    FOREIGN KEY (environment_id) REFERENCES Environments(id),
    FOREIGN KEY (agent_id) REFERENCES Agents(id)
);

CREATE TABLE IF NOT EXISTS SimulationEvents (
    id INTEGER PRIMARY KEY,
    prompt TEXT NOT NULL,
    event_type INTEGER NOT NULL,
    event_description TEXT NOT NULL,
    agent_id INTEGER NULL,
    cycle_id INTEGER NOT NULL,
    time DATETIME NOT NULL, -- to preserve the order, more consistent when sorting
    FOREIGN KEY (agent_id) REFERENCES Agents(id),
    FOREIGN KEY (cycle_id) REFERENCES SimulationCycles(id)
);

CREATE TABLE IF NOT EXISTS SimulationCycles (
    id INTEGER PRIMARY KEY,
    simulation_id INTEGER NOT NULL,
    time DATETIME NOT NULL, -- to preserve the order, more consistent when sorting
    FOREIGN KEY (simulation_id) REFERENCES Simulations(id)
);

CREATE TABLE IF NOT EXISTS Simulations (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    max_cycle_count INTEGER NOT NULL,
    is_price_opt_enabled BOOLEAN NOT NULL,
    status INTEGER NOT NULL,
    environment_id INTEGER NOT NULL,
    business_id INTEGER NOT NULL,
    FOREIGN KEY (environment_id) REFERENCES Environments(id),
    FOREIGN KEY (business_id) REFERENCES Businesses(id)
);

CREATE TABLE IF NOT EXISTS EnvironmentProducts (
    environment_id INTEGER NOT NULL,
    product_id INTEGER NOT NULL,
    PRIMARY KEY (environment_id, product_id),
    FOREIGN KEY (environment_id) REFERENCES Environments(id),
    FOREIGN KEY (product_id) REFERENCES Products(id)
);
