// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package insertopt

import (
	"reflect"

	"github.com/mongodb/mongo-go-driver/core/option"
	"github.com/mongodb/mongo-go-driver/core/writeconcern"
)

var oneBundle = new(OneBundle)
var manyBundle = new(ManyBundle)

// One is options for InsertOne
type One interface {
	one()
	ConvertOneOption() option.InsertOptioner
}

// Many is optinos for InsertMany
type Many interface {
	many()
	ConvertManyOption() option.InsertOptioner
}

// OneBundle is a bundle of One options
type OneBundle struct {
	option One
	next   *OneBundle
}

// Implement the One interface
func (ob *OneBundle) one() {}

// ConvertOneOption implements the One interface
func (ob *OneBundle) ConvertOneOption() option.InsertOptioner { return nil }

// BundleOne bundles One options
func BundleOne(opts ...One) *OneBundle {
	head := oneBundle

	for _, opt := range opts {
		newBundle := OneBundle{
			option: opt,
			next:   head,
		}

		head = &newBundle
	}

	return head
}

// BypassDocumentValidation adds an option allowing the write to opt-out of the document-level validation.
func (ob *OneBundle) BypassDocumentValidation(b bool) *OneBundle {
	bundle := &OneBundle{
		option: BypassDocumentValidation(b),
		next:   ob,
	}

	return bundle
}

// WriteConcern adds an option to specify a write concern
func (ob *OneBundle) WriteConcern(wc *writeconcern.WriteConcern) *OneBundle {
	bundle := &OneBundle{
		option: WriteConcern(wc),
		next:   ob,
	}

	return bundle
}

// Calculates the total length of a bundle, accounting for nested bundles.
func (ob *OneBundle) bundleLength() int {
	if ob == nil {
		return 0
	}

	bundleLen := 0
	for ; ob != nil && ob.option != nil; ob = ob.next {
		if converted, ok := ob.option.(*OneBundle); ok {
			// nested bundle
			bundleLen += converted.bundleLength()
			continue
		}

		bundleLen++
	}

	return bundleLen
}

// Unbundle unwinds and deduplicates the options used to create it and those
// added after creation into a single slice of options.
//
// The deduplicate parameter is used to determine if the bundle is just flattened or
// if we actually deduplicate options.
//
// Since a bundle can be recursive, this method will unwind all recursive bundles.
func (ob *OneBundle) Unbundle(deduplicate bool) ([]option.InsertOptioner, error) {
	options, err := ob.unbundle()
	if err != nil {
		return nil, err
	}

	if !deduplicate {
		return options, nil
	}

	// iterate backwards and make dedup slice
	optionsSet := make(map[reflect.Type]struct{})

	for i := len(options) - 1; i >= 0; i-- {
		currOption := options[i]
		optionType := reflect.TypeOf(currOption)

		if _, ok := optionsSet[optionType]; ok {
			// option already found
			options = append(options[:i], options[i+1:]...)
			continue
		}

		optionsSet[optionType] = struct{}{}
	}

	return options, nil
}

// Helper that recursively unwraps bundle into slice of options
func (ob *OneBundle) unbundle() ([]option.InsertOptioner, error) {
	if ob == nil {
		return nil, nil
	}

	listLen := ob.bundleLength()

	options := make([]option.InsertOptioner, listLen)
	index := listLen - 1

	for listHead := ob; listHead != nil && listHead.option != nil; listHead = listHead.next {
		// if the current option is a nested bundle, Unbundle it and add its options to the current array
		if converted, ok := listHead.option.(*OneBundle); ok {
			nestedOptions, err := converted.unbundle()
			if err != nil {
				return nil, err
			}

			// where to start inserting nested options
			startIndex := index - len(nestedOptions) + 1

			// add nested options in order
			for _, nestedOp := range nestedOptions {
				options[startIndex] = nestedOp
				startIndex++
			}
			index -= len(nestedOptions)
			continue
		}

		options[index] = listHead.option.ConvertOneOption()
		index--
	}

	return options, nil
}

// String implements the Stringer interface
func (ob *OneBundle) String() string {
	if ob == nil {
		return ""
	}

	str := ""
	for head := ob; head != nil && head.option != nil; head = head.next {
		if converted, ok := head.option.(*OneBundle); ok {
			str += converted.String()
			continue
		}

		str += head.option.ConvertOneOption().String() + "\n"
	}

	return str
}

// ManyBundle is a bundle of InsertMany options
type ManyBundle struct {
	option Many
	next   *ManyBundle
}

// BundleMany bundles Many options
func BundleMany(opts ...Many) *ManyBundle {
	head := manyBundle

	for _, opt := range opts {
		newBundle := ManyBundle{
			option: opt,
			next:   head,
		}

		head = &newBundle
	}

	return head
}

// Implement the Many interface
func (mb *ManyBundle) many() {}

// ConvertManyOption implements the Many interface
func (mb *ManyBundle) ConvertManyOption() option.InsertOptioner { return nil }

// BypassDocumentValidation adds an option allowing the write to opt-out of the document-level validation.
func (mb *ManyBundle) BypassDocumentValidation(b bool) *ManyBundle {
	bundle := &ManyBundle{
		option: BypassDocumentValidation(b),
		next:   mb,
	}

	return bundle
}

