PRAGMA foreign_keys = ON;

-- =========================
-- Tablas base
-- =========================

CREATE TABLE IF NOT EXISTS solicitante (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nombre TEXT NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_solicitante_nombre 
ON solicitante(nombre);

CREATE TABLE IF NOT EXISTS colaborador (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nombre TEXT NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_colaborador_nombre 
ON colaborador(nombre);

-- =========================

CREATE TABLE IF NOT EXISTS codigo_SAP (
    codigo TEXT PRIMARY KEY,
    descripcion TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS codigo_ID (
    codigo TEXT PRIMARY KEY,
    descripcion TEXT,
    estado TEXT NOT NULL,
    fecha_pedido DATETIME NOT NULL,
    fecha_actualizacion DATETIME
);

CREATE TABLE IF NOT EXISTS documento (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    codigo TEXT NOT NULL,
    titulo TEXT NOT NULL,
    tipo TEXT NOT NULL,
    ubicacion_path TEXT
);

-- =========================
-- Tabla principal
-- =========================

CREATE TABLE IF NOT EXISTS pendientes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    titulo TEXT NOT NULL,
    descripcion TEXT,

    solicitante_id INTEGER NOT NULL,
    fecha_pedido DATETIME DEFAULT CURRENT_TIMESTAMP,

    asignado_id INTEGER,
    fecha_asignado DATETIME,

    cierre TEXT,
    fecha_cierre DATETIME,

    identificacion_tabla_pendiente TEXT,

    FOREIGN KEY (solicitante_id) REFERENCES solicitante(id),
    FOREIGN KEY (asignado_id) REFERENCES colaborador(id)
);

-- =========================
-- Avances
-- =========================

CREATE TABLE IF NOT EXISTS avance (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    id_pendiente INTEGER NOT NULL,
    descripcion TEXT NOT NULL,
    fecha DATETIME DEFAULT CURRENT_TIMESTAMP,
    mail_path TEXT,

    FOREIGN KEY (id_pendiente) REFERENCES pendientes(id) ON DELETE CASCADE
);

-- =========================
-- Adjunto
-- =========================

CREATE TABLE IF NOT EXISTS adjunto (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    id_pendiente INTEGER NOT NULL,
    descripcion TEXT,
    mail_path TEXT,

    FOREIGN KEY (id_pendiente) REFERENCES pendientes(id) ON DELETE CASCADE
);

-- =========================
-- Tablas intermedias (N:N)
-- =========================

CREATE TABLE IF NOT EXISTS pendientes_documento (
    pendiente_id INTEGER NOT NULL,
    documento_id INTEGER NOT NULL,

    PRIMARY KEY (pendiente_id, documento_id),

    FOREIGN KEY (pendiente_id) REFERENCES pendientes(id) ON DELETE CASCADE,
    FOREIGN KEY (documento_id) REFERENCES documento(id)
);

CREATE TABLE IF NOT EXISTS pendientes_codigo_sap (
    pendiente_id INTEGER NOT NULL,
    codigo_sap_codigo TEXT NOT NULL,

    PRIMARY KEY (pendiente_id, codigo_sap_codigo),

    FOREIGN KEY (pendiente_id) REFERENCES pendientes(id) ON DELETE CASCADE,
    FOREIGN KEY (codigo_sap_codigo) REFERENCES codigo_SAP(codigo)
);

CREATE TABLE IF NOT EXISTS pendientes_codigo_id (
    pendiente_id INTEGER NOT NULL,
    codigo_id_codigo TEXT NOT NULL,

    PRIMARY KEY (pendiente_id, codigo_id_codigo),

    FOREIGN KEY (pendiente_id) REFERENCES pendientes(id) ON DELETE CASCADE,
    FOREIGN KEY (codigo_id_codigo) REFERENCES codigo_ID(codigo)
);