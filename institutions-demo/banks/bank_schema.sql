CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    phone_number VARCHAR(15) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    aadhar VARCHAR(12) NOT NULL,
    dob DATE NOT NULL,
    address TEXT
);


CREATE TABLE accounts (
    account_id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(user_id) ON DELETE CASCADE,
    account_type VARCHAR(50) CHECK (account_type IN ('Savings', 'Current', 'Fixed Deposit', 'Loan')),
    balance DECIMAL(15,2) DEFAULT 0.0,
    nominee VARCHAR(100),
    nominee_relationship VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE transactions (
    transaction_id SERIAL PRIMARY KEY,
    account_id INT REFERENCES accounts(account_id) ON DELETE CASCADE,
    transaction_type VARCHAR(50) CHECK (transaction_type IN ('Credit', 'Debit')),
    amount DECIMAL(15,2) NOT NULL,
    transaction_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    description TEXT
);


CREATE TABLE loans (
    loan_id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(user_id) ON DELETE CASCADE,
    loan_type VARCHAR(50) CHECK (loan_type IN ('Home Loan', 'Car Loan', 'Personal Loan','Education Loan')),
    amount DECIMAL(15,2) NOT NULL,
    interest_rate DECIMAL(5,2) NOT NULL,
    tenure INT NOT NULL,  
    status VARCHAR(50) CHECK (status IN ('Active', 'Closed', 'Defaulted'))
);


CREATE TABLE cards (
    card_id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(user_id) ON DELETE CASCADE,
    card_type VARCHAR(50) CHECK (card_type IN ('Credit', 'Debit')),
    card_number VARCHAR(16) UNIQUE NOT NULL,
    expiry_date DATE NOT NULL,
    cvv VARCHAR(3) NOT NULL,
    status VARCHAR(50) CHECK (status IN ('Active', 'Blocked'))
);
