-- ===== USUARIOS =====

-- name: CrearUsuario :one
INSERT INTO usuarios (nombre, email, password)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUsuariosByID :one
SELECT * FROM usuarios
WHERE id = $1 LIMIT 1;

-- name: ListarUsuarios :many
SELECT * FROM usuarios ORDER BY nombre;

-- name: ActualizarUsuario :one
UPDATE usuarios
SET nombre = $2,
    email = $3
WHERE id = $1
RETURNING *;

-- name: EliminarUsuario :exec
DELETE FROM usuarios WHERE id = $1;

-- name: GetUsuariosByEmail :one
SELECT * FROM usuarios
WHERE email = $1 LIMIT 1;

-- name: GetUsuariosByNombre :one
SELECT * FROM usuarios
WHERE nombre = $1 LIMIT 1;


-- ===== PRENDAS =====

-- name: CrearPrenda :one
INSERT INTO prendas (usuario_id, nombre, categoria, url_imagen)
VALUES ($1, $2, $3, COALESCE(NULLIF($4::text, ''), 'https://placehold.co/400x500/F8E8A2/3C132C?text=Prenda'))
RETURNING *;

-- name: ObtenerPrenda :one
SELECT * FROM prendas WHERE id = $1;

-- name: ListarPrendasPorUsuario :many
SELECT * FROM prendas WHERE usuario_id = $1 ORDER BY fecha_creacion DESC;

-- name: ListarTodasPrendasConAutora :many
SELECT
    p.id, p.usuario_id, p.nombre, p.categoria, p.url_imagen, p.fecha_creacion,
    u.nombre AS autora
FROM prendas p
JOIN usuarios u ON p.usuario_id = u.id
ORDER BY p.fecha_creacion DESC;

-- name: BuscarPrendas :many
SELECT
    p.id, p.usuario_id, p.nombre, p.categoria, p.url_imagen, p.fecha_creacion,
    u.nombre AS autora
FROM prendas p
JOIN usuarios u ON p.usuario_id = u.id
WHERE p.nombre ILIKE '%' || $1 || '%'
   OR p.categoria ILIKE '%' || $1 || '%'
   OR u.nombre ILIKE '%' || $1 || '%'
ORDER BY p.fecha_creacion DESC;

-- name: ActualizarPrenda :one
UPDATE prendas
SET nombre = $3,
    categoria = $4,
    url_imagen = COALESCE(NULLIF($5::text, ''), url_imagen)
WHERE id = $1 AND usuario_id = $2
RETURNING *;

-- name: EliminarPrenda :exec
DELETE FROM prendas WHERE id = $1;

-- name: LimpiarClosetPorUsuario :exec
DELETE FROM prendas WHERE usuario_id = $1;


-- ===== LOOKS (CRUD principal) =====

