package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:QWERTY12345@localhost:5431/job_reservation_db?sslmode=disable"
	}

	db, err := sql.Open("pgx", dsn)

	if err != nil {
		log.Fatalf("Error connection to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed ping to database: %v", err)
	}

	log.Println("Successful connect to database")

	truncateTables(db)

	seedStatusTickets(db)
	seedDepartments(db)
	seedAreas(db)
	seedPhysicalLocation(db)
	seedSpecifiedLocation(db)
	seedEmployees(db)

	log.Println("Finish Seeding!")
}

func truncateTables(db *sql.DB) {
	log.Println("Truncating all tables...")

	tables := []string{
		"rejected_ticket",
		"job",
		"track_status_ticket",
		"ticket",
		"employee",
		"specified_location",
		"physical_location",
		"area",
		"department",
		"status_ticket",
	}

	for _, table := range tables {
		query := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE;", table)
		_, err := db.Exec(query)
		if err != nil {
			log.Fatalf("Failed to truncate table %s: %v", table, err)
		}
	}

	log.Println("All tables truncated successfully.")
}

func seedStatusTickets(db *sql.DB) {
	log.Println("Seed Status Ticket")

	statuses := []struct {
		name      string
		sequence  int
		is_active bool
	}{
		{"Dibatalkan", -100, true},
		{"Approval Section", -2, true},
		{"Approval Department", -1, true},
		{"Menunggu Job", 0, true},
		{"Dikerjakan", 1, true},
		{"Job Selesai", 2, true},
		{"Tiket selesai", 3, true},
	}

	for _, s := range statuses {
		_, err := db.Exec("INSERT INTO status_ticket (name, sequence, is_active) VALUES ($1, $2, $3) ON CONFLICT(name, sequence) DO NOTHING", s.name, s.sequence, s.is_active)
		if err != nil {
			log.Fatalf("Failed to insert status ticket (%v): %v", s.name, err)
		}
	}
	log.Println("Finish to insert data on Status Ticket")
}

func seedDepartments(db *sql.DB) {
	log.Println("Seed Department")

	departments := []struct {
		name        string
		receive_job bool
		is_active   bool
	}{
		{"HRGA", true, true},        // 1
		{"Maintenance", true, true}, // 2
		{"Quality", true, true},     // 3
		{"PE", true, true},          // 4
		{"Office", false, true},     // 5
		{"Marketing", false, true},  // 6
		{"Finance", false, true},    // 7
		{"Operation", false, true},  // 8
	}

	for _, d := range departments {
		_, err := db.Exec("INSERT INTO department (name, receive_job, is_active) VALUES ($1, $2, $3) ON CONFLICT(name) DO NOTHING", d.name, d.receive_job, d.is_active)
		if err != nil {
			log.Fatalf("Failed to insert department (%v): %v", d.name, err)
		}
	}

	log.Println("Finish insert data on Department")
}

func seedAreas(db *sql.DB) {
	log.Println("Seed Area")

	areas := []struct {
		department_id int
		name          string
		is_active     bool
	}{
		{1, "Building", true},
		{1, "Electrical", true},
		{1, "Office", true},
		{2, "Maintenance 1", true},
		{2, "Maintenance 2", true},
		{2, "Maintenance Support", true},
		{3, "Pengukuran", true},
		{3, "Pengujian", true},
	}

	for _, a := range areas {
		_, err := db.Exec("INSERT INTO area (department_id, name, is_active) VALUES ($1, $2, $3) ON CONFLICT(department_id, name) DO NOTHING", a.department_id, a.name, a.is_active)
		if err != nil {
			log.Fatalf("Failed to insert area (%v): %v", a.name, err)
		}
	}

	log.Println("Finish insert data on Area")
}

