package config

var DBFile = "../scheduler.db"
var Port = "7540"
var SQLCreateTables = `CREATE TABLE IF NOT EXISTS scheduler (id INTEGER PRIMARY KEY AUTOINCREMENT, date TEXT, title TEXT, comment TEXT, repeat VARCHAR(128));`
var SQLCreateIndexes = `CREATE INDEX IF NOT EXISTS idx_date ON scheduler (date);`
var WEBDir = "./web"
