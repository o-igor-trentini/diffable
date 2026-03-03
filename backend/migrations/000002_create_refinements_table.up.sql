CREATE TABLE refinements (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    analysis_id     UUID REFERENCES analyses(id) ON DELETE CASCADE,
    instruction     TEXT NOT NULL,
    original_desc   TEXT NOT NULL,
    refined_desc    TEXT,
    model_used      VARCHAR(50),
    tokens_used     INTEGER,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_refinements_analysis_id ON refinements(analysis_id);
