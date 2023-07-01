CREATE DATABASE lsysmon;
\c lsysmon;
CREATE TABLE processes
(
	id serial NOT NULL PRIMARY KEY,
	processes_date TIMESTAMP NOT NULL,
	processes_stat json NOT NULL
);
