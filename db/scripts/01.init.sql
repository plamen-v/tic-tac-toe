--INITIAL SQL
CREATE SEQUENCE IF NOT EXISTS players_id_seq;
CREATE TABLE IF NOT EXISTS players (
    id integer NOT NULL DEFAULT nextval('players_id_seq'),
    login varchar(256) NOT NULL,
    nickname varchar(256) NOT NULL,
    password VARCHAR(60) NOT NULL,
    wins INTEGER NOT NULL DEFAULT 0,
    losses INTEGER NOT NULL DEFAULT 0,
    draws INTEGER NOT NULL DEFAULT 0,
    
    CONSTRAINT players_pk PRIMARY KEY(id),
    UNIQUE(login),
    UNIQUE(nickname)
);
ALTER SEQUENCE players_id_seq
OWNED BY players.id;

CREATE SEQUENCE IF NOT EXISTS rooms_id_seq;
CREATE TABLE IF NOT EXISTS rooms (
    id INTEGER NOT NULL DEFAULT nextval('rooms_id_seq'),
    host_id INTEGER NOT NULL,
    opponent_id INTEGER,
    game_id INTEGER, 
    title VARCHAR(55) NOT NULL,
    info VARCHAR(1024),
    status INTEGER NOT NULL DEFAULT 0,
    
    CONSTRAINT rooms_pk PRIMARY KEY(id),
    CONSTRAINT rooms_fk_host FOREIGN KEY (host_id) REFERENCES players(id),
    CONSTRAINT rooms_fk_opponent FOREIGN KEY (opponent_id) REFERENCES players(id)
);
ALTER SEQUENCE rooms_id_seq
OWNED BY rooms.id;

CREATE SEQUENCE IF NOT EXISTS games_id_seq;
CREATE TABLE IF NOT EXISTS games (
    id INTEGER NOT NULL DEFAULT nextval('games_id_seq'),
    room_id INTEGER NOT NULL,
    player_1_id INTEGER NOT NULL,
    player_1_ready BOOLEAN,
    player_1_mark CHAR(1),
    player_2_id INTEGER NOT NULL,
    player_2_ready BOOLEAN,
    player_2_mark CHAR(1),
    current_player_id INTEGER,
    board TEXT DEFAULT '_________',
    status INTEGER NOT NULL DEFAULT 0,
    winner_id INTEGER,
    loser_id INTEGER,
    
    CONSTRAINT games_pk PRIMARY KEY(id),
    CONSTRAINT games_fk_player_1 FOREIGN KEY (player_1_id) REFERENCES players(id),
    CONSTRAINT games_fk_player_2 FOREIGN KEY (player_2_id) REFERENCES players(id),
    CONSTRAINT games_fk_current_player FOREIGN KEY (current_player_id) REFERENCES players(id),
    CONSTRAINT games_fk_winner FOREIGN KEY (winner_id) REFERENCES players(id),
    CONSTRAINT games_fk_loser FOREIGN KEY (loser_id) REFERENCES players(id),
    CHECK (player_1_mark IN ('X', 'O')),
    CHECK (player_2_mark IN ('X', 'O')),
    UNIQUE (room_id, player_1_id, player_2_id)
);
ALTER SEQUENCE games_id_seq
OWNED BY games.id;

ALTER TABLE rooms DROP CONSTRAINT IF EXISTS rooms_fk_game;
ALTER TABLE rooms ADD CONSTRAINT rooms_fk_game FOREIGN KEY (game_id) REFERENCES games(id);

ALTER TABLE games DROP CONSTRAINT IF EXISTS games_fk_room;
ALTER TABLE games ADD CONSTRAINT games_fk_room FOREIGN KEY (room_id) REFERENCES rooms(id);