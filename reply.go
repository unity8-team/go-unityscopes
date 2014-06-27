package scopes

// #include <stdlib.h>
// #include "shim.h"
import "C"
import (
	"encoding/json"
	"runtime"
	"unsafe"
)

// SearchReply is used to send results of search queries to the client.
type SearchReply struct {
	r C.SharedPtrData
}

func makeSearchReply(replyData *C.uintptr_t) *SearchReply {
	reply := new(SearchReply)
	runtime.SetFinalizer(reply, finalizeSearchReply)
	C.init_search_reply_ptr(&reply.r[0], replyData)
	return reply
}

func finalizeSearchReply(reply *SearchReply) {
	C.destroy_search_reply_ptr(&reply.r[0])
}

// Finished is called to indicate that no further results will be
// pushed to this reply.
//
// This is called automatically if a scope's Search method completes
// without error.
func (reply *SearchReply) Finished() {
	C.search_reply_finished(&reply.r[0])
}

// Error is called to indicate that search query could not be
// completed successfully.
//
// This is called automatically if a scope's Search method completes
// with an error.
func (reply *SearchReply) Error(err error) {
	errString := err.Error()
	C.search_reply_error(&reply.r[0], unsafe.Pointer(&errString))
}

// RegisterCategory registers a new results category with the client.
//
// The template parameter should either be empty (to use the default
// rendering template), or contain a JSON template as described here:
//
// http://developer.ubuntu.com/api/scopes/sdk-14.04/unity.scopes.CategoryRenderer/#details
//
// Categories can be passed to NewCategorisedResult in order to
// construct search results.
func (reply *SearchReply) RegisterCategory(id, title, icon, template string) *Category {
	cat := new(Category)
	runtime.SetFinalizer(cat, finalizeCategory)
	C.search_reply_register_category(&reply.r[0], unsafe.Pointer(&id), unsafe.Pointer(&title), unsafe.Pointer(&icon), unsafe.Pointer(&template), &cat.c[0])
	return cat
}

// Push sends a search result to the client.
func (reply *SearchReply) Push(result *CategorisedResult) error {
	var errorString *C.char = nil
	C.search_reply_push(&reply.r[0], result.result, &errorString)
	return checkError(errorString)
}

// PreviewReply is used to send result previews to the client.
type PreviewReply struct {
	r C.SharedPtrData
}

func makePreviewReply(replyData *C.uintptr_t) *PreviewReply {
	reply := new(PreviewReply)
	runtime.SetFinalizer(reply, finalizePreviewReply)
	C.init_preview_reply_ptr(&reply.r[0], replyData)
	return reply
}

func finalizePreviewReply(reply *PreviewReply) {
	C.destroy_search_reply_ptr(&reply.r[0])
}

// Finished is called to indicate that no further widgets or
// attributes will be pushed to this reply.
//
// This is called automatically if a scope's Preview method completes
// without error.
func (reply *PreviewReply) Finished() {
	C.preview_reply_finished(&reply.r[0])
}

// Error is called to indicate that the preview generation could not
// be completed successfully.
//
// This is called automatically if a scope's Preview method completes
// with an error.
func (reply *PreviewReply) Error(err error) {
	errString := err.Error()
	C.preview_reply_error(&reply.r[0], unsafe.Pointer(&errString))
}

// PushWidgets sends one or more preview widgets to the client.
func (reply *PreviewReply) PushWidgets(widgets ...PreviewWidget) error {
	widget_data := make([]string, len(widgets))
	for i, w := range widgets {
		data, err := w.data()
		if err != nil {
			return err
		}
		widget_data[i] = string(data)
	}
	var errorString *C.char = nil
	C.preview_reply_push_widgets(&reply.r[0], unsafe.Pointer(&widget_data[0]), C.int(len(widget_data)), &errorString)
	return checkError(errorString)
}

// PushAttr pushes a preview attribute to the client.
//
// This will augment the set of attributes in the result available to
// be mapped by preview widgets.  This allows a widget to be sent to
// the client early, and then fill it in later when the information is
// available.
func (reply *PreviewReply) PushAttr(attr string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	json_value := string(data)
	var errorString *C.char = nil
	C.preview_reply_push_attr(&reply.r[0], unsafe.Pointer(&attr), unsafe.Pointer(&json_value), &errorString)
	return checkError(errorString)
}
