package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"charm.land/lipgloss/v2"
	"golang.org/x/term"
)

// confirmDelete makes the user type the database name, like deleting a production database should feel.
func confirmDelete(ctx context.Context, name string, removeVolumes bool) error {
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return fmt.Errorf("%w: stdin is not a terminal", errNotConfirmed)
	}

	what := "it and all of its data"
	if !removeVolumes {
		what = "it (its data volume is kept)"
	}

	printErr(fg(colorWarning).Render("!") + " This permanently deletes " + bold(name) + " and " + what + ".")

	answer, err := prompt(ctx, "  Type "+accent(name)+" to confirm: ")
	if err != nil {
		return err
	}

	if answer != name {
		return fmt.Errorf("%w: %q does not match %q", errNotConfirmed, answer, name)
	}

	return nil
}

// prompt reads one line from stdin, giving up when ctx is canceled (e.g. Ctrl-C).
func prompt(ctx context.Context, question string) (string, error) {
	_, _ = lipgloss.Fprint(os.Stderr, question)

	answers := make(chan string, 1)
	errs := make(chan error, 1)

	// a blocked stdin read cannot be interrupted, so on cancel the reader is abandoned as the process exits
	go func() {
		line, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			errs <- err

			return
		}

		answers <- strings.TrimSpace(line)
	}()

	select {
	case <-ctx.Done():
		printErr("")

		return "", fmt.Errorf("%w: %w", errNotConfirmed, ctx.Err())
	case err := <-errs:
		return "", fmt.Errorf("read confirmation: %w", err)
	case answer := <-answers:
		return answer, nil
	}
}
