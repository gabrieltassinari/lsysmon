CREATE DATABASE lsysmon;
\c lsysmon;
CREATE TABLE processes
(
	id serial NOT NULL PRIMARY KEY,
	processes_date TIMESTAMP NOT NULL,
	processes_stat jsonb NOT NULL
);
CREATE TABLE cpu
(
	id serial NOT NULL PRIMARY KEY,
	cpu_date TIMESTAMP NOT NULL,
	cpu_usage REAL NOT NULL
);
