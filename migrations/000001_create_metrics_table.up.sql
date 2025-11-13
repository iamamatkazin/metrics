-- Создание таблицы метрик
CREATE TABLE metrics (
    id VARCHAR(255) PRIMARY KEY,
    mtype VARCHAR(255) NOT NULL,
    val DOUBLE PRECISION,
    delta INTEGER
);