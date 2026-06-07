CREATE TABLE alerts (
    id UUID PRIMARY KEY,

    tenant_id UUID NOT NULL REFERENCES tenants(id),

    service_id UUID NOT NULL REFERENCES services(id),

    rule_id UUID NOT NULL REFERENCES alert_rules(id),

    title VARCHAR(255) NOT NULL,

    description TEXT,

    severity VARCHAR(20) NOT NULL
    CHECK (severity IN ('INFO', 'WARNING', 'CRITICAL')),

    status VARCHAR(20) NOT NULL
    CHECK (status IN ('OPEN', 'ACKNOWLEDGED', 'RESOLVED')),

    triggered_at TIMESTAMP NOT NULL,

    acknowledged_at TIMESTAMP,
    acknowledged_by UUID REFERENCES users(id),

    resolved_at TIMESTAMP,
    resolved_by UUID REFERENCES users(id),

    notification_sent_at TIMESTAMP,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_alerts_status
ON alerts(status);

CREATE INDEX idx_alerts_tenant
ON alerts(tenant_id);