CREATE TABLE IF NOT EXISTS finances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticker STRING NOT NULL,
    date DATE NOT NULL,
    open DECIMAL(10,2),
    high DECIMAL(10,2),
    low DECIMAL(10,2),
    close DECIMAL(10,2),
    volume BIGINT,
    source STRING,
    scraped_at TIMESTAMPTZ DEFAULT now(),
    UNIQUE (ticker, date)
);
