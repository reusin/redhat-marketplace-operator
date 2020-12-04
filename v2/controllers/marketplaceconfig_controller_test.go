// Copyright 2020 IBM Corp.
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

package controllers

import (
	"github.com/gotidy/ptr"
	. "github.com/redhat-marketplace/redhat-marketplace-operator/v2/test/rectest"

	. "github.com/onsi/ginkgo"
	opsrcApi "github.com/operator-framework/api/pkg/operators/v1"
	operatorsv1alpha1 "github.com/operator-framework/api/pkg/operators/v1alpha1"
	marketplacev1alpha1 "github.com/redhat-marketplace/redhat-marketplace-operator/v2/api/v1alpha1"
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/utils"
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/utils/reconcileutils"
	"github.com/spf13/viper"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var _ = Describe("Testing with Ginkgo", func() {

	var (
		name                 = utils.MARKETPLACECONFIG_NAME
		namespace            = "redhat-marketplace-operator"
		customerID    string = "example-userid"
		razeeName            = "rhm-marketplaceconfig-razeedeployment"
		meterBaseName        = "rhm-marketplaceconfig-meterbase"
		req                  = reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      name,
				Namespace: namespace,
			},
		}

		opts = []StepOption{
			WithRequest(req),
		}
		marketplaceconfig = utils.BuildMarketplaceConfigCR(namespace, customerID)
		razeedeployment   = utils.BuildRazeeCr(namespace, marketplaceconfig.Spec.ClusterUUID, marketplaceconfig.Spec.DeploySecretName)
		meterbase         = utils.BuildMeterBaseCr(namespace)
	)

	var setup = func(r *ReconcilerTest) error {
		s := scheme.Scheme
		_ = opsrcApi.AddToScheme(s)
		_ = operatorsv1alpha1.AddToScheme(s)
		s.AddKnownTypes(marketplacev1alpha1.SchemeGroupVersion, marketplaceconfig)
		s.AddKnownTypes(marketplacev1alpha1.SchemeGroupVersion, razeedeployment)
		s.AddKnownTypes(marketplacev1alpha1.SchemeGroupVersion, meterbase)

		r.Client = fake.NewFakeClient(r.GetGetObjects()...)
		r.Reconciler = &ReconcileMarketplaceConfig{
			client: r.Client,
			scheme: s,
			cc:     reconcileutils.NewLoglessClientCommand(r.Client, s),
		}
		return nil
	}

	var testCleanInstall = func(t GinkgoTInterface) {
		t.Parallel()
		marketplaceconfig.Spec.EnableMetering = ptr.Bool(true)
		marketplaceconfig.Spec.InstallIBMCatalogSource = ptr.Bool(true)
		reconcilerTest := NewReconcilerTest(setup, marketplaceconfig)
		reconcilerTest.TestAll(t,
			ReconcileStep(opts, ReconcileWithUntilDone(true)),
			GetStep(opts,
				GetWithNamespacedName(razeeName, namespace),
				GetWithObj(&marketplacev1alpha1.RazeeDeployment{}),
			),
			GetStep(opts,
				GetWithNamespacedName(meterBaseName, namespace),
				GetWithObj(&marketplacev1alpha1.MeterBase{}),
			),
			GetStep(opts,
				GetWithNamespacedName(utils.OPSRC_NAME, utils.OPERATOR_MKTPLACE_NS),
				GetWithObj(&unstructured.Unstructured{}),
			),
			GetStep(opts,
				GetWithNamespacedName(utils.IBM_CATALOGSRC_NAME, utils.OPERATOR_MKTPLACE_NS),
				GetWithObj(&operatorsv1alpha1.CatalogSource{}),
			),
			GetStep(opts,
				GetWithNamespacedName(utils.OPENCLOUD_CATALOGSRC_NAME, utils.OPERATOR_MKTPLACE_NS),
				GetWithObj(&operatorsv1alpha1.CatalogSource{}),
			),
		)
	}

	It("marketplace config controller", func() {
		defaultFeatures := []string{"razee", "meterbase"}
		viper.Set("features", defaultFeatures)
		viper.Set("IBMCatalogSource", true)
		testCleanInstall(GinkgoT())
	})
})
