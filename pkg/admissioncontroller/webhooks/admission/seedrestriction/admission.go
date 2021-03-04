// Copyright (c) 2021 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package seedrestriction

import (
	"context"

	acadmission "github.com/gardener/gardener/pkg/admissioncontroller/webhooks/admission"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

const (
	// HandlerName is the name of this admission webhook handler.
	HandlerName = "seedrestriction"
	// WebhookPath is the HTTP handler path for this admission webhook handler.
	WebhookPath = "/webhooks/admission/seedrestriction"
)

// New creates a new webhook handler restricting requests by gardenlets. It allows all requests.
func New(logger logr.Logger) *handler {
	return &handler{logger: logger}
}

type handler struct {
	logger logr.Logger
}

var _ admission.Handler = &handler{}

func (p *handler) Handle(_ context.Context, _ admission.Request) admission.Response {
	return acadmission.Allowed("")
}
