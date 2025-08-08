# StockInsight

StockInsight es una plataforma para la gestión, análisis y recomendación de acciones financieras, compuesta por un backend en Go y un frontend en Vue.js. El sistema permite almacenar información de acciones y finanzas, realizar migraciones y semillas de datos, consumir APIs externas, y mostrar recomendaciones y datos relevantes en una interfaz moderna.

## Estructura del Proyecto

### STI-backend

Backend desarrollado en Go, responsable de la lógica de negocio, persistencia de datos y exposición de APIs.

#### Principales Carpetas y Archivos

- **cmd/main.go**  
  Punto de entrada de la aplicación backend.

- **internal/api/router.go**  
  Define las rutas y endpoints expuestos por la API REST.

- **internal/db/**  
  - `connection.go`: Maneja la conexión a la base de datos (CockroachDB).
  - `migrate.go`: Ejecuta migraciones de esquema.
  - `migrations/`: Scripts SQL para crear y modificar tablas.
  - `seeds/`: Datos de ejemplo y utilidades para importar/exportar datos iniciales.

- **internal/finance/**  
  - `domain/`: Define entidades y repositorios de finanzas.
  - `infrastructure/persistence/`: Implementa la persistencia en CockroachDB.
  - `infrastructure/scraper/yahoo.go`: Scraper para obtener datos financieros desde Yahoo Finance.
  - `interfaces/finance_handler.go`: Handlers para los endpoints de finanzas.
  - `use-cases/update_data.go`: Lógica para actualizar datos financieros.

- **internal/stock/**  
  - `domain/`: Entidades y lógica de recomendación de acciones.
  - `infrastructure/api/external_api.go`: Integración con APIs externas de datos de acciones.
  - `infrastructure/repository/persistence_repo.go`: Persistencia de datos de acciones.
  - `interfaces/routes.go`: Rutas específicas para acciones.
  - `interfaces/stock_handler.go`: Handlers para endpoints de acciones.
  - `interfaces/sync_handler.go`: Sincronización de datos.
  - `use_cases/stock_service.go`: Lógica de negocio para acciones.
  - `use_cases/sync_service.go`: Servicio de sincronización de datos.

- **tests/stock_e2e_test.go**  
  Pruebas end-to-end para el módulo de acciones.

- **docs/**  
  Documentación y especificaciones Swagger para la API.

- **docker-compose.yml**  
  Orquestación de servicios para desarrollo y despliegue.

### STI-frontend

Frontend desarrollado en Vue.js, encargado de la visualización y experiencia de usuario.

#### Principales Carpetas y Archivos

- **src/App.vue & main.ts**  
  Configuración principal de la aplicación y punto de entrada.

- **src/components/**  
  Componentes reutilizables como barra de navegación, pie de página y tablas de datos.

- **src/components/data/**  
  - `StockDataTable.vue`, `StockTable.vue`: Tablas para mostrar datos de acciones.
  - `RecommendationCard.vue`, `StockRecommendations.vue`: Visualización de recomendaciones.

- **src/components/forms/FilterForm.vue**  
  Formulario para filtrar datos de acciones.

- **src/components/ui/popover/**  
  Componentes para popovers interactivos.

- **src/components/ui/tabs/**  
  Componentes para navegación por pestañas.

- **src/pages/HomePage.vue**  
  Página principal de la aplicación.

- **src/router/index.ts**  
  Configuración de rutas del frontend.

- **src/stores/**  
  - `stock.ts`: Store para el estado de acciones.
  - `recommendations.ts`: Store para el estado de recomendaciones.

- **src/types/Stock.d.ts**  
  Tipos TypeScript para datos de acciones.

- **public/**  
  Recursos estáticos como favicon.

- **package.json, vite.config.ts, tsconfig.json**  
  Configuración de dependencias, compilación y herramientas.

## ¿Qué hace StockInsight?

- Permite importar, almacenar y consultar información financiera y de acciones.
- Realiza migraciones y semillas de datos para inicializar la base de datos.
- Sincroniza y actualiza datos desde fuentes externas (ej. Yahoo Finance).
- Expone una API REST para interactuar con los datos desde el frontend.
- Proporciona una interfaz web moderna para visualizar recomendaciones, tablas de datos y realizar búsquedas y filtros.
- Facilita la gestión y análisis de portafolios de acciones, mostrando métricas y sugerencias basadas en datos actualizados.

## ¿Cómo funciona?

1. El backend gestiona la persistencia, lógica de negocio y exposición de APIs.
2. El frontend consume la API para mostrar datos y recomendaciones al usuario.
3. Los scripts de migración y semillas permiten inicializar y poblar la base de datos.
4. Los scrapers y servicios de sincronización mantienen los datos actualizados.
5. El usuario interactúa con la interfaz web para consultar información relevante y tomar decisiones informadas.

---

## Endpoints Principales (Backend)

La API REST expone los siguientes endpoints principales:

### Acciones (Stocks)

- `GET /stocks`  
  Obtiene la lista de acciones almacenadas.

### Recomendaciones

- `GET /recommendations`  
  Devuelve recomendaciones de acciones basadas en los datos actuales.

### Ejemplo de Uso (Consumo de API)

```bash
# Obtener todas las acciones
curl http://localhost:8080/api/stocks

# Sincronizar datos de acciones
curl -X POST http://localhost:8080/api/recommendations
```

---

## Despliegue y Ejecución

### Requisitos

- Docker y Docker Compose
- Go 1.20+
- Node.js y pnpm (para frontend)

### Backend

1. Clona el repositorio y accede a la carpeta `STI-backend`.
2. Configura las variables de entorno en `.env`.
3. Ejecuta los servicios con Docker Compose:

   ```powershell
   docker-compose up --build
   ```
4. Ejecuta la aplicación del backend

    ```powershell
        run ./STI-backend/cmd/main.go --serve
    ```
5. La API estará disponible en `http://localhost:8080`.

### Frontend

1. Accede a la carpeta `STI-frontend`.
2. Instala dependencias:

   ```powershell
   pnpm install
   ```

3. Ejecuta el servidor de desarrollo:

   ```powershell
   pnpm run dev
   ```

4. Accede a la interfaz web en `http://localhost:5173`.

---

## Ejemplo de Flujo de Usuario

1. El usuario accede a la web y visualiza la lista de acciones y recomendaciones.
2. Puede filtrar y buscar acciones usando el formulario.
3. Al seleccionar una acción, ve detalles y métricas relevantes.
4. El usuario recibe recomendaciones basadas en los datos recopilados.

---

## Para más información
Consultar `*/STI-backend/README.md` para obtener más información de como levantar el backend.

## Variables .ENV para el backend
Antes de iniciar la aplicación, asegúrate de configurar las siguientes variables de entorno en el archivo `.env` dentro de la carpeta `STI-backend`:

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
## Variables .ENV para el frontend
Antes de iniciar el frontend asegúrate de configurar las siguientes variables de entorno en el archivo `.env` dentro de la carpeta`STI-frontend`:
- `VITE_API_URL`: Url para el api que levanta en backend, normalmente http://localhost:8080.


Ejemplo de archivo `.env`:

```env
VITE_API_URL=http://localhost:8080
```

# Análisis de predicciones por bróker (vistas SQL)

Esta sección documenta el **núcleo del análisis** usado para evaluar objetivamente a los brókers a partir de sus predicciones históricas. Todo el análisis se realiza en **CockroachDB** mediante dos vistas: `broker_predictions` y `broker_evaluation`. El backend (Golang) consulta estas vistas y la UI (Vue 3 + TypeScript + Pinia + Tailwind, **Composition API** con `<script setup>`) presenta los resultados.

---

## Vistas creadas

### 1) `broker_predictions`
Evalúa cada predicción registrada en `stocks` contra el **cierre real más cercano** disponible en `finances` (por `ticker`) y calcula métricas por predicción:

- `prediction_date`: fecha de la predicción (`stocks.created_at`).
- `actual_price`: cierre real más cercano (`finances.close` via *LATERAL JOIN* por proximidad de fecha).
- `prediction_direction`: dirección esperada según `target_from` → `target_to` (`up`, `down`, `neutral`).
- `is_correct`: 1 si la dirección real del precio coincide con la dirección prevista, 0 en caso contrario.
- `error_percentage`: |`target_to` − `actual_price`| / `actual_price` * 100 (redondeado a 2 decimales).

> Nota: la vista filtra predicciones sin `target_to` y tolera faltantes de `target_from`/`actual_price` dejando valores nulos donde corresponde.

**Consulta típica:**

```sql
-- Predicciones evaluadas (con dirección, acierto y error)
SELECT *
FROM broker_predictions
WHERE prediction_date >= '2025-01-01'
ORDER BY prediction_date DESC;
```

### 2) `broker_evaluation`
Agrega las predicciones evaluadas por `brokerage` y genera una métrica compuesta:

- `total_predictions`: total de predicciones evaluadas por bróker.
- `total_hits`: cantidad de aciertos (dirección correcta).
- `accuracy`: % de acierto (con dos decimales).
- `weight_score`: señal compuesta que pondera **precisión** y **volumen** de aciertos:
  
  ```
  weight_score = total_hits * (accuracy / 100)
  ```

Esta métrica es útil para ponderar la “confiabilidad histórica” del bróker en algoritmos de ranking/recomendación.

**Consulta típica:**

```sql
-- Mejores brókers por peso histórico (precisión * volumen)
SELECT *
FROM broker_evaluation
ORDER BY weight_score DESC
LIMIT 10;
```

---

## Cómo se usa desde el backend (Go)

1. **Leer calidad del bróker**: consultar `broker_evaluation` para obtener `accuracy` y `weight_score` por `brokerage`.
2. **Cruzar con predicciones recientes**: consultar `broker_predictions` filtrando por ventana de fechas deseada.
3. **Componer un score final** (en Go) si se requiere un ranking para “hoy”:
   - Por ejemplo: combinar `weight_score` del bróker con señales actuales (p. ej., *upside* calculado contra el último `close` disponible) y ordenarlo descendentemente.
4. **Exponer por API**: endpoint REST del backend entrega las filas ordenadas y los campos explicativos (accuracy, total_hits, error_percentage).

> Las vistas están diseñadas para ser **explicables** y **auditables**: cada recomendación puede rastrearse hasta las métricas calculadas por predicción y por bróker.

---

## Ejemplos de consultas útiles

```sql
-- 1) Ver el historial del error porcentual por bróker y ticker
SELECT brokerage, ticker, prediction_date, error_percentage
FROM broker_predictions
WHERE brokerage = 'Morgan Stanley'
ORDER BY prediction_date DESC;

