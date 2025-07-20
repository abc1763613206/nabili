package migration

import (
	"log"

	"github.com/spf13/viper"
	"github.com/abc1763613206/nabili/internal/db"
)

// AutoMigrateRemoteSources automatically adds missing remote sources to existing configurations
func AutoMigrateRemoteSources() {
	currentDBs := db.List{}
	err := viper.UnmarshalKey("databases", &currentDBs)
	if err != nil {
		log.Printf("AutoMigration: failed to unmarshal databases: %v", err)
		return
	}

	// Define all remote sources
	remoteSources := map[string]*db.DB{
		"bili": {
			Name:      "bili",
			Format:    db.FormatRemote,
			File:      "",
			Languages: []string{"zh-CN"},
			Types:     []db.Type{db.TypeIPv4, db.TypeIPv6},
		},
		"ipsb": {
			Name:      "ipsb",
			Format:    db.FormatRemote,
			File:      "",
			Languages: []string{"en"},
			Types:     []db.Type{db.TypeIPv4, db.TypeIPv6},
		},
		"iqiyi": {
			Name:      "iqiyi",
			Format:    db.FormatRemote,
			File:      "",
			Languages: []string{"zh-CN"},
			Types:     []db.Type{db.TypeIPv4, db.TypeIPv6},
		},
		"baidu": {
			Name:      "baidu",
			Format:    db.FormatRemote,
			File:      "",
			Languages: []string{"zh-CN"},
			Types:     []db.Type{db.TypeIPv4, db.TypeIPv6},
		},
	}

	// Track existing database names
	existingNames := make(map[string]bool)
	for _, db := range currentDBs {
		existingNames[db.Name] = true
	}

	// Add missing remote sources
	added := false
	for name, remoteDB := range remoteSources {
		if !existingNames[name] {
			currentDBs = append(currentDBs, remoteDB)
			added = true
			log.Printf("AutoMigration: Added %s remote source to configuration", name)
		}
	}

	if added {
		viper.Set("databases", currentDBs)
		
		// Save updated configuration
		if err := viper.WriteConfig(); err != nil {
			log.Printf("AutoMigration: failed to write config: %v", err)
			return
		}
	}
}