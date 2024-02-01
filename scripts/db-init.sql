CREATE TABLE sensors (
    id SERIAL PRIMARY KEY,
    name VARCHAR UNIQUE,
    location GEOMETRY(Point,4326)  -- 4326 is the SRID for WGS84 (std GPS coordinate system)
);

CREATE TABLE tags (
    id SERIAL PRIMARY KEY,
    sensor_id INT REFERENCES sensors,
    value VARCHAR
);
