ALTER TABLE `cybernoti`.`Websites` 
ADD INDEX `FK_USER_idx` (`user_id` ASC) VISIBLE;
;
ALTER TABLE `cybernoti`.`Websites` 
ADD CONSTRAINT `FK_USER`
  FOREIGN KEY (`user_id`)
  REFERENCES `cybernoti`.`Users` (`id`)
  ON DELETE NO ACTION
  ON UPDATE NO ACTION;

ALTER TABLE `cybernoti`.`Users` 
ADD COLUMN `deleted_at` DATETIME NULL DEFAULT NULL AFTER `updated_at`,
ADD INDEX `IDX_DELETE` (`deleted_at` ASC) VISIBLE;
;

CREATE TABLE `cybernoti`.`Notifications` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `total` INT NULL DEFAULT 1,
  `device_id` INT NULL,
  `tag_id` INT NULL,
  `created_at` DATETIME NULL DEFAULT now(),
  `updated_at` DATETIME NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `IDX_DEVICEID` (`device_id` ASC) VISIBLE,
  INDEX `IDX_TAGID` (`tag_id` ASC) VISIBLE);
