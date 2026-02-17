# Sistema de Gestión de Streaming

Proyecto académico desarrollado para la asignatura **Programación Orientada a Objetos**.

## Descripción

El presente proyecto consiste en el desarrollo de un sistema de gestión de streaming que simula el funcionamiento básico de una plataforma de transmisión de contenidos audiovisuales.

El sistema permite:

- Registro e inicio de sesión de usuarios.
- Gestión de roles (Administrador y Usuario).
- Exploración y filtrado de un catálogo de contenidos.
- Reproducción multimedia con almacenamiento de progreso.
- Sección “Continúa viendo”.
- Control de acceso mediante suscripciones.
- Administración completa del catálogo por parte del rol administrador.

El diseño del sistema aplica principios fundamentales de la Programación Orientada a Objetos, arquitectura por capas y manejo estructurado de errores.


## Arquitectura del Sistema

El sistema está basado en una arquitectura **cliente–servidor**, compuesta por:

- **Backend:** API REST desarrollada en Go.
- **Frontend:** Interfaz web en HTML5, CSS3 y JavaScript.
- **Comunicación:** JSON sobre HTTP mediante API REST.

Se aplica separación por capas:

manejadores → Capa de presentación (API REST)
repositorios → Acceso a datos (persistencia JSON)
modelos → Entidades del sistema
middleware → Seguridad (JWT y roles

## Módulos del Sistema

### 1.Autenticación y Gestión de Usuarios
- Registro de usuarios.
- Inicio de sesión con JWT.
- Control de acceso por roles.
- Middleware de autenticación y autorización.

### 2.Catálogo de Contenidos
- Listado general.
- Búsqueda por título.
- Filtros por género, tipo y año.
- Visualización de detalle.
- Sección “Mi Lista”.

### 3.Reproducción de Contenidos
- Reproductor multimedia HTML5.
- Guardado automático de progreso.
- Continuación desde el último segundo visualizado.
- Sección “Continúa viendo”.

### 4.Gestión de Suscripciones
- Estado de suscripción (ACTIVE / EXPIRED / CANCELED).
- Restricción de reproducción según estado.

### 5.Administración del Catálogo
- Crear contenidos.
- Editar contenidos.
- Activar / desactivar contenidos.
- Eliminación lógica.
- Acceso exclusivo para administradores.

## Conceptos de Programación Orientada a Objetos Aplicados

- Encapsulación mediante estructuras y control de acceso interno.
- Uso de interfaces para desacoplar la lógica de negocio del acceso a datos.
- Manejo estructurado de errores HTTP (400, 401, 403, 404, 500).
- Separación de responsabilidades por capas.
- Uso de constructores (`NewXxx`) para inicialización controlada.
- 
## Tecnologías Utilizadas

- **Lenguaje Backend:** Go (Golang)
- **Framework HTTP:** Gin
- **Autenticación:** JWT
- **Encriptación de contraseñas:** bcrypt
- **Frontend:** HTML5, CSS3, JavaScript
- **Persistencia:** Archivos JSON
- **Servidor estático:** Gin (static files)

## Estructura del Proyecto

backend/
├── cmd/api
├── interno/
│ ├── manejadores
│ ├── modelos
│ ├── repositorios
│ ├── middleware
├── datos/
├── static/videos/

frontend/
├── pagini/
├── activos/js/
├── activos/css/

## Ejecución del Proyecto

### Backend

```
cd backend
go run cmd/api/main.go
```

Servidor disponible en:

http://localhost:8080


### Frontend

Puede ejecutivo mediante:

Servidor en vivo (código VS), o

Servidor HTTP de Python:

```
cd frontend
python -m http.server 5500
```


Acceso desde:

http://localhost:5500

## Autor

Esteban Nicolás Simbaña Mora
Proyecto académico – Programación Orientada a Objetos
