CREATE TABLE IF NOT EXISTS services (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  desc TEXT NOT NULL DEFAULT "",
  host TEXT NOT NULL, 
  path TEXT NOT NULL, 
  repo TEXT NOT NULL DEFAULT "",
  branch TEXT NOT NULL DEFAULT "master",
  hash TEXT NOT NULL DEFAULT "",
  start TEXT NOT NULL DEFAULT "node",
  args TEXT NOT NULL DEFAULT `[ "-e", "console.log('No Command Specified')" ]`,
  env TEXT NOT NULL DEFAULT "[]",
  install TEXT NOT NULL DEFAULT "node",
  installArgs TEXT NOT NULL DEFAULT `[ "-e", "console.log('No Install Specified')" ]`,
  installEnv TEXT NOT NULL DEFAULT "[]",
  enabled BOOLEAN NOT NULL DEFAULT 1,
  UNIQUE(host, path)
  UNIQUE(name)
);

CREATE TABLE IF NOT EXISTS admins (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user TEXT NOT NULL,
  serial TEXT NOT NULL,
  UNIQUE(serial)
);