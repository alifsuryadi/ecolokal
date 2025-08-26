-- Create enum for user roles
CREATE TYPE user_role AS ENUM ('warga', 'petugas', 'admin');

-- Create enum for pickup status
CREATE TYPE pickup_status AS ENUM ('pending', 'scheduled', 'in_progress', 'completed', 'cancelled');

-- Create enum for transaction type
CREATE TYPE transaction_type AS ENUM ('add', 'redeem');

-- Users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role user_role NOT NULL DEFAULT 'warga',
    phone VARCHAR(20),
    address TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Waste types table
CREATE TABLE waste_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    point_per_kg INTEGER NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Pickup requests table
CREATE TABLE pickup_requests (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    petugas_id INTEGER REFERENCES users(id),
    scheduled_date DATE,
    scheduled_time VARCHAR(50),
    status pickup_status DEFAULT 'pending',
    notes TEXT,
    total_points INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Pickup items table
CREATE TABLE pickup_items (
    id SERIAL PRIMARY KEY,
    pickup_request_id INTEGER REFERENCES pickup_requests(id) ON DELETE CASCADE,
    waste_type_id INTEGER REFERENCES waste_types(id),
    estimated_weight DECIMAL(10,2),
    actual_weight DECIMAL(10,2),
    points_earned INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- User points table
CREATE TABLE user_points (
    id SERIAL PRIMARY KEY,
    user_id INTEGER UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    total_points INTEGER DEFAULT 0,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Transactions table
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    type transaction_type NOT NULL,
    points INTEGER NOT NULL,
    description TEXT,
    reference_id INTEGER, -- bisa refer ke pickup_request_id atau redemption_id
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_pickup_requests_user_id ON pickup_requests(user_id);
CREATE INDEX idx_pickup_requests_status ON pickup_requests(status);
CREATE INDEX idx_transactions_user_id ON transactions(user_id);
CREATE INDEX idx_pickup_items_pickup_request_id ON pickup_items(pickup_request_id);

-- Create trigger for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_waste_types_updated_at BEFORE UPDATE ON waste_types
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_pickup_requests_updated_at BEFORE UPDATE ON pickup_requests
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Insert default admin
INSERT INTO users (name, email, password, role) 
VALUES ('Admin EcoLokal', 'admin@ecolokal.com', '$2a$10$YourHashedPasswordHere', 'admin');

-- Insert default waste types
INSERT INTO waste_types (name, point_per_kg, description) VALUES
('Botol Plastik', 100, 'Botol plastik PET (air mineral, minuman)'),
('Kardus', 50, 'Kardus dan karton bekas'),
('Kertas', 30, 'Kertas HVS, koran, majalah'),
('Logam', 150, 'Kaleng, besi, aluminium'),
('Kaca', 40, 'Botol kaca, pecahan kaca');