package config

import (
	"github.com/fsnotify/fsnotify"
	"log"
)

type ConfigWatcher struct {
	cfgManager *ConfigManager
}

func NewConfigWatcher(cfgManager *ConfigManager) *ConfigWatcher {
	return &ConfigWatcher{cfgManager: cfgManager}
}

func (cfgWatcher ConfigWatcher) WatchGateWayConfig(path string) {
	watcher, _ := fsnotify.NewWatcher()
	defer watcher.Close()

	watcher.Add(path)

	for {
		select {
		case event := <-watcher.Events:
			if event.Has(fsnotify.Write) {
				log.Printf("Event: %s %s", event.Op, event.Name)
				cfg, err := LoadGateWayConfig(path)
				if err != nil {
					log.Println("Error Loading Config : ", path)
				}
				cfgWatcher.cfgManager.UpdateGateWayConfig(cfg)
			}

		case err := <-watcher.Errors:
			log.Println(err)
		}
	}
}

func (cfgWatcher ConfigWatcher) WatchApplicationConfig(path string) {
	watcher, _ := fsnotify.NewWatcher()
	defer watcher.Close()

	watcher.Add(path)

	for {
		select {
		case event := <-watcher.Events:
			if event.Has(fsnotify.Write) {
				log.Printf("Event: %s %s", event.Op, event.Name)
				cfg, err := LoadApplicationConfig()
				if err != nil {
					log.Println("Error Loading Config : ", path)
				}
				cfgWatcher.cfgManager.UpdateApplicationConfig(cfg)
			}

		case err := <-watcher.Errors:
			log.Println(err)
		}
	}
}