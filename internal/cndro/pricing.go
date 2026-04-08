package cndro

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// WritePricing prints the pricing table to w.
func WritePricing(w io.Writer) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "CATEGORY\tMar 15, 2026 - May 3, 2026\tMay 4, 2026 - EVENT")
	fmt.Fprintln(tw, "\tSTANDARD\tLATE")
	fmt.Fprintln(tw, "CORPORATE\t199€\t239€")
	fmt.Fprintln(tw, "INDIVIDUAL\t99€\t119€")
	fmt.Fprintln(tw, "ACADEMIC\t10€\t10€")
	fmt.Fprintln(tw, "\tFULL EVENT\tFULL EVENT")
	return tw.Flush()
}
