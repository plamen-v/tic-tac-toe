
INSERT INTO players(login, nickname, password)
VALUES ('player_1', 'Player 1', '$2a$10$7ZEpGQ9Kw75c5MnOQYYsk.3NYZznnVInfra.H1eFiwsZv4zehK53.'),
       ('player_2', 'Player 2', '$2a$10$e.AlVeSnp1.2b8g0o.pDT.Fqihe3KL2HjPdOGxnpx9lrhymTLJGnK'),
       ('player_3', 'Player 3', '$2a$10$mzrQg650tcwsvxcrtzS8pO6tKNrEt.nVM784L0oVRFLNm2CosxYQi'),
       ('player_4', 'Player 4', '$2a$10$tCC9TIgK29FYzhqA39k5s./Tx8MT0/UjDyNgp7bF1OgaDBANnozD2'),
       ('player_5', 'Player 5', '$2a$10$7ZEpGQ9Kw75c5MnOQYYsk.3NYZznnVInfra.H1eFiwsZv4zehK53.'),
       ('player_6', 'Player 6', '$2a$10$e.AlVeSnp1.2b8g0o.pDT.Fqihe3KL2HjPdOGxnpx9lrhymTLJGnK'),
       ('player_7', 'Player 7', '$2a$10$mzrQg650tcwsvxcrtzS8pO6tKNrEt.nVM784L0oVRFLNm2CosxYQi'),
       ('player_8', 'Player 8', '$2a$10$tCC9TIgK29FYzhqA39k5s./Tx8MT0/UjDyNgp7bF1OgaDBANnozD2'),
       ('player_9', 'Player 9', '$2a$10$7ZEpGQ9Kw75c5MnOQYYsk.3NYZznnVInfra.H1eFiwsZv4zehK53.'),
       ('player_10', 'Player 10', '$2a$10$e.AlVeSnp1.2b8g0o.pDT.Fqihe3KL2HjPdOGxnpx9lrhymTLJGnK'),
       ('player_11', 'Player 11', '$2a$10$mzrQg650tcwsvxcrtzS8pO6tKNrEt.nVM784L0oVRFLNm2CosxYQi'),
       ('player_12', 'Player 12', '$2a$10$tCC9TIgK29FYzhqA39k5s./Tx8MT0/UjDyNgp7bF1OgaDBANnozD2'),
       ('player_13', 'Player 13', '$2a$10$7ZEpGQ9Kw75c5MnOQYYsk.3NYZznnVInfra.H1eFiwsZv4zehK53.'),
       ('player_14', 'Player 14', '$2a$10$e.AlVeSnp1.2b8g0o.pDT.Fqihe3KL2HjPdOGxnpx9lrhymTLJGnK'),
       ('player_15', 'Player 15', '$2a$10$mzrQg650tcwsvxcrtzS8pO6tKNrEt.nVM784L0oVRFLNm2CosxYQi'),
       ('player_16', 'Player 16', '$2a$10$tCC9TIgK29FYzhqA39k5s./Tx8MT0/UjDyNgp7bF1OgaDBANnozD2');

INSERT INTO players_stats (player_id)
SELECT p.id 
FROM players p
WHERE p.login IN ('player_1', 'player_2','player_3','player_4','player_5',
 'player_6','player_7','player_8','player_9', 'player_10','player_11','player_12',
'player_13', 'player_14','player_15','player_16');
