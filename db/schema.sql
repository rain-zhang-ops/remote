-- Initial schema for Control Plane
CREATE TABLE tenants (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    plan TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    tenant_id INTEGER REFERENCES tenants(id),
    email TEXT NOT NULL,
    role TEXT NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE devices (
    id SERIAL PRIMARY KEY,
    tenant_id INTEGER REFERENCES tenants(id),
    user_id INTEGER REFERENCES users(id),
    platform TEXT,
    pubkey TEXT,
    status TEXT,
    last_seen TIMESTAMPTZ
);

CREATE TABLE networks (
    id SERIAL PRIMARY KEY,
    tenant_id INTEGER REFERENCES tenants(id),
    cidr_v4 TEXT,
    cidr_v6 TEXT,
    dns TEXT,
    acl_version TEXT
);

CREATE TABLE memberships (
    device_id INTEGER REFERENCES devices(id),
    network_id INTEGER REFERENCES networks(id),
    tags JSONB,
    PRIMARY KEY (device_id, network_id)
);

CREATE TABLE relays (
    id SERIAL PRIMARY KEY,
    region TEXT,
    anycast_ip INET,
    proto TEXT,
    capacity INTEGER,
    health TEXT
);

CREATE TABLE sessions (
    id SERIAL PRIMARY KEY,
    src_dev INTEGER REFERENCES devices(id),
    dst_dev INTEGER REFERENCES devices(id),
    via TEXT,
    rtt INTEGER,
    loss REAL,
    bytes_up BIGINT,
    bytes_down BIGINT,
    started_at TIMESTAMPTZ,
    ended_at TIMESTAMPTZ
);

CREATE TABLE acls (
    id SERIAL PRIMARY KEY,
    tenant_id INTEGER REFERENCES tenants(id),
    expr TEXT,
    effect TEXT,
    priority INTEGER
);
