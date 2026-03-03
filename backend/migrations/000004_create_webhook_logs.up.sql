CREATE TABLE webhook_logs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_key       VARCHAR(100) NOT NULL,
    payload         JSONB,
    status          VARCHAR(20) NOT NULL DEFAULT 'received',
    analysis_id     UUID REFERENCES analyses(id) ON DELETE SET NULL,
    error_message   TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_webhook_logs_created_at ON webhook_logs(created_at DESC);
CREATE INDEX idx_webhook_logs_status ON webhook_logs(status);
