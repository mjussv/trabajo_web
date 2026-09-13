# Segunda Entrega:

**Requisitos Previos:**

  Para ejecutar el proyecto y correr el archivo de pruebas es necesario contar con las siguientes herramientas instaladas:
    Docker (v29.1.3+), Docker Compose (v2.27.0+), Go (v1.22.2+), sqlc (v1.31.1+) y GNU Make (v4.3+)
  
**Ejecución de Pruebas:**

  La ejecución de los tests está completamente automatizada en el archivo Makefile.

**Comando de ejecución (desde terminal):**

  " make test "

**A tener en cuenta:**

  Este proyecto implementa la capa de persistencia para una aplicación utilizando Go, PostgreSQL 16 y generación de código SQL seguro mediante sqlc.
  
  Si es la primera vez que se descarga sqlc el mismo puede bajarse en el bin de la carpeta "go", lo que podria generar errores; para solucionar esto, utilizar el siguiente comando: " echo 'export PATH=$PATH:~/go/bin'>> ~/.bashrc "
