# SEB-Windows — Lecciones Aprendidas

## 1. syscall.SIGTERM no existe en Windows
- **Error**: Usar `syscall.SIGTERM` no compila en `GOOS=windows`.
- **Solución**: Reemplazar por `os.Kill` o solo usar `os.Interrupt` (Ctrl+C).
- **Regla**: Para señales cross-platform, usar solo constantes del paquete `os`, no `syscall`.

## 2. Variables de entorno de ruta difieren entre plataformas
- **Linux**: `$HOME` → `/home/user`
- **Windows**: `%APPDATA%` → `C:\Users\user\AppData\Roaming`
- **Regla**: Usar `os.Getenv("APPDATA")` en Windows, nunca `os.Getenv("HOME")`.

## 3. Dependencias de plataforma deben aislarse
- `jezek/xgb` (X11) y `unix.AF_NETLINK` solo compilan en Linux.
- **Regla**: Código específico de plataforma va en archivos con sufijo `_linux.go` / `_windows.go`, o se elimina si no aplica.

## 4. Siempre documentar en tasks/
- **Regla**: Actualizar `todo.md` y `lessons.md` tras cada tarea completada o corrección.
