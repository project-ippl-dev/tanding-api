CREATE TABLE documents(
                          id BIGSERIAL PRIMARY KEY,
                          user_id UUID NOT NULL,
                          birth_certificate VARCHAR(255) NOT NULL,
                          family_card VARCHAR(255) NOT NULL,
                          user_identity VARCHAR(255) NOT NULL,
                          belt_certificate VARCHAR(255) NOT NULL,
                          elementary_certificate VARCHAR(255) NOT NULL,
                          middle_certificate VARCHAR(255) NOT NULL,
                          high_certificate VARCHAR(255) NOT NULL,
                          bachelor_certificate VARCHAR(255) NOT NULL,
                          master_certificate VARCHAR(255) NOT NULL,
                          created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                          updated_at TIMESTAMP NULL
);

CREATE INDEX idx_documents_user_id ON documents(user_id);
