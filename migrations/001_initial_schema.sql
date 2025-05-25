-- Удаляем старые таблицы, если они существуют, для чистого старта (ОСТОРОЖНО!)
DROP TABLE IF EXISTS reviews CASCADE;
DROP TABLE IF EXISTS booking_places CASCADE;
DROP TABLE IF EXISTS bookings CASCADE;
DROP TABLE IF EXISTS menu_items CASCADE;
DROP TABLE IF EXISTS places CASCADE;
DROP TABLE IF EXISTS halls CASCADE;
DROP TABLE IF EXISTS establishments CASCADE;
DROP TABLE IF EXISTS users CASCADE;

DROP TYPE IF EXISTS user_role_enum CASCADE;
DROP TYPE IF EXISTS establishment_type_enum CASCADE;
DROP TYPE IF EXISTS place_status_enum CASCADE;
DROP TYPE IF EXISTS booking_status_enum CASCADE;


-- Создание типов ENUM
CREATE TYPE user_role_enum AS ENUM ('user', 'admin'); -- 'superadmin' можно добавить позже
CREATE TYPE establishment_type_enum AS ENUM ('restaurant', 'coworking');
CREATE TYPE place_status_enum AS ENUM ('free', 'booked', 'occupied', 'unavailable');
CREATE TYPE booking_status_enum AS ENUM ('pending', 'confirmed', 'cancelled', 'completed');

