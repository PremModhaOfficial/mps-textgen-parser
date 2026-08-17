-- Auto-generated SQL for UserManagement
-- Schema: um

CREATE SCHEMA IF NOT EXISTS um;

-- Table: um.users
CREATE TABLE IF NOT EXISTS um.users (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	name VARCHAR(255) NOT NULL,
	email VARCHAR(255) NOT NULL UNIQUE,
	password VARCHAR(255) NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_users_email ON um.users(email);

-- Table: um.roles
CREATE TABLE IF NOT EXISTS um.roles (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	name VARCHAR(255) NOT NULL UNIQUE,
	permissions JSONB NOT NULL DEFAULT '[]',
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Junction table: um.user_roles
CREATE TABLE IF NOT EXISTS um.user_roles (
	user_id UUID NOT NULL REFERENCES um.users(id) ON DELETE CASCADE,
	role_id UUID NOT NULL REFERENCES um.roles(id) ON DELETE CASCADE,
	PRIMARY KEY (user_id, role_id)
);