func seedEmployees(db *sql.DB) {
	log.Println("Seed Employee")

	firstNames := []string{"Adi", "Budi", "Cahyo", "Deni", "Eka", "Fajar", "Gita", "Hadi", "Indra", "Joko", "Kartika", "Lia", "Mega", "Nadia", "Oscar", "Putra", "Rina", "Sari", "Tono", "Wati"}
	lastNames := []string{"Wijaya", "Susanto", "Pratama", "Kusumo", "Lestari", "Nugroho", "Wahyuni", "Setiawan", "Hidayat", "Purnama"}

	departments := []struct {
		id   int
		name string
	}{
		{1, "HRGA"}, {2, "Maintenance"}, {3, "Quality"}, {4, "PE"}, {5, "Office"}, {6, "Marketing"}, {7, "Finance"}, {8, "Operation"},
	}
	deptToAreaIDs := map[int][]int{
		1: {1, 2, 3},
		2: {4, 5, 6},
		3: {7, 8},
	}

	npkCounter := 1
	rand.Seed(time.Now().UnixNano())

	for _, dept := range departments {
		log.Printf("... Seeding employees for Department: %s", dept.name)

		// --- Kelompok A: Karyawan dengan Department, Area, dan Position ---
		if areaIDs, ok := deptToAreaIDs[dept.id]; ok {

			// Buat 1 Head of Department
			createEmployee(db, &npkCounter, firstNames, lastNames, "Head of Department", dept.id, randomAreaID(areaIDs))

			// Buat 2 Section
			for i := 0; i < 2; i++ {
				createEmployee(db, &npkCounter, firstNames, lastNames, "Section", dept.id, randomAreaID(areaIDs))
			}

			// Buat 2-3 Leader
			for i := 0; i < rand.Intn(2)+2; i++ {
				createEmployee(db, &npkCounter, firstNames, lastNames, "Leader", dept.id, randomAreaID(areaIDs))
			}

			// Buat 5-10 Staff
			for i := 0; i < rand.Intn(6)+5; i++ {
				createEmployee(db, &npkCounter, firstNames, lastNames, "Staff", dept.id, randomAreaID(areaIDs))
			}
		}

		// --- Kelompok B: Karyawan dengan Department dan Position (tanpa Area) ---
		if _, ok := deptToAreaIDs[dept.id]; !ok {

			// Buat 1 Head of Department
			createEmployee(db, &npkCounter, firstNames, lastNames, "Head of Department", dept.id, sql.NullInt64{}) // area_id = NULL

			// Buat 2 Section
			for i := 0; i < 2; i++ {
				createEmployee(db, &npkCounter, firstNames, lastNames, "Section", dept.id, sql.NullInt64{})
			}

			// Buat 5-10 Staff
			for i := 0; i < rand.Intn(6)+5; i++ {
				createEmployee(db, &npkCounter, firstNames, lastNames, "Staff", dept.id, sql.NullInt64{})
			}
		}
	}
	log.Println("Finish insert data on Employee")
}

func createEmployee(db *sql.DB, npkCounter *int, firstNames, lastNames []string, position string, deptID int, areaID sql.NullInt64) {
	npk := fmt.Sprintf("EMP%04d", *npkCounter)
	*npkCounter++

	fullName := firstNames[rand.Intn(len(firstNames))] + " " + lastNames[rand.Intn(len(lastNames))]

	_, err := db.Exec(`
        INSERT INTO employee (npk, department_id, area_id, name, position, is_active)
        VALUES ($1, $2, $3, $4, $5, true) ON CONFLICT(npk) DO NOTHING`,
		npk, deptID, areaID, fullName, position)

	if err != nil {
		log.Fatalf("Failed to insert employee (%s): %v", fullName, err)
	}
}

func randomAreaID(areaIDs []int) sql.NullInt64 {
	return sql.NullInt64{
		Int64: int64(areaIDs[rand.Intn(len(areaIDs))]),
		Valid: true,
	}
}

func seedPhysicalLocation(db *sql.DB) {
	log.Println("Seed Physical Location")

	physicalLocation := []struct {
		name      string
		is_active bool
	}{
		{"Forging", true},          // 1
		{"Production", true},       // 2
		{"Log", true},              // 3
		{"Building Office", false}, // 4
	}

	for _, p := range physicalLocation {
		_, err := db.Exec("INSERT INTO physical_location (name, is_active) VALUES ($1, $2) ON CONFLICT(name) DO NOTHING", p.name, p.is_active)
		if err != nil {
			log.Fatalf("Failed to insert Physical Location (%v): %v", p.name, err)
		}
	}

	log.Println("Finish insert data on Physical Location")
}

func seedSpecifiedLocation(db *sql.DB) {
	log.Println("Seed Specified Location")

	specifiedLocation := []struct {
		physical_location_id int
		name                 string
		is_active            bool
	}{
		{1, "F 1", true},
		{1, "F 2", true},
		{1, "F 3", true},
		{1, "F 4", true},
		{2, "Machine Production", true},
		{2, "Tool Production", true},
		{2, "Support Production", true},
		{3, "Input Log", true},
		{3, "Output Log", true},
	}

	for _, s := range specifiedLocation {
		_, err := db.Exec("INSERT INTO specified_location (physical_location_id, name, is_active) VALUES ($1, $2, $3) ON CONFLICT(physical_location_id, name) DO NOTHING", s.physical_location_id, s.name, s.is_active)
		if err != nil {
			log.Fatalf("Failed to insert Specified Location (%v): %v", s.name, err)
		}
	}

	log.Println("Finish insert data on Specified Location")
}
