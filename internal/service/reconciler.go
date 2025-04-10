package service

import (
	"context"
	"log"
	"time"
)

type Reconciler struct {
	manager       Manager
	interval      time.Duration
	stopCh        chan struct{}
	isRunning     bool
}

func NewReconciler(manager Manager, interval time.Duration) *Reconciler {
	return &Reconciler{
		manager:   manager,
		interval:  interval,
		stopCh:    make(chan struct{}),
		isRunning: false,
	}
}

func (r *Reconciler) Start() {
	if r.isRunning {
		return
	}

	r.isRunning = true
	go r.reconcileLoop()
}

func (r *Reconciler) Stop() {
	if !r.isRunning {
		return
	}

	r.stopCh <- struct{}{}
	r.isRunning = false
}

func (r *Reconciler) reconcileLoop() {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), r.interval/2)
			if err := r.manager.ReconcileServices(ctx); err != nil {
				log.Printf("Error reconciling services: %v", err)
			}
			cancel()
		case <-r.stopCh:
			return
		}
	}
}