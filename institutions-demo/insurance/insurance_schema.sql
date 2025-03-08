CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    phone_number VARCHAR(15) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    aadhar VARCHAR(12) NOT NULL,
    dob DATE NOT NULL,
    address TEXT
);

CREATE TABLE insurance_providers (
    provider_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    contact_email VARCHAR(255),
    contact_phone VARCHAR(15),
    address TEXT
);

CREATE TABLE insurance_policies (
    policy_id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(user_id) ON DELETE CASCADE,
    provider_id INT REFERENCES insurance_providers(provider_id) ON DELETE CASCADE,
    policy_type VARCHAR(50) CHECK (policy_type IN ('Health', 'Life', 'Vehicle')),
    policy_number VARCHAR(50) UNIQUE NOT NULL,
    sum_assured DECIMAL(15,2), 
    premium_amount DECIMAL(15,2),
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    nominee VARCHAR(100),
    nominee_relationship VARCHAR(50)
);
