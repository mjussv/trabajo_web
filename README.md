# Tercera Entrega:

**Requisitos Previos:**

  Para ejecutar el proyecto y correr el archivo de pruebas es necesario contar con las herramientas de la segunda entrega.
  
**Ejecución de Pruebas:**

  La ejecución de los tests está completamente automatizada en el archivo Makefile.

**Comando de ejecución (desde terminal y carpeta raíz):**

  " make test-api "

**A tener en cuenta:**

  La prueba automatizada es la previamente dicha, que se ejecuta con el " make test-api "; si se quisieran verificar las entregas previas: 
  - La primera entrega requiere que se levante primero el docker y luego ejecutar el main.go, ya que el Makefile (al hacer la prueba " make test-api ") se encarga de apagar el primero, los comandos que se deberian ejecutar son los siguientes:
    - " docker compose up -d postgresql "
    - (opcional para verificar que PostgreSQL esté listo y poder ejecutar el main.go: " docker exec postgres_container pg_isready -U postgres -d outfits_db ", sino esperar unos segundos)
    - " go run main.go "
  - La segunda entrega solo requiere que se ejecute " make test-db "
  - Si se quisieran probar todos los tests que hay disponibles en todo el proyecto se debe ejecutar " make test "
