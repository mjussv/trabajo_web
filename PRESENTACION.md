# OOTD (Outfit of the Day) 🛍️
**Presentación del Proyecto — Grupo 24**
*Programación Web (2026)*

**Integrantes:**
- Graziano, Federica
- Irrazabal, Ainara
- Venturini, Maria Justina

---

## ¿Qué es OOTD?

**OOTD (Outfit of the Day)** es una aplicación web diseñada como un guardarropa digital y comunidad de (?)(inspiración de moda). Su objetivo es resolver un problema cotidiano: organizar las prendas que cada persona tiene en su clóset, combinarlas en *looks* completos según la ocasión, la temporada o el estilo, y descubrir ideas de vestimenta compartidas por otros usuarios de manera simple, rápida y visual.

---

## Funcionalidades Principales

1. **Gestión de Cuentas y Sesión Segura:**
   - Registro e inicio de sesión con validación de datos, almacenamiento seguro de contraseñas mediante **bcrypt** y persistencia de sesión a través de **cookies firmadas (HMAC-SHA256)**.
   - Validación de permisos en el servidor para garantizar que cada usuario solo pueda modificar o eliminar sus propias prendas y looks, y que no pueda darse *Me gusta* a sus propias publicaciones.

2. **Guardarropa Digital Personal (Prendas y Looks):**
   - **Carga y gestión de prendas:** Permite registrar prendas individuales con nombre, categoría (abrigos, calzado, tops, etc.) e imagen, con respaldo visual automático en caso de enlaces externos caídos.
   - **Creación de looks completos:** Permite armar combinaciones con título, descripción, ocasión, temporada, estilo e imagen, vinculando opcionalmente una **prenda estrella** del propio guardarropa.
   - **Filtros y estadísticas en tiempo real:** Cada perfil muestra los contadores actualizados de prendas y looks, permitiendo filtrar la vista entre *Todos*, *Prendas* o *Looks*.

3. **Comunidad, Búsqueda e Interacción Social:**
   - **Explorar Red (`/explorar`):** Buscador global que permite filtrar publicaciones de toda la comunidad por palabra clave (nombre, estilo, ocasión, temporada, prenda o autor) y por tipo de contenido.
   - **Perfiles Públicos (`/perfil/{nombre}`):** Vista de solo lectura para recorrer los looks y prendas publicadas por otros integrantes de la comunidad.
   - **Sistema de "Me gusta":** Interacción instantánea sobre los looks de otros usuarios con actualización inmediata del contador y estado del botón.

---

## Arquitectura y Experiencia de Usuario (SPA sin JavaScript propio)

OOTD fue construida siguiendo una arquitectura en capas limpia en **Go** (`handlers`, `logic`, `db/sqlc` y `views`), respaldada por **PostgreSQL 16** y consultas SQL tipadas generadas con **sqlc**.

---

## Stack Tecnológico

- **Backend:** Go (biblioteca estándar `net/http`, arquitectura modular en capas)
- **Motor de Plantillas:** Templ (renderizado HTML tipado en servidor)
- **Interactividad Frontend:** HTMX 2.0 + CSS3 personalizado
- **Base de Datos y Persistencia:** PostgreSQL 16 + `sqlc` + driver `pgx/v5`
- **Infraestructura y Pruebas Automatizadas:** Docker, Docker Compose, GNU Make, `go test` y Hurl
