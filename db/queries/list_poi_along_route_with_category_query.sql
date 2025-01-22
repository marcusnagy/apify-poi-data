-- name: ListPOIAlongRouteWithCategoryH3 :many
WITH route_line AS (
  SELECT ST_MakeLine(
    ST_SetSRID(ST_Point($1::float8, $2::float8), 4326),
    ST_SetSRID(ST_Point($3::float8, $4::float8), 4326)
  ) AS geom
),
corridor_poly AS (
  -- Buffer in meters around the route, then convert back to geometry
  SELECT ST_Buffer(geom::geography, $5::float8)::geometry AS poly
  FROM route_line
),
cells AS (
  -- Fill the buffered corridor polygon with H3 cells
  SELECT unnest(h3_polyfill(poly, 9)) AS h3_cell
  FROM corridor_poly
)
SELECT gm.*
FROM poi_data_schema.google_maps gm
JOIN cells c ON gm.h3_index = c.h3_cell
WHERE EXISTS (
  SELECT 1
  FROM unnest(gm.categories) cat
  WHERE cat ILIKE '%' || $6 || '%'
);