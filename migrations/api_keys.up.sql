CREATE TABLE api_keys (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    service_id UUID NOT NULL REFERENCES services(id),
    key_hash VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_api_keys_tenant ON api_keys(tenant_id);
CREATE INDEX idx_api_keys_active ON api_keys(is_active) WHERE is_active = TRUE;
