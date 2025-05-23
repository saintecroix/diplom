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

CREATE TABLE IF NOT EXISTS diplom_raw.wagon_operations (
    id SERIAL PRIMARY KEY,
    "Группа наименование" VARCHAR(255),
    "Вагон" VARCHAR(20),
    "Дата операции" TIMESTAMP WITHOUT TIME ZONE,
    "Операция" VARCHAR(255),
    "Ст. опер." VARCHAR(100),
    "Дрг. опер." VARCHAR(100),
    "Дата отпраления" TIMESTAMP WITHOUT TIME ZONE,
    "Станция отпр." VARCHAR(255),
    "Дрг. отпр." VARCHAR(100),
    "Ст. назнач." VARCHAR(255),
    "Дрг. назн." VARCHAR(100),
    "Код груза" VARCHAR(5),
    "Груз наименование" VARCHAR(255),
    "Дата доставки (расчетная)" DATE,
    "Расст., км" REAL,
    "Простой (дни)" NUMERIC,
    "груж." VARCHAR(10),
    "Время загрузки данных" TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() at time zone 'utc+3')
);
ALTER TABLE diplom_raw.wagon_operations OWNER TO keril;

CREATE TABLE IF NOT EXISTS diplom_raw.wg_group (
    id SERIAL PRIMARY KEY,
    "Группа наименование" varchar(255),
    "Вагоны" TEXT,
    "Время изменнеия данных" TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() at time zone 'utc+3')
    );
ALTER TABLE diplom_raw.wg_group OWNER TO keril;

CREATE TABLE IF NOT EXISTS diplom_raw.naming (
    id SERIAL PRIMARY KEY,
    "Оригинальное наименование" varchar(255),
    "Альтернативные имена" TEXT,
    "Время изменения данных" TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() at time zone 'utc+3')
    );
ALTER TABLE diplom_raw.naming OWNER TO keril;