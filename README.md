# Cuarta entrega:

**Requisitos Previos:**

  Para ejecutar el proyecto y correr el archivo de pruebas es necesario contar con las siguientes herramientas: Go (v1.22.2+), sqlc (v1.31.1+), templ (v0.2.747+), Docker (v29.1.3+), Docker Compose (v2.27.0+), GNU Make (v4.3+) y Hurl (v4.3.0+) 
  
**Ejecución de Pruebas:**

  La ejecución de los tests está completamente automatizada en el archivo Makefile.

**Comando de ejecución (desde terminal y carpeta raíz del proyecto "trabajo-web"):**

    - " cd trabajo-web "
    - " make run " 

**A tener en cuenta:**

  La prueba automatizada es la previamente dicha, que se ejecuta con el " make run ", pero el Makefile ofrece otros test (con el objetivo de corroborar la consistencia del proyecto), que se deben ejecutar de la misma forma que se aclara en el apartado "Comando de ejecución" (es decir, desde la terminal, y dentro de la carpeta raiz del proyecto "trabajo-web"):
  
    - " make test-db ": se encarga de realizar los tests correspondientes a la 2da entrega (aquellos que son sobre la base de datos)
    - " make test ": ejecuta todos los tests que encuentre en el proyecto  

**Aclaración para la corrección:**

  En esta entrega, por cuestiones de tiempo y complejidad, implementamos con la lógica y vistas de Templ las entidades principales: "usuarios" y "prendas", como se nos aconsejó en la clase consulta del 6/11; nos faltaría hacer lo mismo con "looks" y "me_gusta_looks", por lo que vamos a seguir trabajando en eso en un branch aparte, con el fin de tenerlo listo y completo para la próxima entrega.
