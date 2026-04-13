ALTER TABLE students 
ADD CONSTRAINT students_class_check 
CHECK (class IN ('10th', '11th', '12th'));
