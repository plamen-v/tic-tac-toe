--INITIAL SQL
CREATE TABLE IF NOT EXISTS players (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    login varchar(256) NOT NULL,
    password VARCHAR(60) NOT NULL,
    nickname varchar(30) NOT NULL,
       
    UNIQUE(login),
    UNIQUE(nickname)
);

CREATE TABLE IF NOT EXISTS players_stats (
    player_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    wins INTEGER NOT NULL DEFAULT 0,
    losses INTEGER NOT NULL DEFAULT 0,
    draws INTEGER NOT NULL DEFAULT 0,

    CONSTRAINT players_stats_fk_player FOREIGN KEY (player_id) REFERENCES players(id)
);


CREATE TABLE IF NOT EXISTS rooms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    host_id UUID NOT NULL,
    host_request_new_game BOOLEAN NOT NULL DEFAULT false,
    guest_id UUID,
    guest_request_new_game BOOLEAN NOT NULL DEFAULT false,
    game_id UUID,     
    title VARCHAR(30) NOT NULL,
    description VARCHAR(150),
    phase INTEGER NOT NULL DEFAULT 0,
    
    CONSTRAINT rooms_fk_host FOREIGN KEY (host_id) REFERENCES players(id),
    CONSTRAINT rooms_fk_guest FOREIGN KEY (guest_id) REFERENCES players(id),
    UNIQUE(guest_id)
);

CREATE TABLE IF NOT EXISTS games (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    host_id UUID NOT NULL,
    host_mark CHAR(1) NOT NULL,
    guest_id UUID NOT NULL,
    guest_mark CHAR(1) NOT NULL,
    current_player_id UUID NOT NULL,
    board TEXT NOT NULL DEFAULT '_________',
    winner_id UUID,
    phase INTEGER NOT NULL DEFAULT 0,
    
    CONSTRAINT games_fk_host FOREIGN KEY (host_id) REFERENCES players(id),
    CONSTRAINT games_fk_guest FOREIGN KEY (guest_id) REFERENCES players(id),
    CONSTRAINT games_fk_current_player FOREIGN KEY (current_player_id) REFERENCES players(id),
    CONSTRAINT games_fk_winner FOREIGN KEY (winner_id) REFERENCES players(id),
    CHECK (host_mark IN ('X', 'O')),
    CHECK (guest_mark IN ('X', 'O'))
);

ALTER TABLE rooms DROP CONSTRAINT IF EXISTS rooms_fk_game;
ALTER TABLE rooms ADD CONSTRAINT rooms_fk_game FOREIGN KEY (game_id) REFERENCES games(id);
