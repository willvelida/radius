/*
Copyright 2023 The Radius Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package queueprovider

import (
	"context"
	"errors"
	"sync"

	"github.com/radius-project/radius/pkg/components/queue"
)

var (
	ErrUnsupportedQueueProvider = errors.New("unsupported queue provider")
)

// QueueProvider is the provider to create and manage queue client.
type QueueProvider struct {
	options QueueProviderOptions

	queueClient queue.Client
	once        sync.Once
}

// New creates new QueueProvider instance.
func New(opts QueueProviderOptions) *QueueProvider {
	return &QueueProvider{
		queueClient: nil,
		options:     opts,
	}
}

// GetClient creates or gets queue client.
func (p *QueueProvider) GetClient(ctx context.Context) (queue.Client, error) {
	if p.queueClient != nil {
		return p.queueClient, nil
	}

	err := ErrUnsupportedQueueProvider
	p.once.Do(func() {
		if fn, ok := clientFactory[p.options.Provider]; ok {
			p.queueClient, err = fn(ctx, p.options)
		}
	})

	return p.queueClient, err
}

// SetClient sets the queue client for the QueueProvider. This should be used by tests that need to mock the queue client.
func (p *QueueProvider) SetClient(client queue.Client) {
	p.queueClient = client
}
