CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    surname VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(20) UNIQUE NOT NULL,
    role VARCHAR(50) NOT NULL,
    password VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS user_tokens (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    role VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    access_token VARCHAR(255) UNIQUE NOT NULL,
    refresh_token VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE if not exists photos (
    id SERIAL primary key,
    route varchar not null
);

CREATE OR REPLACE FUNCTION update_user_token_updated_at()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER user_tokens_update_updated_at
    BEFORE UPDATE ON user_tokens
    FOR EACH ROW
EXECUTE FUNCTION update_user_token_updated_at();


CREATE TABLE IF NOT EXISTS restaurants (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    address VARCHAR(255) NOT NULL,
    description TEXT,
    city VARCHAR(100),
    status BOOLEAN DEFAULT TRUE,
    phone VARCHAR(100) NOT NULL,
    owner_id INTEGER NOT NULL,
    mode_from TIMESTAMP NOT NULL,
    mode_to TIMESTAMP NOT NULL,
    icon_id integer,
    FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE,
    foreign key (icon_id) references photos(id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS restaurant_photos (
     id SERIAL PRIMARY KEY,
    photo_id integer not null,
     restaurant_id INTEGER NOT NULL,
     FOREIGN KEY (restaurant_id) REFERENCES restaurants(id) ON DELETE CASCADE,
    foreign key (photo_id) references photos(id) on delete cascade
);

CREATE table if not exists services (
  id SERIAL PRIMARY KEY,
  name VARCHAR NOT NULL UNIQUE
);

create table if not exists restaurant_service (
    id serial primary key,
    service_id integer not null,
    restaurant_id integer not null,
    foreign key (service_id) references services(id) on delete cascade,
    foreign key (restaurant_id) references restaurants(id) on DELETE cascade
);

CREATE TABLE IF NOT EXISTS favorite_restaurants (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    restaurant_id INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (restaurant_id) REFERENCES restaurants(id) on delete cascade
);

CREATE TABLE IF NOT EXISTS tables (
     id SERIAL PRIMARY KEY,
     name VARCHAR(255) NOT NULL,
     type VARCHAR(100) NOT NULL,
     description TEXT,
     capacity INTEGER NOT NULL,
     photo_id integer,
     restaurant_id INTEGER NOT NULL,
     FOREIGN KEY (restaurant_id) REFERENCES restaurants(id) ON DELETE CASCADE,
    FOREIGN KEY (photo_id) references photos(id) on delete set null
);

CREATE TABLE IF NOT EXISTS orders (
     id SERIAL PRIMARY KEY,
     restaurant_id INTEGER NOT NULL,
     total_sum FLOAT NOT NULL,
     user_id INTEGER,
     table_id INTEGER,
     date TIMESTAMP NOT NULL,
     status VARCHAR(100) NOT NULL,
     FOREIGN KEY (restaurant_id) REFERENCES restaurants(id) ON DELETE CASCADE,
     FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
     FOREIGN KEY (table_id) REFERENCES tables(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS foods (
      id SERIAL PRIMARY KEY,
      name VARCHAR(255) NOT NULL,
      type VARCHAR(100) NOT NULL,
      description TEXT,
      price FLOAT NOT NULL,
      available BOOLEAN NOT NULL DEFAULT TRUE,
      photo_id integer,
      restaurant_id INTEGER NOT NULL,
      FOREIGN KEY (restaurant_id) REFERENCES restaurants(id) ON DELETE CASCADE,
    foreign key (photo_id) references photos(id) on delete set null
);

CREATE TABLE IF NOT EXISTS order_foods (
     id SERIAL PRIMARY KEY,
     order_id INTEGER NOT NULL,
     food_id INTEGER NOT NULL,
     FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
     FOREIGN KEY (food_id) REFERENCES foods(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS restaurant_reviews (
    id SERIAL PRIMARY KEY,
    stars INTEGER NOT NULL,
    description VARCHAR,
    user_id INTEGER NOT NULL,
    restaurant_id INTEGER NOT NULL,
    date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (restaurant_id) REFERENCES restaurants(id) ON DELETE CASCADE
)