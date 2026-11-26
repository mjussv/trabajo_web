# Entrega tp6
**Requisitos Previos:**

  Para ejecutar la aplicación y correr las suites de pruebas de forma automatizada es necesario contar con las siguientes herramientas instaladas:

  - Go (v1.25.0+)
  - sqlc (v1.31.1+)
  - templ (v0.3.1100+)
  - Docker (v29.1.3+)
  - Docker Compose (v2.27.0+)
  - GNU Make (v4.3+)
  - Hurl (v4.3.0+)

  *(Nota: el motor de base de datos **PostgreSQL 16** y su esquema inicial se levantan y configuran automáticamente mediante Docker Compose, por lo que no requiere instalación local. Asimismo, **HTMX v2.0.0** se importa directamente desde CDN en el layout base).*

  Si es la primera vez que se instalan `sqlc` o `templ` con `go install`, asegurarse de tener el directorio `bin` de Go en el `PATH`:

    - `echo 'export PATH=$PATH:~/go/bin' >> ~/.bashrc`

---

**Ejecución de la Aplicación:**

  Todo el ciclo de vida de ejecución está automatizado en el archivo `Makefile` dentro de la carpeta `trabajo-web`. Al iniciar, el script limpia contenedores y volúmenes previos, genera el código actualizado de `sqlc` y `templ`, levanta PostgreSQL en Docker, espera a que la base esté lista para aceptar conexiones y arranca el servidor web. Al detener la ejecución (`Ctrl + C`) o ante cualquier fallo, se garantiza la detención y limpieza automática de contenedores y volúmenes.

**Comandos para levantar la aplicación (desde la terminal y dentro de la carpeta `trabajo-web`):**

    - `cd trabajo-web`
    - `make run` 

**URL de acceso en el navegador:**

    - `http://localhost:8080`

---

**Ejecución de Pruebas:**

  A diferencia de la ejecución interactiva de la aplicación, las pruebas automatizadas levantan un entorno limpio, ejecutan las validaciones correspondientes y realizan la limpieza final de contenedores y volúmenes al terminar. Todos los comandos deben ejecutarse desde la carpeta `trabajo-web`:

    - `make test-hurl`: levanta la base de datos y la aplicación completa, espera a que el servidor responda en `http://localhost:8080` y ejecuta la suite de pruebas de integración HTTP y fragmentos HTMX definida en `requests.hurl` (verificando códigos de estado esperados, cabeceras, cookies de sesión, casos de error 400/401/403/404 y contenido HTML/OOB devuelto).
    - `make test`: ejecuta todos los tests unitarios y de integración en Go (`go test -v ./...`), incluyendo la suite completa sobre la capa de persistencia y consultas `sqlc` (`TestSuiteCompleta`).
    - `make test-db`: ejecuta únicamente los tests correspondientes a la capa de base de datos (`go test -v ./db/sqlc`).

---

**Aclaraciones para el docente:**

  - **Cobertura completa del dominio en la interfaz web:** A diferencia de la cuarta entrega (donde solo se habían trabajado usuarios y prendas), en la versión final se encuentran implementadas de punta a punta las cuatro entidades del modelo (`usuarios`, `prendas`, `looks` y `me_gusta_looks`) con sus respectivas vistas Templ, lógica de negocio y consultas `sqlc`.
