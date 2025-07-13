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

	// 1. Master Data Independen
	seedDepartment(db)
	seedPhysicalLocation(db)
	seedPosition(db)
	seedSectionStatusTickets(db)
	seedWorkflow(db)

	// 2. Master Data Dependen
	seedAreas(db)
	seedSpecifiedLocation(db)
	seedStatusTicket(db)
	seedWorkflowStep(db)

	// 3. Main Data
	seedPositionToWorkflowMapping(db)
	seedEmployees(db)

	log.Println("Finish Seeding!")
}

func truncateTables(db *sql.DB) {
	log.Println("Truncating all tables...")

	tables := []string{
		"job",
		"track_status_ticket",
		"ticket",
		"workflow_step",
		"position_to_workflow_mapping",
		"employee",
		"status_ticket",
		"area",
		"department",
		"specified_location",
		"physical_location",
		"workflow",
		"section_status_ticket",
		"position",
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

// 1. Master Data Independen

// 1.1.
func seedDepartment(db *sql.DB) {
	log.Println("Seed Department")

	departments := []struct {
		name        string
		receive_job bool
	}{
		{"HRGA", true},        // 1
		{"Maintenance", true}, // 2
		{"Quality", true},     // 3
		{"PE", true},          // 4
		{"Office", false},     // 5
		{"Marketing", false},  // 6
		{"Finance", false},    // 7
		{"Operation", false},  // 8
	}

	for _, d := range departments {
		_, err := db.Exec("INSERT INTO department (name, receive_job, is_active) VALUES ($1, $2, true) ON CONFLICT(name) DO NOTHING", d.name, d.receive_job)
		if err != nil {
			log.Fatalf("Failed to insert department (%v): %v", d.name, err)
		}
	}

	log.Println("Finish insert data on Department")
}

// 1.2.
func seedPhysicalLocation(db *sql.DB) {
	log.Println("Seed Physical Location")

	physicalLocation := []struct {
		name      string
		is_active bool
	}{
		{"Forging", true},         // 1
		{"Production", true},      // 2
		{"Log", true},             // 3
		{"Building Office", true}, // 4
	}

	for _, p := range physicalLocation {
		_, err := db.Exec("INSERT INTO physical_location (name, is_active) VALUES ($1, $2) ON CONFLICT(name) DO NOTHING", p.name, p.is_active)
		if err != nil {
			log.Fatalf("Failed to insert Physical Location (%v): %v", p.name, err)
		}
	}

	log.Println("Finish insert data on Physical Location")
}

// 1.3.
func seedPosition(db *sql.DB) {
	log.Println("Seeding position...")
	positions := []string{"Department", "Section", "Frontman", "Leader"}
	for _, p := range positions {
		_, err := db.Exec("INSERT INTO position (name, is_active) VALUES ($1, true) ON CONFLICT(name) DO NOTHING", p)
		if err != nil {
			log.Fatalf("Failed to seed position: %v", err)
		}
	}
}

// 1.4.
func seedSectionStatusTickets(db *sql.DB) {
	log.Println("Seeding section_status_ticket...")

	sections := []struct {
		name     string
		sequence int
	}{
		{"Delete Section", 1},
		{"Approval Section", 2},
		{"Actual Section", 3},
	}

	for _, s := range sections {
		_, err := db.Exec("INSERT INTO section_status_ticket (name, sequence, is_active) VALUES ($1, $2, true) ON CONFLICT(name) DO NOTHING", s.name, s.sequence)
		if err != nil {
			log.Fatalf("Failed to seed section_status_ticket: %v", err)
		}
	}
}

// 1.5.
func seedWorkflow(db *sql.DB) {
	log.Println("Seeding workflow...")
	workflows := []string{"Direct to Job Workflow", "Department Approval Workflow", "Full Approval Workflow"}
	for _, w := range workflows {
		_, err := db.Exec("INSERT INTO workflow (name, is_active) VALUES ($1, true) ON CONFLICT(name) DO NOTHING", w)
		if err != nil {
			log.Fatalf("Failed to seed workflow: %v", err)
		}
	}
}

// 2. Master Data Dependen

// 2.1.
func seedAreas(db *sql.DB) {
	log.Println("Seed Area")

	areas := []struct {
		department_id int
		name          string
	}{
		{1, "Building"},
		{1, "Electrical"},
		{1, "Office"},
		{2, "Maintenance 1"},
		{2, "Maintenance 2"},
		{2, "Maintenance Support"},
		{3, "Pengukuran"},
		{3, "Pengujian"},
	}

	for _, a := range areas {
		_, err := db.Exec("INSERT INTO area (department_id, name, is_active) VALUES ($1, $2, true) ON CONFLICT(department_id, name) DO NOTHING", a.department_id, a.name)
		if err != nil {
			log.Fatalf("Failed to insert area (%v): %v", a.name, err)
		}
	}

	log.Println("Finish insert data on Area")
}

// 2.2.
func seedSpecifiedLocation(db *sql.DB) {
	log.Println("Seed Specified Location")

	specifiedLocation := []struct {
		physical_location_id int
		name                 string
	}{
		{1, "F 1"},
		{1, "F 2"},
		{1, "F 3"},
		{1, "F 4"},
		{2, "Machine Production"},
		{2, "Tool Production"},
		{2, "Support Production"},
		{3, "Input Log"},
		{3, "Output Log"},
	}

	for _, s := range specifiedLocation {
		_, err := db.Exec("INSERT INTO specified_location (physical_location_id, name, is_active) VALUES ($1, $2, true) ON CONFLICT(physical_location_id, name) DO NOTHING", s.physical_location_id, s.name)
		if err != nil {
			log.Fatalf("Failed to insert Specified Location (%v): %v", s.name, err)
		}
	}

	log.Println("Finish insert data on Specified Location")
}

// 2.3.
func seedStatusTicket(db *sql.DB) {
	log.Println("Seed Status Ticket")

	statuses := []struct {
		name       string
		section_id int
		sequence   int
	}{
		// DELETE SECTION
		{"Dibatalkan", 1, 0},

		// APPROVAL SECTION
		{"Approval Section", 2, 1},
		{"Approval Department", 2, 2},

		// ACTUAL SECTION
		{"Menunggu Job", 3, 3},
		{"Dikerjakan", 3, 4},
		{"Job Selesai", 3, 5},
		{"Tiket selesai", 3, 6},
	}

	for _, s := range statuses {
		_, err := db.Exec("INSERT INTO status_ticket (name, section_id, sequence, is_active) VALUES ($1, $2, $3, true) ON CONFLICT ON CONSTRAINT status_ticket_unique_name DO NOTHING", s.name, s.section_id, s.sequence)
		if err != nil {
			log.Fatalf("Failed to insert status ticket (%v): %v", s.name, err)
		}
	}
	log.Println("Finish to insert data on Status Ticket")
}

// 2.4.
func seedWorkflowStep(db *sql.DB) {
	log.Println("Seeding workflow_step...")

	steps := []struct {
		workflow_id      int
		status_ticket_id int
		step_sequence    int
	}{
		// Full Approval Workflow (ID=3)
		{3, 2, 0}, // Approval Section
		{3, 3, 1}, // Approval Department
		{3, 4, 2}, // Menunggu Job

		// Department Approval Workflow (ID=2)
		{2, 3, 0}, // Approval Department
		{2, 4, 1}, // Menunggu Job

		// Direct to Job Workflow (ID=1)
		{1, 4, 0}, // Menunggu Job
	}
	for _, s := range steps {
		_, err := db.Exec(`
            INSERT INTO workflow_step (workflow_id, status_ticket_id, step_sequence, is_active) VALUES ($1, $2, $3, true)
            ON CONFLICT (workflow_id, step_sequence) DO NOTHING`, s.workflow_id, s.status_ticket_id, s.step_sequence)
		if err != nil {
			log.Fatalf("Failed to seed workflow_step: %v", err)
		}
	}
}

// 3. Main Data

// 3.1.
func seedPositionToWorkflowMapping(db *sql.DB) {
	log.Println("Seeding position_to_workflow_mapping...")
	mappings := []struct {
		position_id int
		workflow_id int
	}{
		{1, 1}, // Department -> Direct
		{2, 2}, // Section -> Dept Approval
		{3, 3}, // Frontman -> Full Approval
		{4, 3}, // Leader -> Full Approval
	}
	for _, m := range mappings {
		_, err := db.Exec("INSERT INTO position_to_workflow_mapping (position_id, workflow_id) VALUES ($1, $2) ON CONFLICT(position_id) DO NOTHING", m.position_id, m.workflow_id)
		if err != nil {
			log.Fatalf("Failed to seed position_to_workflow_mapping: %v", err)
		}
	}
}

// 3.2.
func seedEmployees(db *sql.DB) {
	log.Println("Seed Employee")

	firstNames := []string{"Adi", "Budi", "Cahyo", "Deni", "Eka", "Fajar", "Gita", "Hadi", "Indra", "Joko", "Kartika", "Lia", "Mega", "Nadia", "Oscar", "Putra", "Rina", "Sari", "Tono", "Wati"}
	lastNames := []string{"Wijaya", "Susanto", "Pratama", "Kusumo", "Lestari", "Nugroho", "Wahyuni", "Setiawan", "Hidayat", "Purnama"}
	npkCounter := 1
	rand.Seed(time.Now().UnixNano())

	deptIDs := []int{1, 2, 3, 4, 5, 6, 7, 8}
	areaIDsByDept := map[int][]int{
		1: {1, 2, 3},
		2: {4, 5, 6},
		3: {7, 8},
	}
	positionIDs := []int{1, 2, 3, 4}

	for _, dept := range deptIDs {
		log.Printf("... Seeding employees for Department: %v", dept)

		// --- Kelompok A: Karyawan dengan Department, Area, dan Position ---
		if areaIDs, ok := areaIDsByDept[dept]; ok {

			// Buat 1 Head of Department
			createEmployee(db, &npkCounter, firstNames, lastNames, positionIDs[0], dept, randomAreaID(areaIDs))

			// Buat 2 Section
			for i := 0; i < 2; i++ {
				createEmployee(db, &npkCounter, firstNames, lastNames, positionIDs[1], dept, randomAreaID(areaIDs))
			}

			// Buat 2-3 Frontmant
			for i := 0; i < rand.Intn(2)+2; i++ {
				createEmployee(db, &npkCounter, firstNames, lastNames, positionIDs[2], dept, randomAreaID(areaIDs))
			}

			// Buat 5-10 Leader
			for i := 0; i < rand.Intn(6)+5; i++ {
				createEmployee(db, &npkCounter, firstNames, lastNames, positionIDs[3], dept, randomAreaID(areaIDs))
			}
		}

		// --- Kelompok B: Karyawan dengan Department dan Position (tanpa Area) ---
		if _, ok := areaIDsByDept[dept]; !ok {

			// Buat 1 Head of Department
			createEmployee(db, &npkCounter, firstNames, lastNames, positionIDs[0], dept, sql.NullInt64{})

			// Buat 2 Section
			for i := 0; i < 2; i++ {
				createEmployee(db, &npkCounter, firstNames, lastNames, positionIDs[1], dept, sql.NullInt64{})
			}

			// Buat 2-3 Frontman
			for i := 0; i < rand.Intn(2)+2; i++ {
				createEmployee(db, &npkCounter, firstNames, lastNames, positionIDs[2], dept, sql.NullInt64{})
			}

			// Buat 5-10 Leader
			for i := 0; i < rand.Intn(6)+5; i++ {
				createEmployee(db, &npkCounter, firstNames, lastNames, positionIDs[3], dept, sql.NullInt64{})
			}
		}
	}
	log.Println("Finish insert data on Employee")
}

func createEmployee(db *sql.DB, npkCounter *int, firstNames, lastNames []string, position int, deptID int, areaID sql.NullInt64) {
	npk := fmt.Sprintf("EMP%04d", *npkCounter)
	*npkCounter++

	fullName := firstNames[rand.Intn(len(firstNames))] + " " + lastNames[rand.Intn(len(lastNames))]

	_, err := db.Exec(`
        INSERT INTO employee (npk, department_id, area_id, name, position_id, is_active)
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
