DROP TRIGGER IF EXISTS user_tokens_update_updated_at ON user_tokens;

DROP FUNCTION IF EXISTS update_user_token_updated_at();

DROP TABLE IF EXISTS restaurant_reviews;
DROP TABLE IF EXISTS order_foods;
DROP TABLE IF EXISTS foods;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS tables;
DROP TABLE IF EXISTS favorite_restaurants;
DROP TABLE IF EXISTS restaurant_service;
DROP TABLE IF EXISTS services;
DROP TABLE IF EXISTS restaurant_photos;
DROP TABLE IF EXISTS restaurants;
DROP TABLE IF EXISTS user_tokens;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS photos;
