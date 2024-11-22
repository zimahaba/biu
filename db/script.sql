
CREATE TABLE user_credentials (
    id            SERIAL NOT NULL,
    username      TEXT NOT NULL,
    password      TEXT NOT NULL,
    refresh_token text,
    app_user_id   SERIAL NOT NULL
);

CREATE TABLE app_user (
    id                  SERIAL NOT NULL,
    name                TEXT   NOT NULL,
    email               TEXT   NOT NULL,
    birthday            DATE
);

ALTER TABLE app_user         ADD CONSTRAINT app_user_pkey         PRIMARY KEY (id);
ALTER TABLE user_credentials ADD CONSTRAINT user_credentials_pkey PRIMARY KEY (id);

ALTER TABLE user_credentials ADD CONSTRAINT user_credentials_app_user_fkey FOREIGN KEY (app_user_id) REFERENCES app_user (id);

ALTER TABLE user_credentials ADD CONSTRAINT user_credentials_username_unq UNIQUE (username);
ALTER TABLE app_user         ADD CONSTRAINT user_email_unq                UNIQUE (email);