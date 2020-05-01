package main

import (
	"archive/zip"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/pkg/errors"
)

// terraformVersion is a constant of the version of the terraform binary
// to download and execute. We want to control this versioning to ensure compatibility
// with modules, providers, etc
const terraformVersion = "0.12.24"

func main() {
	// flag for apply vs. destroy
	var destroy bool
	// path for terraform binary
	var tfPath string
	flag.StringVar(&tfPath, "tfPath", "", "path for terraform binary")
	flag.BoolVar(&destroy, "destroy", false, "destroy")
	flag.Parse()

	if !terraformInstalled(tfPath) {
		log.Println("Installing terraform")
		if err := installTerraform(tfPath); err != nil {
			log.Fatalln("Unable to install terraform", err)
		}
	}

	// terraform init
	if err := execute("terraform", "init", "-input=false"); err != nil {
		log.Fatalln("Failed to terraform init", err)
	}

	// terraform apply/destroy
	action := "apply"
	if destroy {
		action = "destroy"
	}

	if err := execute("terraform", action, "-input=false", "-auto-approve"); err != nil {
		log.Fatalln("Failed to terraform apply/destroy", err)
	}

	// finished
	if destroy {
		log.Println("Changes destroyed successfully")
	} else {
		log.Println("Changes applied successfully")
	}
}

// Checks to see if terraform already exists at path
// Note: at this point assuming if terraform already exists, that it is the
// correct version, os, arch. User may have previously installed a version
// that we don't support. May want to add handling in the future.
func terraformInstalled(tfPath string) bool {
	path, err := exec.LookPath("terraform")
	if err != nil {
		return false
	}

	tfPath = filepath.Join(tfPath, "terraform")
	if tfPath == path {
		return true
	}

	// have terraform at a different path
	return false
}

// Install terraform: download file, unzip binary into path
func installTerraform(tfPath string) error {
	// info needed for binary
	opsys := runtime.GOOS
	arch := runtime.GOARCH

	filename := fmt.Sprintf("terraform_%s_%s_%s.zip", terraformVersion, opsys, arch)
	url := fmt.Sprintf("https://releases.hashicorp.com/terraform/%s/%s", terraformVersion, filename)

	if err := download(url, filename); err != nil {
		return errors.Wrap(err, "Unable to download zip")
	}

	if err := unzip(filename, tfPath); err != nil {
		return errors.Wrap(err, "Unable to unzip binary")
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
// Logs are standard-out-ed. Errors are captured and wrapped in an error object
func execute(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		return errors.Wrap(err, "Failed cmd.Start()")
	}

	// cmd.Wait() should be called only finishing reading from stdoutIn and stderrIn.
	// wg ensures that we finish
	var errBuf bytes.Buffer
	var errStdout, errErrBuf error

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		errStdout = capture(os.Stdout, stdoutIn)
		wg.Done()
	}()

	errErrBuf = capture(&errBuf, stderrIn)
	wg.Wait()

	if err := cmd.Wait(); err != nil {
		// see if we captured any error info
		if errBuf.Len() > 0 {
			retError := errors.New(errBuf.String())
			err = errors.Wrap(retError, err.Error())
		}
		return errors.Wrap(err, "Failed cmd.Wait()")
	}

	if errStdout != nil || errErrBuf != nil {
		return errors.New("Failed to capture stdout or error buffer")
	}

	if errBuf.Len() > 0 {
		return errors.New(errBuf.String())
	}

	return nil
}

// https://blog.kowalczyk.info/article/wOYk/advanced-command-execution-in-go-with-osexec.html
// With some modifications
func capture(w io.Writer, r io.Reader) error {
	var out []byte
	buf := make([]byte, 1024, 1024)
	for {
		n, err := r.Read(buf[:])
		if n > 0 {
			d := buf[:n]
			out = append(out, d...)
			if _, err := w.Write(d); err != nil {
				return err
			}
		}
		if err != nil {
			// Read returns io.EOF at the end of file, which is not an error for us
			if err == io.EOF {
				err = nil
			}
			return err
		}
	}
}
