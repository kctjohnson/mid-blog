CREATE TABLE
  blogger (
	  id INT NOT NULL AUTO_INCREMENT, 
	  first_name VARCHAR(255) NOT NULL, 
	  last_name VARCHAR(255) NOT NULL,
	  email VARCHAR(255) NOT NULL,
	  age VARCHAR(255) NOT NULL,
	  gender ENUM ('Male', 'Female') NOT NULL,
	  bio VARCHAR(255) NOT NULL,
	  create_date TIMESTAMP NOT NULL,
	  PRIMARY KEY (id)
  );

CREATE TABLE
  post (
	  id INT NOT NULL AUTO_INCREMENT, 
	  blogger_id INT NOT NULL,
	  title VARCHAR(255) NOT NULL,
	  content MEDIUMTEXT NOT NULL,
	  likes INT NOT NULL,
	  dislikes int NOT NULL,
	  create_date TIMESTAMP NOT NULL,
	  PRIMARY KEY (id),
	  CONSTRAINT fk_post_blogger FOREIGN KEY (blogger_id) REFERENCES blogger (id) ON DELETE CASCADE
  );

CREATE TABLE
  user (
	  id INT NOT NULL AUTO_INCREMENT,
	  username VARCHAR(255) NOT NULL,
	  password VARCHAR(255) NOT NULL,
	  create_date TIMESTAMP,
	  PRIMARY KEY (id)
  );

CREATE TABLE
comment (
	id INT NOT NULL AUTO_INCREMENT,
	user_id INT NOT NULL,
	post_id INT NOT NULL,
	content MEDIUMTEXT NOT NULL,
	likes INT NOT NULL,
	dislikes INT NOT NULL,
	create_date TIMESTAMP NOT NULL,
	PRIMARY KEY (id),
  CONSTRAINT fk_comment_user FOREIGN KEY (user_id) REFERENCES user (id) ON DELETE CASCADE,
  CONSTRAINT fk_comment_post FOREIGN KEY (post_id) REFERENCES post (id) ON DELETE CASCADE
);
