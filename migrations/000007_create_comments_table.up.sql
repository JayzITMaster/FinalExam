CREATE TABLE IF NOT EXISTS comments (
  admin_id serial,
  form_id serial NOT NULL,
  comments1 text ,
  comments2 text ,
  comments3 text ,
  comments35 text ,
  comments4 text ,
  comments5 text ,
  comments6 text ,
  comments7 text ,
  comments8 text ,
  comments9 text ,
  comments10 text ,
  edit_made TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
);