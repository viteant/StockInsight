CREATE TABLE IF NOT EXISTS stocks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticker STRING NOT NULL,
    company STRING NOT NULL,
    brokerage STRING NOT NULL,
    action STRING NOT NULL,
    rating_from STRING,
    rating_to STRING,
    normalize_rating_from STRING,
    normalize_rating_to STRING,
    target_from DECIMAL(10,2),
    target_to DECIMAL(10,2),
    created_at TIMESTAMPTZ DEFAULT now()
);