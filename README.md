# ⚡ Flux

A production-ready full-stack framework for **Go**, inspired by **Laravel + Inertia.js**.
Write a Go controller, create a Vue page — Flux handles everything in between: routing,
page rendering, data serialization, validation, auth and SPA navigation.

```
Browser → Vue 3 → Flux Inertia Adapter → Go Router → Middleware → Controller → Service → Model (GORM) → MySQL
```

| Backend | Frontend |
|---|---|
| Go 1.24+, Gin (HTTP engine only) | Vue 3 + TypeScript + Vite |
| GORM + MySQL | Tailwind CSS 4 + shadcn-style components |
| JWT auth + httpOnly session cookie | Pinia |
| slog structured logging | Inertia-style page rendering (custom adapter) |
| DTO pattern, constructor DI, MVC | Axios |

## Quick start

```bash
# 1. Configure
copy .env.example .env          # set APP_KEY (any long random string) and DB_*

# 2. Install frontend deps
npm install

# 3. Build the CLI
go build -o flux.exe ./cmd/flux

# 4. Database (created automatically if missing)
./flux.exe migrate
./flux.exe seed                 # admin@flux.local / password + demo users

# 5. Run (Go server + Vite dev server together)
./flux.exe serve
```

Open http://localhost:8080 (or your `APP_PORT`) and sign in with `admin@flux.local` / `password`.

## The developer experience

The goal: you only write a controller method and a Vue page.

```go
// app/controllers/user_controller.go
func (c *UserController) Index(req *framework.Request) {
    users, pagination, err := c.service.GetUsers(req.Context(), page, 10, search)
    if err != nil {
        req.Error(err)
        return
    }
    req.View("Users/Index", dto.UserIndexDTO{Users: users, Pagination: pagination})
}
```

```vue
<!-- resources/js/Pages/Users/Index.vue — resolved automatically from "Users/Index" -->
<script setup lang="ts">
defineProps<{ users: UserDTO[]; pagination: PaginationDTO }>()  // props injected automatically
</script>
```

No API routes, no fetch calls, no manual wiring. `req.View("Users/Index", dto)` serializes
the DTO, and the Flux client renders `Pages/Users/Index.vue` with the DTO fields as props.
Navigation happens over XHR (Inertia protocol) — no full page reloads.

## Architecture

**Layers (no Repository pattern):** Controller → Service → Model, with DTOs at the boundary.

- **Controllers** (`app/controllers/`) — structs holding services. Read input, call the service, render. Never import Gin or GORM.
- **Services** (`app/services/`) — all business logic. Talk to GORM directly, return DTOs.
- **Models** (`app/models/`) — GORM models. Database concerns only; never serialized to Vue.
- **DTOs** (`app/dto/`) — every page has a props DTO; every write endpoint has a validated input DTO.

### Request object (controllers never touch Gin)

```go
req.Input("name")      // body field (JSON or form), falls back to query
req.Query("page")      // query string          req.QueryDefault("page", "1")
req.Param("id")        // route param            req.ParamUint("id")
req.Header("X-Foo")    // request header
req.Bind(&dto)         // decode body into struct
req.Validate(&dto)     // bind + validate; auto-responds 422 and returns false on failure
req.Context()          // context.Context for services/GORM
req.File("avatar")     // uploaded *multipart.FileHeader   req.SaveFile(f, dst)
req.User() / req.UserID()  // authenticated user (set by middleware.Auth)
```

### Response object

```go
req.View("Users/Index", dto)   // render Pages/Users/Index.vue with DTO as props
req.JSON(200, data)            // raw JSON
req.Success(data)              // {"success":true,"data":...}
req.Error(err)                 // status-aware (HTTPError, gorm.ErrRecordNotFound → 404, else 500)
req.Redirect("/users")         // 303 on non-GET so Inertia forms behave like Laravel
req.Download(path, "file.pdf") // attachment
req.ServeFile(path)            // inline file (req.File reads uploads; ServeFile sends files)
req.NoContent()                // 204
```

### Routing

```go
route.Get("/dashboard", ctrl.Index)
route.Post("/logout", auth.Logout)
route.Put("/users/:id", users.Update)
route.Delete("/users/:id", users.Destroy)
route.Group("/admin", middleware.Auth(app, authService))  // sub-router with middleware

route.Resource("/users", userController)
// GET    /users            → Index
// GET    /users/create     → Create
// POST   /users            → Store
// GET    /users/:id        → Show
// GET    /users/:id/edit   → Edit
// PUT    /users/:id        → Update   (PATCH too)
// DELETE /users/:id        → Destroy
```

