CREATE TABLE icecream(
    name  varchar NOT NULL,
    image_closed varchar NOT NULL,
    image_open varchar NOT NULL,
    description varchar,
    story varchar,
    sourcing_values varchar[],
    ingredients varchar[],
    allergy_info varchar,
    dietary_certifications varchar,
    product_id varchar NOT NULL
);