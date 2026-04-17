# SEB-POC (Safe Exam Browser Proof of Concept)

Este repositorio contiene implementaciones experimentales de un cliente para **Safe Exam Browser (SEB)**, diseñado como una prueba de concepto (PoC) para entender cómo los sistemas de gestión de aprendizaje (LMS) como Moodle validan la integridad de las sesiones de examen.

El proyecto se divide en dos implementaciones nativas independientes:

- **/seb-windows**: Cliente desarrollado en Go diseñado para entornos Windows, utilizando `chromedp` para el control del navegador y la inyección de headers criptográficos exigidos por el protocolo SEB.
- **/seb-linux**: Versión equivalente para sistemas operativos basados en Linux, enfocada en la replicación del protocolo de handshake y la validación de integridad.

## Objetivo del Proyecto
Este repositorio tiene fines **estrictamente educativos y de investigación**. El objetivo es analizar cómo los navegadores de examen aplican restricciones de seguridad y cómo se generan los hashes de configuración (`ConfigKey`) y de sesión (`BrowserExamKey`) para garantizar que el entorno de examen no ha sido alterado.

## Disclaimer (Descargo de Responsabilidad)
**ESTE SOFTWARE SE PROPORCIONA "TAL CUAL", SIN GARANTÍA DE NINGÚN TIPO.**

1. El autor no se hace responsable de cualquier uso indebido, daño, o consecuencia legal derivada del uso de este software.
2. Este proyecto **no es un producto oficial de Safe Exam Browser** ni está afiliado a Moodle Trust o ninguna institución educativa.
3. El uso de herramientas de evasión o manipulación en entornos de evaluación académica puede violar las políticas de tu institución. **Utiliza este código únicamente en entornos de prueba controlados y con fines de aprendizaje.**
4. Cualquier intento de usar este software para realizar trampas o acceder sin autorización a plataformas académicas es responsabilidad exclusiva del usuario.

## Requisitos Previos
- **Go 1.26+**
- Navegador Google Chrome o Chromium instalado.
- Entorno de desarrollo compatible (Git Bash/MSYS2 para Windows, terminal estándar para Linux).

## Estructura
```text
SEB-POC/
├── seb-windows/    # Implementación nativa para Windows
└── seb-linux/      # Implementación nativa para Linux
```

## Licencia
Este proyecto se libera bajo la licencia MIT. Consulta el archivo `LICENSE` (si aplica) para más detalles.
