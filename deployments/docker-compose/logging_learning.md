# centralized Logging with Fluentd & OpenSearch

This guide explains the logging infrastructure added to the Unified Observability Stack.

## Architecture

```mermaid
graph LR
    App[Application] -->|Forward (24224)| Fluentd
    AppHTTP[Scripts/Curl] -->|HTTP (9880)| Fluentd
    Fluentd -->|Index: fluentd-*| OpenSearch
    OpenSearch -->|Query| Dashboards
```

-   **Fluentd**: Aggregates logs. Exposed on ports `24224` (Forward) and `9880` (HTTP).
-   **OpenSearch**: Stores logs. Port `9200`.
-   **OpenSearch Dashboards**: Visualization UI. Port `5601`.
-   **Application Logs**: `webdemosvc` and `natsdemosvc` automatically send logs to Fluentd via the Docker logging driver.

## Quick Start

1.  **Start the Stack**:
    ```bash
    cd deployments/docker-compose
    docker compose up -d
    ```

2.  **Verify Service Logs**:
    -   Open [http://localhost:5601](http://localhost:5601).
    -   Go to **Discover**.
    -   Create index pattern `fluentd-*`.
    -   Search for `container_name:webdemosvc` or `container_name:natsdemosvc`.

3.  **Send a Manual Test Log**:
    ```bash
    curl -X POST -d 'json={"logger":"test", "message":"Hello World"}' http://localhost:9880/app.test
    ```

3.  **Visualize**:
    -   Open [http://localhost:5601](http://localhost:5601).
    -   Go to **Stack Management** -> **Index Patterns**.
    -   Create pattern `fluentd-*`.
    -   Go to **Discover** to see your logs.

## Configuration

-   **Fluentd Config**: `deployments/docker-compose/config/fluentd/fluent.conf`
-   **Docker Compose**: `deployments/docker-compose/docker-compose.yml`
