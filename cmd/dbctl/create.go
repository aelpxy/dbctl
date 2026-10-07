package main

import (
	"cmp"
	"context"
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/aelpxy/dbctl/internal/database"
	"github.com/aelpxy/dbctl/internal/docker"
	"github.com/aelpxy/dbctl/internal/generate"
	"github.com/aelpxy/dbctl/internal/hostnet"
	"github.com/spf13/cobra"
)

const defaultWaitTimeout = 2 * time.Minute

type createOptions struct {
	output   *outputFormat
	password string
	name     string
	port     int
	timeout  time.Duration
	wait     bool
}

// TODO: add support for memory and cpu args
func (a *app) newCreateCmd() *cobra.Command {
	opts := createOptions{}

	cmd := &cobra.Command{
		Use:   "create <database-type>[:image-tag]",
		Short: "Create a new database",
		Long:  "Create a new database and wait until it is ready.\n\nSupported databases: " + strings.Join(a.databases.Names(), ", "),
		Example: `  dbctl create postgres
  dbctl create redis:latest
  dbctl create mysql --password mypassword --port 3306 --name mydb
  dbctl create pgvector -o json`,
		Aliases:   []string{"mk"},
		Args:      cobra.ExactArgs(1),
		ValidArgs: a.databases.Names(),
		RunE: a.withDocker(func(ctx context.Context, c *docker.Client, args []string) error {
			return a.runCreate(ctx, c, args[0], &opts)
		}),
	}

	cmd.Flags().StringVarP(&opts.password, "password", "P", "", "Specify a custom password for the database (default: generated password).")
	cmd.Flags().IntVarP(&opts.port, "port", "p", 0, "Specify a custom port for the database (default: free port).")
	cmd.Flags().StringVarP(&opts.name, "name", "n", "", "Specify a custom name for the database (default: generated name).")
	cmd.Flags().BoolVar(&opts.wait, "wait", true, "Wait until the database is ready to accept connections.")
	cmd.Flags().DurationVar(&opts.timeout, "timeout", defaultWaitTimeout, "How long to wait for the database to become ready.")
	opts.output = addOutputFlag(cmd)

	return cmd
}

func (a *app) runCreate(ctx context.Context, c *docker.Client, arg string, opts *createOptions) error {
	def, tag, err := a.databases.Parse(arg)
	if err != nil {
		return fmt.Errorf("parse database type: %w", err)
	}

	if !def.Supports(runtime.GOARCH) {
		return fmt.Errorf("%s does not support %s systems", def.Name(), runtime.GOARCH)
	}

	create, err := resolveCreate(ctx, def, def.Image(tag), opts)
	if err != nil {
		return err
	}

	if err := pullImage(ctx, c, create.Image); err != nil {
		return err
	}

	s := startSpinner("Creating " + create.Name + " database...")
	id, err := c.Create(ctx, create)

	s.Stop()

	if err != nil {
		return fmt.Errorf("create database: %w", err)
	}

	if opts.wait {
		if err := waitReady(ctx, c, id, opts.timeout); err != nil {
			return err
		}
	}

	return printCreated(id, create, opts.output)
}

func resolveCreate(ctx context.Context, def *database.Definition, image string, opts *createOptions) (*docker.CreateOptions, error) {
	create := &docker.CreateOptions{
		Definition: def,
		Image:      image,
		Name:       cmp.Or(opts.name, generate.Name()),
		Password:   cmp.Or(opts.password, generate.Password()),
		HostIP:     hostnet.OutboundIP(ctx).String(),
		HostPort:   opts.port,
	}

	if create.HostPort != 0 {
		return create, nil
	}

	port, err := hostnet.FreePort(ctx)
	if err != nil {
		return nil, fmt.Errorf("find a free port: %w", err)
	}

	create.HostPort = port

	return create, nil
}

func pullImage(ctx context.Context, c *docker.Client, image string) error {
	exists, err := c.ImageExists(ctx, image)
	if err != nil {
		return fmt.Errorf("check image: %w", err)
	}

	if exists {
		progressf("Image %s already exists, skipping pull\n", image)

		return nil
	}

	s := startSpinner("Pulling " + image + "...")
	defer s.Stop()

	if err := c.Pull(ctx, image); err != nil {
		return fmt.Errorf("pull image: %w", err)
	}

	s.Stop()
	progressf("Image %s pulled successfully\n", image)

	return nil
}

func waitReady(ctx context.Context, c *docker.Client, id string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	s := startSpinner("Waiting for the database to be ready...")
	err := c.WaitReady(ctx, id)

	s.Stop()

	if err != nil {
		return fmt.Errorf("wait for database: %w", err)
	}

	return nil
}

func printCreated(id string, opts *docker.CreateOptions, output *outputFormat) error {
	conn, err := newConnection(opts.Definition, docker.ContainerPrefix+opts.Name, opts.HostIP, opts.HostPort, opts.Password)
	if err != nil {
		return err
	}

	conn.ID = id
	conn.Image = opts.Image

	if output.isJSON() {
		return writeJSON(conn)
	}

	table := newTable("Key", "Value")
	table.AppendBulk([][]string{
		{"Container ID", shortID(conn.ID)},
		{"Name", conn.Name},
		{"Database Type", conn.Type},
		{"Image Tag", conn.Image},
		{"Port", strconv.Itoa(conn.Port)},
		{"Password", conn.Password},
	})
	table.Render()

	fmt.Printf("\nConnection String: %s\n", conn.URL)

	return nil
}