-- 2) Top brókers por precisión (no sólo por weight_score)
SELECT brokerage, total_predictions, total_hits, accuracy
FROM broker_evaluation
ORDER BY accuracy DESC, total_predictions DESC
LIMIT 10;

-- 3) Filtrar predicciones recientes (últimos 30 días)
SELECT *
FROM broker_predictions
WHERE prediction_date >= current_date - INTERVAL '30 days'
ORDER BY prediction_date DESC;
```

---

## Notas y consideraciones

- El emparejamiento con precios usa el **cierre más cercano en fecha** para robustez ante huecos de calendario.
- `weight_score` es **simple y transparente**; puede sustituirse o complementarse con otras métricas si el producto lo requiere a futuro.
- Toda la lógica de evaluación vive en SQL, lo que facilita su **mantenimiento**, **auditoría** y **performance** para el backend.

## Licencia

MIT

## Autor
**Vicente Chiriguaya M.**
[LinkedIn](https://www.linkedin.com/in/vchiriguaya) | [GitHub](https://github.com/viteant)  
Arquitecto de software disfrazado de full stack. Me obsesiona que las cosas funcionen, pero también que tengan sentido. Trabajo entre diseño de sistemas, automatización e inteligencia artificial, con preferencia por stacks limpios, estructuras predecibles y código que no sorprenda... salvo para bien.
