# StockInsight Backend

Este proyecto es el backend de StockInsight, desarrollado en Go y utilizando una base de datos compatible con CockroachDB.

## Requisitos

- [Go 1.21+](https://golang.org/dl/)
- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)

## Instalación

1. **Clona el repositorio:**

   ```sh
   git clone https://github.com/tu-usuario/StockInsight.git
   cd StockInsight/go-backend
   ```
2. **Configura el entorno:**

Antes de iniciar la aplicación, asegúrate de configurar las siguientes variables de entorno en el archivo `.env` dentro de la carpeta `backend`:

- `ENVIRONMENT`: Entorno de ejecución, puede ser `env` o `prod`.
- `DATABASE_URI`: URI de conexión a la base de datos (ejemplo: `root:@localhost:26257/stockinsights?sslmode=disable`)
- `API_ENDPOINT`: Endpoint de la API externa de stocks.
- `API_TOKEN`: Token de autenticación para la API.

Ejemplo de archivo `.env`:

```env
ENVIRONMENT=env
DATABASE_URI="root:@localhost:26257/stockinsights?sslmode=disable"
API_ENDPOINT="https://api.example.com/list"
API_TOKEN="tu_token"
```

Asegúrate de que el archivo .env esté correctamente configurado antes de ejecutar migraciones y levantar el backend

3. **Configura la base de datos:**

   Inicia los servicios con Docker Compose:

   ```sh
   docker-compose up -d
   ```

   Esto levantará la base de datos y otros servicios necesarios.

4. **Instala las dependencias de Go:**

   ```sh
   go mod download
   ```

5. **Ejecuta las migraciones de la base de datos:**

   Asegúrate de que la base de datos esté corriendo y ejecuta el comando para migrar:

   ```sh
   go run cmd/main.go --migrate
   ```

## Uso general

```bash
go run main.go [comando] [flags]
```

## Comandos disponibles

### `--migrate`

Ejecuta las migraciones de la base de datos.

```bash
go run main.go --migrate
```

#### Opcional: `--reset`

Borra la base de datos antes de migrar.

```bash
go run main.go --migrate --reset
```

---

### `--sync`

Sincroniza los datos de **stocks** desde la API externa.

```bash
go run main.go --sync
```

---

### `--update-finance`

Actualiza los datos históricos financieros de Yahoo Finance por cada ticker.

```bash
go run main.go --update-finance
```

---

### `--serve`

Inicializa el servidor de la aplicación.

```bash
go run main.go --serve
```

> Si no se pasa ningún flag, también se ejecuta este comando por defecto.

---

### 📤 `--export` y `--table`

Exporta datos desde la base a un archivo `.json`.

>Ejemplo de exportación de stocks y finanzas con datos de ejemplos incluídos:

```bash
go run cmd/main.go --export=internal/db/seeds/stocks_seed.json --table="stocks"
go run cmd/main.go --export=internal/db/seeds/finances_seed.json --table="finances"
```

---

### 📥 `--import` y `--table`

Importa datos desde un archivo `.json` hacia la base de datos.

```bash
go run cmd/main.go --import=internal/db/seeds/import.json --table="stocks"
go run cmd/main.go --import=internal/db/seeds/finances_seed.json --table="finances"
```

---

## Endpoints disponibles

El servidor expone los siguientes endpoints REST bajo el path base `/api`:

### `GET /api/stocks`

Obtiene una lista paginada de acciones.

**Parámetros de consulta disponibles:**

- `page`: número de página (por defecto: 1)
- `limit`: cantidad por página (por defecto: 20)
- `ticker`: filtra por símbolo (ILIKE)
- `company`: filtra por nombre de empresa (ILIKE)
- `brokerage`: filtra por nombre del bróker (ILIKE)
- `target_from_min`: valor mínimo para target_from
- `target_from_max`: valor máximo para target_from
- `target_to_min`: valor mínimo para target_to
- `target_to_max`: valor máximo para target_to
- `date_from`: fecha mínima de creación (`YYYY-MM-DD`)
- `date_to`: fecha máxima de creación (`YYYY-MM-DD`)
- `orderBy`: campo por el cual ordenar (`ticker`, `company`, `created_at`, etc.)
- `orderDir`: dirección del orden (`asc` o `desc`)

### `GET /api/recommendations`

Obtiene una lista de recomendaciones agrupadas por tipo (`buy`, `hold`, `sell`) basada en el puntaje (`weight_score`) de los brokers.

---

## Documentación Swagger (OpenAPI)

Este proyecto incluye documentación autogenerada con [Swaggo](https://github.com/swaggo/swag).

### 📦 Instalación de Swaggo

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

Agrega `$GOPATH/bin` a tu `PATH` si no puedes ejecutar `swag`.

### 📄 Generar la documentación

Ejecuta en la raíz del proyecto:

```bash
swag init --dir ./cmd,./internal --output ./docs
```

Esto escaneará los comentarios de tus handlers y generará los archivos `docs/swagger.json` y `docs/swagger.yaml`.

### 🔍 Ver documentación

Una vez que corras el servidor (`--serve`), abre en el navegador:

```
http://localhost:3000/swagger/index.html
```

## Testing

### 🧪 Pruebas E2E

El proyecto incluye pruebas de extremo a extremo para validar el comportamiento real de la API, incluyendo validación de filtros dinámicos.

Ejemplo:

```go
func TestGetStocksE2E_WithDynamicFilters(t *testing.T)
```

Este test:

- Consulta `GET /api/stocks?limit=5`
- Usa el primer resultado como base
- Realiza una segunda consulta aplicando filtros con esos datos
- Valida que los filtros hayan sido aplicados correctamente

Para ejecutar las pruebas:

```bash
go test ./tests -v
```

Asegúrate de tener una base de datos con datos válidos antes de ejecutar.


## Estructura del proyecto

- `cmd/main.go`: Punto de entrada de la aplicación.
- `internal/db/`: Conexión, migraciones y seeds de la base de datos.
- `internal/finance/`: Lógica de finanzas.
- `internal/stock/`: Lógica de stocks.

## Notas

- Puedes modificar la configuración de la base de datos en el archivo `docker-compose.yml`.
- Para detener los servicios de Docker:

  ```sh
  docker-compose down
  ```

## Licencia

MIT

## Autor
**Vicente Chiriguaya M.**
[LinkedIn](https://www.linkedin.com/in/vchiriguaya) | [GitHub](https://github.com/viteant)  
Arquitecto de software disfrazado de full stack. Me obsesiona que las cosas funcionen, pero también que tengan sentido. Trabajo entre diseño de sistemas, automatización e inteligencia artificial, con preferencia por stacks limpios, estructuras predecibles y código que no sorprenda... salvo para bien.


