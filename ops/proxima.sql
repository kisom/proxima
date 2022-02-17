CREATE TABLE updates (
-- {
-- 	"drift":-0.588729962,
-- 	"elapsed":5.588729962,
-- 	"mission": {
-- 		"action":"accelerating",
-- 		"clock": {
-- 			"launched":"2022-02-16 10:01",
-- 			"observer":"2022-02-16 10:01",
-- 			"observer_et":"6",
-- 			"relative":"2022-02-16 10:01",
-- 			"relative_et":"6"
-- 		},
-- 		"state": {
-- 			"v":"11229",
-- 			"x":"56025"
-- 		}
-- 	}
-- }
    id		    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created		BIGINT NOT NULL,
    drift		REAL NOT NULL, -- drift between observer and earth
    elapsed 	REAL NOT NULL, -- mission elapsed time
    action		TEXT NOT NULL, -- current mission phase
    launched	TEXT NOT NULL, -- timestamp
    observer	TEXT NOT NULL, -- observer's time
    observer_et	REAL NOT NULL, -- observer's elapsed time
    relative	TEXT NOT NULL, -- earth clock
    relative_et	REAL NOT NULL, -- earth elapsed time
    -- velocity and distance are big numbers, so we keep them in text
    -- form so they can be loaded into a big.Rat.
    velocity	TEXT NOT NULL, -- speed in m/s
    distance	TEXT NOT NULL  -- distance traveled in meters
);
