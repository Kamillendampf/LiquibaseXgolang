package LiquibaseXgolang

import (
	"fmt"
	"log"
	"os/exec"
)

type Mode string

type Config struct {
	ChangelogFile string
	Username      string
	Password      string
	URL           string
	PathToCLI     string
}

func runCLI(cfg Config, args []string) error {
	if cfg.PathToCLI != "" {
		log.Print("Check for CLI Path Exists")
		_, err := exec.LookPath(cfg.PathToCLI)
		if err != nil {
			return fmt.Errorf("Liquibase CLI not found (%s): %v", cfg.PathToCLI, err)
		}
	}

	cmd := exec.Command(cfg.PathToCLI, args...)
	out, err := cmd.CombinedOutput()
	fmt.Println(string(out))
	return err
}

type Liquibase struct {
	cfg Config
}

func New(cfg Config) *Liquibase {
	if cfg.PathToCLI == "" {
		log.Print("No custom Value for Command - use default command 'liquibase' \n")
		cfg.PathToCLI = "liquibase"
	}

	return &Liquibase{cfg: cfg}
}

func (l *Liquibase) Update() error {
	args := []string{
		"--changeLogFile=" + l.cfg.ChangelogFile,
		"--url=" + l.cfg.URL,
		"--username=" + l.cfg.Username,
		"--password=" + l.cfg.Password,
		"update",
	}
	return runCLI(l.cfg, args)
}

func (l *Liquibase) Rollback(tag string) error {
	args := []string{
		"--changeLogFile=" + l.cfg.ChangelogFile,
		"--url=" + l.cfg.URL,
		"--username=" + l.cfg.Username,
		"--password=" + l.cfg.Password,
		"rollback", tag,
	}
	return runCLI(l.cfg, args)
}

func (l *Liquibase) Tag(name string) error {
	args := []string{
		"--changeLogFile=" + l.cfg.ChangelogFile,
		"--url=" + l.cfg.URL,
		"--username=" + l.cfg.Username,
		"--password=" + l.cfg.Password,
		"tag", name,
	}
	return runCLI(l.cfg, args)
}

func (l *Liquibase) Status() error {
	args := []string{
		"--changeLogFile=" + l.cfg.ChangelogFile,
		"--url=" + l.cfg.URL,
		"--username=" + l.cfg.Username,
		"--password=" + l.cfg.Password,
		"status",
	}
	return runCLI(l.cfg, args)
}

func (l *Liquibase) Validate() error {
	args := []string{
		"--changeLogFile=" + l.cfg.ChangelogFile,
		"--url=" + l.cfg.URL,
		"--username=" + l.cfg.Username,
		"--password=" + l.cfg.Password,
		"validate",
	}
	return runCLI(l.cfg, args)
}

func (l *Liquibase) ClearChecksums() error {
	args := []string{
		"--changeLogFile=" + l.cfg.ChangelogFile,
		"--url=" + l.cfg.URL,
		"--username=" + l.cfg.Username,
		"--password=" + l.cfg.Password,
		"clearCheckSums",
	}
	return runCLI(l.cfg, args)
}

func (l *Liquibase) ReleaseLocks() error {
	args := []string{
		"--changeLogFile=" + l.cfg.ChangelogFile,
		"--url=" + l.cfg.URL,
		"--username=" + l.cfg.Username,
		"--password=" + l.cfg.Password,
		"releaseLocks",
	}
	return runCLI(l.cfg, args)
}

func (l *Liquibase) History() error {
	args := []string{
		"--changeLogFile=" + l.cfg.ChangelogFile,
		"--url=" + l.cfg.URL,
		"--username=" + l.cfg.Username,
		"--password=" + l.cfg.Password,
		"history",
	}
	return runCLI(l.cfg, args)
}
