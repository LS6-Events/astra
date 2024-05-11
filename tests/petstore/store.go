package petstore

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

var Pets = []Pet{
	{ID: "a0652c3a-142f-438a-a553-381a7c135d9b", Name: "Dog", PhotoURLs: []string{}, Status: "available", Tags: nil},
	{ID: "95188ede-b9ab-44ba-95b1-71b5d5c31f7f", Name: "Cat", PhotoURLs: []string{}, Status: "pending", Tags: nil},
}

var petsLock = &sync.Mutex{}

func newPetID() string {
	return uuid.New().String()
}

func AddPet(pet Pet) {
	petsLock.Lock()
	defer petsLock.Unlock()
	pet.ID = newPetID()
	Pets = append(Pets, pet)
}

func RemovePet(id string) {
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

func PetByID(id string) (*Pet, error) {
	for _, pet := range Pets {
		if pet.ID == id {
			return &pet, nil
		}
	}
	return nil, fmt.Errorf("not found: pet %s", id)
}
