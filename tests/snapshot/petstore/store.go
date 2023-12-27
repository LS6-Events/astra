package petstore

import (
	"fmt"
	"sync"
	"sync/atomic"
)

var Pets = []Pet{
	{ID: 1, Name: "Dog", PhotoURLs: []string{}, Status: "available", Tags: nil},
	{ID: 2, Name: "Cat", PhotoURLs: []string{}, Status: "pending", Tags: nil},
}

var petsLock = &sync.Mutex{}
var lastPetID int64 = 2

func newPetID() int64 {
	return atomic.AddInt64(&lastPetID, 1)
}

func AddPet(pet Pet) {
	petsLock.Lock()
	defer petsLock.Unlock()
	pet.ID = newPetID()
	Pets = append(Pets, pet)
}

func RemovePet(id int64) {
	petsLock.Lock()
	defer petsLock.Unlock()
	var newPets []Pet
	for _, pet := range Pets {
		if pet.ID != id {
			newPets = append(newPets, pet)
		}
	}
	Pets = newPets
}

func PetByID(id int64) (*Pet, error) {
	for _, pet := range Pets {
		if pet.ID == id {
			return &pet, nil
		}
	}
	return nil, fmt.Errorf("not found: pet %d", id)
}
