package main

import (
	"bufio"
	"net"
	"os"
	"strings"

	"fmt"

	"sync"

	"github.com/bo0rsh201/borssh/common"
	"github.com/ryanuber/go-glob"
)

type initialSyncer struct {
	printer  *ProgressPrinter
	compiler compiler
}

func (s *initialSyncer) checkAndSync() error {

	masks := s.compiler.config.InitialSync
	if len(masks) < 1 {
		return nil
	}
	hosts, err := readKnownHosts()
	if err != nil {
		return err
	}

	localHash, err := getLocalHash()
	if err != nil {
		return err
	}

	numWorkers := 10
	syncCh := make(chan string, numWorkers)
	wg := sync.WaitGroup{}
	l := sync.Mutex{}
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				hostname, ok := <-syncCh
				if !ok {
					return
				}
				err := s.sync(hostname, localHash)
				if err != nil {
					// package "color" is not thread safe
					l.Lock()
					s.printer.yellow("\t%s: %s", hostname, err.Error())
					l.Unlock()
				}
			}
		}()
	}
	for _, h := range hosts {
		for _, m := range masks {
			if glob.Glob(m, h) {
				syncCh <- h
			}
		}
	}
	close(syncCh)
	wg.Wait()
	return nil
}

func (s *initialSyncer) sync(hostname string, localHash string) error {

	ex, err := common.NewExecutor(hostname)
	if err != nil {
		return err
	}

	// check
	match, err := ex.DoesVersionMatch(localHash, paths.GetHashPath(false))
	if err != nil {
		return err
	}
	if match {
		return nil
	}
	// rsync
	cmd, err := ex.Rsync(
		paths.GetBaseDir(true),
		paths.GetHomeDir(false)+"/",
		"--delete",
		"--copy-unsafe-links",
		"--exclude",
		common.COMPILED_HASH_FILE,
	)
	if err != nil {
		return err
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(
			"Failed to rsync base dir: '%s' output: '%s'", err.Error(),
			string(out),
		)
	}

	err = s.compiler.install(ex)
	if err != nil {
		return err
	}
	cmd, err = ex.Rsync(
		paths.GetHashPath(true),
		paths.GetHashPath(false),
	)
	if err != nil {
		return err
	}
	out, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(
			"Failed to rsync compiled hash: '%s' output: '%s'",
			err.Error(),
			string(out),
		)
	}
	return nil
}

func readKnownHosts() ([]string, error) {

	var res []string
	fp, err := os.Open(paths.GetKnownHostsPath(true))
	if err == os.ErrNotExist {
		return []string{}, nil
	}
	if err != nil {
		return res, err
	}
	s := bufio.NewScanner(fp)
	s.Split(bufio.ScanLines)
	for s.Scan() {
		line := strings.TrimLeft(s.Text(), "\t ")
		if strings.HasPrefix(line, "#") || strings.HasPrefix(line, "@") || len(line) == 0 {
			continue
		}
		end := strings.IndexAny(line, "\t ")
		if end <= 0 {
			continue
		}
		parts := strings.Split(line[:end], ",")
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if net.ParseIP(p) == nil {
				res = append(res, p)
				break
			}
		}
	}
	return res, s.Err()
}
