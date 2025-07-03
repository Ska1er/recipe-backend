CREATE TABLE ingredients (
	id 					serial PRIMARY KEY, 
	name 				varchar(100) NOT NULL,
	measure 		varchar(50) NOT NULL,
	created_at 	timestamp with time zone DEFAULT NOW(),
	updated_at 	timestamp with time zone DEFAULT NOW()
);

CREATE TABLE recipes (
	id 						serial PRIMARY KEY, 
	name 					varchar(100) NOT NULL,
	description 	text NULL,
	image					varchar(200) NOT NULL,
	steps					jsonb NOT NULL,
	is_custom 		boolean NOT NULL DEFAULT false,
	is_visible		boolean NOT NULL DEFAULT false,
	cooking_time	int NOT NULL,
	difficulty		varchar(100) NOT NULL,
	created_at 		timestamp with time zone DEFAULT NOW(),
	updated_at 		timestamp with time zone DEFAULT NOW()
);

CREATE TABLE recipe_ingredients (
	id						serial PRIMARY KEY, 
	sort					integer NOT NULL,
	value 				numeric NOT NULL,
	recipe_id 		integer NOT NULL,
	ingredient_id integer NOT NULL,

	CONSTRAINT FK_Recipe FOREIGN KEY (recipe_id) REFERENCES recipes(id),
	CONSTRAINT FK_Ingredient FOREIGN KEY (ingredient_id) REFERENCES ingredients(id)
);
