CREATE TABLE IF NOT EXISTS usuarios (
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(100) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    fecha_creacion TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS prendas (
    id SERIAL PRIMARY KEY,
    usuario_id INT NOT NULL REFERENCES usuarios(id) ON DELETE CASCADE,
    nombre VARCHAR(100) NOT NULL,
    categoria VARCHAR(50) NOT NULL,
    url_imagen TEXT NOT NULL DEFAULT '',
    fecha_creacion TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS looks (
    id SERIAL PRIMARY KEY,
    usuario_id INT NOT NULL REFERENCES usuarios(id) ON DELETE CASCADE,
    nombre VARCHAR(255) NOT NULL,
    descripcion VARCHAR(500),
    ocasion VARCHAR(100) NOT NULL,
    temporada VARCHAR(100) NOT NULL,
    prenda_estrella_id INT REFERENCES prendas(id) ON DELETE SET NULL,
    url_imagen TEXT NOT NULL,
    estilo VARCHAR(100),
    fecha_creacion TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS me_gusta_looks (
    usuario_id INT NOT NULL REFERENCES usuarios(id) ON DELETE CASCADE,
    look_id INT NOT NULL REFERENCES looks(id) ON DELETE CASCADE,
    fecha_creacion TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (usuario_id, look_id)
);

CREATE INDEX IF NOT EXISTS idx_me_gusta_look ON me_gusta_looks(look_id);
CREATE INDEX IF NOT EXISTS idx_looks_usuario ON looks(usuario_id);
CREATE INDEX IF NOT EXISTS idx_prendas_usuario ON prendas(usuario_id);
CREATE INDEX IF NOT EXISTS idx_looks_prenda_estrella ON looks(prenda_estrella_id);
