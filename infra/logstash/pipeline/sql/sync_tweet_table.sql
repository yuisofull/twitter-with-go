SELECT *, updated_at AS ts FROM temp_db.tweets
WHERE updated_at > :sql_last_value AND updated_at < NOW()
ORDER BY id