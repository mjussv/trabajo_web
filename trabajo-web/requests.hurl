# 1. Vista inicial sin sesión y manejo de filtro inválido (Ruta: GET /)
GET http://localhost:8080/?filtro=invalido
HTTP 200
[Asserts]
header "Content-Type" contains "text/html"
body contains "Iniciar Sesión"
body contains "Registrar Nuevo Usuario"

# 2. Registro de Usuario 1 y sesión con cookie firmada (Ruta: POST /usuarios/crear-web)
POST http://localhost:8080/usuarios/crear-web
[FormParams]
nombre: prueba
email: prueba@example.com
password: clave123
HTTP 303
[Asserts]
header "Location" == "/"
cookie "session_user_id" exists

GET http://localhost:8080/
HTTP 200
[Asserts]
body contains "prueba"
body contains "Cargar Nueva Prenda"

# 3. Error de validación al crear prenda vacía (HTTP 400)
POST http://localhost:8080/prendas
[FormParams]
nombre:
categoria: Abrigos
HTTP 400
[Asserts]
body contains "Error al guardar la prenda"

# 4. Crear Prenda vía HTMX (Ruta: POST /prendas)
POST http://localhost:8080/prendas
[FormParams]
nombre: Camperon de Cuero
categoria: Abrigos
url_imagen: https://example.com/campera.jpg
HTTP 200
[Asserts]
header "Content-Type" contains "text/html"
body contains "Camperon de Cuero"
body contains "hx-swap-oob=\"true\""

# 5. Crear Look vía HTMX asociando la prenda estrella 1 (Ruta: POST /looks)
POST http://localhost:8080/looks
[FormParams]
nombre: Urbano Noche
descripcion: Ideal para salir
ocasion: Noche
temporada: Invierno
estilo: Streetwear
prenda_estrella_id: 1
url_imagen: https://example.com/look.jpg
HTTP 200
[Asserts]
header "Content-Type" contains "text/html"
body contains "Urbano Noche"
body contains "Camperon de Cuero"

# 6. Error de validación al dar Me Gusta a un look propio (HTTP 403)
POST http://localhost:8080/me-gusta/toggle-web
[FormParams]
look_id: 1
HTTP 403
[Asserts]
body contains "no podés darle me gusta a tu propio look"

# 7. Cerrar sesión del Usuario 1 (Ruta: POST /logout)
POST http://localhost:8080/logout
HTTP 303
[Asserts]
header "Location" == "/"

# 8. Error al iniciar sesión con contraseña incorrecta (HTTP 401)
POST http://localhost:8080/login
[FormParams]
nombre: prueba
password: clave_incorrecta
HTTP 401
[Asserts]
body contains "Usuario o contraseña incorrectos."

# 9. Error al registrar email duplicado (HTTP 400)
POST http://localhost:8080/usuarios/crear-web
[FormParams]
nombre: otro_usuario
email: prueba@example.com
password: clave123
HTTP 400
[Asserts]
body contains "ese nombre de usuario o email ya está en uso."

# 10. Registrar Usuario 2 para probar interacción entre usuarios
POST http://localhost:8080/usuarios/crear-web
[FormParams]
nombre: amiga
email: amiga@example.com
password: clave123
HTTP 303

# 11. Usuario 2 explora y le da Me Gusta al look de Usuario 1 (Rutas: GET /explorar y POST /me-gusta/toggle-web)
GET http://localhost:8080/explorar?q=Urbano&tipo=invalido
HTTP 200
[Asserts]
body contains "Urbano Noche"

POST http://localhost:8080/me-gusta/toggle-web
[FormParams]
look_id: 1
HTTP 200
[Asserts]
body contains "❤️ 1"

# 12. Usuario 2 quita el Me Gusta (toggle)
POST http://localhost:8080/me-gusta/toggle-web
[FormParams]
look_id: 1
HTTP 200
[Asserts]
body contains "🤍 0"

# 13. Usuario 2 intenta borrar una prenda y un look ajenos (HTTP 404)
DELETE http://localhost:8080/prendas/1
HTTP 404

DELETE http://localhost:8080/looks/1
HTTP 404

# 14. Perfiles públicos y protección de archivos estáticos (Rutas: GET /perfil/{nombre} y GET /static/)
GET http://localhost:8080/perfil/prueba
HTTP 200
[Asserts]
body contains "@prueba"
body contains "Urbano Noche"

GET http://localhost:8080/perfil/amiga
HTTP 303
[Asserts]
header "Location" == "/"

GET http://localhost:8080/perfil/usuario_inexistente
HTTP 404

GET http://localhost:8080/static/
HTTP 404

GET http://localhost:8080/static/css/styles.css
HTTP 200
[Asserts]
header "Content-Type" contains "text/css"

# 15. Cerrar sesión Usuario 2 e ingresar como Usuario 1 (Ruta: POST /login válido)
POST http://localhost:8080/logout
HTTP 303

POST http://localhost:8080/login
[FormParams]
nombre: prueba
password: clave123
HTTP 303

# 16. Borrar la última prenda (Ruta: DELETE /prendas/{id})
DELETE http://localhost:8080/prendas/1
HTTP 200
[Asserts]
body contains "No tenés prendas cargadas en tu clóset aún."
body contains "hx-swap-oob=\"true\""

# 17. Borrar el último look (Ruta: DELETE /looks/{id})
DELETE http://localhost:8080/looks/1
HTTP 200
[Asserts]
body contains "Todavía no armaste ningún look."
body contains "hx-swap-oob=\"true\""
