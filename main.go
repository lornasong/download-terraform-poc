package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/pkg/errors"
)

func main() {
	// parse info needed for binary
	var version, opsys, arch, tfPath string
	if err := parse(&version, &opsys, &arch, &tfPath); err != nil {
		log.Fatalln("Unable to parse binary information", err)
	}

	// download and unzip binary
	filename := fmt.Sprintf("terraform_%s_%s_%s.zip", version, opsys, arch)
	url := fmt.Sprintf("https://releases.hashicorp.com/terraform/%s/%s", version, filename)
	download(url, filename)

	if err := unzip(filename, tfPath); err != nil {
		log.Fatalln("Unable to unzip binary", err)
	}

	// terraform init and apply
	if err := execute("terraform", "init"); err != nil {
		log.Fatalln("Failed to terraform init", err)
	}

	if err := execute("terraform", "apply", "-auto-approve=true"); err != nil {
		log.Fatalln("Failed to terraform apply", err)
	}

	log.Println("Changes applied successfully")
}

func parse(version, opsys, arch, tfPath *string) error {
	flag.StringVar(version, "tfv", "", "terraform version")
	flag.StringVar(opsys, "os", "", "operating system")
	flag.StringVar(arch, "arch", "", "architecture")
	flag.StringVar(tfPath, "tfPath", "", "path for terraform binary")
	flag.Parse()

	switch *opsys {
	case "darwin":
	case "freebsd":
	case "linux":
	case "openbsd":
	case "solaris":
	case "windows":
	default:
		return fmt.Errorf("Unknown operating system %s", *opsys)
	}

	switch *arch {
	case "amd64":
	case "386":
	case "arm":
	default:
		return fmt.Errorf("Unknown architecture %s", *arch)
	}

	return nil
}

func download(url, filename string) error {
	log.Printf("Downloading %s from %s\n", filename, url)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, resp.Body); err != nil {
		return err
	}
	return nil
}

// unzip unzips the terraform binary. Assumes only one file inside the zip file.
// Simplied version of https://stackoverflow.com/a/24792688 for the sake of the poc
func unzip(src, dest string) error {
	log.Printf("Unziping %s to %s\n", src, dest)

	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		path := filepath.Join(dest, f.Name)
		fc, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer fc.Close()

		_, err = io.Copy(fc, rc)
		if err != nil {
			return err
		}
	}
	return nil
}

// execute executes command and logs out to console. Reworked from:
// https://blog.kowalczyk.info/article/wOYk/advanced-command-execution-in-go-with-osexec.html
// We have to do a little extra in order to stream logs in realish time vs. get a dump after execution.
// We also have to do a little extra to capture error on why an execution failed.
func execute(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		return errors.Wrap(err, "Failed cmd.Start()")
	}

	// cmd.Wait() should be called only finishing reading from stdoutIn and stderrIn.
	// wg ensures that we finish
	var stdout, stderr []byte
	var errStdout, errStderr error

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		stdout, errStdout = capture(os.Stdout, stdoutIn)
		wg.Done()
	}()
	stderr, errStderr = capture(os.Stderr, stderrIn)
	wg.Wait()

	if err := cmd.Wait(); err != nil {
		return errors.Wrap(err, "Failed cmd.Wait()")
	}

	if errStdout != nil || errStderr != nil {
		return errors.New("Failed to capture stdout or stderr")
	}

	log.Println(string(stdout))

	if len(stderr) > 0 {
		log.Println(string(stderr))
	}

	return nil
}

// https://blog.kowalczyk.info/article/wOYk/advanced-command-execution-in-go-with-osexec.html
func capture(w io.Writer, r io.Reader) ([]byte, error) {
	var out []byte
	buf := make([]byte, 1024, 1024)
	for {
		n, err := r.Read(buf[:])
		if n > 0 {
			d := buf[:n]
			out = append(out, d...)
			_, err := w.Write(d)
			if err != nil {
				return out, err
			}
		}
		if err != nil {
			// Read returns io.EOF at the end of file, which is not an error for us
			if err == io.EOF {
				err = nil
			}
			return out, err
		}
	}
}
