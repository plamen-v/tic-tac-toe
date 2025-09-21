
INSERT INTO players(login, nickname, password)
VALUES ('player_1', 'Player 1', '$2a$10$7ZEpGQ9Kw75c5MnOQYYsk.3NYZznnVInfra.H1eFiwsZv4zehK53.'),
       ('player_2', 'Player 2', '$2a$10$e.AlVeSnp1.2b8g0o.pDT.Fqihe3KL2HjPdOGxnpx9lrhymTLJGnK'),
       ('player_3', 'Player 3', '$2a$10$mzrQg650tcwsvxcrtzS8pO6tKNrEt.nVM784L0oVRFLNm2CosxYQi'),
       ('player_4', 'Player 4', '$2a$10$tCC9TIgK29FYzhqA39k5s./Tx8MT0/UjDyNgp7bF1OgaDBANnozD2');

INSERT INTO players_stats (player_id)
SELECT p.id 
FROM players p
WHERE p.login IN ('player_1', 'player_2','player_3','player_4')
