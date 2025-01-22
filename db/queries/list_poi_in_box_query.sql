-- name: ListPOIInBox :many
SELECT *
FROM poi_data_schema.google_maps
WHERE ST_Contains(
    ST_MakeEnvelope($1::float8, $2::float8, $3::float8, $4::float8, 4326),
    geom
);