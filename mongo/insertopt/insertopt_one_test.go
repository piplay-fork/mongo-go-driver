package insertopt

import (
	"testing"

	"reflect"

	"github.com/mongodb/mongo-go-driver/core/option"
	"github.com/mongodb/mongo-go-driver/internal/testutil/helpers"
)

func createNestedInsertOneBundle1(t *testing.T) *OneBundle {
	nestedBundle := BundleOne(BypassDocumentValidation(false))
	testhelpers.RequireNotNil(t, nestedBundle, "nested bundle was nil")

	outerBundle := BundleOne(BypassDocumentValidation(true), BypassDocumentValidation(true), nestedBundle)
	testhelpers.RequireNotNil(t, outerBundle, "outer bundle was nil")

	return outerBundle
}

// Test doubly nested bundle
func createNestedInsertOneBundle2(t *testing.T) *OneBundle {
	b1 := BundleOne(BypassDocumentValidation(false))
	testhelpers.RequireNotNil(t, b1, "nested bundle was nil")

	b2 := BundleOne(BypassDocumentValidation(false), b1)
	testhelpers.RequireNotNil(t, b2, "nested bundle was nil")

	outerBundle := BundleOne(BypassDocumentValidation(true), BypassDocumentValidation(true), b2)
	testhelpers.RequireNotNil(t, outerBundle, "outer bundle was nil")

	return outerBundle
}

// Test two top level nested bundles
func createNestedInsertOneBundle3(t *testing.T) *OneBundle {
	b1 := BundleOne(BypassDocumentValidation(false))
	testhelpers.RequireNotNil(t, b1, "nested bundle was nil")

	b2 := BundleOne(BypassDocumentValidation(false), b1)
	testhelpers.RequireNotNil(t, b2, "nested bundle was nil")

	b3 := BundleOne(BypassDocumentValidation(true))
	testhelpers.RequireNotNil(t, b3, "nested bundle was nil")

	b4 := BundleOne(BypassDocumentValidation(false), b3)
	testhelpers.RequireNotNil(t, b4, "nested bundle was nil")

	outerBundle := BundleOne(b4, BypassDocumentValidation(true), b2)
	testhelpers.RequireNotNil(t, outerBundle, "outer bundle was nil")

	return outerBundle
}

