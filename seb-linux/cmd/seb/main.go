package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"seb-linux/internal/browser"
	"seb-linux/internal/config"
	"seb-linux/internal/crypto"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	log.SetPrefix("[SEB-CORE] ")

	sebFile := "config.seb"
	if len(os.Args) > 1 {
		sebFile = os.Args[1]
	}

	log.Printf("Iniciando SEB-Linux con: %s", sebFile)

	spoofedBEK := "1e9b2524a1b021966a337ab2881a6c42ddf510be13f56cf80bba3fc9fcb476eb"

	cfg, err := config.ParseSEBFile(sebFile)
	if err != nil {
		log.Fatalf("Fallo crítico leyendo .seb: %v", err)
	}

	calculatedConfigKey, err := crypto.DeriveConfigKey(sebFile)
	if err != nil {
		log.Fatalf("Fallo derivando ConfigKey matemático: %v", err)
	}
	cfg.ConfigKey = calculatedConfigKey
	log.Printf("[CRYPTO] ConfigKey derivado exitosamente: %s", cfg.ConfigKey)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		log.Println("Señal de apagado recibida. Limpiando entorno...")
		cancel()
	}()

	log.Println("Lanzando navegador (modo libre, sin lockdown)...")
	if err := browser.RunExamSession(ctx, cfg, spoofedBEK); err != nil {
		log.Fatalf("Navegador finalizó con error: %v", err)
	}
}
