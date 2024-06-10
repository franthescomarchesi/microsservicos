CREATE TABLE clients (
    id VARCHAR(255), 
    name VARCHAR(255), 
    email VARCHAR(255), 
    created_at date, 
    updated_at date
);
CREATE TABLE accounts (
    id VARCHAR(255), 
    client_id VARCHAR(255), 
    balance int, 
    created_at date, 
    updated_at date
);
CREATE TABLE transactions (
    id VARCHAR(255), 
    account_id_from VARCHAR(255), 
    account_id_to VARCHAR(255), 
    amount int, 
    created_at date
);
INSERT INTO clients (
    id, 
    name, 
    email, 
    created_at, 
    updated_at
) VALUES (
    '1b1799eb-3ba2-4152-8881-1e433c038aa3',
    'John Doe',
    'john@j.com',
    '2024-05-29',
    '2024-05-29'
), (
    '44a75715-bc34-467b-b9cf-046db6b2a453',
    'Jane Doe',
    'jane@j.com',
    '2024-05-29',
    '2024-05-29'
);
INSERT INTO accounts (
    id, 
    client_id, 
    balance, 
    created_at, 
    updated_at) 
VALUES (
    '90912d8b-1f55-4fd5-b992-2ecc21d548ac',
    '1b1799eb-3ba2-4152-8881-1e433c038aa3',
    '1000',
    '2024-05-29',
    '2024-05-29'
), (
    '5dc6a210-f5d8-485e-8516-db0fa856e520',
    '44a75715-bc34-467b-b9cf-046db6b2a453',
    '1000',
    '2024-05-29',
    '2024-05-29'
);
