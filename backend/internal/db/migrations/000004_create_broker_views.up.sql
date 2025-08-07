-- Vista broker_predictions usando el cierre más cercano por ticker
CREATE VIEW broker_predictions AS
SELECT
    s.brokerage,
    s.ticker,
    s.created_at AS prediction_date,
    s.target_to,
    s.target_from,
    f.close AS actual_price,

    CASE
        WHEN s.target_from IS NOT NULL AND s.target_to IS NOT NULL THEN
            CASE
                WHEN s.target_to > s.target_from THEN 'up'
                WHEN s.target_to < s.target_from THEN 'down'
                ELSE 'neutral'
                END
        ELSE NULL
        END AS prediction_direction,

    CASE
        WHEN s.target_from IS NOT NULL AND s.target_to IS NOT NULL AND f.close IS NOT NULL THEN
            CASE
                WHEN SIGN(s.target_to - s.target_from) = SIGN(f.close - s.target_from) THEN 1
                ELSE 0
                END
        ELSE NULL
        END AS is_correct,

    CASE
        WHEN s.target_to IS NOT NULL AND f.close IS NOT NULL AND f.close != 0 THEN
            ROUND(ABS(s.target_to - f.close) / f.close * 100, 2)
        ELSE NULL
        END AS error_percentage

FROM stocks s
         JOIN LATERAL (
    SELECT close
        FROM finances f
        WHERE f.ticker = s.ticker
        ORDER BY ABS(f.date - s.created_at::date) ASC
        LIMIT 1
        ) f ON true
        WHERE s.target_to IS NOT NULL;


CREATE VIEW broker_evaluation AS
SELECT
    brokerage,
    COUNT(*) FILTER (WHERE is_correct IS NOT NULL) AS total_predictions,
    SUM(is_correct) AS total_hits,
    ROUND(100.0 * SUM(is_correct)::float / NULLIF(COUNT(*) FILTER (WHERE is_correct IS NOT NULL)::float, 0.0), 2) AS accuracy,
    ROUND(SUM(is_correct)::float * (
        100.0 * SUM(is_correct)::float / NULLIF(COUNT(*) FILTER (WHERE is_correct IS NOT NULL)::float, 0.0)
    ) / 100.0, 2) AS weight_score
FROM
    broker_predictions
GROUP BY
    brokerage
ORDER BY
    weight_score DESC;
