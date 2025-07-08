package config

import (
	"case_study_api/entities"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&entities.Ticket{},
		&entities.Event{},
		&entities.User{},
	)
}

func SeedData(db *gorm.DB) error {
	log.Println("Starting database seeding...")

	var count int64
	db.Model(&entities.User{}).Where("role = ?", "admin").Count(&count)
	if count > 0 {
		return fmt.Errorf("failed to create admin user")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Failed to hash admin password")
		return fmt.Errorf("failed to create admin user: %v", err)
	}

	admin := entities.User{
		Name:     "Super Admin",
		Email:    "admin@system.com",
		Password: string(hashed),
		Role:     "admin",
	}
	if err := db.Create(&admin).Error; err != nil {
		log.Printf("Failed to create admin user: %v", err)
		return fmt.Errorf("failed to create admin user: %v", err)
	} else {
		log.Println("✅ Admin user seeded: admin@system.com / admin123")
	}

	// Seed Events
	log.Println("Seeding events...")
	var eventCount int64
	db.Model(&entities.Event{}).Count(&eventCount)
	if eventCount == 0 {
		events := []entities.Event{
			{
				Title:       "Tech Innovation Summit 2025",
				Description: "Join industry leaders for the latest in technology and innovation. Network with professionals and discover cutting-edge solutions.",
				Location:    "Jakarta Convention Center",
				Category:    "Technology",
				Status:      "upcoming",
				Date:        time.Now().AddDate(0, 0, 30),
				EndDate:     time.Now().AddDate(0, 0, 31),
				Capacity:    500,
				Price:       250000,
				SoldTickets: 0,
				CreatedBy:   admin.ID,
				IsActive:    true,
			},
			{
				Title:       "Digital Marketing Masterclass",
				Description: "Learn the latest digital marketing strategies and tools from industry experts. Hands-on workshops included.",
				Location:    "Bali International Convention Centre",
				Category:    "Marketing",
				Status:      "upcoming",
				Date:        time.Now().AddDate(0, 0, 45),
				EndDate:     time.Now().AddDate(0, 0, 45),
				Capacity:    200,
				Price:       150000,
				SoldTickets: 0,
				CreatedBy:   admin.ID,
				IsActive:    true,
			},
			{
				Title:       "Startup Funding Workshop",
				Description: "Connect with investors and learn how to secure funding for your startup. Pitch sessions and networking opportunities.",
				Location:    "Surabaya Business Center",
				Category:    "Business",
				Status:      "upcoming",
				Date:        time.Now().AddDate(0, 0, 60),
				EndDate:     time.Now().AddDate(0, 0, 60),
				Capacity:    100,
				Price:       200000,
				SoldTickets: 0,
				CreatedBy:   admin.ID,
				IsActive:    true,
			},
			{
				Title:       "AI & Machine Learning Conference",
				Description: "Explore the future of artificial intelligence and machine learning. Expert speakers and hands-on demos.",
				Location:    "Bandung Tech Park",
				Category:    "Technology",
				Status:      "upcoming",
				Date:        time.Now().AddDate(0, 0, 75),
				EndDate:     time.Now().AddDate(0, 0, 76),
				Capacity:    300,
				Price:       300000,
				SoldTickets: 0,
				CreatedBy:   admin.ID,
				IsActive:    true,
			},
			{
				Title:       "E-commerce Excellence Summit",
				Description: "Discover strategies to boost your online business. Learn from successful e-commerce entrepreneurs.",
				Location:    "Yogyakarta Convention Hall",
				Category:    "Business",
				Status:      "upcoming",
				Date:        time.Now().AddDate(0, 0, 90),
				EndDate:     time.Now().AddDate(0, 0, 90),
				Capacity:    250,
				Price:       175000,
				SoldTickets: 0,
				CreatedBy:   admin.ID,
				IsActive:    true,
			},
			{
				Title:       "Creative Design Workshop",
				Description: "Unleash your creativity with the latest design tools and techniques. Perfect for designers and artists.",
				Location:    "Medan Creative Center",
				Category:    "Design",
				Status:      "upcoming",
				Date:        time.Now().AddDate(0, 0, 105),
				EndDate:     time.Now().AddDate(0, 0, 105),
				Capacity:    150,
				Price:       125000,
				SoldTickets: 0,
				CreatedBy:   admin.ID,
				IsActive:    true,
			},
			{
				Title:       "Data Science Bootcamp",
				Description: "Intensive bootcamp covering data analysis, visualization, and machine learning. Hands-on projects included.",
				Location:    "Semarang University Campus",
				Category:    "Technology",
				Status:      "upcoming",
				Date:        time.Now().AddDate(0, 0, 120),
				EndDate:     time.Now().AddDate(0, 0, 122),
				Capacity:    80,
				Price:       400000,
				SoldTickets: 0,
				CreatedBy:   admin.ID,
				IsActive:    true,
			},
			{
				Title:       "Mobile App Development Summit",
				Description: "Learn to build stunning mobile applications for iOS and Android. Expert-led workshops and networking.",
				Location:    "Makassar Tech Hub",
				Category:    "Technology",
				Status:      "upcoming",
				Date:        time.Now().AddDate(0, 0, 135),
				EndDate:     time.Now().AddDate(0, 0, 136),
				Capacity:    180,
				Price:       275000,
				SoldTickets: 0,
				CreatedBy:   admin.ID,
				IsActive:    true,
			},
			{
				Title:       "Blockchain & Cryptocurrency Forum",
				Description: "Explore the world of blockchain technology and cryptocurrency. Investment strategies and technical deep-dives.",
				Location:    "Denpasar Financial District",
				Category:    "Finance",
				Status:      "upcoming",
				Date:        time.Now().AddDate(0, 0, 150),
				EndDate:     time.Now().AddDate(0, 0, 150),
				Capacity:    220,
				Price:       350000,
				SoldTickets: 0,
				CreatedBy:   admin.ID,
				IsActive:    true,
			},
			{
				Title:       "Leadership Excellence Conference",
				Description: "Develop your leadership skills with renowned speakers and interactive workshops. Perfect for managers and executives.",
				Location:    "Palembang Business Center",
				Category:    "Leadership",
				Status:      "upcoming",
				Date:        time.Now().AddDate(0, 0, 165),
				EndDate:     time.Now().AddDate(0, 0, 165),
				Capacity:    300,
				Price:       225000,
				SoldTickets: 0,
				CreatedBy:   admin.ID,
				IsActive:    true,
			},
		}

		for _, event := range events {
			if err := db.Create(&event).Error; err != nil {
				log.Printf("Failed to create event '%s': %v", event.Title, err)
			} else {
				log.Printf("✅ Event seeded: %s", event.Title)
			}
		}
		log.Println("✅ All events seeded successfully!")
	} else {
		log.Println("Events already exist, skipping seeding...")
	}

	return nil
}

func ResetDatabase(db *gorm.DB) error {
	log.Println("Starting database reset...")

	// Drop all tables
	if err := db.Migrator().DropTable(
		&entities.Ticket{},
		&entities.Event{},
		&entities.User{},
	); err != nil {
		return fmt.Errorf("failed to drop tables: %v", err)
	}
	log.Println("✓ All tables dropped")

	// Recreate tables
	if err := AutoMigrate(db); err != nil {
		return fmt.Errorf("failed to migrate tables: %v", err)
	}
	log.Println("✓ Tables recreated")

	// Seed fresh data
	if err := SeedData(db); err != nil {
		return fmt.Errorf("failed to seed data: %v", err)
	}

	log.Println("Database reset completed successfully!")
	return nil
}
