# 1. Vista inicial sin sesión
GET http://localhost:8080/
HTTP 200
[Asserts]
header "Content-Type" contains "text/html"
body contains "Iniciar Sesión"
body contains "Registrar Nuevo Usuario"

# 2. Registro de Usuario 1 y sesión con cookie firmada
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

# 3. Crear Prenda vía HTMX (devuelve fragmento HTML)
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

# 4. Crear Look vía HTMX
POST http://localhost:8080/looks
[FormParams]
nombre: Urbano Noche
descripcion: Ideal para salir
ocasion: Noche
temporada: Invierno
estilo: Streetwear
url_imagen: https://example.com/look.jpg
HTTP 200
[Asserts]
header "Content-Type" contains "text/html"
body contains "Urbano Noche"

# 5. Explorar red
GET http://localhost:8080/explorar?q=Urbano&tipo=todos
HTTP 200
[Asserts]
body contains "Urbano Noche"

# 6. Cerrar sesión
GET http://localhost:8080/logout
HTTP 303
[Asserts]
header "Location" == "/"
