CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    phone_number VARCHAR(15) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    aadhar VARCHAR(12) NOT NULL,
    dob DATE NOT NULL,
    address TEXT
);

CREATE TABLE epf_accounts (
    epf_id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    uan VARCHAR(20) UNIQUE NOT NULL, -- EPF Account Number
    employer_name VARCHAR(255) NOT NULL,  -- Employer name
    balance DECIMAL(15,2) DEFAULT 0.00,
    monthly_contribution DECIMAL(10,2) DEFAULT 0.00,
    last_contribution_date DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE nps_accounts (
    nps_id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    pran_number VARCHAR(20) UNIQUE NOT NULL, -- PRAN (Permanent Retirement Account Number)
    pension_tier VARCHAR(10) CHECK (pension_tier IN ('Tier 1', 'Tier 2')), -- Tier type
    balance DECIMAL(15,2) DEFAULT 0.00,
    monthly_contribution DECIMAL(10,2) DEFAULT 0.00,
    last_contribution_date DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE ppf_accounts (
    ppf_id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    ppf_number VARCHAR(20) UNIQUE NOT NULL, -- PPF Account Number
    balance DECIMAL(15,2) DEFAULT 0.00,
    yearly_contribution DECIMAL(10,2) DEFAULT 0.00,
    last_contribution_date DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

