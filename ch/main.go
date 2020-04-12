package main

import "github.com/sirupsen/logrus"

func main() {
	defer func() {
		if e := recover(); e != nil {
			err := e.(error)
			logrus.Error(err)
		}
	}()
	c := make(chan struct{})
	close(c)
	close(c)
}
