package main

import "errors"

func (c *commands) register(name string, f func(*state, command) error) error {
	if name == "" {
		return errors.New("Missing name parameter")
	}

	c.handlersByName[name] = f

	return nil
}

func (c *commands) run(s *state, cmd command) error {
	funcToRun := c.handlersByName[cmd.name]
	if funcToRun == nil {
		return errors.New("function not found")
	}
	return funcToRun(s, cmd)
}
