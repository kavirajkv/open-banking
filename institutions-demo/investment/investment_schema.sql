CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    phone_number VARCHAR(15) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    aadhar VARCHAR(12) NOT NULL,
    dob DATE NOT NULL,
    address TEXT
);

CREATE TABLE investment_brokers (
    broker_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    broker_type VARCHAR(50) CHECK (broker_type IN ('Stock Broker', 'AMC')),
    contact_email VARCHAR(255),
    contact_phone VARCHAR(15),
    website TEXT
);

CREATE TABLE demat_accounts (
    demat_id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(user_id) ON DELETE CASCADE,
    broker_id INT REFERENCES investment_brokers(broker_id) ON DELETE CASCADE,
    demat_number VARCHAR(20) UNIQUE NOT NULL, -- NSDL/CDSL Demat ID
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE stocks_and_bonds (
    investment_id SERIAL PRIMARY KEY,
    demat_id INT REFERENCES demat_accounts(demat_id) ON DELETE CASCADE,
    investment_type VARCHAR(50) CHECK (investment_type IN ('Stock', 'Bond')),
    asset_name VARCHAR(255) NOT NULL, -- Stock Name / Bond Name
    asset_symbol VARCHAR(50), -- Stock Symbol (e.g., TCS, INFY)
    quantity DECIMAL(15,2) DEFAULT 0.0,
    purchase_price DECIMAL(15,2) DEFAULT 0.0,
    current_price DECIMAL(15,2) DEFAULT 0.0,
    purchase_date DATE NOT NULL
);

CREATE TABLE mf_accounts (
    mf_account_id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(user_id) ON DELETE CASCADE,
    broker_id INT REFERENCES investment_brokers(broker_id) ON DELETE CASCADE, -- AMC details
    folio_number VARCHAR(50) UNIQUE NOT NULL, -- Mutual Fund Folio Number
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE mutual_funds (
    investment_id SERIAL PRIMARY KEY,
    mf_account_id INT REFERENCES mf_accounts(mf_account_id) ON DELETE CASCADE,
    fund_name VARCHAR(255) NOT NULL, -- Mutual Fund Name
    fund_type VARCHAR(50) CHECK (fund_type IN ('Equity', 'Debt', 'Hybrid')),
    quantity DECIMAL(15,2) DEFAULT 0.0, -- Number of Units
    purchase_price DECIMAL(15,2) DEFAULT 0.0,
    current_price DECIMAL(15,2) DEFAULT 0.0,
    purchase_date DATE NOT NULL
);
