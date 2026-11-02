# ===== 1. CREAR USUARIO =====
POST http://localhost:8080/usuarios
Content-Type: application/json
{
    "nombre": "Usuario Hurl",
    "email": "hurl_test@example.com"
}
HTTP 201
[Captures]
usuario_id: jsonpath "$.id"

# ===== 2. OBTENER USUARIOS =====
GET http://localhost:8080/usuarios
HTTP 200

# ===== 3. CREAR PRENDA =====
POST http://localhost:8080/prendas
Content-Type: application/json
{
    "usuario_id": {{usuario_id}},
    "nombre": "Remera Hurl",
    "categoria": "Remeras"
}
HTTP 201
[Captures]
prenda_id: jsonpath "$.id"

# ===== 4. ACTUALIZAR PRENDA (PUT) =====
PUT http://localhost:8080/prendas/{{prenda_id}}
Content-Type: application/json
{
    "usuario_id": {{usuario_id}},
    "nombre": "Remera Hurl Editada",
    "categoria": "Remeras"
}
HTTP 200

# ===== 5. CREAR LOOK =====
POST http://localhost:8080/looks
Content-Type: application/json
{
    "usuario_id": {{usuario_id}},
    "nombre": "Look Hurl",
    "descripcion": "Descripción test",
    "ocasion": "Casual",
    "temporada": "Verano",
    "prenda_estrella_id": {{prenda_id}},
    "url_imagen": "https://example.com/img.jpg",
    "estilo": "Urbano"
}
HTTP 201
[Captures]
look_id: jsonpath "$.id"

# ===== 6. ACTUALIZAR LOOK (PUT) =====
PUT http://localhost:8080/looks/{{look_id}}
Content-Type: application/json
{
    "usuario_id": {{usuario_id}},
    "nombre": "Look Hurl Editado",
    "descripcion": "Nueva descripcion",
    "ocasion": "Noche",
    "temporada": "Invierno",
    "prenda_estrella_id": {{prenda_id}},
    "url_imagen": "https://example.com/img2.jpg",
    "estilo": "Elegante"
}
HTTP 200

# ===== 7. ELIMINAR LOOK =====
DELETE http://localhost:8080/looks/{{look_id}}
HTTP 204

# ===== 8. ELIMINAR PRENDA =====
DELETE http://localhost:8080/prendas/{{prenda_id}}
HTTP 204
