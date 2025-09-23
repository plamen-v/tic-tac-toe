--INITIAL SQL
CREATE SEQUENCE IF NOT EXISTS players_id_seq;
CREATE TABLE IF NOT EXISTS players (
    id integer NOT NULL DEFAULT nextval('players_id_seq'),
    login varchar(256) NOT NULL,
    password VARCHAR(60) NOT NULL,
    nickname varchar(30) NOT NULL,
       
    CONSTRAINT players_pk PRIMARY KEY(id),
    UNIQUE(login),
    UNIQUE(nickname)
);
ALTER SEQUENCE players_id_seq
OWNED BY players.id;

CREATE TABLE IF NOT EXISTS players_stats (
    player_id INTEGER PRIMARY KEY,
    wins INTEGER NOT NULL DEFAULT 0,
    losses INTEGER NOT NULL DEFAULT 0,
    draws INTEGER NOT NULL DEFAULT 0,

    CONSTRAINT players_stats_fk_player FOREIGN KEY (player_id) REFERENCES players(id)
);


CREATE TABLE IF NOT EXISTS rooms (
    host_id INTEGER PRIMARY KEY,
    host_continue BOOLEAN NOT NULL DEFAULT false,
    guest_id INTEGER,
    guest_continue BOOLEAN NOT NULL DEFAULT false,
    game_id INTEGER,     
    title VARCHAR(30) NOT NULL,
    description VARCHAR(150),
    phase INTEGER NOT NULL DEFAULT 0,
    
    CONSTRAINT rooms_fk_host FOREIGN KEY (host_id) REFERENCES players(id),
    CONSTRAINT rooms_fk_guest FOREIGN KEY (guest_id) REFERENCES players(id),
    UNIQUE(guest_id)
);

CREATE SEQUENCE IF NOT EXISTS games_id_seq;
CREATE TABLE IF NOT EXISTS games (
    id BIGINT NOT NULL DEFAULT nextval('games_id_seq'),
    host_id INTEGER NOT NULL,
    host_mark CHAR(1) NOT NULL,
    guest_id INTEGER NOT NULL,
    guest_mark CHAR(1) NOT NULL,
    current_player_id INTEGER NOT NULL,
    board TEXT NOT NULL DEFAULT '_________',
    winner_id INTEGER,
    phase INTEGER NOT NULL DEFAULT 0,
    
    CONSTRAINT games_pk PRIMARY KEY(id),
    CONSTRAINT games_fk_host FOREIGN KEY (host_id) REFERENCES players(id),
    CONSTRAINT games_fk_guest FOREIGN KEY (guest_id) REFERENCES players(id),
    CONSTRAINT games_fk_current_player FOREIGN KEY (current_player_id) REFERENCES players(id),
    CONSTRAINT games_fk_winner FOREIGN KEY (winner_id) REFERENCES players(id),
    CHECK (host_mark IN ('X', 'O')),
    CHECK (guest_mark IN ('X', 'O'))
);
ALTER SEQUENCE games_id_seq
OWNED BY games.id;

ALTER TABLE rooms DROP CONSTRAINT IF EXISTS rooms_fk_game;
ALTER TABLE rooms ADD CONSTRAINT rooms_fk_game FOREIGN KEY (game_id) REFERENCES games(id);
