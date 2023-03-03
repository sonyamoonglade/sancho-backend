package meta_cache

import (
	"github.com/sonyamoonglade/sancho-backend/internal/domain"
	"go.uber.org/atomic"
)

type MetaCache struct {
	data *atomic.Pointer[domain.BusinessMeta]
	set  *atomic.Bool
}

func NewMetaCache() *MetaCache {
	return &MetaCache{
		data: atomic.NewPointer[domain.BusinessMeta](nil),
		set:  atomic.NewBool(false),
	}
}

func (m *MetaCache) Get() *domain.BusinessMeta {
	if ok := m.set.Load(); ok {
		return m.data.Load()
	}

	panic("meta is not set")
}

func (m *MetaCache) Set(meta domain.BusinessMeta) {
	m.data.Store(&meta)
	m.set.Store(true)
}
