# ==========================================
# 1. USUARIOS
# ==========================================

# 1.1 Crear Usuario (Éxito 201)
POST http://localhost:8080/usuarios
Content-Type: application/json
{
    "nombre": "Usuario Hurl",
    "email": "hurl_test@example.com"
}
HTTP 201
[Captures]
usuario_id: jsonpath "$.id"

# 1.2 Listar Usuarios (Éxito 200)
GET http://localhost:8080/usuarios
HTTP 200

# 1.3 Obtener Usuario por ID (Éxito 200)
GET http://localhost:8080/usuarios/{{usuario_id}}
HTTP 200

# 1.4 Actualizar Usuario (Éxito 200)
PUT http://localhost:8080/usuarios/{{usuario_id}}
Content-Type: application/json
{
    "nombre": "Usuario Hurl Editado",
    "email": "hurl_editado@example.com"
}
HTTP 200

# 1.5 Obtener Usuario Inexistente (Error 404)
GET http://localhost:8080/usuarios/99999
HTTP 404

# ==========================================
# 2. PRENDAS
# ==========================================

# 2.1 Crear Prenda (Éxito 201)
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

# 2.2 Listar Prendas por Usuario (Éxito 200 - Requiere usuario_id en Query Param)
GET http://localhost:8080/prendas?usuario_id={{usuario_id}}
HTTP 200

# 2.3 Listar Prendas Sin usuario_id (Error 400 Bad Request)
GET http://localhost:8080/prendas
HTTP 400

# 2.4 Obtener Prenda por ID (Éxito 200)
GET http://localhost:8080/prendas/{{prenda_id}}
HTTP 200

# 2.5 Actualizar Prenda (Éxito 200)
PUT http://localhost:8080/prendas/{{prenda_id}}
Content-Type: application/json
{
    "usuario_id": {{usuario_id}},
    "nombre": "Remera Hurl Editada",
    "categoria": "Remeras"
}
HTTP 200

# ==========================================
# 3. LOOKS
# ==========================================

# 3.1 Crear Look (Éxito 201)
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

# 3.2 Listar Looks (Éxito 200)
GET http://localhost:8080/looks
HTTP 200

# 3.3 Obtener Look por ID (Éxito 200)
GET http://localhost:8080/looks/{{look_id}}
HTTP 200

# 3.4 Actualizar Look (Éxito 200)
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

# ==========================================
# 4. ME GUSTA LOOKS
# ==========================================

# 4.1 Dar Me Gusta (Éxito 201)
POST http://localhost:8080/me-gusta
Content-Type: application/json
{
    "usuario_id": {{usuario_id}},
    "look_id": {{look_id}}
}
HTTP 201

# 4.2 Quitar Me Gusta (Éxito 204)
DELETE http://localhost:8080/me-gusta
Content-Type: application/json
{
    "usuario_id": {{usuario_id}},
    "look_id": {{look_id}}
}
HTTP 204

# ==========================================
# 5. ELIMINACIONES Y VERIFICACIÓN FINAL
# ==========================================

# 5.1 Eliminar Look (Éxito 204)
DELETE http://localhost:8080/looks/{{look_id}}
HTTP 204

# 5.2 Eliminar Prenda (Éxito 204)
DELETE http://localhost:8080/prendas/{{prenda_id}}
HTTP 204

# 5.3 Eliminar Usuario (Éxito 204)
DELETE http://localhost:8080/usuarios/{{usuario_id}}
HTTP 204

# 5.4 Verificar que el Look eliminado ya no existe (Error 404)
GET http://localhost:8080/looks/{{look_id}}
HTTP 404
