package store

import (
	"encoding/json"
	"fmt"
	"time"

	"podium/internal/models"
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

// Add these functions to your existing boltdb.go file

func (s *BoltStore) CreateService(service models.Service) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("services"))
		if err != nil {
			return fmt.Errorf("failed to create services bucket: %w", err)
		}
		
		data, err := json.Marshal(service)
		if err != nil {
			return fmt.Errorf("failed to marshal service: %w", err)
		}
		
		return b.Put([]byte(service.ID), data)
	})
}

func (s *BoltStore) GetService(id string) (models.Service, error) {
	var service models.Service
	
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("services"))
		if b == nil {
			return fmt.Errorf("services bucket not found")
		}
		
		data := b.Get([]byte(id))
		if data == nil {
			return fmt.Errorf("service not found: %s", id)
		}
		
		return json.Unmarshal(data, &service)
	})
	
	return service, err
}

func (s *BoltStore) UpdateService(service models.Service) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("services"))
		if b == nil {
			return fmt.Errorf("services bucket not found")
		}
		
		data, err := json.Marshal(service)
		if err != nil {
			return fmt.Errorf("failed to marshal service: %w", err)
		}
		
		return b.Put([]byte(service.ID), data)
	})
}

func (s *BoltStore) DeleteService(id string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("services"))
		if b == nil {
			return fmt.Errorf("services bucket not found")
		}
		
		return b.Delete([]byte(id))
	})
}

func (s *BoltStore) ListServices() ([]models.Service, error) {
	var services []models.Service
	
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("services"))
		if b == nil {
			return nil
		}
		
		return b.ForEach(func(k, v []byte) error {
			var service models.Service
			if err := json.Unmarshal(v, &service); err != nil {
				return fmt.Errorf("failed to unmarshal service: %w", err)
			}
			
			services = append(services, service)
			return nil
		})
	})
	
	return services, err
}

func (s *BoltStore) GetServiceByName(name string) (models.Service, error) {
	var service models.Service
	var found bool
	
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("services"))
		if b == nil {
			return fmt.Errorf("services bucket not found")
		}
		
		return b.ForEach(func(k, v []byte) error {
			var s models.Service
			if err := json.Unmarshal(v, &s); err != nil {
				return fmt.Errorf("failed to unmarshal service: %w", err)
			}
			
			if s.Name == name {
				service = s
				found = true
				return nil
			}
			return nil
		})
	})
	
	if !found {
		return service, fmt.Errorf("service not found: %s", name)
	}
	
	return service, err
}


func (s *BoltStore) GetContainersByStatus(status string) ([]models.Container, error) {
	var containers []models.Container
	
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("containers"))
		
		return b.ForEach(func(k, v []byte) error {
			var container models.Container
			if err := json.Unmarshal(v, &container); err != nil {
				return fmt.Errorf("failed to unmarshal container: %w", err)
			}
			
			if container.Status == status {
				containers = append(containers, container)
			}
			return nil
		})
	})
	
	return containers, err
}