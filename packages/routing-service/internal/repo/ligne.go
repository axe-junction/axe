	package repo

	import (
		"context"
		"log"

		"github.com/axe-junction/axe-server/internal/models"
		"github.com/google/uuid"
		"gorm.io/gorm"
	)

	type LigneRepo struct {
		db *gorm.DB
	}

	func NewLigneRepo(db *gorm.DB) models.LigneRepo {
		return &LigneRepo{db: db}
	}

	func (repo *LigneRepo) GetLignes() (models.Lignes, error) {
		var lignes models.Lignes
		if err := repo.db.Preload("Stations").Find(&lignes).Error; err != nil {
			return nil, err
		}
		return lignes, nil
	}

	func (repo *LigneRepo) GetByID(id uuid.UUID) (models.Ligne, error) {
		var ligne models.Ligne
		if err := repo.db.Preload("Stations").First(&ligne, "id = ?", id).Error; err != nil {
			return models.Ligne{}, err
		}
		return ligne, nil
	}

	func (repo *LigneRepo) GetByType(typee string) (models.Lignes, error) {
		var lignes models.Lignes
		if err := repo.db.Preload("Stations").Where("type = ?", typee).Find(&lignes).Error; err != nil {
			return nil, err
		}
		return lignes, nil
	}
	func (repo *LigneRepo) GetNearby(latitude, longitude float64) (models.Lignes, error) {
		radius := 20
		var lignes models.Lignes
		for i := radius; i > 0; i++ {
			if err := repo.db.Preload("Stations").Where("ST_DWithin(ST_MakePoint(longitude, latitude)::geography, ST_MakePoint(?, ?)::geography, ?)", longitude, latitude, radius).Find(&lignes).Error; err != nil {
				return nil, err
			}
			if len(lignes) == 0 {
				radius += 100
				continue
			}
			break
		}
		return lignes, nil
	}
	func (repo *LigneRepo) GetRouteBetweenStations(ctx context.Context, startID, endID uuid.UUID) (models.Lignes, error) {
		log.Printf("🔍 SQL: Searching routes between station %s and %s", startID, endID)

		var lignes models.Lignes

		// Find routes that have both start and end stations
		query := `
			SELECT DISTINCT l.* FROM lignes l
			JOIN stops s1 ON l.id = s1.route_id AND s1.station_id = ?
			JOIN stops s2 ON l.id = s2.route_id AND s2.station_id = ?
			WHERE s1.sequence < s2.sequence
		`

		log.Printf("🔍 Executing query: %s", query)
		log.Printf("🔍 Parameters: startID=%s, endID=%s", startID, endID)

		if err := repo.db.Raw(query, startID, endID).Scan(&lignes).Error; err != nil {
			log.Printf("❌ SQL Error: %v", err)
			return nil, err
		}

		log.Printf("📊 Found %d routes from raw query", len(lignes))

		// For each route, load only the stations between start and end (inclusive)
		for i := range lignes {
			var stations []models.Station

			stationQuery := `
				SELECT s.* FROM stations s
				JOIN stops st ON s.id = st.station_id
				WHERE st.route_id = ? 
				AND st.sequence >= (SELECT sequence FROM stops WHERE route_id = ? AND station_id = ?)
				AND st.sequence <= (SELECT sequence FROM stops WHERE route_id = ? AND station_id = ?)
				ORDER BY st.sequence
			`

			err := repo.db.Raw(stationQuery, lignes[i].ID, lignes[i].ID, startID, lignes[i].ID, endID).Scan(&stations).Error
			if err != nil {
				log.Printf("❌ Error loading stations between %s and %s for route %s: %v", startID, endID, lignes[i].ID, err)
				continue
			}

			lignes[i].Stations = stations
			log.Printf("✅ Loaded %d stations between start and end for route %s (%s)", len(stations), lignes[i].ID, lignes[i].Name)

			// Debug: show the station sequence
			for j, station := range stations {
				log.Printf("   Station %d: %s", j+1, station.Name)
			}
		}

		return lignes, nil
	}

	var _ models.LigneRepo = (*LigneRepo)(nil)
