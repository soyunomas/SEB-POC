# SEB-Windows

Navegador que se identifica como Safe Exam Browser ante Moodle y otras plataformas LMS. Inyecta los headers criptográficos (`X-SafeExamBrowser-ConfigKeyHash`, `X-SafeExamBrowser-RequestHash`) y el User-Agent necesarios para que el servidor acepte la sesión.

## Requisitos

- Go 1.26+
- Google Chrome o Chromium instalado en el sistema

## Compilar

```bash
go build -ldflags="-s -w" -o bin/seb-windows.exe ./cmd/seb
```

O con Make (requiere Git Bash / MSYS2):

```bash
make build
```

## Uso

```bash
seb-windows.exe config.seb
```

Si no se indica archivo, usa `config.seb` del directorio actual.

El navegador abrirá Chrome con la URL definida en `startURL` del archivo `.seb`, inyectando automáticamente los headers SEB en cada petición.

## Archivo .seb

Es un XML Plist de Apple con la configuración del examen. Se obtiene desde Moodle (enlace `sebs://` o descarga directa). Ejemplo mínimo:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>startURL</key><string>https://tu-moodle.com/mod/quiz/view.php?id=12345</string>
    <key>sendBrowserExamKey</key><true/>
</dict>
</plist>
```

## Perfil de Chrome

Se guarda en `%APPDATA%\seb-windows-profile`. Las extensiones instaladas manualmente desde `chrome://extensions` persisten entre sesiones.
