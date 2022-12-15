CREATE TABLE "sync" (
	"key"	TEXT,
	"value"	INTEGER NOT NULL DEFAULT 0,
	PRIMARY KEY("key")
);

CREATE TABLE "stat_documents" (
	"timestamp_day"	INTEGER NOT NULL DEFAULT 0,
	"count"	INTEGER NOT NULL DEFAULT 0,
	PRIMARY KEY("timestamp_day")
);

CREATE TABLE "stat_patients" (
	"timestamp_day"	INTEGER NOT NULL DEFAULT 0,
	"count"	INTEGER NOT NULL DEFAULT 0,
	PRIMARY KEY("timestamp_day")
);
