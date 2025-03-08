CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    phone_number VARCHAR(15) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    aadhar VARCHAR(12) UNIQUE NOT NULL,
    dob DATE NOT NULL,
    address TEXT NOT NULL
);

CREATE TABLE accounts (
    account_id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(user_id) ON DELETE CASCADE,
    account_number VARCHAR(20) UNIQUE NOT NULL,  -- Unique Account Number
    ifsc_code VARCHAR(11) NOT NULL,  -- Added IFSC Code
    micr_code VARCHAR(9),  -- Optional MICR Code
    branch_name VARCHAR(100),
    account_type VARCHAR(50) CHECK (account_type IN ('Savings', 'Current', 'Fixed Deposit', 'Loan')),
    account_sub_type VARCHAR(50) CHECK (account_sub_type IN ('REGULAR', 'NRI', 'SALARY', 'SAVINGS', 'CURRENT')),
    balance DECIMAL(15,2) DEFAULT 0.0,
    nominee VARCHAR(100),
    nominee_relationship VARCHAR(50),
    nominee_status VARCHAR(20) CHECK (nominee_status IN ('REGISTERED', 'NOT_REGISTERED')),
    account_status VARCHAR(20) CHECK (account_status IN ('ACTIVE', 'INACTIVE', 'CLOSED')),
    opening_date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE transactions (
    transaction_id SERIAL PRIMARY KEY,
    account_id INT REFERENCES accounts(account_id) ON DELETE CASCADE,
    transaction_type VARCHAR(10) CHECK (transaction_type IN ('Credit', 'Debit')),
    mode VARCHAR(50) CHECK (mode IN ('UPI', 'NEFT', 'RTGS', 'IMPS', 'CARD', 'CASH', 'CHEQUE', 'OTHER')),
    amount DECIMAL(15,2) NOT NULL,
    transaction_timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    description TEXT
);

CREATE TABLE loans (
    loan_id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(user_id) ON DELETE CASCADE,
    account_id INT REFERENCES accounts(account_id) ON DELETE CASCADE,  
    loan_type VARCHAR(50) CHECK (loan_type IN ('Home Loan', 'Car Loan', 'Personal Loan', 'Education Loan')),
    amount DECIMAL(15,2) NOT NULL,
    interest_rate DECIMAL(5,2) NOT NULL,
    tenure INT NOT NULL,  
    status VARCHAR(50) CHECK (status IN ('Active', 'Closed', 'Defaulted')),
    sanction_date DATE NOT NULL,
    due_amount DECIMAL(15,2) DEFAULT 0.0
);

CREATE TABLE cards (
    card_id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(user_id) ON DELETE CASCADE,
    account_id INT REFERENCES accounts(account_id) ON DELETE CASCADE,  
    card_type VARCHAR(50) CHECK (card_type IN ('Credit', 'Debit')),
    card_number VARCHAR(16) UNIQUE NOT NULL,  
    expiry_date DATE NOT NULL,
    cvv VARCHAR(3) NOT NULL,
    status VARCHAR(50) CHECK (status IN ('Active', 'Blocked')),
    credit_limit DECIMAL(15,2) DEFAULT 0.0, 
    available_balance DECIMAL(15,2) DEFAULT 0.0  
);
