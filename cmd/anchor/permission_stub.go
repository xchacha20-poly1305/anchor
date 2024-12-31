//go:build !windows && !unix

package main

func checkPermission() error {
	return nil
}
