## 8. Runtime

### Handler type
```go
type Handler func(*Context) Node
```
- Shared by `Page`, `Action`, `Render`, `Handle`.

### HTTP / rendering failure behavior (consolidated)
- `Page` endpoint + non-`GET` method: `405 Method Not Allowed`.
- `Action` endpoint + non-`POST` method: `405 Method Not Allowed`.
- `Action` request `ParseForm` failure: `400 Bad Request`.
- Node render failure (page/fragment): `500 Internal Server Error`.
- Shell render failure (page): `500 Internal Server Error`.
- Missing root page (`/` not registered): `GET /` returns `500 Internal Server Error`.

### Response content type
- Both full page and fragment writes set:
  - `Content-Type: text/html; charset=utf-8`.
