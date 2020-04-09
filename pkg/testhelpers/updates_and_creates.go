package testhelpers

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	kpackfakes "github.com/pivotal/kpack/pkg/client/clientset/versioned/fake"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime"
	k8sfakes "k8s.io/client-go/kubernetes/fake"
	clientgotesting "k8s.io/client-go/testing"
)

func TestKpackActions(t *testing.T, clientset *kpackfakes.Clientset, expectUpdates []clientgotesting.UpdateActionImpl, expectCreates []runtime.Object, expectDeletes []string) {
	t.Helper()
	actions, err := ActionRecorderList{clientset}.ActionsByVerb()
	require.NoError(t, err)

	for i, want := range expectCreates {
		if i >= len(actions.Creates) {
			t.Errorf("Missing create: %#v", want)
			continue
		}

		got := actions.Creates[i].GetObject()

		if diff := cmp.Diff(want, got, cmpopts.EquateEmpty()); diff != "" {
			t.Errorf("Unexpected create (-want, +got): %s", diff)
		}
	}

	for i, want := range expectDeletes {
		if i >= len(actions.Deletes) {
			t.Errorf("Missing delete: %#v", want)
			continue
		}

		got := actions.Deletes[i].GetName()

		if diff := cmp.Diff(want, got, cmpopts.EquateEmpty()); diff != "" {
			t.Errorf("Unexpected delete (-want, +got): %s", diff)
		}
	}

	if got, want := len(actions.Deletes), len(expectDeletes); got > want {
		for _, extra := range actions.Deletes[want:] {
			t.Errorf("Extra delete: %#v", extra.GetName())
		}
	}

	for i, want := range expectUpdates {
		if i >= len(actions.Updates) {
			wo := want.GetObject()
			t.Errorf("Missing update for %#v", wo)
			continue
		}

		got := actions.Updates[i].GetObject()

		if diff := cmp.Diff(want.GetObject(), got, cmpopts.EquateEmpty()); diff != "" {
			t.Errorf("Unexpected update (-want, +got): %s", diff)
		}
	}

	if got, want := len(actions.Updates), len(expectUpdates); got > want {
		for _, extra := range actions.Updates[want:] {
			t.Errorf("Extra update: %#v", extra.GetObject())
		}
	}

	for _, extra := range actions.DeleteCollections {
		t.Errorf("Extra delete-collection: %#v", extra)
	}

	for _, extra := range actions.Patches {
		t.Errorf("Extra patch: %#v; raw: %s", extra, string(extra.GetPatch()))
	}
}

func TestK8sActions(t *testing.T, clientset *k8sfakes.Clientset, expectUpdates []clientgotesting.UpdateActionImpl, expectCreates []runtime.Object, expectDeletes []string) {
	t.Helper()
	actions, err := ActionRecorderList{clientset}.ActionsByVerb()
	require.NoError(t, err)

	for i, want := range expectCreates {
		if i >= len(actions.Creates) {
			t.Errorf("Missing create: %#v", want)
			continue
		}

		got := actions.Creates[i].GetObject()

		if diff := cmp.Diff(want, got, cmpopts.EquateEmpty()); diff != "" {
			t.Errorf("Unexpected create (-want, +got): %s", diff)
		}
	}

	for i, want := range expectDeletes {
		if i >= len(actions.Deletes) {
			t.Errorf("Missing delete: %#v", want)
			continue
		}

		got := actions.Deletes[i].GetName()

		if diff := cmp.Diff(want, got, cmpopts.EquateEmpty()); diff != "" {
			t.Errorf("Unexpected delete (-want, +got): %s", diff)
		}
	}

	if got, want := len(actions.Deletes), len(expectDeletes); got > want {
		for _, extra := range actions.Deletes[want:] {
			t.Errorf("Extra delete: %#v", extra.GetName())
		}
	}

	for i, want := range expectUpdates {
		if i >= len(actions.Updates) {
			wo := want.GetObject()
			t.Errorf("Missing update for %#v", wo)
			continue
		}

		got := actions.Updates[i].GetObject()

		if diff := cmp.Diff(want.GetObject(), got, cmpopts.EquateEmpty()); diff != "" {
			t.Errorf("Unexpected update (-want, +got): %s", diff)
		}
	}

	if got, want := len(actions.Updates), len(expectUpdates); got > want {
		for _, extra := range actions.Updates[want:] {
			t.Errorf("Extra update: %#v", extra.GetObject())
		}
	}

	for _, extra := range actions.DeleteCollections {
		t.Errorf("Extra delete-collection: %#v", extra)
	}

	for _, extra := range actions.Patches {
		t.Errorf("Extra patch: %#v; raw: %s", extra, string(extra.GetPatch()))
	}
}
