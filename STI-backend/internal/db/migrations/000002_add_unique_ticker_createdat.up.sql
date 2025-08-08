ALTER TABLE stocks
ADD CONSTRAINT unique_ticker_created_at UNIQUE (ticker, created_at);