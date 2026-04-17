package monitor

import (
	"log"
	"os"
	"strings"
)

// CheckHardwareVirtualization lee la tabla DMI del Kernel.
// Precepto #38 (Data Locality): Lectura directa del FileSystem en memoria (sysfs).
func CheckHardwareVirtualization() {
	vendorData, err := os.ReadFile("/sys/class/dmi/id/sys_vendor")
	if err == nil {
		// Normalizamos a minúsculas y limpiamos saltos de línea
		vendor := strings.ToLower(strings.TrimSpace(string(vendorData)))
		log.Printf("[HARDWARE-AUDIT] 🖥️  Fabricante del Sistema: %s", vendor)

		// Firmas conocidas de Hypervisores
		if strings.Contains(vendor, "qemu") || 
		   strings.Contains(vendor, "virtualbox") || 
		   strings.Contains(vendor, "innotek") || 
		   strings.Contains(vendor, "vmware") {
			log.Fatalf("🚨[ANTI-CHEAT] Ejecución en Máquina Virtual detectada (%s). Abortando examen por riesgo de fraude.", vendor)
		}
	} else {
		log.Printf("[HARDWARE-AUDIT] ⚠️ No se pudo leer la tabla DMI. Posible ofuscación.")
	}
}
