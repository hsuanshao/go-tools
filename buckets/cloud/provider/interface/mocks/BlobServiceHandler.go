// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	entity "github.com/hsuanshao/go-tools/buckets/entity"
	ctx "github.com/hsuanshao/go-tools/ctx"

	mock "github.com/stretchr/testify/mock"

	time "time"
)

// BlobServiceHandler is an autogenerated mock type for the BlobServiceHandler type
type BlobServiceHandler struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *BlobServiceHandler) Close() {
	_m.Called()
}

// Delete provides a mock function with given fields: _a0, contentType, objPathes
func (_m *BlobServiceHandler) Delete(_a0 ctx.CTX, contentType entity.ContentType, objPathes []string) (bool, error) {
	ret := _m.Called(_a0, contentType, objPathes)

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(ctx.CTX, entity.ContentType, []string) (bool, error)); ok {
		return rf(_a0, contentType, objPathes)
	}
	if rf, ok := ret.Get(0).(func(ctx.CTX, entity.ContentType, []string) bool); ok {
		r0 = rf(_a0, contentType, objPathes)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(ctx.CTX, entity.ContentType, []string) error); ok {
		r1 = rf(_a0, contentType, objPathes)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GenReadPresignedURL provides a mock function with given fields: _a0, objURL, duration
func (_m *BlobServiceHandler) GenReadPresignedURL(_a0 ctx.CTX, objURL string, duration time.Duration) (string, error) {
	ret := _m.Called(_a0, objURL, duration)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(ctx.CTX, string, time.Duration) (string, error)); ok {
		return rf(_a0, objURL, duration)
	}
	if rf, ok := ret.Get(0).(func(ctx.CTX, string, time.Duration) string); ok {
		r0 = rf(_a0, objURL, duration)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(ctx.CTX, string, time.Duration) error); ok {
		r1 = rf(_a0, objURL, duration)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetObjectList provides a mock function with given fields: _a0, prefix, delim
func (_m *BlobServiceHandler) GetObjectList(_a0 ctx.CTX, prefix string, delim string) ([]string, error) {
	ret := _m.Called(_a0, prefix, delim)

	var r0 []string
	var r1 error
	if rf, ok := ret.Get(0).(func(ctx.CTX, string, string) ([]string, error)); ok {
		return rf(_a0, prefix, delim)
	}
	if rf, ok := ret.Get(0).(func(ctx.CTX, string, string) []string); ok {
		r0 = rf(_a0, prefix, delim)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	if rf, ok := ret.Get(1).(func(ctx.CTX, string, string) error); ok {
		r1 = rf(_a0, prefix, delim)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Health provides a mock function with given fields: _a0, csp
func (_m *BlobServiceHandler) Health(_a0 ctx.CTX, csp entity.CloudServiceProvider) (entity.HealthStatus, error) {
	ret := _m.Called(_a0, csp)

	var r0 entity.HealthStatus
	var r1 error
	if rf, ok := ret.Get(0).(func(ctx.CTX, entity.CloudServiceProvider) (entity.HealthStatus, error)); ok {
		return rf(_a0, csp)
	}
	if rf, ok := ret.Get(0).(func(ctx.CTX, entity.CloudServiceProvider) entity.HealthStatus); ok {
		r0 = rf(_a0, csp)
	} else {
		r0 = ret.Get(0).(entity.HealthStatus)
	}

	if rf, ok := ret.Get(1).(func(ctx.CTX, entity.CloudServiceProvider) error); ok {
		r1 = rf(_a0, csp)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsObjectExists provides a mock function with given fields: _a0, objURL
func (_m *BlobServiceHandler) IsObjectExists(_a0 ctx.CTX, objURL string) (bool, error) {
	ret := _m.Called(_a0, objURL)

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(ctx.CTX, string) (bool, error)); ok {
		return rf(_a0, objURL)
	}
	if rf, ok := ret.Get(0).(func(ctx.CTX, string) bool); ok {
		r0 = rf(_a0, objURL)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(ctx.CTX, string) error); ok {
		r1 = rf(_a0, objURL)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Override provides a mock function with given fields: _a0, ct, objPath, objNewRaw, objmetadata
func (_m *BlobServiceHandler) Override(_a0 ctx.CTX, ct entity.ContentType, objPath string, objNewRaw []byte, objmetadata map[string]string) (string, error) {
	ret := _m.Called(_a0, ct, objPath, objNewRaw, objmetadata)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(ctx.CTX, entity.ContentType, string, []byte, map[string]string) (string, error)); ok {
		return rf(_a0, ct, objPath, objNewRaw, objmetadata)
	}
	if rf, ok := ret.Get(0).(func(ctx.CTX, entity.ContentType, string, []byte, map[string]string) string); ok {
		r0 = rf(_a0, ct, objPath, objNewRaw, objmetadata)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(ctx.CTX, entity.ContentType, string, []byte, map[string]string) error); ok {
		r1 = rf(_a0, ct, objPath, objNewRaw, objmetadata)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutPresignedURL provides a mock function with given fields: _a0, objURL, mime, duration, metaData
func (_m *BlobServiceHandler) PutPresignedURL(_a0 ctx.CTX, objURL string, mime entity.ContentType, duration time.Duration, metaData map[string]string) (string, error) {
	ret := _m.Called(_a0, objURL, mime, duration, metaData)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(ctx.CTX, string, entity.ContentType, time.Duration, map[string]string) (string, error)); ok {
		return rf(_a0, objURL, mime, duration, metaData)
	}
	if rf, ok := ret.Get(0).(func(ctx.CTX, string, entity.ContentType, time.Duration, map[string]string) string); ok {
		r0 = rf(_a0, objURL, mime, duration, metaData)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(ctx.CTX, string, entity.ContentType, time.Duration, map[string]string) error); ok {
		r1 = rf(_a0, objURL, mime, duration, metaData)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ReadObjectContent provides a mock function with given fields: _a0, objectPath
func (_m *BlobServiceHandler) ReadObjectContent(_a0 ctx.CTX, objectPath string) ([]byte, map[string]string, error) {
	ret := _m.Called(_a0, objectPath)

	var r0 []byte
	var r1 map[string]string
	var r2 error
	if rf, ok := ret.Get(0).(func(ctx.CTX, string) ([]byte, map[string]string, error)); ok {
		return rf(_a0, objectPath)
	}
	if rf, ok := ret.Get(0).(func(ctx.CTX, string) []byte); ok {
		r0 = rf(_a0, objectPath)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(ctx.CTX, string) map[string]string); ok {
		r1 = rf(_a0, objectPath)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(map[string]string)
		}
	}

	if rf, ok := ret.Get(2).(func(ctx.CTX, string) error); ok {
		r2 = rf(_a0, objectPath)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// Upload provides a mock function with given fields: _a0, ct, objpath, objraw, objmetadata
func (_m *BlobServiceHandler) Upload(_a0 ctx.CTX, ct entity.ContentType, objpath string, objraw []byte, objmetadata map[string]string) (string, string, error) {
	ret := _m.Called(_a0, ct, objpath, objraw, objmetadata)

	var r0 string
	var r1 string
	var r2 error
	if rf, ok := ret.Get(0).(func(ctx.CTX, entity.ContentType, string, []byte, map[string]string) (string, string, error)); ok {
		return rf(_a0, ct, objpath, objraw, objmetadata)
	}
	if rf, ok := ret.Get(0).(func(ctx.CTX, entity.ContentType, string, []byte, map[string]string) string); ok {
		r0 = rf(_a0, ct, objpath, objraw, objmetadata)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(ctx.CTX, entity.ContentType, string, []byte, map[string]string) string); ok {
		r1 = rf(_a0, ct, objpath, objraw, objmetadata)
	} else {
		r1 = ret.Get(1).(string)
	}

	if rf, ok := ret.Get(2).(func(ctx.CTX, entity.ContentType, string, []byte, map[string]string) error); ok {
		r2 = rf(_a0, ct, objpath, objraw, objmetadata)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

type mockConstructorTestingTNewBlobServiceHandler interface {
	mock.TestingT
	Cleanup(func())
}

// NewBlobServiceHandler creates a new instance of BlobServiceHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewBlobServiceHandler(t mockConstructorTestingTNewBlobServiceHandler) *BlobServiceHandler {
	mock := &BlobServiceHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
