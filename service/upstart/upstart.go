// Copyright 2012, 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package upstart

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
	"text/template"
	"time"

	"github.com/juju/errors"
	"github.com/juju/utils"

	"github.com/juju/juju/service/common"
)

var servicesRe = regexp.MustCompile("^([a-zA-Z0-9-_:]+)\\.conf$")

// ListServices returns the name of all installed services on the
// local host.
func ListServices(initDir string) ([]string, error) {
	fis, err := ioutil.ReadDir(initDir)
	if err != nil {
		return nil, errors.Trace(err)
	}

	var services []string
	for _, fi := range fis {
		if groups := servicesRe.FindStringSubmatch(fi.Name()); len(groups) > 0 {
			services = append(services, groups[1])
		}
	}
	return services, nil
}

// ListCommand returns a command that will list the services on a host.
func ListCommand() string {
	// TODO(ericsnow) Do "ls /etc/init/*.conf" instead?
	return `sudo initctl list | awk '{print $1}' | sort | uniq`
}

var startedRE = regexp.MustCompile(`^.* start/running, process (\d+)\n$`)

// InitDir holds the default init directory name.
var InitDir = "/etc/init"

var InstallStartRetryAttempts = utils.AttemptStrategy{
	Total: 1 * time.Second,
	Delay: 250 * time.Millisecond,
}

// Service provides visibility into and control over an upstart service.
type Service struct {
	common.Service
}

func NewService(name string, conf common.Conf) *Service {
	if conf.InitDir == "" {
		conf.InitDir = InitDir
	}
	return &Service{
		Service: common.Service{
			Name: name,
			Conf: conf,
		},
	}
}

// Name implements service.Service.
func (s Service) Name() string {
	return s.Service.Name
}

// Conf implements service.Service.
func (s Service) Conf() common.Conf {
	return s.Service.Conf
}

// confPath returns the path to the service's configuration file.
func (s *Service) confPath() string {
	return path.Join(s.Service.Conf.InitDir, s.Service.Name+".conf")
}

// Validate returns an error if the service is not adequately defined.
func (s *Service) Validate() error {
	if err := s.Service.Validate(); err != nil {
		return errors.Trace(err)
	}

	if s.Service.Conf.InitDir == "" {
		return errors.New("missing InitDir")
	}

	if s.Service.Conf.Transient {
		if len(s.Service.Conf.Env) > 0 {
			return errors.NotSupportedf("Conf.Env (when transient)")
		}
		if len(s.Service.Conf.Limit) > 0 {
			return errors.NotSupportedf("Conf.Limit (when transient)")
		}
		if s.Service.Conf.Out != "" {
			return errors.NotSupportedf("Conf.Out (when transient)")
		}
		if s.Service.Conf.ExtraScript != "" {
			return errors.NotSupportedf("Conf.ExtraScript (when transient)")
		}
	} else {
		if s.Service.Conf.AfterStopped != "" {
			return errors.NotSupportedf("Conf.AfterStopped (when not transient)")
		}
		if s.Service.Conf.ExecStopPost != "" {
			return errors.NotSupportedf("Conf.ExecStopPost (when not transient)")
		}
	}

	return nil
}

// render returns the upstart configuration for the service as a slice of bytes.
func (s *Service) render() ([]byte, error) {
	if err := s.Validate(); err != nil {
		return nil, err
	}
	conf := s.Conf()
	if conf.Transient {
		conf.ExecStopPost = "rm " + s.confPath()
	}
	return Serialize(s.Name(), conf)
}

// Installed returns whether the service configuration exists in the
// init directory.
func (s *Service) Installed() bool {
	_, err := os.Stat(s.confPath())
	return err == nil
}

// Exists returns whether the service configuration exists in the
// init directory with the same content that this Service would have
// if installed.
func (s *Service) Exists() bool {
	// In any error case, we just say it doesn't exist with this configuration.
	// Subsequent calls into the Service will give the caller more useful errors.
	_, same, _, err := s.existsAndSame()
	if err != nil {
		return false
	}
	return same
}

func (s *Service) existsAndSame() (exists, same bool, conf []byte, err error) {
	expected, err := s.render()
	if err != nil {
		return false, false, nil, errors.Trace(err)
	}
	current, err := ioutil.ReadFile(s.confPath())
	if err != nil {
		if os.IsNotExist(err) {
			// no existing config
			return false, false, expected, nil
		}
		return false, false, nil, errors.Trace(err)
	}
	return true, bytes.Equal(current, expected), expected, nil
}

