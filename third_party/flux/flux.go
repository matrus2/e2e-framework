package flux

import (
	"fmt"
	"github.com/vladimirvivien/gexe"
	log "k8s.io/klog/v2"
	"strings"
)

type Opts struct {
	namespace string
	mode      string
	branch    string
	tag       string
	args      []string
}

type Manager struct {
	e          *gexe.Echo
	kubeConfig string
}

type Option func(*Opts)

func (m *Manager) processOpts(opts ...Option) *Opts {
	option := &Opts{}
	for _, op := range opts {
		op(option)
	}
	return option
}

func (m *Manager) InstallFlux(opts ...Option) error {
	o := m.processOpts(opts...)
	o.mode = "install"
	return m.run(o)
}

func (m *Manager) UninstallFlux(opts ...Option) error {
	o := m.processOpts(opts...)
	o.mode = "uninstall"
	return m.run(o)
}

func WithNamespace(namespace string) Option {
	return func(opts *Opts) {
		opts.namespace = namespace
	}
}

func WithTag(namespace string) Option {
	return func(opts *Opts) {
		opts.tag = namespace
	}
}

func WithBranch(namespace string) Option {
	return func(opts *Opts) {
		opts.branch = namespace
	}
}

func WithArgs(args []string) Option {
	return func(opts *Opts) {
		opts.args = args
	}
}

func (m *Manager) run(opts *Opts) (err error) {
	if m.e.Prog().Avail("flux") == "" {
		err = fmt.Errorf("'flux' command is missing. Please ensure the tool exists before using the flux manager")
		return
	}
	command, err := m.getCommand(opts)
	if err != nil {
		return
	}
	log.V(4).InfoS("Running Flux Operation", "command", command)
	proc := m.e.RunProc(command)
	result := proc.Result()
	log.V(4).Info("Flux Command output \n", result)
	if proc.IsSuccess() {
		return nil
	} else {
		return proc.Err()
	}
}

func New(kubeConfig string) *Manager {
	return &Manager{e: gexe.New(), kubeConfig: kubeConfig}
}

func (m *Manager) getCommand(opt *Opts) (string, error) {
	commandParts := []string{"flux", opt.mode}
	if opt.namespace != "" {
		commandParts = append(commandParts, "--namespace", opt.namespace)
	}
	if opt.branch != "" {
		commandParts = append(commandParts, "--branch", opt.branch)
	}
	if opt.tag != "" {
		commandParts = append(commandParts, "--tag", opt.tag)
	}

	commandParts = append(commandParts, opt.args...)
	commandParts = append(commandParts, "--kubeconfig", m.kubeConfig)
	return strings.Join(commandParts, " "), nil
}

func (m *Manager) CreateGitRepo(name string, url string, opts ...Option) error {
	o := m.processOpts(opts...)
	command := []string{"create source git"}

	if name == "" {
		return fmt.Errorf("GitRepository name is required")
	}
	if url == "" {
		return fmt.Errorf("GitRepository url is required")
	}

	command = append(command, name, "--url", url)

	if o.branch == "" || o.tag == "" {
		return fmt.Errorf("GitRepository tag or branch is required")
	}

	return m.run(o)
}

func (m *Manager) CreateHelmRepo(name string, url string, opts ...Option) error {
	o := m.processOpts(opts...)
	command := []string{"create source helm"}

	if name == "" {
		return fmt.Errorf("GitRepository name is required")
	}
	if url == "" {
		return fmt.Errorf("GitRepository url is required")
	}

	command = append(command, name, "--url", url)

	if o.branch == "" || o.tag == "" {
		return fmt.Errorf("GitRepository tag or branch is required")
	}

	return m.run(o)
}

func (m *Manager) CreateKustomization(name string, source string, opts ...Option) error {
	o := m.processOpts(opts...)
	command := []string{"create ks"}

	if name == "" {
		return fmt.Errorf("Kustomization name is required")
	}
	if source == "" {
		return fmt.Errorf("Kustomization source is required")
	}

	command = append(command, name, "--source", source)

	if o.branch == "" || o.tag == "" {
		return fmt.Errorf("GitRepository tag or branch is required")
	}

	return m.run(o)
}

func (m *Manager) CreateHelmRelease(name string, source string, chartName string, opts ...Option) error {
	o := m.processOpts(opts...)
	command := []string{"create hr"}

	if name == "" {
		return fmt.Errorf("Kustomization name is required")
	}
	if source == "" {
		return fmt.Errorf("Kustomization source is required")
	}

	command = append(command, name, "--source", source)

	if o.branch == "" || o.tag == "" {
		return fmt.Errorf("GitRepository tag or branch is required")
	}

	return m.run(o)
}
