-- noinspection SqlDialectInspectionForFile

CREATE TABLE if not exists Accounts (
   id SERIAL NOT NULL,
   admin_user_id    integer,
   account_name           TEXT    NOT NULL UNIQUE,
   PRIMARY KEY(id)
);


CREATE TABLE if not exists Images (
   id SERIAL NOT NULL,
   format         TEXT    NOT NULL,
   imageData      BYTEA,
   PRIMARY KEY(id)
);



CREATE TABLE if not exists Players (
   id SERIAL NOT NULL,
   account_id integer NOT NULL,
   player_name    TEXT    NOT NULL,
   email          TEXT,
   phone          TEXT,
   paid           integer NOT NULL DEFAULT 0,
   image_id INTEGER NOT NULL DEFAULT 0,
   PRIMARY KEY(id)
);

CREATE TABLE if not exists Users (
   id SERIAL PRIMARY KEY NOT NULL,
   account_id integer NOT NULL,
   subject           TEXT    NOT NULL UNIQUE,
   email             TEXT    NOT NULL UNIQUE,
   provider          TEXT,
   given_name        TEXT    NOT NULL,
   preferred_handle  TEXT    NOT NULL,
   picture_url       TEXT,
   stack          integer NOT NULL DEFAULT 0,
   last_login  timestamp without time zone DEFAULT CURRENT_TIMESTAMP

);

CREATE TABLE if not exists UserSessions (
   session_id  SERIAL PRIMARY KEY NOT NULL,
   account_id integer NOT NULL,
   subject        TEXT    NOT NULL UNIQUE,
   email          TEXT    NOT NULL UNIQUE,
   provider       TEXT,
   given_name     TEXT    NOT NULL UNIQUE,
   picture_url    TEXT
);

CREATE TABLE if not exists Tables (
   id  SERIAL PRIMARY KEY NOT NULL,
   button integer NOT NULL,
   num_seats  integer NOT NULL,
   bigblind integer NOT NULL,
   channel TEXT
);

CREATE TABLE if not exists  Seats (
   table_id integer NOT NULL references Users(id),
   num integer NOT NULL,
   empty boolean NOT NULL default true,
   session_id  TEXT NOT NULL,
   sitting_out boolean NOT NULL default false,
   PRIMARY KEY (table_id, num)
);


CREATE TABLE if not exists  TablePlayerState (
    user_id         integer NOT NULL references Users(id),
    table_id        integer,
    seatnum         integer,
    stack_change    integer,
    starting_stack  integer,
    CONSTRAINT TablePlayerState_pkey PRIMARY KEY(user_id, table_id, seatnum)
);


CREATE TABLE if not exists Transactions (
   id               SERIAL PRIMARY KEY NOT NULL,
   user_id          integer NOT NULL references Users(id),
   table_id         integer NOT NULL,
   typ              integer NOT NULL,
   amount           integer NOT NULL,
   transaction_date timestamp NOT NULL DEFAULT NOW()
);

CREATE TABLE if not exists TableRecovery (
   id               SERIAL PRIMARY KEY NOT NULL,
   table_id         integer NOT NULL
);

