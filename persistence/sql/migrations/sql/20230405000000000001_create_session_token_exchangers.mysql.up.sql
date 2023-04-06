CREATE TABLE IF NOT EXISTS session_token_exchangers (
    id CHAR(36) NOT NULL PRIMARY KEY,
    nid CHAR(36) NOT NULL,
    flow_id CHAR(36) NOT NULL,
    session_id CHAR(36) DEFAULT NULL,
    code VARCHAR(64) NOT NULL,

    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Relevant query:
--   SELECT * from session_token_exchangers
--   WHERE flow_id = ? AND nid = ? AND code = ? AND session_id IS NOT NULL AND code <> '';
-- CREATE INDEX session_token_exchangers_nid_flow_id_code_idx ON session_token_exchangers (nid, flow_id, code);
