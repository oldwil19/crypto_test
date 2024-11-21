-- Crear base de datos y asignar propietario
CREATE DATABASE cryptodb WITH OWNER = crypto_user;

-- Conectarse a la base de datos

-- Crear tabla users
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY, -- Usar UUID como identificador único
    username TEXT UNIQUE NOT NULL, -- Nombre de usuario único
    password_hash TEXT NOT NULL, -- Contraseña encriptada
    balance NUMERIC(15, 2) DEFAULT 1000.00, -- Saldo inicial virtual en USD
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Crear tabla transactions para operaciones simuladas de trading
CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY, -- ID autoincremental
    user_id UUID NOT NULL REFERENCES users(id), -- Relación con la tabla users
    coin TEXT NOT NULL, -- Moneda (e.g., BTC, SOL)
    amount DECIMAL(18, 8) NOT NULL, -- Cantidad de la criptomoneda
    price DECIMAL(18, 8) NOT NULL, -- Precio en USD
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- Fecha de la transacción
);

-- Crear tabla balances para seguimiento de saldos
CREATE TABLE IF NOT EXISTS balances (
    user_id UUID PRIMARY KEY REFERENCES users(id), -- Relación con la tabla users
    usd_balance DECIMAL(18, 8) NOT NULL DEFAULT 1000.00, -- Saldo inicial en USD
    btc_balance DECIMAL(18, 8) NOT NULL DEFAULT 0, -- Saldo de BTC
    sol_balance DECIMAL(18, 8) NOT NULL DEFAULT 0 -- Saldo de SOL
);