// Running returns true if the Service appears to be running.
func (s *Service) Running() bool {
	cmd := exec.Command("status", "--system", s.Service.Name)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}
	return startedRE.Match(out)
}

// Start starts the service.
func (s *Service) Start() error {
	if s.Running() {
		return nil
	}
	err := runCommand("start", "--system", s.Service.Name)
	if err != nil {
		// Double check to see if we were started before our command ran.
		if s.Running() {
			return nil
		}
	}
	return err
}

func runCommand(args ...string) error {
	out, err := exec.Command(args[0], args[1:]...).CombinedOutput()
	if err == nil {
		return nil
	}
	out = bytes.TrimSpace(out)
	if len(out) > 0 {
		return fmt.Errorf("exec %q: %v (%s)", args, err, out)
	}
	return fmt.Errorf("exec %q: %v", args, err)
}

// Stop stops the service.
func (s *Service) Stop() error {
	if !s.Running() {
		return nil
	}
	return runCommand("stop", "--system", s.Service.Name)
}

// StopAndRemove stops the service and then deletes the service
// configuration from the init directory.
func (s *Service) StopAndRemove() error {
	if !s.Installed() {
		return nil
	}
	if err := s.Stop(); err != nil {
		return err
	}
	return os.Remove(s.confPath())
}

// Remove deletes the service configuration from the init directory.
func (s *Service) Remove() error {
	if !s.Installed() {
		return nil
	}
	return os.Remove(s.confPath())
}

// Install installs and starts the service.
func (s *Service) Install() error {
	exists, same, conf, err := s.existsAndSame()
	if err != nil {
		return errors.Trace(err)
	}
	if same {
		return nil
	}
	if exists {
		if err := s.StopAndRemove(); err != nil {
			return errors.Annotate(err, "upstart: could not remove installed service")
		}

	}
	if err := ioutil.WriteFile(s.confPath(), conf, 0644); err != nil {
		return errors.Trace(err)
	}

	// On slower disks, upstart may take a short time to realise
	// that there is a service there.
	for attempt := InstallStartRetryAttempts.Start(); attempt.Next(); {
		if err = s.Start(); err == nil {
			break
		}
	}
	return err
}

// InstallCommands returns shell commands to install and start the service.
func (s *Service) InstallCommands() ([]string, error) {
	conf, err := s.render()
	if err != nil {
		return nil, err
	}
	return []string{
		fmt.Sprintf("cat >> %s << 'EOF'\n%sEOF\n", s.confPath(), conf),
		"start " + s.Service.Name,
	}, nil
}

// Serialize renders the conf as raw bytes.
func Serialize(name string, conf common.Conf) ([]byte, error) {
	var buf bytes.Buffer
	if conf.Transient {
		if err := transientConfT.Execute(&buf, conf); err != nil {
			return nil, err
		}
	} else {
		if err := confT.Execute(&buf, conf); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

// TODO(ericsnow) Use a different solution than templates?

// BUG: %q quoting does not necessarily match libnih quoting rules
// (as used by upstart); this may become an issue in the future.
var confT = template.Must(template.New("").Parse(`
description "{{.Desc}}"
author "Juju Team <juju@lists.ubuntu.com>"
start on runlevel [2345]
stop on runlevel [!2345]
respawn
normal exit 0
{{range $k, $v := .Env}}env {{$k}}={{$v|printf "%q"}}
{{end}}
{{range $k, $v := .Limit}}limit {{$k}} {{$v}}
{{end}}
script
{{if .ExtraScript}}{{.ExtraScript}}{{end}}
{{if .Out}}
  # Ensure log files are properly protected
  touch {{.Out}}
  chown syslog:syslog {{.Out}}
  chmod 0600 {{.Out}}
{{end}}
  exec {{.ExecStart}}{{if .Out}} >> {{.Out}} 2>&1{{end}}
end script
`[1:]))

var transientConfT = template.Must(template.New("").Parse(`
description "{{.Desc}}"
author "Juju Team <juju@lists.ubuntu.com>"
start on stopped {{.AfterStopped}}
script
  {{.ExecStart}}
end script
{{if .ExecStopPost}}
post-stop script
  {{.ExecStopPost}}
end script
{{end}}
`[1:]))
