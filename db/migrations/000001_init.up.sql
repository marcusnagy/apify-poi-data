-- 1) Ensure the PostGIS & H3 extensions exist
CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS h3;

-- 2) Create schema if not already present
CREATE SCHEMA IF NOT EXISTS poi_data_schema;

-- 3) Create the google_maps table
CREATE TABLE IF NOT EXISTS poi_data_schema.google_maps (
    id SERIAL PRIMARY KEY,
    search_string TEXT,
    rank INT,
    search_page_url TEXT,
    is_advertisement BOOLEAN,
    title TEXT,
    sub_title TEXT,
    price TEXT,
    category_name TEXT,
    address TEXT,
    neighborhood TEXT,
    street TEXT,
    city TEXT,
    postal_code TEXT,
    state TEXT,
    country_code TEXT,
    website TEXT,
    phone TEXT,
    phone_unformatted TEXT,
    claim_this_business BOOLEAN,
    location_lat DOUBLE PRECISION,
    location_lng DOUBLE PRECISION,
    total_score DOUBLE PRECISION,
    permanently_closed BOOLEAN,
    temporarily_closed BOOLEAN,
    place_id TEXT UNIQUE,     -- unique constraint on place_id
    categories TEXT[],
    fid TEXT,
    cid TEXT,
    reviews_count INT,
    images_count INT,
    image_categories TEXT[],
    scraped_at TIMESTAMPTZ,
    google_food_url TEXT,
    hotel_ads JSONB,
    opening_hours JSONB,
    people_also_search JSONB,
    places_tags JSONB,
    reviews_tags JSONB,
    additional_info JSONB,
    gas_prices JSONB,
    url TEXT,
    image_url TEXT,
    kgmid TEXT,
    h3_index TEXT,  -- new column for H3
    geom GEOMETRY,  -- geometry column
    search_page_loaded_url TEXT,
    description TEXT,
    located_in TEXT,
    plus_code TEXT,
    menu TEXT,
    reserve_table_url TEXT,
    hotel_stars TEXT,
    hotel_description TEXT,
    check_in_date TEXT,
    check_out_date TEXT,
    similar_hotels_nearby JSONB,
    hotel_review_summary JSONB,
    popular_times_live_text TEXT,
    popular_times_live_percent INT,
    popular_times_histogram JSONB,
    questions_and_answers JSONB,
    updates_from_customers JSONB,
    web_results JSONB,
    parent_place_url TEXT,
    table_reservation_links JSONB,
    booking_links JSONB,
    order_by JSONB,
    images TEXT,
    image_urls TEXT[],
    reviews JSONB,
    user_place_note JSONB,
    restaurant_data JSONB,
    owner_updates JSONB
);

-- 4) Add a PostGIS geometry column (WGS84)
ALTER TABLE poi_data_schema.google_maps
  ADD COLUMN IF NOT EXISTS geom geometry(Point, 4326);

-- 5) Add a column for H3 indexing ( for storing H3 index)
ALTER TABLE poi_data_schema.google_maps
  ADD COLUMN IF NOT EXISTS h3_index h3index;

-- 6) Create a spatial GIST index on the geometry column
CREATE INDEX IF NOT EXISTS idx_google_maps_geom
  ON poi_data_schema.google_maps
  USING GIST (geom);

-- 7) Create a standard b-tree (or hash) index for the H3 column
CREATE INDEX IF NOT EXISTS idx_google_maps_h3_index
  ON poi_data_schema.google_maps (h3_index);

-- 8) Create a GIN index for faster JSONB queries on additional_info
CREATE INDEX IF NOT EXISTS idx_google_maps_additional_info_gin
  ON poi_data_schema.google_maps
  USING GIN (additional_info jsonb_path_ops);

-- 9) Example tripadvisor table (currently empty, but schema is set up)
CREATE TABLE IF NOT EXISTS poi_data_schema.tripadvisor (
    -- define columns here...
);