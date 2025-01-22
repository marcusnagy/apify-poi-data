-- name: ListPOIInBoxWithCategoryH3 :many
WITH envelope_poly AS (
  SELECT ST_MakeEnvelope(
    $1::float8,  -- minX (western longitude)
    $2::float8,  -- minY (southern latitude)
    $3::float8,  -- maxX (eastern longitude)
    $4::float8,  -- maxY (northern latitude)
    4326
  ) AS poly
),
cells AS (
  SELECT unnest(h3_polyfill(poly, 9)) AS h3_cell  -- pick H3 res=9 for example
  FROM envelope_poly
)
SELECT gm.*
FROM poi_data_schema.google_maps gm
JOIN cells c ON gm.h3_index = c.h3_cell
WHERE EXISTS (
  SELECT 1
  FROM unnest(gm.categories) cat
  WHERE cat ILIKE '%' || $5 || '%'
);