# Brutalist Todo App

A simple todo application built with Go, HTMX, and DaisyUI featuring a brutalist design theme.

## Stack

- **Go** - Backend server with embedded templates
- **HTMX** - Dynamic interactions without writing JavaScript
- **DaisyUI** - Tailwind CSS component library
- **Brutalist Theme** - Bold, raw design aesthetic

## Features

- Add, toggle, and delete todos
- Real-time updates with HTMX (no page reloads)
- Responsive brutalist design
- Single binary deployment

## Run Locally

```bash
go run main.go
```

Then open http://localhost:3000

## Deploy on Dalang.io

1. Push this repo to GitHub
2. Go to [Dalang.io Dashboard](https://dalang.io/dashboard/services/create/github)
3. Connect your GitHub account
4. Select this repository
5. Deploy!

The app includes a Dockerfile for automatic container builds.

## Environment Variables

- `PORT` - Server port (default: 3000)

<!-- push test 2026-10-29T06:30:26Z: checks that a push rebuilds the deployment -->

<!-- push test 2026-10-29T06:35:10Z: second push, hook now points at staging -->

<!-- push test 2026-10-29T06:52:54Z: third push, after the first build was cut off by an API restart -->
