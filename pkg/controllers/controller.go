package controllers

import (
	"context"
	"fmt"
)

// Controller represents a basic controller interface
type Controller interface {
	Start(ctx context.Context) error
	Stop() error
}

// BaseController implements basic controller functionality
type BaseController struct {
	Name string
}

// NewBaseController creates a new base controller
func NewBaseController(name string) *BaseController {
	return &BaseController{
		Name: name,
	}
}

// Start starts the controller
func (c *BaseController) Start(ctx context.Context) error {
	fmt.Printf("Starting controller: %s\n", c.Name)
	return nil
}

// Stop stops the controller
func (c *BaseController) Stop() error {
	fmt.Printf("Stopping controller: %s\n", c.Name)
	return nil
}
