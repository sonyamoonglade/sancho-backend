package meta_cache

import (
	"sync"
	"testing"

	"github.com/sonyamoonglade/sancho-backend/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestMetaCache(t *testing.T) {

	t.Run("should save and get meta", func(t *testing.T) {
		meta := &domain.BusinessMeta{
			DeliveryPunishmentThreshold: 100,
			DeliveryPunishmentValue:     100,
		}
		cache := NewMetaCache()
		cache.Set(*meta)

		receivedMeta := cache.Get()
		require.EqualValues(t, *meta, *receivedMeta)
		require.NotNil(t, receivedMeta)
	})

	t.Run("get without save should panic", func(t *testing.T) {
		cache := NewMetaCache()
		require.Panics(t, func() {
			cache.Get()
		})
	})

	t.Run("concurrent get", func(t *testing.T) {
		cache := NewMetaCache()
		wg := new(sync.WaitGroup)
		meta := &domain.BusinessMeta{
			DeliveryPunishmentThreshold: 100,
			DeliveryPunishmentValue:     100,
		}
		cache.Set(*meta)

		wg.Add(100)
		for i := 0; i < 100; i++ {
			go func(wg *sync.WaitGroup) {
				receivedMeta := cache.Get()
				require.EqualValues(t, *meta, *receivedMeta)
				require.NotNil(t, receivedMeta)
				defer wg.Done()
			}(wg)
		}

		wg.Wait()
	})

	t.Run("concurrent get and set", func(t *testing.T) {
		cache := NewMetaCache()
		meta := domain.BusinessMeta{
			DeliveryPunishmentThreshold: 400,
			DeliveryPunishmentValue:     500,
		}
		cache.Set(meta)
		wg := new(sync.WaitGroup)
		wg.Add(200)
		for i := 0; i < 100; i++ {
			// reader
			go func(wg *sync.WaitGroup) {
				receivedMeta := cache.Get()
				require.NotNil(t, receivedMeta)
				defer wg.Done()
			}(wg)
			// writer
			go func(wg *sync.WaitGroup) {
				meta := domain.BusinessMeta{
					DeliveryPunishmentThreshold: 200,
					DeliveryPunishmentValue:     300,
				}
				cache.Set(meta)
				defer wg.Done()
			}(wg)
		}

		wg.Wait()
	})

}
