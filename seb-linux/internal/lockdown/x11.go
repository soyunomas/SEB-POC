package lockdown

import (
	"context"
	"log"
	"runtime"

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

func StartX11Lockdown(ctx context.Context) error {
	c, err := xgb.NewConn()
	if err != nil {
		return err
	}

	setup := xproto.Setup(c)
	screen := setup.DefaultScreen(c)

	dummyWin, err := xproto.NewWindowId(c)
	if err != nil {
		return err
	}

	xproto.CreateWindow(c, screen.RootDepth, dummyWin, screen.Root,
		0, 0, 1, 1, 0, 0, 0, 0,[]uint32{})

	clipboardAtom := internAtom(c, "CLIPBOARD")
	primaryAtom := internAtom(c, "PRIMARY")

	// Reclamación inicial
	xproto.SetSelectionOwner(c, dummyWin, clipboardAtom, xproto.TimeCurrentTime)
	xproto.SetSelectionOwner(c, dummyWin, primaryAtom, xproto.TimeCurrentTime)

	log.Println("[LOCKDOWN-X11] Escudo Anti-Copy/Paste Activo. Portapapeles asegurado.")

	go func() {
		<-ctx.Done()
		// Al salir, liberamos el portapapeles asignando el dueño a "Ninguno" (0)
		xproto.SetSelectionOwner(c, 0, clipboardAtom, xproto.TimeCurrentTime)
		xproto.SetSelectionOwner(c, 0, primaryAtom, xproto.TimeCurrentTime)
		c.Close()
	}()

	go func() {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()

		for {
			ev, err := c.WaitForEvent()
			if err != nil || ev == nil {
				break
			}

			switch e := ev.(type) {
			case xproto.SelectionRequestEvent:
				// Alguien intenta LEER el portapapeles. Se lo denegamos.
				reply := xproto.SelectionNotifyEvent{
					Time:      e.Time,
					Requestor: e.Requestor,
					Selection: e.Selection,
					Target:    e.Target,
					Property:  0, // None
				}
				xproto.SendEvent(c, false, e.Requestor, xproto.EventMaskNoEvent, string(reply.Bytes()))
				// Log silenciado aquí para no saturar la terminal si un daemon insiste
				
			case xproto.SelectionClearEvent:
				// FÍSICA X11: Alguien (ej. el usuario haciendo Ctrl+C) acaba de robarnos 
				// la propiedad del portapapeles. ¡Lo robamos de vuelta en el acto!
				xproto.SetSelectionOwner(c, dummyWin, e.Selection, xproto.TimeCurrentTime)
				log.Println("[LOCKDOWN-X11] 🛡️ Intento de Copiar detectado y neutralizado (Reclamación forzada).")
			}
		}
	}()

	return nil
}

func internAtom(c *xgb.Conn, name string) xproto.Atom {
	reply, err := xproto.InternAtom(c, false, uint16(len(name)), name).Reply()
	if err != nil {
		return 0
	}
	return reply.Atom
}
