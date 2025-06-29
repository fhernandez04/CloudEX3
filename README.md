## Exercise 3 for Cloud Computing

### The challenge

1. Use Docker Compose to orchestrate your **five** containers and deploy **NGINX**
to control the traffic between your services based on the request method. The
containers are as follows:
    - A container to handle each operation for `/api/books`: GET, POST, PUT, and DELETE
    - A container to handle requests to `/*`, i.e., the rendering of the webpage.
    - A container for NGINX.
2. Use a Multi-stage Dockerfile to minimize the size of the image.
3. Publish to a public container registry (e.g., Docker Hub) your images.

## Explanations
- We split the code into microservices
- We create one Dockerfile per service (GET, POST, PUT, DELETE, HOMEPAGE)
- A docker-compose.yml file uses these Dockerfiles/images to spin up containers
- nginx.conf routes incoming traffic to each container based on the request type

What we're doing is splitting each of our endpoints (which isn't ideal in real-world scenarios, but serves well for simulating microservices) into separate containers. In our case, the GET endpoint for api/books is responsible for initializing and populating the database. That’s why all other containers wait for it to be up before connecting to the database using an environment variable named DATABASE_URI.

Request traffic to the server is managed by NGINX, which filters incoming HTTP requests based on whether they are GET, POST, PUT, DELETE, or generic (frontend), and forwards them to the corresponding container. NGINX runs in its own container.

Each service container exposes its internal port 3030 within the Docker internal network, meaning it cannot be accessed directly from outside. This way, we delegate all external access to NGINX.

Inside its container, NGINX listens on port 3030 for HTTP connections. However, externally (from your machine or any client), you don't access port 3030 directly — you use port 8080 on your host.

**Port 3030**: is the port NGINX listens to inside its container.

**Port 8080**: is the port on your host machine through which you access the service.

A request first reaches port 8080 on your machine, which is mapped to port 3030 of the NGINX container. Inside the container, NGINX receives the request on port 3030. Based on the route and HTTP method (GET, POST, PUT, DELETE, etc.), NGINX decides which microservice container to forward the request to.

The **nginx.conf** file defines which container and port each request should be forwarded to. For example, if the GET container exposes port 3030, the file would include:

    upstream get_service {
        server get:3030;
    }

Summary:
    GET Container: exposes port 3030
    NGINX Container: listens on port 3030 internally but is accessed externally through port 8080
    It routes the request to the exposed port of the appropriate microservice (e.g., GET)

--------------------------------------------------------------------------------------

## Explicaciones

- Dividimos el código en microservicios
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
                      
