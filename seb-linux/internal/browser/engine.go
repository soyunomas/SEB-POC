package browser

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"seb-linux/internal/config"

	"github.com/chromedp/cdproto/browser"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/fetch"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func RunExamSession(ctx context.Context, cfg *config.SEBConfig, internalBEK string) error {
	// Perfil persistente en HOME — instala ahí tus plugins y se mantienen
	profileDir := filepath.Join(os.Getenv("HOME"), ".seb-linux-profile")
	os.MkdirAll(profileDir, 0755)

	opts := []chromedp.ExecAllocatorOption{
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.Flag("headless", false),
		chromedp.Flag("user-data-dir", profileDir),
		chromedp.Flag("start-maximized", true),
		chromedp.Flag("disable-background-networking", true),
		chromedp.Flag("disable-breakpad", true),
		chromedp.Flag("disable-component-update", true),
		chromedp.Flag("disable-default-apps", true),
		chromedp.Flag("disable-hang-monitor", true),
		chromedp.Flag("disable-popup-blocking", true),
		chromedp.Flag("disable-prompt-on-repost", true),
		chromedp.Flag("disable-sync", true),
		chromedp.Flag("enable-automation", false),
	}

	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	taskCtx, cancelTask := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancelTask()

	// JS: objeto SEB + forzar misma ventana + desbloquear copy-paste
	sebJS := fmt.Sprintf(`
		Object.defineProperty(window, 'SafeExamBrowser', {
			value: {
				version: '3.3.0',
				security: { browserExamKey: '%s', configKey: '%s' }
			}, writable: false, enumerable: true, configurable: false
		});

		// Forzar que todo se abra en la misma ventana (no popups)
		var _origOpen = window.open.bind(window);
		Object.defineProperty(window, 'open', {
			value: function(url, name, features) {
				if (url && url !== '' && url !== 'about:blank') {
					location.href = url;
				}
				return window;
			},
			writable: false, configurable: false
		});

		// Desbloquear copy-paste: capturar en fase capture y frenar propagación
		['copy','cut','paste','selectstart','contextmenu'].forEach(function(evt) {
			document.addEventListener(evt, function(e) { e.stopImmediatePropagation(); }, true);
		});

		// CSS: forzar selección de texto en todo
		var _sebStyle = document.createElement('style');
		_sebStyle.textContent = '* { user-select: text !important; -webkit-user-select: text !important; -moz-user-select: text !important; } [unselectable] { user-select: text !important; }';
		(document.head || document.documentElement).appendChild(_sebStyle);

		// Limpiar atributos inline que bloquean copia
		function _sebUnlock() {
			document.querySelectorAll('[oncopy],[onpaste],[onselectstart],[oncontextmenu],[unselectable]').forEach(function(el) {
				el.removeAttribute('oncopy');
				el.removeAttribute('onpaste');
				el.removeAttribute('onselectstart');
				el.removeAttribute('oncontextmenu');
				el.removeAttribute('unselectable');
			});
			// Quitar target _blank de formularios y enlaces
			document.querySelectorAll('form[target], a[target="_blank"]').forEach(function(el) {
				el.removeAttribute('target');
			});
		}
		new MutationObserver(_sebUnlock).observe(document.documentElement, {childList:true, subtree:true});
		document.addEventListener('DOMContentLoaded', _sebUnlock);
		_sebUnlock();
	`, internalBEK, cfg.ConfigKey)

	// Listener para inyectar headers SEB solo en documentos HTML
	chromedp.ListenTarget(taskCtx, func(ev interface{}) {
		switch ev := ev.(type) {
		case *fetch.EventRequestPaused:
			go func() {
				c := chromedp.FromContext(taskCtx)
				execCtx := cdp.WithExecutor(taskCtx, c.Target)

				cleanURL := ev.Request.URL
				if idx := strings.Index(cleanURL, "#"); idx != -1 {
					cleanURL = cleanURL[:idx]
				}

				hashBEK := sha256.Sum256([]byte(cleanURL + internalBEK))
				hashCK := sha256.Sum256([]byte(cleanURL + cfg.ConfigKey))

				var newHeaders []*fetch.HeaderEntry
				for k, v := range ev.Request.Headers {
					newHeaders = append(newHeaders, &fetch.HeaderEntry{Name: k, Value: fmt.Sprintf("%v", v)})
				}
				newHeaders = append(newHeaders,
					&fetch.HeaderEntry{Name: "X-SafeExamBrowser-RequestHash", Value: hex.EncodeToString(hashBEK[:])},
					&fetch.HeaderEntry{Name: "X-SafeExamBrowser-ConfigKeyHash", Value: hex.EncodeToString(hashCK[:])},
				)

				log.Printf("[HEADER] %s", cleanURL)

				if err := fetch.ContinueRequest(ev.RequestID).WithHeaders(newHeaders).Do(execCtx); err != nil {
					log.Printf("[HEADER-ERR] %v", err)
					fetch.ContinueRequest(ev.RequestID).Do(execCtx)
				}
			}()
		}
	})

	userAgent := "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 SafeExamBrowser/3.3.0"

	err := chromedp.Run(taskCtx,
		fetch.Enable().WithPatterns([]*fetch.RequestPattern{
			{RequestStage: fetch.RequestStageRequest, ResourceType: network.ResourceTypeDocument},
		}),
		emulation.SetUserAgentOverride(userAgent),
		// Otorgar permisos de clipboard
		browser.SetPermission(&browser.PermissionDescriptor{
			Name: "clipboard-read",
		}, browser.PermissionSettingGranted),
		browser.SetPermission(&browser.PermissionDescriptor{
			Name:                     "clipboard-write",
			AllowWithoutSanitization: true,
		}, browser.PermissionSettingGranted),
		chromedp.ActionFunc(func(ctx context.Context) error {
			_, err := page.AddScriptToEvaluateOnNewDocument(sebJS).Do(ctx)
			return err
		}),
	)
	if err != nil {
		return fmt.Errorf("fallo configurando CDP: %v", err)
	}

	log.Printf("[BROWSER] Navegando a: %s", cfg.StartURL)
	log.Println("[INFO] Perfil persistente en ~/.seb-linux-profile — instala plugins desde chrome://extensions")
	if err := chromedp.Run(taskCtx, chromedp.Navigate(cfg.StartURL)); err != nil {
		log.Printf("[BROWSER-WARNING] %v", err)
	}

	<-taskCtx.Done()
	return nil
}
