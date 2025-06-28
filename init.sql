\connect hiroyuki_diet


CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TYPE gender AS ENUM ('man', 'woman');
CREATE TYPE meal_type AS ENUM ('breakfast', 'lunch', 'dinner', 'snacking');
CREATE TYPE skin_part AS ENUM ('head', 'face', 'body');
CREATE TYPE field AS ENUM ('login', 'signin', 'home', 'meal', 'meal_form', 'meal_edit', 'data', 'profile', 'profile_edit', 'exercise', 'exercise_complete', 'achievement', 'achievement_complete', 'chibi_hiroyuki');


