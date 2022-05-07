package main

import (
	"encoding/json"
	"log"
	"os"
)

const dataFile = "data/data.json"

type FuckEntry struct {
	Name  string `json:"username"`
	Total int    `json:"total"`
	Fuck  int    `json:"fuck"`
}

type FuckManager struct {
	Entries map[int64]*FuckEntry
}

func NewFuckManager() *FuckManager {
	return &FuckManager{Entries: map[int64]*FuckEntry{}}
}

func (manager *FuckManager) LoadData() error {
	if _, err := os.Stat(dataFile); err != nil {
		return nil
	}

	bytes, err := os.ReadFile(dataFile)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bytes, &manager.Entries); err != nil {
		return err
	}

	log.Printf("loaded %d fuck entries\n", len(manager.Entries))
	return nil
}

func (manager *FuckManager) SaveData() error {
	bytes, err := json.Marshal(manager.Entries)
	if err != nil {
		return err
	}

	if err := os.WriteFile(dataFile, bytes, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func (manager *FuckManager) AddMessage(userId int64, name string, fuck int) {
	if _, exists := manager.Entries[userId]; !exists {
		manager.Entries[userId] = &FuckEntry{name, 0, 0}
	}

	manager.Entries[userId].Name = name
	manager.Entries[userId].Total++
	manager.Entries[userId].Fuck += fuck
}
