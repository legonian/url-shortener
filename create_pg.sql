--to have access to gen_random_bytes(int)
DROP EXTENSION IF EXISTS pgcrypto;
CREATE EXTENSION pgcrypto;

--create random char generator function
CREATE OR REPLACE FUNCTION generate_uid(size INT) RETURNS TEXT AS $$
DECLARE
  characters TEXT := 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
  bytes BYTEA := gen_random_bytes(size);
  l INT := length(characters);
  i INT := 0;
  output TEXT := '';
BEGIN
  WHILE i < size LOOP
    output := output || substr(characters, get_byte(bytes, i) % l + 1, 1);
    i := i + 1;
  END LOOP;
  RETURN output;
END;
$$ LANGUAGE plpgsql VOLATILE;

--create table with URLs
DROP TABLE IF EXISTS short_urls;
CREATE TABLE IF NOT EXISTS short_urls (
  url_id int primary key generated always as identity,
  short_url varchar(10) not null UNIQUE DEFAULT generate_uid(6),
  full_url varchar(100) not null,
  views_count int not null default 0 check (0 <= views_count),
  create_date timestamp not null default current_timestamp
);

--function to add short URL
DROP FUNCTION IF EXISTS add_url(character varying);
CREATE FUNCTION add_url(IN f_url varchar(100), OUT short_url varchar(10), OUT full_url varchar(100), OUT views_count int) AS $$
  INSERT INTO short_urls (full_url)
  VALUES (f_url)
  RETURNING short_url, full_url, views_count;
$$ LANGUAGE SQL;

--function to read short URL
DROP FUNCTION IF EXISTS get_full_url(character varying);
CREATE FUNCTION get_full_url(IN s_url varchar(10), OUT short_url varchar(10), OUT full_url varchar(100), OUT views_count int) AS $$
  UPDATE short_urls SET views_count = views_count+1 WHERE short_url = s_url;
  
  SELECT short_url, full_url, views_count FROM short_urls
  WHERE short_url = s_url
$$ LANGUAGE SQL;

--usage:
--select * from add_url('https://www.google.com/')
--select * from get_full_url('O1yQQM')