func TestInsertInsertOneOpt(t *testing.T) {
	var bundle1 *OneBundle
	bundle1 = bundle1.BypassDocumentValidation(true).BypassDocumentValidation(false)
	testhelpers.RequireNotNil(t, bundle1, "created bundle was nil")
	bundle1Opts := []option.Optioner{
		BypassDocumentValidation(true).ConvertInsertOption(),
		BypassDocumentValidation(false).ConvertInsertOption(),
	}
	bundle1DedupOpts := []option.Optioner{
		BypassDocumentValidation(false).ConvertInsertOption(),
	}

	bundle2 := BundleOne(BypassDocumentValidation(true))
	bundle2Opts := []option.Optioner{
		BypassDocumentValidation(true).ConvertInsertOption(),
	}

	bundle3 := BundleOne().
		BypassDocumentValidation(false).
		BypassDocumentValidation(true)

	bundle3Opts := []option.Optioner{
		OptBypassDocumentValidation(false).ConvertInsertOption(),
		OptBypassDocumentValidation(true).ConvertInsertOption(),
	}

	bundle3DedupOpts := []option.Optioner{
		OptBypassDocumentValidation(true).ConvertInsertOption(),
	}

	nilBundle := BundleOne()
	var nilBundleOpts []option.Optioner

	nestedBundle1 := createNestedInsertOneBundle1(t)
	nestedBundleOpts1 := []option.Optioner{
		OptBypassDocumentValidation(true).ConvertInsertOption(),
		OptBypassDocumentValidation(true).ConvertInsertOption(),
		OptBypassDocumentValidation(false).ConvertInsertOption(),
	}
	nestedBundleDedupOpts1 := []option.Optioner{
		OptBypassDocumentValidation(false).ConvertInsertOption(),
	}

	nestedBundle2 := createNestedInsertOneBundle2(t)
	nestedBundleOpts2 := []option.Optioner{
		OptBypassDocumentValidation(true).ConvertInsertOption(),
		OptBypassDocumentValidation(true).ConvertInsertOption(),
		OptBypassDocumentValidation(false).ConvertInsertOption(),
		OptBypassDocumentValidation(false).ConvertInsertOption(),
	}
	nestedBundleDedupOpts2 := []option.Optioner{
		OptBypassDocumentValidation(false).ConvertInsertOption(),
	}

	nestedBundle3 := createNestedInsertOneBundle3(t)
	nestedBundleOpts3 := []option.Optioner{
		OptBypassDocumentValidation(false).ConvertInsertOption(),
		OptBypassDocumentValidation(true).ConvertInsertOption(),
		OptBypassDocumentValidation(true).ConvertInsertOption(),
		OptBypassDocumentValidation(false).ConvertInsertOption(),
		OptBypassDocumentValidation(false).ConvertInsertOption(),
	}
	nestedBundleDedupOpts3 := []option.Optioner{
		OptBypassDocumentValidation(false).ConvertInsertOption(),
	}

	t.Run("MakeOptions", func(t *testing.T) {
		head := bundle1

		bundleLen := 0
		for head != nil && head.option != nil {
			bundleLen++
			head = head.next
		}

		if bundleLen != len(bundle1Opts) {
			t.Errorf("expected bundle length %d. got: %d", len(bundle1Opts), bundleLen)
		}
	})

	t.Run("TestAll", func(t *testing.T) {
		opts := []OneOption{
			BypassDocumentValidation(true),
		}
		params := make([]One, len(opts))
		for i := range opts {
			params[i] = opts[i]
		}
		bundle := BundleOne(params...)

		deleteOpts, _, err := bundle.Unbundle(true)
		testhelpers.RequireNil(t, err, "got non-nill error from unbundle: %s", err)

		if len(deleteOpts) != len(opts) {
			t.Errorf("expected unbundled opts len %d. got %d", len(opts), len(deleteOpts))
		}

		for i, opt := range opts {
			if !reflect.DeepEqual(opt.ConvertInsertOption(), deleteOpts[i]) {
				t.Errorf("opt mismatch. expected %#v, got %#v", opt, deleteOpts[i])
			}
		}
	})

	t.Run("Nil Option Bundle", func(t *testing.T) {
		sess := InsertSessionOpt{}
		opts, _, err := BundleOne(BypassDocumentValidation(true), BundleOne(nil), sess, nil).unbundle()
		testhelpers.RequireNil(t, err, "got non-nil error from unbundle: %s", err)

		if len(opts) != 1 {
			t.Errorf("expected bundle length 1. got: %d", len(opts))
		}

		opts, _, err = BundleOne(nil, sess, BundleOne(nil), BypassDocumentValidation(true)).unbundle()
		testhelpers.RequireNil(t, err, "got non-nil error from unbundle: %s", err)

		if len(opts) != 1 {
			t.Errorf("expected bundle length 1. got: %d", len(opts))
		}
	})

	t.Run("Unbundle", func(t *testing.T) {
		var cases = []struct {
			name         string
			dedup        bool
			bundle       *OneBundle
			expectedOpts []option.Optioner
		}{
			{"NilBundle", false, nilBundle, nilBundleOpts},
			{"Bundle1", false, bundle1, bundle1Opts},
			{"Bundle1Dedup", true, bundle1, bundle1DedupOpts},
			{"Bundle2", false, bundle2, bundle2Opts},
			{"Bundle2Dedup", true, bundle2, bundle2Opts},
			{"Bundle3", false, bundle3, bundle3Opts},
			{"Bundle3Dedup", true, bundle3, bundle3DedupOpts},
			{"NestedBundle1_DedupFalse", false, nestedBundle1, nestedBundleOpts1},
			{"NestedBundle1_DedupTrue", true, nestedBundle1, nestedBundleDedupOpts1},
			{"NestedBundle2_DedupFalse", false, nestedBundle2, nestedBundleOpts2},
			{"NestedBundle2_DedupTrue", true, nestedBundle2, nestedBundleDedupOpts2},
			{"NestedBundle3_DedupFalse", false, nestedBundle3, nestedBundleOpts3},
			{"NestedBundle3_DedupTrue", true, nestedBundle3, nestedBundleDedupOpts3},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				options, _, err := tc.bundle.Unbundle(tc.dedup)
				testhelpers.RequireNil(t, err, "got non-nill error from unbundle: %s", err)

				if len(options) != len(tc.expectedOpts) {
					t.Errorf("options length does not match expected length. got %d expected %d", len(options),
						len(tc.expectedOpts))
				} else {
					for i, opt := range options {
						if !reflect.DeepEqual(opt, tc.expectedOpts[i]) {
							t.Errorf("expected: %s\nreceived: %s", opt, tc.expectedOpts[i])
						}
					}
				}
			})
		}
	})
}