-- name: CrearLook :one
INSERT INTO looks (
    usuario_id, nombre, descripcion, ocasion, temporada, prenda_estrella_id, url_imagen, estilo
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: ObtenerLookConDetalle :one
SELECT
    l.id, l.usuario_id, l.nombre, l.descripcion, l.ocasion, l.temporada, l.url_imagen, l.estilo, l.fecha_creacion,
    u.nombre AS autora,
    p.nombre AS prenda_estrella
FROM looks l
JOIN usuarios u ON l.usuario_id = u.id
LEFT JOIN prendas p ON l.prenda_estrella_id = p.id
WHERE l.id = $1;

-- name: ListarLooksConDetalle :many
SELECT
    l.id, l.usuario_id, l.nombre, l.descripcion, l.ocasion, l.temporada, l.url_imagen, l.estilo, l.fecha_creacion,
    u.nombre AS autora,
    p.nombre AS prenda_estrella,
    (SELECT COUNT(*) FROM me_gusta_looks mg WHERE mg.look_id = l.id) AS total_me_gusta,
    EXISTS(SELECT 1 FROM me_gusta_looks mg WHERE mg.look_id = l.id AND mg.usuario_id = $1) AS le_dio_me_gusta
FROM looks l
JOIN usuarios u ON l.usuario_id = u.id
LEFT JOIN prendas p ON l.prenda_estrella_id = p.id
ORDER BY l.fecha_creacion DESC;

-- name: BuscarLooks :many
SELECT
    l.id, l.usuario_id, l.nombre, l.descripcion, l.ocasion, l.temporada, l.url_imagen, l.estilo, l.fecha_creacion,
    u.nombre AS autora,
    p.nombre AS prenda_estrella,
    (SELECT COUNT(*) FROM me_gusta_looks mg WHERE mg.look_id = l.id) AS total_me_gusta,
    EXISTS(SELECT 1 FROM me_gusta_looks mg WHERE mg.look_id = l.id AND mg.usuario_id = $1) AS le_dio_me_gusta
FROM looks l
JOIN usuarios u ON l.usuario_id = u.id
LEFT JOIN prendas p ON l.prenda_estrella_id = p.id
WHERE l.nombre ILIKE '%' || $2 || '%'
   OR COALESCE(l.descripcion, '') ILIKE '%' || $2 || '%'
   OR l.ocasion ILIKE '%' || $2 || '%'
   OR l.temporada ILIKE '%' || $2 || '%'
   OR COALESCE(l.estilo, '') ILIKE '%' || $2 || '%'
   OR COALESCE(p.nombre, '') ILIKE '%' || $2 || '%'
   OR u.nombre ILIKE '%' || $2 || '%'
ORDER BY l.fecha_creacion DESC;

-- name: ListarLooksPorUsuario :many
SELECT
    l.id, l.usuario_id, l.nombre, l.descripcion, l.ocasion, l.temporada, l.url_imagen, l.estilo, l.fecha_creacion,
    u.nombre AS autora,
    p.nombre AS prenda_estrella,
    (SELECT COUNT(*) FROM me_gusta_looks mg WHERE mg.look_id = l.id) AS total_me_gusta,
    EXISTS(SELECT 1 FROM me_gusta_looks mg WHERE mg.look_id = l.id AND mg.usuario_id = $2) AS le_dio_me_gusta
FROM looks l
JOIN usuarios u ON l.usuario_id = u.id
LEFT JOIN prendas p ON l.prenda_estrella_id = p.id
WHERE l.usuario_id = $1
ORDER BY l.fecha_creacion DESC;

-- name: ListarLooksPorOcasion :many
SELECT * FROM looks WHERE ocasion = $1 ORDER BY fecha_creacion DESC;

-- name: ListarLooksPorPrenda :many
SELECT l.*
FROM looks l
JOIN prendas p ON l.prenda_estrella_id = p.id
WHERE p.nombre ILIKE '%' || $1 || '%'
ORDER BY l.fecha_creacion DESC;

-- name: ActualizarLook :one
UPDATE looks
SET nombre = $3,
    descripcion = $4,
    ocasion = $5,
    temporada = $6,
    prenda_estrella_id = $7,
    url_imagen = $8,
    estilo = $9
WHERE id = $1 AND usuario_id = $2
RETURNING *;

-- name: EliminarLook :exec
DELETE FROM looks WHERE id = $1;


-- ===== ME GUSTA (like tipo Instagram) =====

-- name: DarMeGusta :exec
INSERT INTO me_gusta_looks (usuario_id, look_id) VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: QuitarMeGusta :execrows
DELETE FROM me_gusta_looks WHERE usuario_id = $1 AND look_id = $2;

-- name: ContarMeGustaPorLook :one
SELECT COUNT(*) FROM me_gusta_looks WHERE look_id = $1;

-- name: UsuarioLeDioMeGusta :one
SELECT EXISTS (
    SELECT 1 FROM me_gusta_looks WHERE usuario_id = $1 AND look_id = $2
);
