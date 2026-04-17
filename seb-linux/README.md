# SEB-POC - Cliente para Linux

Esta es la implementación nativa para **Linux** del proyecto principal [SEB-POC](../README.md). Se trata de una prueba de concepto (PoC) desarrollada en **Go** que replica el comportamiento y las firmas de red de Safe Exam Browser (SEB) en sistemas basados en Linux.

Su objetivo principal es la investigación de los mecanismos de seguridad, específicamente la replicación del protocolo de *handshake* y la validación de integridad exigida por plataformas de e-learning como Moodle.

## 🚀 Características Principales

- **Replicación de Handshake**: Simulación del proceso de conexión seguro esperado por el LMS.
- **Validación de Integridad (Spoofing)**: Generación e inyección de los hashes criptográficos obligatorios del protocolo SEB:
  - `ConfigKey` (Verificación estática de la configuración del examen).
  - `BrowserExamKey` (Hash dinámico que valida la sesión y el entorno).
- **Control del Navegador**: Interacción con una instancia local de Chromium/Google Chrome para simular el entorno de evaluación bajo los parámetros requeridos.

## 📋 Requisitos Previos

Para ejecutar y compilar este cliente en tu máquina Linux, necesitarás:

- **Go 1.26** o superior.
- **Google Chrome o Chromium** instalado y accesible en tu `$PATH`.
- Una terminal estándar de Linux (Bash, Zsh, etc.).

## ⚙️ Instalación y Uso

1. **Clonar el repositorio y acceder al directorio:**
   ```bash
   git clone https://github.com/soyunomas/SEB-POC.git
   cd SEB-POC/seb-linux
   ```

2. **Descargar y sincronizar dependencias:**
   ```bash
   go mod tidy
   ```

3. **Ejecutar la Prueba de Concepto:**
   Puedes probar el código directamente en tu entorno de desarrollo:
   ```bash
   go run .
   ```
   O compilar el binario nativo para tu distribución Linux:
   ```bash
   go build -o seb-poc-linux
   ./seb-poc-linux
   ```
   *(Nota: Revisa el código fuente para añadir argumentos o URLs específicas del LMS si tu script de Go lo requiere mediante flags o variables de entorno).*

## 🏗️ Arquitectura Técnica en Linux

A diferencia de la versión de Windows (que puede lidiar con APIs propias del sistema para el modo kiosko), esta versión se enfoca puramente en la **capa de red y el protocolo de autenticación HTTP**. 

Durante su ejecución, la aplicación:
1. Inicia una sesión controlada del navegador.
2. Intercepta y manipula las cabeceras HTTP de salida.
3. Calcula e inyecta dinámicamente las cabeceras `x-safeexambrowser-configkey` y `x-safeexambrowser-requesthash` en las peticiones web, permitiendo superar con éxito los *checks* del LMS.

## ⚠️ Advertencia Ética y Legal

Al igual que el proyecto principal, este código se proporciona **exclusivamente con fines de investigación, ciberseguridad y educación**. 

No debe utilizarse para vulnerar sistemas académicos reales, realizar trampas o evadir restricciones de exámenes en entornos de producción. El uso indebido de esta herramienta para acceder sin autorización o falsificar telemetría académica es responsabilidad exclusiva del usuario.
