package utils

import (
	"encoding/json"
	"flag"
	"fmt"
	"hot-coffee/config"
	"os"
)

func DirectoryExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func Help(s []string) {
	for _, v := range s {
		if v == "--help" || v == "-help" || v == "-h" {
			fmt.Println(`Coffee Shop Management System

Usage:
	hot-coffee [--port <N>] [--dir <S>] 
	hot-coffee --help
			
Options:
	--help       Show this screen.
	--port N     Port number.
	--dir S      Path to the data directory.`)
			os.Exit(0)
		}
	}
}

func CreateDir(path string) error {
	// Ensure the directory exists
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return err
	}

	inventoryPath := fmt.Sprintf("%s/inventory.json", path)
	menuItemsPath := fmt.Sprintf("%s/menu_items.json", path)
	ordersPath := fmt.Sprintf("%s/orders.json", path)
	configPath := fmt.Sprintf("%s/config.json", path)

	// Create or initialize inventory.json
	if _, err := os.Stat(inventoryPath); os.IsNotExist(err) {
		err := saveJson(config.DefaultInventoryData, inventoryPath)
		if err != nil {
			return err
		}
	}

	// Create or initialize menu_items.json
	if _, err := os.Stat(menuItemsPath); os.IsNotExist(err) {
		err := saveJson(config.DefaultMenuItemsData, menuItemsPath)
		if err != nil {
			return err
		}
	}

	// Create or initialize orders.json
	if _, err := os.Stat(ordersPath); os.IsNotExist(err) {
		err := saveJson(config.DefaultOrdersData, ordersPath)
		if err != nil {
			return err
		}
	}

	// Create or initialize config.json
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		err := saveJson(config.DefaultConfigData, configPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func saveJson(items any, path string) error {
	jsonData, err := json.MarshalIndent(items, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, jsonData, 0o644)
}

func Flagchecker() (string, string) {
	dir := flag.String("dir", "data", "Path to the data directory")
	port := flag.String("port", "8080", "Port Number")
	flag.Parse()

	config.BaseDir = *dir
	return *dir, *port
}