-- Таблица пользователей
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    phone_number VARCHAR(20) NOT NULL UNIQUE,
    email VARCHAR(255) UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role user_role_enum NOT NULL DEFAULT 'user', -- ДОБАВЛЕНО поле роли
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Таблица заведений
CREATE TABLE IF NOT EXISTS establishments (
    id BIGSERIAL PRIMARY KEY,
    owner_user_id BIGINT REFERENCES users(id) ON DELETE SET NULL, -- ДОБАВЛЕНО: ID владельца. SET NULL если админ удален. Или ON DELETE CASCADE, если заведение должно удалиться.
    name VARCHAR(255) NOT NULL,
    type establishment_type_enum NOT NULL,
    address TEXT,
    working_hours VARCHAR(255),
    description TEXT,
    photos TEXT[], 
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_establishments_owner ON establishments (owner_user_id);

-- Таблица залов
CREATE TABLE IF NOT EXISTS halls (
    id BIGSERIAL PRIMARY KEY,
    establishment_id BIGINT NOT NULL REFERENCES establishments(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    capacity INT,
    has_air_conditioner BOOLEAN DEFAULT FALSE,
    photos TEXT[],
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_halls_establishment_id ON halls (establishment_id);

-- Таблица мест
CREATE TABLE IF NOT EXISTS places (
    id BIGSERIAL PRIMARY KEY,
    hall_id BIGINT NOT NULL REFERENCES halls(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL, 
    type VARCHAR(100), -- Тип места (стол на 2-х, рабочее место и т.д.)
    status place_status_enum NOT NULL DEFAULT 'free',
    visual_info JSONB, -- { "x": 10, "y": 20, "width": 5, "height": 5, "shape": "rect" / "circle" }
    icon_free_url TEXT,
    icon_booked_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(hall_id, name) 
);
CREATE INDEX idx_places_hall_id ON places (hall_id);

-- Таблица элементов меню
CREATE TABLE IF NOT EXISTS menu_items (
    id BIGSERIAL PRIMARY KEY,
    establishment_id BIGINT NOT NULL REFERENCES establishments(id) ON DELETE CASCADE,
    category_name VARCHAR(100) NOT NULL, -- Было category, переименовал для ясности
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    description TEXT,
    photo_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_menu_items_establishment_category ON menu_items (establishment_id, category_name);

-- Таблица бронирований
CREATE TABLE IF NOT EXISTS bookings (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    establishment_id BIGINT NOT NULL REFERENCES establishments(id) ON DELETE CASCADE,
    hall_id BIGINT NOT NULL REFERENCES halls(id) ON DELETE CASCADE,
    booking_time TIMESTAMPTZ NOT NULL,
    duration_minutes INT, 
    people_count INT NOT NULL CHECK (people_count > 0),
    status booking_status_enum NOT NULL DEFAULT 'pending',
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_bookings_user_time ON bookings (user_id, booking_time);
CREATE INDEX idx_bookings_establishment_time ON bookings (establishment_id, booking_time);

-- Связующая таблица для бронирований и мест (многие-ко-многим)
CREATE TABLE IF NOT EXISTS booking_places (
    booking_id BIGINT NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
    place_id BIGINT NOT NULL REFERENCES places(id) ON DELETE RESTRICT, 
    PRIMARY KEY (booking_id, place_id)
);

-- Таблица отзывов
CREATE TABLE IF NOT EXISTS reviews (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    establishment_id BIGINT NOT NULL REFERENCES establishments(id) ON DELETE CASCADE,
    rating INT NOT NULL CHECK (rating >= 1 AND rating <= 5),
    text TEXT NOT NULL,
    photo_url TEXT,
    is_moderated BOOLEAN DEFAULT FALSE,
    is_approved BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_reviews_establishment_approved ON reviews (establishment_id, is_approved, created_at DESC);

-- Функция для автоматического обновления updated_at
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Триггеры для таблиц
CREATE TRIGGER set_timestamp_users BEFORE UPDATE ON users FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();
CREATE TRIGGER set_timestamp_establishments BEFORE UPDATE ON establishments FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();
CREATE TRIGGER set_timestamp_halls BEFORE UPDATE ON halls FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();
CREATE TRIGGER set_timestamp_places BEFORE UPDATE ON places FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();
CREATE TRIGGER set_timestamp_menu_items BEFORE UPDATE ON menu_items FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();
CREATE TRIGGER set_timestamp_bookings BEFORE UPDATE ON bookings FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();

-- ДОБАВЛЕНИЕ ТЕСТОВЫХ АДМИНИСТРАТОРОВ И ЗАВЕДЕНИЙ
-- Пароль для всех админов будет 'adminpass' (хэш для него нужно сгенерировать)
-- Хэш для 'adminpass': $2a$10$EazevzL3aL1jV1eVbbJ43OsdP7sYqjHpl2KxUKn21j3pYJgHb8K2K (сгенерирован с bcrypt.DefaultCost)

INSERT INTO users (name, phone_number, email, password_hash, role) VALUES
('Admin Vista', '+79990000001', 'admin_vista@minity.app', '$2a$10$EazevzL3aL1jV1eVbbJ43OsdP7sYqjHpl2KxUKn21j3pYJgHb8K2K', 'admin'),
('Admin MostHub', '+79990000002', 'admin_mosthub@minity.app', '$2a$10$EazevzL3aL1jV1eVbbJ43OsdP7sYqjHpl2KxUKn21j3pYJgHb8K2K', 'admin'),
('Regular User', '+79990000003', 'user@minity.app', '$2a$10$EazevzL3aL1jV1eVbbJ43OsdP7sYqjHpl2KxUKn21j3pYJgHb8K2K', 'user'); -- Пароль 'adminpass' для простоты теста

-- Получаем ID администраторов (в реальном приложении это делается не так прямолинейно в миграции)
-- Для простоты предположим, что они получили ID 1 и 2.
-- В реальной миграции лучше использовать SELECT id FROM users WHERE email = '...'

INSERT INTO establishments (owner_user_id, name, type, address, working_hours, description, photos) VALUES
(
    (SELECT id FROM users WHERE email = 'admin_vista@minity.app'),
    'Ресторан "Vista"',
    'restaurant',
    'ул. Небесная, 100, Этаж 50, Город Панорам',
    '12:00 - 00:00 (Ежедневно)',
    'Панорамный ресторан "Vista" предлагает изысканные блюда европейской кухни и потрясающий вид на город. Идеальное место для романтического ужина или деловой встречи.',
    ARRAY['https://i.ibb.co/QRYvJTp/vista.jpg']
),
(
    (SELECT id FROM users WHERE email = 'admin_mosthub@minity.app'),
    'Коворкинг "Most IT Hub"',
    'coworking',
    'пр. IT Разработчиков, 123, Офис 101, Техноград',
    '09:00 - 21:00 (Пн-Сб), Вс - выходной',
    'Современный и технологичный коворкинг в центре города для продуктивной работы. Высокоскоростной интернет, удобные рабочие места, переговорные комнаты и зоны отдыха.',
    ARRAY['https://i.ibb.co/DPjqjR0h/MOST-IT-Hub-Almaty-I-2.jpg']
);

-- Добавим тестовые залы и места для Vista (ID заведения предположительно 1)
INSERT INTO halls (establishment_id, name, description, capacity) VALUES
( (SELECT id FROM establishments WHERE name = 'Ресторан "Vista"'), 'Основной зал Vista', 'Просторный зал с видом на город', 50),
( (SELECT id FROM establishments WHERE name = 'Ресторан "Vista"'), 'VIP-комната Vista', 'Уединенная комната для особых случаев', 10);

INSERT INTO places (hall_id, name, type, status, visual_info) VALUES
( (SELECT id FROM halls WHERE name = 'Основной зал Vista'), 'V1', 'table-2', 'free', '{"x": 10, "y": 10}'::jsonb),
( (SELECT id FROM halls WHERE name = 'Основной зал Vista'), 'V2', 'table-4', 'booked', '{"x": 30, "y": 10}'::jsonb),
( (SELECT id FROM halls WHERE name = 'Основной зал Vista'), 'V3', 'table-2', 'free', '{"x": 50, "y": 10}'::jsonb);

-- Добавим тестовые залы и места для Most IT Hub (ID заведения предположительно 2)
INSERT INTO halls (establishment_id, name, description, capacity) VALUES
( (SELECT id FROM establishments WHERE name = 'Коворкинг "Most IT Hub"'), 'Open Space MostHub', 'Большое открытое пространство с рабочими местами', 40);

INSERT INTO places (hall_id, name, type, status, visual_info) VALUES
( (SELECT id FROM halls WHERE name = 'Open Space MostHub'), 'MW1', 'desk', 'free', '{"x": 5, "y": 10}'::jsonb),
( (SELECT id FROM halls WHERE name = 'Open Space MostHub'), 'MW2', 'desk', 'free', '{"x": 20, "y": 10}'::jsonb),
( (SELECT id FROM halls WHERE name = 'Open Space MostHub'), 'MW3', 'desk-large', 'booked', '{"x": 35, "y": 10}'::jsonb);

-- Добавим тестовое меню для Vista
INSERT INTO menu_items (establishment_id, category_name, name, price, description) VALUES
( (SELECT id FROM establishments WHERE name = 'Ресторан "Vista"'), 'Закуски', 'Брускетта с томатами и базиликом', 350.00, 'Хрустящий хлеб, свежие томаты, базилик, оливковое масло Extra Virgin.'),
( (SELECT id FROM establishments WHERE name = 'Ресторан "Vista"'), 'Закуски', 'Карпаччо из говядины', 550.00, 'Тонко нарезанная говяжья вырезка с рукколой, пармезаном и бальзамическим кремом.'),
( (SELECT id FROM establishments WHERE name = 'Ресторан "Vista"'), 'Горячие блюда', 'Стейк Рибай', 1200.00, 'Сочный стейк из мраморной говядины с овощами гриль.');