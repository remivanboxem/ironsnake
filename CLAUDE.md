# CLAUDE.md - IronSnake Project Guide

## Project Overview

**IronSnake** is an educational Learning Management System (LMS) designed for programming courses at universities. The platform enables teachers to create programming assignments, quizzes, and exercises, while students can submit solutions with real-time code execution feedback in secure, sandboxed Docker containers.

### Core Capabilities
- **Course Management**: YAML-based course definitions with file-system persistence
- **Real-time Code Execution**: Sandboxed Docker containers with strict resource limits
- **Multiple Problem Types**: Code exercises, multiple choice questions, fill-in-the-blank
- **Access Control**: Role-based access with admin/tutor/student hierarchy
- **Syllabus Support**: mdBook-compatible syllabus with table of contents
- **i18n Ready**: Paraglide-based internationalization (English/French)

### Target Users
| Role | Capabilities |
|------|-------------|
| **Students** | Submit work, view feedback, track progress |
| **Teachers** | Create tasks, review submissions, provide feedback |
| **Administrators** | Bulk course/user management, configure execution pods |

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────────┐
│                              Client Browser                              │
└─────────────────────────────────┬───────────────────────────────────────┘
                                  │ HTTP/HTTPS
                                  ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                         Nginx Reverse Proxy (:80)                        │