`Resource` uses reflection: implement any subset of the seven actions with signature
`func(*framework.Request)` and only those are registered.

### Validation

DTOs carry `validate` tags (go-playground/validator) and `json` tags (field names in errors):

```go
type CreateUserDTO struct {
    Name     string `json:"name"     validate:"required,min=2,max=255"`
    Email    string `json:"email"    validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}
```

```go
var input dto.CreateUserDTO
if !req.Validate(&input) {
    return // 422 with {"errors": {"email": ["The email must be a valid email address."]}} already sent
}
```

On the Vue side, `useForm` catches the 422 automatically:

```ts
const form = useForm({ name: '', email: '', password: '' })
form.post('/users')          // on failure: form.errors.email === "The email must be a valid..."
```

### Authentication

Complete flow included: **login, logout, register, forgot/reset password, email
verification** — built on JWT (signed with `APP_KEY`) delivered two ways:

- **Session:** httpOnly cookie (`AUTH_SESSION_COOKIE`) for browser navigation.
- **API:** `Authorization: Bearer <token>` header.

Verification and reset links are short-lived purpose-scoped JWTs. Mail uses
`MAIL_DRIVER=log` (writes to `storage/logs/mail.log`) in development or `smtp` in production.

### Frontend adapter (`resources/js/flux/`)

- `router.ts` — Inertia protocol client: XHR visits with `X-Inertia` headers, history
  push/pop, 409 asset-version reloads, 401 → login, 422 → error bags.
- `Link.vue` — SPA `<a>` replacement.
- `useForm.ts` — reactive forms with `processing` + `errors`.
- `usePage.ts` / Pinia `stores/auth.ts` — shared props (`auth.user`, `appName`) everywhere.

## CLI

```
flux new <name>               scaffold a project (set FLUX_TEMPLATE_REPO)
flux serve                    Go server + Vite dev server together
flux build                    vite build + go build → bin/app
flux migrate                  run pending migrations (creates the database if missing)
flux seed                     seed demo data
flux make:controller Post     app/controllers/post_controller.go (full resource)
flux make:model Post          app/models/post.go
flux make:service Post        app/services/post_service.go
flux make:dto Post            app/dto/post.go (page DTO + create/update DTOs)
flux make:middleware Rate     app/middleware/rate.go
flux make:migration create_posts_table
flux make:page Posts/Index    resources/js/Pages/Posts/Index.vue
```

## Folder structure

```
app/
  controllers/     HTTP layer (structs, framework.Request only)
  services/        business logic (GORM, return DTOs)
  models/          GORM models
  dto/             page props + validated inputs
  middleware/      Auth, Guest, Verified, ...
  routes/          web.go — the route map + DI wiring
framework/         Flux core: router, request/response, inertia, jwt, validation, mail, slog
config/            .env loader (stdlib only)
database/          connection, migrator, migrations/, seeders/
resources/
  views/app.html   root template (page object + Vite tags injected)
  css/app.css      Tailwind 4 + shadcn design tokens
  js/
    Pages/         one .vue per req.View() component name
    Components/    ui/ (shadcn-style button, input, card, table, badge, label)
    Layouts/       DashboardLayout, AuthLayout
    flux/          the Inertia-style adapter
    stores/        Pinia (auth)
    app.ts         bootstrap    router.ts  navigation facade
storage/logs/      flux.log, mail.log
public/build/      production bundle (flux build)
cmd/flux/          the CLI
```

## Production

```bash
flux build                    # bundle + binary
APP_ENV=production ./bin/app  # serves assets from public/build via the Vite manifest
```

Set a strong `APP_KEY`, real `DB_*` credentials and `MAIL_DRIVER=smtp`. In production the
logger switches to JSON, Gin runs in release mode, and asset versioning triggers client
reloads after each deploy.

## Adding a feature (the Flux loop)

```bash
flux make:model Post && flux make:migration create_posts_table
flux make:dto Post && flux make:service Post && flux make:controller Post
flux make:page Posts/Index
```

then register `route.Resource("/posts", controllers.NewPostController(...))` in
`app/routes/web.go`, fill in the stubs, and `flux migrate`.
