package monitor

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"

	"golang.org/x/sys/unix"
)

// Constantes físicas del Kernel de Linux para Netlink Proc Connector
const (
	CN_IDX_PROC          = 1
	CN_VAL_PROC          = 1
	PROC_CN_MCAST_LISTEN = 1
	PROC_EVENT_EXEC      = 0x00000002
)

var forbiddenApps = [][]byte{[]byte("discord"),
	[]byte("obs"),
	[]byte("teamviewer"),
	[]byte("anydesk"),[]byte("skype"),[]byte("zoom"),
}

// StartProcessMonitor abre un socket Netlink para escuchar ejecuciones de procesos.
func StartProcessMonitor(ctx context.Context) error {
	// 1. Abrir Socket Netlink crudo
	fd, err := unix.Socket(unix.AF_NETLINK, unix.SOCK_DGRAM, unix.NETLINK_CONNECTOR)
	if err != nil {
		return fmt.Errorf("requiere permisos Root (sudo) para escuchar al Kernel: %v", err)
	}

	addr := &unix.SockaddrNetlink{
		Family: unix.AF_NETLINK,
		Groups: CN_IDX_PROC,
	}
	if err := unix.Bind(fd, addr); err != nil {
		unix.Close(fd)
		return fmt.Errorf("fallo al hacer bind en Netlink: %v", err)
	}

	// 2. Armar y enviar el comando PROC_CN_MCAST_LISTEN (Estructuras en binario puro)
	if err := sendListenCommand(fd); err != nil {
		unix.Close(fd)
		return fmt.Errorf("fallo enviando orden de escucha al Kernel: %v", err)
	}

	log.Println("[KERNEL-MONITOR] Enlace Netlink establecido. Escuchando Ring-0...")

	go func() {
		// FÍSICA: Bloquear al hilo del SO para ingesta de eventos de baja latencia
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		defer unix.Close(fd)

		// PRECEPTO #2 y #10: Buffers preasignados fuera del bucle. ZERO ALLOC.
		nlBuf := make([]byte, 4096)
		cmdBuf := make([]byte, 512)

		for {
			select {
			case <-ctx.Done():
				return
			default:
				// Syscall bloqueante. Consume 0 CPU mientras espera.
				n, _, err := unix.Recvfrom(fd, nlBuf, 0)
				if err != nil || n < 60 {
					continue
				}

				// Offset 36 es donde está el What (tipo de evento) del struct proc_event
				what := binary.LittleEndian.Uint32(nlBuf[36:40])
				if what == PROC_EVENT_EXEC {
					// En un evento EXEC, el TGID (Process ID) está en el offset 56
					tgid := binary.LittleEndian.Uint32(nlBuf[56:60])
					checkProcess(tgid, cmdBuf)
				}
			}
		}
	}()

	return nil
}

func sendListenCommand(fd int) error {
	buf := make([]byte, 40) // NLMSGHDR (16) + CN_MSG (20) + OP (4)

	// Netlink Header
	binary.LittleEndian.PutUint32(buf[0:4], 40)
	binary.LittleEndian.PutUint16(buf[4:6], unix.NLMSG_DONE)
	binary.LittleEndian.PutUint32(buf[12:16], uint32(os.Getpid()))

	// Connector Header
	binary.LittleEndian.PutUint32(buf[16:20], CN_IDX_PROC)
	binary.LittleEndian.PutUint32(buf[20:24], CN_VAL_PROC)
	binary.LittleEndian.PutUint16(buf[32:34], 4)

	// Op (PROC_CN_MCAST_LISTEN)
	binary.LittleEndian.PutUint32(buf[36:40], PROC_CN_MCAST_LISTEN)

	addr := &unix.SockaddrNetlink{Family: unix.AF_NETLINK}
	return unix.Sendto(fd, buf, 0, addr)
}

func checkProcess(pid uint32, cmdBuf[]byte) {
	// Acceso rápido a /proc. Sin allocations complejas.
	path := "/proc/" + strconv.FormatUint(uint64(pid), 10) + "/cmdline"
	
	f, err := os.Open(path)
	if err != nil {
		return // El proceso pudo morir muy rápido
	}
	defer f.Close()

	n, _ := f.Read(cmdBuf)
	if n == 0 {
		return
	}

	// Reemplazar bytes nulos por espacios para evaluación
	cleanCmd := bytes.ReplaceAll(cmdBuf[:n], []byte{0},[]byte{' '})
	cleanCmd = bytes.ToLower(cleanCmd)

	for _, app := range forbiddenApps {
		// Zero-Alloc Contains
		if bytes.Contains(cleanCmd, app) {
			log.Printf("[KERNEL-MONITOR] 🚨 ¡INFRACCIÓN CRÍTICA! Proceso prohibido detectado: %s", app)
			log.Printf("   -> CMD: %s", cleanCmd)
			
			// ACCIÓN PUNITIVA: Matar el proceso infractor instantáneamente
			unix.Kill(int(pid), unix.SIGKILL)
			log.Printf("[KERNEL-MONITOR] 💀 Proceso %d exterminado con SIGKILL.", pid)
		}
	}
}
