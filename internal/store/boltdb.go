package store

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/odlemon/podium/internal/models"
	bolt "go.etcd.io/bbolt"
)

type BoltStore struct {
	db *bolt.DB
}

func NewBoltStore(path string) (*BoltStore, error) {
	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("failed to open bolt db: %w", err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("containers"))
		if err != nil {
			return fmt.Errorf("failed to create containers bucket: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create buckets: %w", err)
	}

	return &BoltStore{db: db}, nil
}

func (s *BoltStore) Close() error {
	return s.db.Close()
}

func (s *BoltStore) CreateContainer(container models.Container) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("containers"))
		
		data, err := json.Marshal(container)
		if err != nil {
			return fmt.Errorf("failed to marshal container: %w", err)
		}
		
		return b.Put([]byte(container.ID), data)
	})
}

func (s *BoltStore) GetContainer(id string) (models.Container, error) {
	var container models.Container
	
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("containers"))
		data := b.Get([]byte(id))
		
		if data == nil {
			return fmt.Errorf("container not found: %s", id)
		}
		
		return json.Unmarshal(data, &container)
	})
	
	return container, err
}

func (s *BoltStore) ListContainers() ([]models.Container, error) {
	var containers []models.Container
	
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("containers"))
		
		return b.ForEach(func(k, v []byte) error {
			var container models.Container
			if err := json.Unmarshal(v, &container); err != nil {
				return fmt.Errorf("failed to unmarshal container: %w", err)
			}
			
			containers = append(containers, container)
			return nil
		})
	})
	
	return containers, err
}

func (s *BoltStore) UpdateContainer(container models.Container) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("containers"))
		
		data, err := json.Marshal(container)
		if err != nil {
			return fmt.Errorf("failed to marshal container: %w", err)
		}
		
		return b.Put([]byte(container.ID), data)
	})
}

func (s *BoltStore) DeleteContainer(id string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("containers"))
		return b.Delete([]byte(id))
	})
}