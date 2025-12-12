STATUS:
	go run main.go status

CREATE-SCRIPT-CREATE-TABLE:
	go run main.go new-script create --version-type TIME --description create_table

CREATE-SCRIPT-ALTER-SCHEMA:
	go run main.go new-script create --version-type TIME --description alter_schema

CREATE-SCRIPT-ALTER-SCHEMA-ALLOW-NO-CHANGE:
	go run main.go new-script create --version-type TIME --description alter_schema --allow-empty-script true

CREATE-SCRIPT-ALTER-COLUMN:
	go run main.go new-script create --version-type TIME --description alter_column

UPDATE-SCRIPT:
	go run main.go new-script update

MIGRATE-ALL:
	go run main.go migrate all

MIGRATE-PREVIEW-INC:
	go run main.go preview inc

MIGRATE-INC:
	go run main.go migrate inc

#MIGRATE-PREVIEW-DEC:
#	go run main.go preview dec

# DO NOT ENABLE THIS TO AVOID UNINTENDED DESTRUCTIVE OPERATIONS.
# 这里不打开是避免误操作导致问题
#MIGRATE-DEC:
#	go run main.go migrate dec
