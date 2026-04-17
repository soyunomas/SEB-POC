# SEB-Windows — TODO

## Fase 0: Fundación ✅
- [x] Estructura Go, Makefile, `go.mod`

## Fase 1: Criptografía y Parseo ✅
- [x] Parser de `.seb` (XML Plist)
- [x] Derivación del Config Key (JSON canónico, SHA-256)
- [x] Cálculo del BEK (SHA-256 del ejecutable)

## Fase 2: Navegador + Headers ✅
- [x] Chrome vía CDP (chromedp)
- [x] Inyección de headers `X-SafeExamBrowser-ConfigKeyHash` y `X-SafeExamBrowser-RequestHash`
- [x] User-Agent con `SafeExamBrowser/3.3.0`
- [x] Objeto JS `window.SafeExamBrowser` inyectado en cada página
- [x] Desbloqueo de copy-paste y selección de texto

## Migración Linux → Windows ✅
- [x] Módulo renombrado `seb-linux` → `seb-windows`
- [x] Eliminadas dependencias Linux: `jezek/xgb` (X11), Netlink (`golang.org/x/sys` directo)
- [x] User-Agent cambiado a `Windows NT 10.0; Win64; x64`
- [x] Perfil Chrome: `$HOME/.seb-linux-profile` → `%APPDATA%\seb-windows-profile`
- [x] Señal de apagado: `syscall.SIGTERM` → `os.Kill` (compatible Windows)
- [x] Lockdown X11 eliminado → stub vacío (sin restricciones)
- [x] Monitor Netlink + hardware DMI eliminados → stub vacío (sin restricciones)
- [x] Makefile genera `seb-windows.exe`
- [x] Extensión Chrome adaptada (nombre + UA Windows)
- [x] Compilación cruzada verificada: `GOOS=windows GOARCH=amd64` → PE32+ OK
