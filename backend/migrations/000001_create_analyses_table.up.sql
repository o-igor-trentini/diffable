CREATE TYPE analysis_type AS ENUM ('single_commit','commit_range','pull_request');
CREATE TYPE analysis_status AS ENUM ('pending','processing','completed','failed');

CREATE TABLE analyses (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    analysis_type   analysis_type NOT NULL,
    status          analysis_status NOT NULL DEFAULT 'pending',
    workspace       VARCHAR(255),
    repo_slug       VARCHAR(255),
    commit_hash     VARCHAR(40),
    from_hash       VARCHAR(40),
    to_hash         VARCHAR(40),
    pr_id           INTEGER,
    raw_diff        TEXT,
    diff_hash       VARCHAR(64),
    generated_desc  TEXT,
    model_used      VARCHAR(50),
    tokens_used     INTEGER,
    error_message   TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_analyses_diff_hash ON analyses(diff_hash) WHERE diff_hash IS NOT NULL;
CREATE INDEX idx_analyses_created_at ON analyses(created_at DESC);
CREATE INDEX idx_analyses_type ON analyses(analysis_type);