// Ordered adds an option that if true and insert fails, returns without performing remaining writes, otherwise continues
func (mb *ManyBundle) Ordered(b bool) *ManyBundle {
	bundle := &ManyBundle{
		option: Ordered(b),
		next:   mb,
	}

	return bundle
}

// WriteConcern adds an option to specify a write concern
func (mb *ManyBundle) WriteConcern(wc *writeconcern.WriteConcern) *ManyBundle {
	bundle := &ManyBundle{
		option: WriteConcern(wc),
		next:   mb,
	}

	return bundle
}

// Calculates the total length of a bundle, accounting for nested bundles.
func (mb *ManyBundle) bundleLength() int {
	if mb == nil {
		return 0
	}

	bundleLen := 0
	for ; mb != nil && mb.option != nil; mb = mb.next {
		if converted, ok := mb.option.(*ManyBundle); ok {
			// nested bundle
			bundleLen += converted.bundleLength()
			continue
		}

		bundleLen++
	}

	return bundleLen
}

// Unbundle unwinds and deduplicates the options used to create it and those
// added after creation into a single slice of options.
//
// The deduplicate parameter is used to determine if the bundle is just flattened or
// if we actually deduplicate options.
//
// Since a bundle can be recursive, this method will unwind all recursive bundles.
func (mb *ManyBundle) Unbundle(deduplicate bool) ([]option.InsertOptioner, error) {
	options, err := mb.unbundle()
	if err != nil {
		return nil, err
	}

	if !deduplicate {
		return options, nil
	}

	// iterate backwards and make dedup slice
	optionsSet := make(map[reflect.Type]struct{})

	for i := len(options) - 1; i >= 0; i-- {
		currOption := options[i]
		optionType := reflect.TypeOf(currOption)

		if _, ok := optionsSet[optionType]; ok {
			// option already found
			options = append(options[:i], options[i+1:]...)
			continue
		}

		optionsSet[optionType] = struct{}{}
	}

	return options, nil
}

// Helper that recursively unwraps bundle into slice of options
func (mb *ManyBundle) unbundle() ([]option.InsertOptioner, error) {
	if mb == nil {
		return nil, nil
	}

	listLen := mb.bundleLength()

	options := make([]option.InsertOptioner, listLen)
	index := listLen - 1

	for listHead := mb; listHead != nil && listHead.option != nil; listHead = listHead.next {
		// if the current option is a nested bundle, Unbundle it and add its options to the current array
		if converted, ok := listHead.option.(*ManyBundle); ok {
			nestedOptions, err := converted.unbundle()
			if err != nil {
				return nil, err
			}

			// where to start inserting nested options
			startIndex := index - len(nestedOptions) + 1

			// add nested options in order
			for _, nestedOp := range nestedOptions {
				options[startIndex] = nestedOp
				startIndex++
			}
			index -= len(nestedOptions)
			continue
		}

		options[index] = listHead.option.ConvertManyOption()
		index--
	}

	return options, nil
}

// String implements the Stringer interface
func (mb *ManyBundle) String() string {
	if mb == nil {
		return ""
	}

	str := ""
	for head := mb; head != nil && head.option != nil; head = head.next {
		if converted, ok := head.option.(*ManyBundle); ok {
			str += converted.String()
			continue
		}

		str += head.option.ConvertManyOption().String() + "\n"
	}

	return str
}

// BypassDocumentValidation allows the write to opt-out of the document-level validation.
func BypassDocumentValidation(b bool) OptBypassDocumentValidation {
	return OptBypassDocumentValidation(b)
}

// Ordered if true and insert fails, returns without performing remaining writes, otherwise continues
func Ordered(b bool) OptOrdered {
	return OptOrdered(b)
}

// WriteConcern specifies a write concern
func WriteConcern(wc *writeconcern.WriteConcern) OptWriteConcern {
	return OptWriteConcern{
		WriteConcern: wc,
	}
}

// OptBypassDocumentValidation allows the write to opt-out of the document-level validation.
type OptBypassDocumentValidation option.OptBypassDocumentValidation

// OptOrdered if true and insert fails, returns without performing remaining writes, otherwise continues
type OptOrdered option.OptOrdered

// OptWriteConcern specifies a write concern
type OptWriteConcern option.OptWriteConcern

func (OptBypassDocumentValidation) one() {}

// ConvertOneOption implements the One interface
func (opt OptBypassDocumentValidation) ConvertOneOption() option.InsertOptioner {
	return option.OptBypassDocumentValidation(opt)
}

func (OptWriteConcern) one() {}

// ConvertOneOption implements the One interface
func (opt OptWriteConcern) ConvertOneOption() option.InsertOptioner {
	return option.OptWriteConcern(opt)
}

func (OptWriteConcern) many() {}

// ConvertManyOption implements the Many interface
func (opt OptWriteConcern) ConvertManyOption() option.InsertOptioner {
	return option.OptWriteConcern(opt)
}

func (OptBypassDocumentValidation) many() {}

// ConvertManyOption implements the Many interface
func (opt OptBypassDocumentValidation) ConvertManyOption() option.InsertOptioner {
	return option.OptBypassDocumentValidation(opt)
}

func (OptOrdered) many() {}

// ConvertManyOption implements the Many interface
func (opt OptOrdered) ConvertManyOption() option.InsertOptioner {
	return option.OptOrdered(opt)
}
