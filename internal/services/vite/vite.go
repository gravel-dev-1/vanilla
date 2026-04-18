package vite

import (
	"context"
	"embed"
	"io/fs"
	"os"
	"os/exec"

	"gbfw/internal/services/env"
)

//go:embed build/*
var productionFS embed.FS

const ServiceName = "Vite"

type Service struct {
	FS fs.FS
}

func (s *Service) Start(context.Context) (err error) {
	s.FS, err = fs.Sub(productionFS, "build")
	if env.IsDev {
		cmd := exec.Command(string(env.Get("JS_RUNTIME", "node")), "node_modules/vite/bin/vite", "--host")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Start()
		s.FS = FS([]fs.FS{os.DirFS("internal/services/vite/dev"), os.DirFS("public")})
	}
	return err
}

func (s *Service) String() string                        { return ServiceName }
func (s *Service) State(context.Context) (string, error) { return "", nil }
func (s *Service) Terminate(context.Context) error       { return nil }

type FS []fs.FS

func (f FS) Open(name string) (file fs.File, err error) {
	for _, filesystem := range f {
		if file, err = filesystem.Open(name); file != nil {
			return file, err
		}
	}
	return nil, fs.ErrNotExist
}
