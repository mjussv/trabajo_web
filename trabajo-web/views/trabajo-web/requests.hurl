# ==========================================
# 1. VERIFICAR VISTA WEB INICIAL (SIN SESIÓN)
# ==========================================
GET http://localhost:8080/
HTTP 200
[Asserts]
header "Content-Type" contains "text/html"
body contains "Iniciar Sesión"
body contains "Registrar Nuevo Usuario"

# ==========================================
# 2. REGISTRO Y LOGIN WEB CON COOKIE DE SESIÓN
# ==========================================
POST http://localhost:8080/usuarios/crear-web
[Form]
nombre: JustinaWeb
email: justina_web@example.com
password: clave123
HTTP 303
[Asserts]
header "Location" == "/"
cookie "session_user_id" exists

# Acceder a / usando la cookie guardada automáticamente por Hurl
GET http://localhost:8080/
HTTP 200
[Asserts]
header "Content-Type" contains "text/html"
body contains "JustinaWeb"
body contains "Cargar Nueva Prenda"

# ==========================================
# 3. API REST - USUARIOS
# ==========================================
POST http://localhost:8080/usuarios
Content-Type: application/json
{
    "nombre": "Usuario Hurl",
    "email": "hurl_test@example.com",
    "password": "secreto123"
}
HTTP 201
[Captures]
usuario_id: jsonpath "$.id"
[Asserts]
header "Content-Type" contains "application/json"
jsonpath "$.id" exists
jsonpath "$.nombre" == "Usuario Hurl"
jsonpath "$.email" == "hurl_test@example.com"

GET http://localhost:8080/usuarios
HTTP 200
[Asserts]
header "Content-Type" contains "application/json"
jsonpath "$" isCollection

GET http://localhost:8080/usuarios/{{usuario_id}}
HTTP 200
[Asserts]
jsonpath "$.id" == {{usuario_id}}
jsonpath "$.nombre" == "Usuario Hurl"

GET http://localhost:8080/usuarios/999999
HTTP 404

# ==========================================
# 4. API REST - PRENDAS
# ==========================================
POST http://localhost:8080/prendas?usuario_id={{usuario_id}}
Content-Type: application/json
{
    "usuario_id": {{usuario_id}},
    "nombre": "Remera Hurl",
    "categoria": "Remeras"
}
HTTP 201
[Captures]
prenda_id: jsonpath "$.id"
[Asserts]
jsonpath "$.id" exists
jsonpath "$.nombre" == "Remera Hurl"
jsonpath "$.categoria" == "Remeras"
jsonpath "$.usuario_id" == {{usuario_id}}

GET http://localhost:8080/prendas?usuario_id={{usuario_id}}
HTTP 200
[Asserts]
jsonpath "$" isCollection
jsonpath "$[0].nombre" == "Remera Hurl"

PUT http://localhost:8080/prendas/{{prenda_id}}?usuario_id={{usuario_id}}
Content-Type: application/json
{
    "usuario_id": {{usuario_id}},
    "nombre": "Remera Hurl Editada",
    "categoria": "Remeras"
}
HTTP 200
[Asserts]
jsonpath "$.id" == {{prenda_id}}
jsonpath "$.nombre" == "Remera Hurl Editada"

# ==========================================
# 5. API REST - LOOKS Y LIMPIEZA
# ==========================================
POST http://localhost:8080/looks?usuario_id={{usuario_id}}
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
[Asserts]
jsonpath "$.id" exists
jsonpath "$.nombre" == "Look Hurl"
jsonpath "$.ocasion" == "Casual"
jsonpath "$.temporada" == "Verano"
jsonpath "$.estilo" == "Urbano"

PUT http://localhost:8080/looks/{{look_id}}?usuario_id={{usuario_id}}
Content-Type: application/json
{
    "nombre": "Look Hurl Editado",
    "descripcion": "Nueva descripcion",
    "ocasion": "Noche",
    "temporada": "Invierno",
    "prenda_estrella_id": {{prenda_id}},
    "url_imagen": "https://example.com/img2.jpg",
    "estilo": "Elegante"
}
HTTP 200
[Asserts]
jsonpath "$.id" == {{look_id}}
jsonpath "$.nombre" == "Look Hurl Editado"
jsonpath "$.ocasion" == "Noche"
jsonpath "$.temporada" == "Invierno"
jsonpath "$.estilo" == "Elegante"

DELETE http://localhost:8080/looks/{{look_id}}?usuario_id={{usuario_id}}
HTTP 204

DELETE http://localhost:8080/prendas/{{prenda_id}}?usuario_id={{usuario_id}}
HTTP 204
