CREATE DATABASE my_db;

\c my_db

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'keril') THEN
        CREATE USER keril with password 'pass';
END IF;
END $$;

GRANT ALL PRIVILEGES ON DATABASE my_db TO keril;

ALTER DATABASE my_db OWNER TO keril;

CREATE SCHEMA IF NOT EXISTS diplom_raw AUTHORIZATION keril;

CREATE TABLE IF NOT EXISTS diplom_raw.trips (
    id SERIAL PRIMARY KEY,
    "Дата и время начала рейса" TIMESTAMP WITHOUT TIME ZONE,
    "Номер вагона" VARCHAR(8),
    "Дорога отправления" VARCHAR(100),
    "Дорога назначения" VARCHAR(100),
    "Номер накладной" VARCHAR(50),
    "Станция отправления" VARCHAR(255),
    "Станция назначения" VARCHAR(255),
    "Наименование груза" VARCHAR(255),
    "Грузоотправитель" VARCHAR(255),
    "Грузополучатель" VARCHAR(255),
    "Тип парка (М/Т)" VARCHAR(50),
    "Тип парка (П/Г)" VARCHAR(50),
    "Время загрузки данных" TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() at time zone 'utc+3')
    );
ALTER TABLE diplom_raw.trips OWNER TO keril;

CREATE TABLE IF NOT EXISTS diplom_raw.road(
    id SERIAL PRIMARY KEY,
    r_name VARCHAR(3),
    load_data TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() at time zone 'utc+3')
    );
ALTER TABLE diplom_raw.road OWNER TO keril;

CREATE TABLE IF NOT EXISTS diplom_raw.station(
    id SERIAL PRIMARY KEY,
    s_name VARCHAR(255),
    load_data TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() at time zone 'utc+3'));
ALTER TABLE diplom_raw.station OWNER TO keril;

CREATE TABLE IF NOT EXISTS diplom_raw.gruz(
    id SERIAL PRIMARY KEY,
    g_name VARCHAR(255),
    load_data TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() at time zone 'utc+3')
);
ALTER TABLE diplom_raw.gruz OWNER TO keril;

CREATE TABLE IF NOT EXISTS diplom_raw.company(
    id SERIAL PRIMARY KEY,
    c_name VARCHAR(255),
    load_data TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() at time zone 'utc+3')
);
ALTER TABLE diplom_raw.company OWNER TO keril;

CREATE TABLE IF NOT EXISTS diplom_raw.wagon_group (
    id SERIAL PRIMARY KEY,
    group_name varchar(255),
    load_data TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() at time zone 'utc+3')
    );
ALTER TABLE diplom_raw.wagon_group OWNER TO keril;

CREATE TABLE IF NOT EXISTS diplom_raw.wagon (
    id SERIAL PRIMARY KEY,
    "Номер вагона" varchar(255),
    group_id INT REFERENCES diplom_raw.wagon_group(id),
    load_data TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() at time zone 'utc+3')
    );
ALTER TABLE diplom_raw.wagon OWNER TO keril;

CREATE TABLE IF NOT EXISTS diplom_raw.naming (
    id SERIAL PRIMARY KEY,
    "Оригинальное наименование" varchar(255),
    "Альтернативные имена" TEXT,
    "Время изменения данных" TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() at time zone 'utc+3')
    );
ALTER TABLE diplom_raw.naming OWNER TO keril;