package cndro

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"

	"golang.org/x/term"
)

// TicketsURL is the conference tickets page.
const TicketsURL = "https://cloudnativedays.ro/tickets"

// Hyperlink wraps linkText as an OSC 8 hyperlink to url (clickable in many modern terminals).
func Hyperlink(url, linkText string) string {
	return fmt.Sprintf("\x1b]8;;%s\x1b\\%s\x1b]8;;\x1b\\", url, linkText)
}

// WriteTicketsURLLine prints the tickets URL; uses a terminal hyperlink when stdout is a TTY
// and NO_COLOR is unset.
func WriteTicketsURLLine(w io.Writer) error {
	line := "\t👉 " + TicketsURL + " 👈"
	if f, ok := w.(*os.File); ok && os.Getenv("NO_COLOR") == "" {
		fd := int(f.Fd())
		if fd >= 0 && term.IsTerminal(fd) {
			line = "\t👉 " + Hyperlink(TicketsURL, TicketsURL) + " 👈"
		}
	}
	_, err := fmt.Fprintln(w, line)
	return err
}

// OpenTicketsURL opens TicketsURL in the system default browser.
func OpenTicketsURL() error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", TicketsURL)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", TicketsURL)
	default:
		cmd = exec.Command("xdg-open", TicketsURL)
	}
	return cmd.Start()
}
