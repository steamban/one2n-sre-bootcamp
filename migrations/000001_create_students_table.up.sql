CREATE TABLE IF NOT EXISTS students (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    first_name TEXT NOT NULL CHECK(length(first_name) <= 50),
    last_name TEXT NOT NULL CHECK(length(last_name) <= 50),
    age INTEGER NOT NULL CHECK(age > 0 AND age < 150),
    gender TEXT NOT NULL CHECK(gender IN ('Male', 'Female', 'Other')),
    email TEXT NOT NULL UNIQUE CHECK(length(email) <= 255),
    phone TEXT NOT NULL CHECK(length(phone) <= 15),
    class TEXT NOT NULL CHECK(length(class) <= 20),
    rank TEXT CHECK(length(rank) == 1 AND rank IN ('A', 'B', 'C', 'D', 'E', 'F')),
    address_line1 TEXT NOT NULL CHECK(length(address_line1) <= 100),
    address_line2 TEXT CHECK(length(address_line2) <= 100),
    city TEXT NOT NULL CHECK(length(city) <= 50),
    state TEXT NOT NULL CHECK(length(state) <= 50),
    pincode TEXT NOT NULL CHECK(length(pincode) == 6),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME
);

-- Trigger to update updated_at on changes
CREATE TRIGGER IF NOT EXISTS update_student_updated_at 
AFTER UPDATE ON students
FOR EACH ROW
BEGIN
    UPDATE students SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
END;
