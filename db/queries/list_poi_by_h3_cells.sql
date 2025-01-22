-- name: ListPOIsByH3Cells :many
WITH parent_cells AS (
    SELECT unnest($2::text[])::h3index AS parent_cell
)
SELECT g.*
FROM poi_data_schema.google_maps g
JOIN LATERAL (
    SELECT h3_cell_to_children(pc.parent_cell, $1::int) AS child_index
    FROM parent_cells pc
) children
ON g.h3_index = children.child_index;