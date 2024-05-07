USE `temp_db`;
DROP TABLE IF EXISTS `temp_table`;
CREATE TABLE `temp_table` (
  `id` INT NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`id`),
  `field` VARCHAR(255) NOT NULL,
  `update_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
INSERT INTO temp_table (id, field)
VALUES (1, 'dummy field');
INSERT INTO temp_table (id, field)
VALUES (2, 'dummy field');
INSERT INTO temp_table (id, field)
VALUES (3, 'dummy field');
INSERT INTO temp_table (id, field)
VALUES (4, 'dummy field');
INSERT INTO temp_table (id, field)
VALUES (5, 'dummy field');
INSERT INTO temp_table (id, field)
VALUES (6, 'dummy field');
INSERT INTO temp_table (id, field)
VALUES (7, 'dummy field');
INSERT INTO temp_table (id, field)
VALUES (8, 'dummy field');
INSERT INTO temp_table (id, field)
VALUES (9, 'dummy field');
INSERT INTO temp_table (id, field)
VALUES (10, 'dummy field');


CREATE TABLE `favorites` (
  `user_id` int(11) NOT NULL,
  `tweet_id` int(11) NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`user_id`, `tweet_id`)
);
CREATE TABLE `feeds` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
);

CREATE TABLE tweets (
  id INT AUTO_INCREMENT PRIMARY KEY,
  user_id INT,
  text_content TEXT,
  images json,
  videos json,
    status int(11) NOT NULL DEFAULT '1',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
   -- Add other properties as needed
);
INSERT INTO tweets (user_id, text_content, images, videos)
VALUES (
    1,
    'This is a sample tweet about AI',
    '{"images":["ai1.jpg", "ai2.jpg"]}',
    '{"videos":["ai_video1.mp4"]}'
  ),
  (
    2,
    'This is a sample tweet about Machine Learning',
    '{"images":["ml1.jpg", "ml2.jpg"]}',
    '{"videos":["ml_video1.mp4"]}'
  ),
  (
    3,
    'This is a sample tweet about Data Science',
    '{"images":["ds1.jpg", "ds2.jpg"]}',
    '{"videos":["ds_video1.mp4"]}'
  ),
  (
    4,
    'This is a sample tweet about Python',
    '{"images":["python1.jpg", "python2.jpg"]}',
    '{"videos":["python_video1.mp4"]}'
  ),
  (
    5,
    'This is a sample tweet about SQL',
    '{"images":["sql1.jpg", "sql2.jpg"]}',
    '{"videos":["sql_video1.mp4"]}'
  );

CREATE TABLE `feeds_tweets` (
  `feed_id` int(11) NOT NULL,
  `tweet_id` int(11) NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`feed_id`, `tweet_id`)
);


CREATE TABLE `followers` (
  `follower_id` int(11) NOT NULL,
  `followee_id` int(11) NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`follower_id`, `followee_id`)
);


CREATE TABLE `images` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `url` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `height` int(11) DEFAULT NULL,
  `width` int(11) DEFAULT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
);


CREATE TABLE `users` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `email` varchar(50) NOT NULL,
  `fb_id` varchar(50) DEFAULT NULL,
  `gg_id` varchar(50) DEFAULT NULL,
  `password` varchar(50) NOT NULL,
  `salt` varchar(50) DEFAULT NULL,
  `last_name` varchar(50) NOT NULL,
  `first_name` varchar(50) NOT NULL,
  `dob` date DEFAULT NULL,
  `phone` varchar(20) DEFAULT NULL,
  `role` enum('user', 'admin', 'shipper') NOT NULL DEFAULT 'user',
  `avatar` json DEFAULT NULL,
  `status` int(11) NOT NULL DEFAULT '1',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `apple_id` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `email` (`email`)
);