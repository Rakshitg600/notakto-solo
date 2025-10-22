-- WARNING: This schema is for context only and is not meant to be run.
-- Table order and constraints may not be valid for execution.

CREATE TABLE player (
  uid character varying NOT NULL DEFAULT nextval('player_uid_seq'::regclass),
  name character varying NOT NULL,
  email character varying NOT NULL UNIQUE,
  profile_pic text,
  CONSTRAINT player_pkey PRIMARY KEY (uid)
);
CREATE TABLE session (
  session_id character varying NOT NULL DEFAULT nextval('session_session_id_seq'::regclass),
  uid character varying NOT NULL,
  expired boolean DEFAULT false,
  created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT session_pkey PRIMARY KEY (session_id),
  CONSTRAINT session_uid_fkey FOREIGN KEY (uid) REFERENCES player(uid)
);
CREATE TABLE sessionstate (
  session_id character varying NOT NULL,
  boards ARRAY,
  current_player integer,
  winner character varying,
  board_size integer,
  number_of_boards integer,
  difficulty integer,
  game_history ARRAY,
  gameover boolean DEFAULT false,
  CONSTRAINT sessionstate_pkey PRIMARY KEY (session_id),
  CONSTRAINT sessionstate_session_id_fkey FOREIGN KEY (session_id) REFERENCES session(session_id)
);
CREATE TABLE wallet (
  uid character varying NOT NULL,
  coins integer DEFAULT 0,
  xp integer DEFAULT 0,
  CONSTRAINT wallet_pkey PRIMARY KEY (uid),
  CONSTRAINT wallet_uid_fkey FOREIGN KEY (uid) REFERENCES player(uid)
);