│  ┌─────────────────────────────┬─────────────────────────────────────┐  │
│  │    /api/*  → core:8080     │        /*  → front:3000             │  │
│  └─────────────────────────────┴─────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────┘
         │                                           │
         ▼                                           ▼
┌─────────────────────┐                   ┌─────────────────────────────┐
│   Go Backend API    │                   │    SvelteKit Frontend       │
│   (core:8080)       │                   │    (front:3000)             │
├─────────────────────┤                   ├─────────────────────────────┤
│ • Course handlers   │                   │ • SSR + SPA routing         │
│ • Code execution    │                   │ • CodeMirror editors        │
│ • YAML parsing      │                   │ • Reactive UI (Svelte 5)    │
│ • DB operations     │                   │ • Markdown rendering        │
└──────────┬──────────┘                   └─────────────────────────────┘
           │
     ┌─────┴─────┐
     ▼           ▼
┌─────────┐  ┌─────────────────────┐
│Postgres │  │  Docker Daemon      │
│  (:5432)│  │  (code execution)   │
└─────────┘  └─────────────────────┘
```

---

## Technology Stack

| Layer | Technology | Version | Purpose |
|-------|------------|---------|---------|
| **Backend** | Go | 1.24 | API server, business logic |
| **ORM** | GORM | 1.25 | Database operations |
| **Frontend** | SvelteKit | 2.x | Full-stack web framework |
| **UI Framework** | Svelte | 5.x | Reactive components |
| **Database** | PostgreSQL | 18 | User/course metadata storage |
| **Code Editor** | CodeMirror | 6.x | In-browser code editing |
| **Styling** | Tailwind CSS | 4.x | Utility-first CSS |
| **UI Components** | bits-ui | 2.x | Accessible component library |
| **Markdown** | markdown-it + KaTeX | - | Rich text with math support |
| **i18n** | Paraglide | 2.x | Internationalization |
| **Containerization** | Docker | - | Service orchestration + code sandboxing |
| **Reverse Proxy** | Nginx | - | Request routing, static files |
| **Task Runner** | just | - | Build automation |
| **Documentation** | Storybook | 10.x | Component documentation |

---

## Project Structure

```
ironsnake/
├── CLAUDE.md                # This file - project guide
├── docker-compose.yml       # Production Docker configuration
├── docker-compose.dev.yaml  # Development Docker configuration (hot reload)
├── justfile                 # Task runner commands
│
├── core/                    # ═══ GO BACKEND ═══
│   ├── main.go              # Entry point, HTTP server, route registration
│   ├── course_handlers.go   # REST handlers for courses/tasks
│   ├── code_executor.go     # Docker-based code execution engine
│   ├── models.go            # GORM models (User, Course, Role, etc.)
│   ├── database.go          # PostgreSQL connection, migrations
│   ├── seed.go              # Development data seeding
│   ├── types.go             # API response DTOs
│   ├── Dockerfile           # Production multi-stage build
│   ├── Dockerfile.dev       # Development with hot reload
│   ├── go.mod               # Go module definition
│   └── courseparser/        # ─── YAML Course Parser Package ───
│       ├── courseparser.go  # Main course/task loader
│       ├── config.go        # config.yaml parsing
│       ├── access.go        # access.yaml parsing
│       ├── task.go          # task.yaml parsing with polymorphic problems
│       ├── syllabus.go      # mdBook syllabus parsing (book.toml, SUMMARY.md)
│       ├── types.go         # Problem interfaces and structs
│       ├── errors.go        # Custom error types
│       └── courseparser_test.go
│
├── front/                   # ═══ SVELTEKIT FRONTEND ═══
│   ├── src/
│   │   ├── app.html         # HTML shell template
│   │   ├── app.d.ts         # Global TypeScript declarations
│   │   ├── hooks.server.ts  # Server-side hooks
│   │   ├── hooks.ts         # Client-side hooks
│   │   ├── routes/          # ─── File-based Routing ───
│   │   │   ├── +layout.svelte         # Root layout (header, theme)
│   │   │   ├── +page.svelte           # Home: course list, recent tasks
│   │   │   └── courses/
│   │   │       └── [id]/
│   │   │           ├── +page.svelte   # Course detail: tasks, syllabus
│   │   │           └── tasks/
│   │   │               └── [taskId]/
│   │   │                   └── +page.svelte  # Task: problems, code editor
│   │   └── lib/
│   │       ├── index.ts             # Lib entry point
│   │       ├── utils.ts             # Utility functions
│   │       ├── services/            # ─── API Client Layer ───
│   │       │   ├── api-client.ts    # Base fetch wrapper with error handling
│   │       │   ├── course.service.ts # Course/task API calls
│   │       │   ├── code.service.ts  # Code execution API
│   │       │   └── index.ts         # Service exports
│   │       ├── types/               # ─── TypeScript Interfaces ───
│   │       │   ├── course.ts        # Course, Task, Problem types
│   │       │   └── index.ts
│   │       ├── components/          # ─── UI Components ───
│   │       │   ├── Header.svelte    # App header with navigation
│   │       │   └── ui/              # Reusable UI primitives (bits-ui based)
│   │       │       ├── button/
│   │       │       ├── card/
│   │       │       ├── dropdown-menu/
│   │       │       ├── avatar/
│   │       │       ├── input/
│   │       │       └── markdown/    # Markdown renderer with KaTeX
│   │       └── paraglide/           # Generated i18n runtime
│   ├── messages/                    # i18n translation files
│   │   ├── en.json
│   │   └── fr.json
│   ├── static/                      # Static assets
│   ├── storybook-static/            # Built Storybook docs
│   ├── stories/                     # Storybook stories
│   ├── package.json
│   ├── svelte.config.js
│   ├── vite.config.ts
│   ├── tsconfig.json
│   ├── Dockerfile                   # Production build
│   └── Dockerfile.dev               # Development with HMR
│
├── proxy/                   # ═══ NGINX REVERSE PROXY ═══
│   ├── nginx.conf           # Routing rules: /api/* → backend, /* → frontend
│   └── Dockerfile
│
├── courses/                 # ═══ COURSE CONTENT (YAML) ═══
│   └── CS01/                # Example course directory
│       ├── config.yaml      # Course metadata (name, admins, tutors)
│       ├── access.yaml      # Task accessibility rules
│       ├── README.md        # Course description
│       ├── syllabus/        # mdBook-compatible syllabus
│       │   ├── book.toml    # Book metadata
│       │   ├── SUMMARY.md   # Table of contents
│       │   └── src/         # Syllabus content markdown
│       └── tasks/           # Task definitions
│           ├── task01/
│           │   ├── task.yaml         # Task config with problems
│           │   ├── run               # Grading script
│           │   └── *.py              # Support files
│           └── task02/
│
└── documentation/           # ═══ PROJECT DOCUMENTATION ═══
    └── OVERVIEW.md          # Feature requirements by role
```

---

## Quick Commands

```bash
# Development
just run-dev          # Start dev environment with hot reload
just logs             # Follow container logs
just stop             # Stop all containers
just build-dev        # Rebuild dev images

# Production
just build            # Build production images
just run              # Start production environment

# Cleanup
just clean            # Remove orphaned containers
```

---

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/courses` | List all courses |
| GET | `/api/courses/:id` | Get course detail with tasks and syllabus |
| GET | `/api/courses/:id/tasks/:taskId` | Get full task with problems |
| POST | `/api/run` | Execute code in sandboxed container |

### Request/Response Examples

**List Courses**
```bash
GET /api/courses
```
```json
[
  {
    "id": "CS01",
    "code": "CS01",
    "name": "[CS01] Introduction to Computer Science",
    "accessible": false,
    "admins": ["jmachin", "nicolasm"],
    "tutors": ["damiend"],
    "taskCount": 7
  }
]
```

**Get Course Detail**
```bash
GET /api/courses/CS01
```
```json
{
  "id": "CS01",
  "name": "[CS01] Introduction to Computer Science",
  "tasks": [
    {
      "id": "task01",
      "name": "Convertisseur binaire vers Base64",
      "author": "Jean Machin",
      "environmentType": "docker",
      "problems": [
        { "id": "binary_to_base64", "type": "code", "name": "..." }
      ]
    }
  ],
  "syllabus": {
    "title": "Course Title",
    "author": "Author Name",
    "summary": [...]
  }
}
```

**Execute Code**
```bash
POST /api/run
Content-Type: application/json

{ "code": "print('Hello, World!')", "language": "python" }
```
```json
{ "output": "Hello, World!\n", "error": "", "exitCode": 0 }
```

---

## Key Architectural Patterns

### 1. Filesystem-Backed Course Content (Content-as-Code)
Courses are defined in YAML files under `courses/`, not in the database. This design enables:
- **Version Control**: Course content can be tracked in Git
- **Easy Migration**: Courses are portable directories
- **Offline Authoring**: Teachers can edit YAML locally
- **Separation of Concerns**: User data (DB) vs course content (files)

```
courses/CS01/
├── config.yaml      # Course metadata
├── access.yaml      # Task scheduling/permissions
└── tasks/task01/
    ├── task.yaml    # Task definition with problems
    └── grading/     # Test scripts (future)
```

### 2. Polymorphic Problem Types with Interface Pattern
Problems use a `type` discriminator field in YAML, parsed into Go interfaces:

```go
// core/courseparser/types.go
type Problem interface {
    GetType() string
    GetName() string
    GetHeader() string
}

type CodeProblem struct {
    BaseProblem `yaml:",inline"`
    Language    string `yaml:"language"`
    Default     string `yaml:"default"`
}

type MultipleChoiceProblem struct {
    BaseProblem `yaml:",inline"`
    Choices     []Choice `yaml:"choices"`
    Limit       int      `yaml:"limit"`
}

type MatchProblem struct {
    BaseProblem `yaml:",inline"`
    Answer      string `yaml:"answer"`
}
```

This enables:
- Type-safe problem handling in Go
- Easy addition of new problem types
- Polymorphic YAML deserialization via `UnmarshalYAML`

### 3. Sandboxed Code Execution (Docker-in-Docker Pattern)
User code runs in ephemeral Docker containers with strict isolation:

```go
// core/code_executor.go - Security Configuration
docker run --rm \
  --network none \              # No network access
  --memory 128m \               # Memory limit
  --cpus 0.5 \                  # CPU throttle
  --pids-limit 64 \             # Fork bomb protection
  --read-only \                 # Immutable filesystem
  --tmpfs /tmp:size=10m \       # Limited scratch space
  --security-opt no-new-privileges \  # No privilege escalation
  -v /code.py:/code.py:ro \     # Read-only code mount
  python:3.14-slim python /code.py
```

**Execution Flow:**
1. Frontend sends code to `/api/run`
2. Backend writes code to temp file
3. Docker container executes with limits
4. Output captured and returned (stdout, stderr, exit code)
5. Container destroyed automatically (`--rm`)

### 4. Service Layer Pattern (Frontend)
API calls are abstracted through typed service classes:

```typescript
// front/src/lib/services/course.service.ts
export const courseService = {
    async getAllCourses(): Promise<Course[]> {
        return apiGet<Course[]>('/courses');
    },
    async getCourseById(id: string): Promise<CourseDetail> {
        return apiGet<CourseDetail>(`/courses/${id}`);
    },
    async getTaskById(courseId: string, taskId: string): Promise<TaskDetail> {
        return apiGet<TaskDetail>(`/courses/${courseId}/tasks/${taskId}`);
    }
};

// front/src/lib/services/code.service.ts
export const codeService = {
    async runCode(code: string, language: string): Promise<RunCodeResponse> {
        return apiPost<RunCodeResponse, RunCodeRequest>('/run', { code, language });
    }
};
```

### 5. Component-Driven UI with Headless Components
UI is built using bits-ui (headless) + Tailwind:

```svelte
<!-- Composable card usage -->
<Card.Root class="hover:shadow-lg">
    <Card.Header>
        <Card.Title>{task.name}</Card.Title>
        <Card.Description>{task.author}</Card.Description>
    </Card.Header>
    <Card.Content>...</Card.Content>
</Card.Root>
```

### 6. Reactive State with Svelte 5 Runes
The frontend uses Svelte 5's new reactivity system:

```svelte
<script lang="ts">
    let task: TaskDetail | null = $state(null);
    let loading = $state(true);
    let editors: Map<string, EditorView> = new SvelteMap();
    
    // Derived state
    let currentMode = $derived(mode.current);
    
    // Effects for side effects
    $effect(() => {
        // React to theme changes
        editors.forEach((editor, problemId) => {
            // Update editor themes
        });
    });
</script>
```

---

## Data Model

### Database Schema (PostgreSQL via GORM)

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│      User       │     │  CourseTeacher  │     │     Course      │
├─────────────────┤     ├─────────────────┤     ├─────────────────┤
│ id (uuid)       │◄────│ teacher_id      │     │ id (uuid)       │
│ username        │     │ course_id       │────►│ code            │
│ email           │     │ role_id         │     │ name            │
│ password_hash   │     └────────┬────────┘     │ description     │
│ first_name      │              │              │ academic_year   │
│ last_name       │              ▼              │ created_by      │
│ created_at      │     ┌─────────────────┐     │ created_at      │
└─────────────────┘     │      Role       │     └─────────────────┘
                        ├─────────────────┤
                        │ id (uuid)       │
                        │ name (unique)   │
                        └─────────────────┘
```

### Filesystem Schema (Course Content)

```yaml
# config.yaml - Course metadata
name: "[CS01] Introduction to Computer Science"
accessible: false
admins: [jmachin, nicolasm]
tutors: [damiend]
groups_student_choice: false
allow_unregister: true
registration: true

# access.yaml - Task scheduling
dispenser_data:
  config:
    task01:
      accessibility: 2026-01-25 19:15:03/2026-01-29 19:15:07/2026-01-28 19:15:04
      evaluation_mode: best
      no_stored_submissions: 3
      submission_limit:
        amount: 1
        period: 1  # minutes

# task.yaml - Task definition
author: Jean Machin
name: Binary to Base64 Converter
environment_type: docker
environment_id: python3
environment_parameters:
  limits:
    time: "3"
    memory: "100"
problems:
  problem_id:
    type: code
    name: Problem Name
    header: |
      Markdown description with :math:`LaTeX` support
    language: python
    default: |
      def solution():
          pass
```

---

## Frontend Route Structure

| Route | Component | Purpose |
|-------|-----------|---------|
| `/` | `+page.svelte` | Dashboard: recent tasks, course list |
| `/courses/[id]` | `courses/[id]/+page.svelte` | Course detail: tasks, syllabus, staff |
| `/courses/[id]/tasks/[taskId]` | `tasks/[taskId]/+page.svelte` | Task view: problems, code editor, output |

### Task Page Features
- **CodeMirror Integration**: Syntax-highlighted Python editor
- **Theme Sync**: Editor theme follows system dark/light mode
- **Per-Problem State**: Each problem has its own editor instance
- **Async Execution**: Non-blocking code runs with loading states
- **Markdown Rendering**: Problem headers with KaTeX math support

---

## Development Workflow

### Starting the Development Environment
```bash
just run-dev          # Starts all services with hot reload
just logs             # View container logs
```

Services in development mode:
- **Frontend**: Vite HMR on port 3000
- **Backend**: Air hot reload on port 8080
- **Database**: PostgreSQL on port 5432
- **Proxy**: Nginx on port 80

### Adding a New Course
1. Create directory: `courses/COURSE_CODE/`
2. Add `config.yaml`:
   ```yaml
   name: "[CODE] Course Name"
   accessible: true
   admins: [admin_username]
   tutors: [tutor_username]
   ```
3. Add `access.yaml`:
   ```yaml
   dispenser_data:
     config:
       task01:
         accessibility: true
   ```
4. Create `tasks/task01/task.yaml` with problems
5. Backend automatically picks up changes (no restart needed)

### Adding a New Problem Type
1. **Backend - Define type** in `core/courseparser/types.go`:
   ```go
   type NewProblemType struct {
       BaseProblem `yaml:",inline"`
       CustomField string `yaml:"custom_field"`
   }
   ```
2. **Backend - Add parsing** in `core/courseparser/task.go` `UnmarshalYAML`:
   ```go
   case "new_type":
       var p NewProblemType
       if err := valueNode.Decode(&p); err != nil {
           return err
       }
       problem = &p
   ```
3. **Backend - Add response type** in `core/types.go`:
   ```go
   type ProblemDetailResponse struct {
       // ... existing fields
       CustomField string `json:"customField,omitempty"`
   }
   ```
4. **Frontend - Add TypeScript type** in `front/src/lib/types/course.ts`:
   ```typescript
   export interface ProblemDetail extends Problem {
       // ... existing fields
       customField?: string;
   }
   ```
5. **Frontend - Add rendering** in task page `+page.svelte`

### Adding a New API Endpoint
1. **Create handler** in `core/course_handlers.go` or new file:
   ```go
   func newHandler(w http.ResponseWriter, r *http.Request) {
       // Implementation
   }
   ```
2. **Register route** in `core/main.go`:
   ```go
   http.HandleFunc("/new-endpoint", newHandler)
   ```
3. **Add service method** in `front/src/lib/services/`:
   ```typescript
   async newMethod(): Promise<Response> {
       return apiGet<Response>('/new-endpoint');
   }
   ```

### Running Tests
```bash
# Backend tests
cd core && go test ./...

# Frontend type checking
cd front && pnpm check

# Frontend linting
cd front && pnpm lint

# Storybook
cd front && pnpm storybook
```

---

## Docker Configuration

### Production Architecture
```yaml
# docker-compose.yml
services:
  postgres:     # PostgreSQL 18, persistent volume
  core:         # Go binary, multi-stage build
  front:        # Node.js SvelteKit, SSR
  proxy:        # Nginx routing
```

### Development Architecture
```yaml
# docker-compose.dev.yaml
services:
  postgres:     # Shared with production
  core:         # Source mount + Air hot reload
    volumes:
      - ./core:/app           # Live code
      - ./courses:/app/courses # Course content
      - /var/run/docker.sock  # Docker-in-Docker for code exec
      - /tmp/ironsnake-code   # Shared temp for code files
  front:        # Source mount + Vite HMR
    volumes:
      - ./front:/app
  proxy:        # Same as production
```

### Code Execution Docker-in-Docker
The core container needs Docker socket access to spawn execution containers:
```yaml
volumes:
  - /var/run/docker.sock:/var/run/docker.sock  # Host Docker daemon
  - /tmp/ironsnake-code:/tmp/ironsnake-code    # Shared code directory
```

---

## Environment Variables

```bash
# Database Configuration
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres      # CHANGE IN PRODUCTION
POSTGRES_DB=ironsnake
POSTGRES_HOST=postgres
POSTGRES_PORT=5432
POSTGRES_SSLMODE=disable

# Container Names (optional)
POSTGRES_CONTAINER_NAME=ironsnake-postgres
CORE_CONTAINER_NAME=ironsnake-core
FRONT_CONTAINER_NAME=ironsnake-front
PROXY_CONTAINER_NAME=ironsnake-proxy

# Ports
CORE_PORT=8080
FRONT_PORT=3000
PROXY_PORT=80

# Frontend
NODE_ENV=development|production
ORIGIN=http://localhost:80
PUBLIC_API_URL=http://localhost:8080  # Dev only
```

---

## Code Execution Security Model

The code executor implements defense-in-depth:

| Layer | Protection | Implementation |
|-------|------------|----------------|
| **Network** | No external access | `--network none` |
| **Memory** | 128MB limit | `--memory 128m` |
| **CPU** | 50% single core | `--cpus 0.5` |
| **Time** | 10s timeout | Go context deadline |
| **Processes** | 64 process limit | `--pids-limit 64` |
| **Filesystem** | Read-only root | `--read-only` |
| **Temp Space** | 10MB tmpfs | `--tmpfs /tmp:size=10m` |
| **Privileges** | No escalation | `--security-opt no-new-privileges` |
| **Cleanup** | Auto-remove | `--rm` |

```go
// Timeout handling in code_executor.go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

if ctx.Err() == context.DeadlineExceeded {
    return RunCodeResponse{
        Error:    "Execution timed out after 10s",
        ExitCode: 124,  // Standard timeout code
    }
}
```

---

## Key File Reference

| Purpose | Location |
|---------|----------|
| Backend entry point | [core/main.go](core/main.go) |
| REST handlers | [core/course_handlers.go](core/course_handlers.go) |
| Code execution | [core/code_executor.go](core/code_executor.go) |
| Course loader | [core/courseparser/courseparser.go](core/courseparser/courseparser.go) |
| Problem types | [core/courseparser/types.go](core/courseparser/types.go) |
| Task parser | [core/courseparser/task.go](core/courseparser/task.go) |
| DB models | [core/models.go](core/models.go) |
| API DTOs | [core/types.go](core/types.go) |
| Frontend types | [front/src/lib/types/course.ts](front/src/lib/types/course.ts) |
| API client | [front/src/lib/services/api-client.ts](front/src/lib/services/api-client.ts) |
| Course service | [front/src/lib/services/course.service.ts](front/src/lib/services/course.service.ts) |
| Home page | [front/src/routes/+page.svelte](front/src/routes/+page.svelte) |
| Course page | [front/src/routes/courses/[id]/+page.svelte](front/src/routes/courses/[id]/+page.svelte) |
| Task page | [front/src/routes/courses/[id]/tasks/[taskId]/+page.svelte](front/src/routes/courses/[id]/tasks/[taskId]/+page.svelte) |
| Nginx config | [proxy/nginx.conf](proxy/nginx.conf) |
| Production compose | [docker-compose.yml](docker-compose.yml) |
| Development compose | [docker-compose.dev.yaml](docker-compose.dev.yaml) |

---

## Known Limitations & Technical Debt

| Area | Limitation | Impact |
|------|------------|--------|
| **Authentication** | No auth implemented | API endpoints are open |
| **Language Support** | Python only | Code executor hardcoded to Python |
| **Grading** | No test validation | Code runs but isn't graded |
| **Persistence** | No submission storage | User work not saved |
| **Rate Limiting** | None | DoS vulnerability on `/api/run` |
| **File Upload** | Not implemented | Can't submit multi-file solutions |
| **WebSocket** | Not used | No real-time collaboration |
| **Caching** | None | Course YAML parsed on every request |

---

## Future Roadmap Considerations

### Authentication & Authorization
- JWT or session-based authentication
- Role-based access control (RBAC) in API
- Course enrollment management
- OAuth2 integration (university SSO)

### Multi-Language Code Execution
```go
// Extend language support
languageConfigs := map[string]LanguageConfig{
    "python":     {Image: "python:3.14-slim", Cmd: "python"},
    "javascript": {Image: "node:22-alpine", Cmd: "node"},
    "java":       {Image: "openjdk:21-slim", Compile: "javac", Run: "java"},
    "rust":       {Image: "rust:1.75-slim", Compile: "rustc", Run: "./a.out"},
}
```

### Grading System
- Test case definitions in `task.yaml`
- Input/output validation
- Partial credit scoring
- Anti-cheating measures (code similarity)

### Submission Persistence
```sql
CREATE TABLE submissions (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    task_id VARCHAR(255),
    course_id VARCHAR(255),
    code TEXT,
    output TEXT,
    grade NUMERIC,
    submitted_at TIMESTAMP
);
```

### Performance Optimizations
- Course YAML caching with file watchers
- Container pool for faster cold starts
- API response caching
- Database connection pooling
