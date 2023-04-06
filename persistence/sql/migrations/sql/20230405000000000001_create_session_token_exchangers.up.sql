CREATE TABLE IF NOT EXISTS session_token_exchangers (
    "id" UUID NOT NULL PRIMARY KEY,
    "nid" UUID NOT NULL,
    "flow_id" UUID NOT NULL,
    "session_id" UUID DEFAULT NULL,
    "code" VARCHAR(64) NOT NULL,

    "created_at" timestamp NOT NULL,
    "updated_at" timestamp NOT NULL
);

-- Relevant query:
--   SELECT * from session_token_exchangers
--   WHERE flow_id = ? AND nid = ? AND code = ? AND session_id IS NOT NULL AND code <> '';
CREATE INDEX IF NOT EXISTS session_token_exchangers_nid_flow_id_code_idx ON session_token_exchangers (nid, flow_id, code);
