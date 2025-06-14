## Exercise 3 for Cloud Computing

### The challenge

To succeed in this assignment, you must do the following:

1. Use Docker Compose to orchestrate your **five** containers and deploy **NGINX**
to control the traffic between your services based on the request method. The
containers are as follows:
    - A container to handle each operation for `/api/books`: GET, POST, PUT, and DELETE
    - A container to handle requests to `/*`, i.e., the rendering of the webpage.
    - A container for NGINX.
2. Use a Multi-stage Dockerfile to minimize the size of the image.
3. Publish to a public container registry (e.g., Docker Hub) your images.

## Explicaciones

- Creamos un Dockerfile por servicio (GET, POST, PUT, DELETE, HOMEPAGE)
- docker-compose.yml que emplea estos dockerfiles/imagenes para generar contenedores
- nginx.conf redirige el tráfico de peticiones a cada contenedor según la naturaleza de
  la misma

Lo que estamos haciendo es dividir cada uno de nuestros endpoints (lo cual no es óptimo,
pero lo hacemos para simular microservicios) en un contenedor distinto. En nuestro caso el
endpoint GET de api/books es el que inicializa la base de datos y la rellena, por eso el resto
de contenedores espera a que este esté levantado para conectarse a dicha base de datos por
medio de una variable de entorno llamada DATABASE_URI.

El tráfico de peticiones al servidor se gestiona con NGINX, donde filtramos las peticiones
según si son GET, POST, PUT, DELETE o genéricas al frontend, redirigiéndolas a su respectivo
contenedor asociado. NGINX tiene su propio contenedor.

Cada contenedor expone su puerto 3030 dentro de la red interna de Docker, lo que significa
que no puede ser accedido desde el exterior, de esta forma delegamos en NGINX el acceso.

NGINX, internamente, escucha en el puerto 3030 para recibir conexiones HTTP. Sin embargo, 
desde fuera del contenedor (es decir, desde tu máquina o cualquier cliente externo), no accederás 
a ese puerto 3030 directamente, sino al puerto 8080 de tu host.
 
- **Puerto 3030**: es el puerto donde NGINX escucha dentro del contenedor.
- **Puerto 8080**: es el puerto de tu máquina por donde accedes al servicio.

Esa petición llega primero al puerto 8080 de tu máquina, que está mapeado al puerto 3030 
del contenedor de NGINX. Dentro del contenedor, NGINX recibe la petición en el puerto 3030. 
Según la ruta y el método HTTP de la petición (GET, POST, PUT, DELETE, etc.), NGINX decide a 
qué contenedor de microservicio debe reenviarla.

El archivo **nginx.conf** decide a que puerto del contenedor seleccionado reenviar la petición,
por tanto, si GET expone el puerto 3030 por ejemplo, en este archivo se define:

    upstream get_service {
        server get:3030;
    }

Resumen:
    Contenedor GET: expone puerto 3030
    Contenedor NGINX: escucha en el 3030 pero el servicio es accedido desde el 8080
                      redirige al puerto expuesto por GET